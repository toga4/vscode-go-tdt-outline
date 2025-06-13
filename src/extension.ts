// src/extension.ts
import * as vscode from "vscode";
import * as cp from "node:child_process";
import * as path from "node:path";
import * as fs from "node:fs";

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

// Output Channel for logging
let outputChannel: vscode.OutputChannel;

export function activate(context: vscode.ExtensionContext) {
	// Create Output Channel
	outputChannel = vscode.window.createOutputChannel("Go TDD Outline");
	context.subscriptions.push(outputChannel);

	outputChannel.appendLine("Go TDD Outline extension is being activated...");

	try {
		const goTestOutlineProvider = new GoTestOutlineProvider(context);

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
	private parserPath: string;
	private parserExists = false;

	constructor(context: vscode.ExtensionContext) {
		// Get path to bundled Go analysis tool
		// Change executable filename based on OS
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

		// Check if parser file exists
		this.parserExists = fs.existsSync(this.parserPath);
		if (!this.parserExists) {
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
			return [];
		}

		// Check for cancellation
		if (token.isCancellationRequested) {
			return [];
		}

		// Execute Go analysis tool
		return new Promise((resolve, reject) => {
			// execFile is preferable for security
			const child = cp.execFile(
				this.parserPath,
				[document.fileName],
				{ timeout: 10000 }, // 10 second timeout
				(err, stdout, stderr) => {
					// Check for cancellation
					if (token.isCancellationRequested) {
						return resolve([]);
					}

					if (err) {
						// Display appropriate error message based on error type
						if (err.code === "ETIMEDOUT") {
							vscode.window.showErrorMessage(
								"Go TDD Outline: Parser timed out. File may be too large.",
							);
						} else if (err.code === "ENOENT") {
							vscode.window.showErrorMessage(
								"Go TDD Outline: Parser file not found.",
							);
							this.parserExists = false;
						} else {
							vscode.window.showErrorMessage(
								`Go TDD Outline: Error occurred during analysis: ${err.message}`,
							);
						}
						log(`Error executing go-outline-parser: ${err}`, "error");
						return resolve([]); // Return empty array on error
					}
					if (stderr?.trim()) {
						log(`go-outline-parser stderr: ${stderr}`, "warn");
					}

					try {
						// Parse output JSON
						const goSymbols: GoSymbol[] = JSON.parse(stdout);
						if (!goSymbols) {
							return resolve([]);
						}

						// Convert Go symbol information to VS Code symbol information
						const vsCodeSymbols = this.convertToVSCodeSymbols(goSymbols);
						resolve(vsCodeSymbols);
					} catch (e) {
						log(`Error parsing JSON from go-outline-parser: ${e}`, "error");
						log(`Raw output: ${stdout}`, "error");
						// Notify user of JSON parse error
						vscode.window.showErrorMessage(
							"Go TDD Outline: Failed to parse parser output. Please check Go file syntax.",
						);
						return resolve([]); // Return empty array on error
					}
				},
			);

			// Terminate process on cancellation
			token.onCancellationRequested(() => {
				if (!child?.killed) {
					child.kill();
				}
			});
		});
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
}

// Helper function for log output
function log(message: string, level: "info" | "error" | "warn" = "info") {
	const timestamp = new Date().toISOString();
	const logMessage = `[${timestamp}] [${level.toUpperCase()}] ${message}`;

	if (outputChannel) {
		outputChannel.appendLine(logMessage);
	}

	switch (level) {
		case "error":
			console.error(logMessage);
			break;
		case "warn":
			console.warn(logMessage);
			break;
		default:
			console.log(logMessage);
	}
}

export function deactivate() {
	if (outputChannel) {
		outputChannel.appendLine("Go TDD Outline extension is being deactivated.");
		outputChannel.dispose();
	}
}
