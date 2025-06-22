import * as cp from "node:child_process";
import * as fs from "node:fs";
import * as path from "node:path";
import { promisify } from "node:util";
import * as vscode from "vscode";

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
}

const execFileAsync = promisify(cp.execFile);

export interface ExtensionApi {
  documentSymbolProvider: GoTddOutlineProvider;
}

export async function activate(context: vscode.ExtensionContext): Promise<ExtensionApi> {
  console.log("activate", context);

  // Create Output Channel
  const outputChannel = vscode.window.createOutputChannel("Go TDD Outline");
  context.subscriptions.push(outputChannel);

  outputChannel.appendLine("Go TDD Outline extension is being activated...");

  try {
    const goTddOutlineProvider = new GoTddOutlineProvider(context, outputChannel);

    // Register DocumentSymbolProvider for Go language files
    const disposable = vscode.languages.registerDocumentSymbolProvider(
      { language: "go", scheme: "file" },
      goTddOutlineProvider,
    );
    context.subscriptions.push(disposable);

    outputChannel.appendLine("Go TDD Outline extension activated successfully.");

    return {
      documentSymbolProvider: goTddOutlineProvider,
    };
  } catch (error) {
    const errorMessage = error instanceof Error ? error.message : String(error);
    outputChannel.appendLine(`Failed to activate Go TDD Outline: ${errorMessage}`);
    throw error;
  }
}

export class GoTddOutlineProvider implements vscode.DocumentSymbolProvider {
  private readonly parserPath: string;
  private readonly config: ExtensionConfig;
  private readonly outputChannel: vscode.OutputChannel;
  private parserExists = false;

  constructor(context: vscode.ExtensionContext, outputChannel: vscode.OutputChannel) {
    this.outputChannel = outputChannel;

    // Load configuration
    this.config = this.loadConfiguration();

    // Get path to bundled Go analysis tool
    const parserFile = process.platform === "win32" ? "parser.exe" : "parser";
    this.parserPath = path.join(context.extensionPath, "bin", parserFile);

    // Check if parser file exists
    this.parserExists = fs.existsSync(this.parserPath);
    if (!this.parserExists) {
      this.outputChannel.appendLine(
        `Error: Parser file not found. Please reinstall the extension. Path: ${this.parserPath}`,
      );
    }
  }

  private loadConfiguration(): ExtensionConfig {
    const config = vscode.workspace.getConfiguration("go-tdt-outline");
    return {
      timeout: config.get<number>("timeout") ?? 10000,
      maxFileSize: config.get<number>("maxFileSize") ?? 1024 * 1024, // 1MB
    };
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

    // Check file size
    const fileStats = fs.statSync(document.fileName);
    if (fileStats.size > this.config.maxFileSize) {
      this.outputChannel.appendLine(
        `Warning: File too large (${Math.round(fileStats.size / 1024)}KB). Skipping analysis.`,
      );
      return [];
    }

    try {
      const { stdout, stderr } = await execFileAsync(this.parserPath, [document.fileName], {
        timeout: this.config.timeout,
      });

      if (stderr?.trim()) {
        this.outputChannel.appendLine(`Parser stderr: ${stderr}`);
      }

      const goSymbols: GoSymbol[] = JSON.parse(stdout);
      if (!goSymbols || !Array.isArray(goSymbols)) {
        return [];
      }

      const vsCodeSymbols = this.convertToVSCodeSymbols(goSymbols);
      return vsCodeSymbols;
    } catch (error) {
      this.outputChannel.appendLine(`Error: ${error}`);
      return [];
    }
  }

  /**
   * Recursively convert JSON objects from Go tool to VSCode DocumentSymbol[]
   */
  private convertToVSCodeSymbols(goSymbols: GoSymbol[]): vscode.DocumentSymbol[] {
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

export function deactivate() {
  // Extension cleanup is handled automatically by VSCode
  // OutputChannel disposal is managed by context.subscriptions
}
