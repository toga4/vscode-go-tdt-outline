// src/extension.ts
import * as vscode from "vscode";
import * as cp from "node:child_process";
import * as path from "node:path";
import * as fs from "node:fs/promises";
import { promisify } from "node:util";

// Type definition for JSON output from Go analysis tool
interface GoSymbol {
	name: string;
	detail: string;
	kind: number;
	range: {
		start: { line: number; character: number };
		end: { line: number; character: number };
	};
	children: GoSymbol[];
}

// Configuration interface
interface ExtensionConfig {
	timeout: number;
	maxFileSize: number;
	enableDebugLog: boolean;
}

// Error type for child process execution
interface ExecError extends Error {
	code?: string;
	killed?: boolean;
	signal?: string;
}

const execFileAsync = promisify(cp.execFile);

export async function activate(context: vscode.ExtensionContext) {
	// Create Output Channel
	const outputChannel = vscode.window.createOutputChannel("Go TDD Outline");
	context.subscriptions.push(outputChannel);

	outputChannel.appendLine("Go TDD Outline extension is being activated...");

	try {
		const goTestOutlineProvider = new GoTestOutlineProvider(
			context,
			outputChannel,
		);

		// Register DocumentSymbolProvider for Go language files
		const disposable = vscode.languages.registerDocumentSymbolProvider(
			{ language: "go", scheme: "file" },
			goTestOutlineProvider,
		);
		context.subscriptions.push(disposable);

		outputChannel.appendLine(
			"Go TDD Outline extension activated successfully.",
		);
	} catch (error) {
		const errorMessage = error instanceof Error ? error.message : String(error);
		outputChannel.appendLine(
			`Failed to activate Go TDD Outline: ${errorMessage}`,
		);
		vscode.window.showErrorMessage(
			`Go TDD Outline: Failed to activate extension: ${errorMessage}`,
		);
		throw error;
	}
}

class GoTestOutlineProvider implements vscode.DocumentSymbolProvider {
	private readonly parserPath: string;
	private readonly config: ExtensionConfig;
	private readonly outputChannel: vscode.OutputChannel;
	private parserExists = false;
	private readonly cache = new Map<
		string,
		{
			symbols: vscode.DocumentSymbol[];
			version: number;
			timestamp: number;
		}
	>();
	private readonly CACHE_EXPIRY_MS = 30000; // 30 seconds

	constructor(
		context: vscode.ExtensionContext,
		outputChannel: vscode.OutputChannel,
	) {
		this.outputChannel = outputChannel;

		// Load configuration
		this.config = this.loadConfiguration();

		// Get path to bundled Go analysis tool
		const parserFile =
			process.platform === "win32"
				? "go-outline-parser.exe"
				: "go-outline-parser";
		this.parserPath = path.join(
			context.extensionPath,
			"out",
			"parser",
			parserFile,
		);

		// Check if parser file exists asynchronously
		this.checkParserExists();
	}

	private loadConfiguration(): ExtensionConfig {
		const config = vscode.workspace.getConfiguration("goTddOutline");
		return {
			timeout: config.get<number>("timeout") ?? 10000,
			maxFileSize: config.get<number>("maxFileSize") ?? 1024 * 1024, // 1MB
			enableDebugLog: config.get<boolean>("enableDebugLog") ?? false,
		};
	}

	private async checkParserExists(): Promise<void> {
		try {
			await fs.access(this.parserPath);
			this.parserExists = true;
			this.log("Parser found successfully", "info");
		} catch {
			this.parserExists = false;
			vscode.window.showErrorMessage(
				`Go TDD Outline: Parser file not found. Please reinstall the extension.\nPath: ${this.parserPath}`,
			);
		}
	}

	async provideDocumentSymbols(
		document: vscode.TextDocument,
		token: vscode.CancellationToken,
	): Promise<vscode.DocumentSymbol[]> {
		// Early return if parser doesn't exist
		if (!this.parserExists) {
			this.log("Parser not available", "warn");
			return [];
		}

		// Check for cancellation
		if (token.isCancellationRequested) {
			return [];
		}

		// Check file size
		const fileStats = await fs.stat(document.fileName);
		if (fileStats.size > this.config.maxFileSize) {
			vscode.window.showWarningMessage(
				`Go TDD Outline: File too large (${Math.round(fileStats.size / 1024)}KB). Skipping analysis.`,
			);
			return [];
		}

		// Check cache
		const cacheKey = `${document.fileName}:${document.version}`;
		const cached = this.cache.get(cacheKey);
		if (cached && Date.now() - cached.timestamp < this.CACHE_EXPIRY_MS) {
			this.log("Using cached result", "info");
			return cached.symbols;
		}

		try {
			this.log(`Analyzing file: ${document.fileName}`, "info");

			const { stdout, stderr } = await execFileAsync(
				this.parserPath,
				[document.fileName],
				{
					timeout: this.config.timeout,
					signal: token.isCancellationRequested ? undefined : undefined,
				},
			);

			if (stderr?.trim()) {
				this.log(`Parser stderr: ${stderr}`, "warn");
			}

			const goSymbols: GoSymbol[] = JSON.parse(stdout);
			if (!goSymbols || !Array.isArray(goSymbols)) {
				this.log("No symbols found or invalid response format", "info");
				return [];
			}

			const vsCodeSymbols = this.convertToVSCodeSymbols(goSymbols);

			// Update cache
			this.cache.set(cacheKey, {
				symbols: vsCodeSymbols,
				version: document.version,
				timestamp: Date.now(),
			});

			// Clean old cache entries
			this.cleanCache();

			this.log(`Found ${vsCodeSymbols.length} test functions`, "info");
			return vsCodeSymbols;
		} catch (error) {
			return this.handleError(error as ExecError);
		}
	}

	private handleError(error: ExecError): vscode.DocumentSymbol[] {
		if (error.code === "ETIMEDOUT") {
			vscode.window.showErrorMessage(
				"Go TDD Outline: Parser timed out. File may be too large.",
			);
		} else if (error.code === "ENOENT") {
			vscode.window.showErrorMessage("Go TDD Outline: Parser file not found.");
			this.parserExists = false;
		} else if (error.name === "SyntaxError") {
			vscode.window.showErrorMessage(
				"Go TDD Outline: Failed to parse parser output. Please check Go file syntax.",
			);
			this.log(`JSON parse error: ${error.message}`, "error");
		} else {
			vscode.window.showErrorMessage(
				`Go TDD Outline: Error occurred during analysis: ${error.message}`,
			);
		}

		this.log(`Error executing parser: ${error}`, "error");
		return [];
	}

	private cleanCache(): void {
		const now = Date.now();
		for (const [key, value] of this.cache.entries()) {
			if (now - value.timestamp > this.CACHE_EXPIRY_MS) {
				this.cache.delete(key);
			}
		}
	}

	/**
	 * Recursively convert JSON objects from Go tool to VSCode DocumentSymbol[]
	 */
	private convertToVSCodeSymbols(
		goSymbols: GoSymbol[],
	): vscode.DocumentSymbol[] {
		return goSymbols.map((s) => {
			const range = new vscode.Range(
				new vscode.Position(s.range.start.line, s.range.start.character),
				new vscode.Position(s.range.end.line, s.range.end.character),
			);

			const symbol = new vscode.DocumentSymbol(
				s.name,
				s.detail,
				s.kind as vscode.SymbolKind, // SymbolKind numbers are aligned with Go side
				range,
				range, // Use same range for selectionRange
			);

			if (s.children && s.children.length > 0) {
				symbol.children = this.convertToVSCodeSymbols(s.children);
			}

			return symbol;
		});
	}

	private log(
		message: string,
		level: "info" | "error" | "warn" = "info",
	): void {
		if (!this.config.enableDebugLog && level === "info") {
			return;
		}

		const timestamp = new Date().toISOString();
		const logMessage = `[${timestamp}] [${level.toUpperCase()}] ${message}`;

		this.outputChannel.appendLine(logMessage);

		switch (level) {
			case "error":
				console.error(logMessage);
				break;
			case "warn":
				console.warn(logMessage);
				break;
			default:
				if (this.config.enableDebugLog) {
					console.log(logMessage);
				}
		}
	}
}

export function deactivate() {
	// Extension cleanup is handled automatically by VSCode
	// OutputChannel disposal is managed by context.subscriptions
}
