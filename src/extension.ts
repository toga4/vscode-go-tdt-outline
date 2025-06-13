// src/extension.ts
import * as vscode from "vscode";
import * as cp from "node:child_process";
import * as path from "node:path";
import * as fs from "node:fs";

// Goの解析ツールが出力するJSONの型定義
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
	// Output Channelを作成
	outputChannel = vscode.window.createOutputChannel("Go TDD Outline");
	context.subscriptions.push(outputChannel);

	outputChannel.appendLine("Go TDD Outline extension is being activated...");

	try {
		const goTestOutlineProvider = new GoTestOutlineProvider(context);

		// Go言語ファイルに対してDocumentSymbolProviderを登録する
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
			`Go TDD Outline: 拡張機能の有効化に失敗しました: ${errorMessage}`,
		);
		throw error;
	}
}

class GoTestOutlineProvider implements vscode.DocumentSymbolProvider {
	private parserPath: string;
	private parserExists = false;

	constructor(context: vscode.ExtensionContext) {
		// 拡張機能にバンドルされたGo製解析ツールのパスを取得
		// OSによって実行ファイル名を変える
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

		// パーサーファイルの存在確認
		this.parserExists = fs.existsSync(this.parserPath);
		if (!this.parserExists) {
			vscode.window.showErrorMessage(
				`Go TDD Outline: パーサーファイルが見つかりません。拡張機能を再インストールしてください。\nパス: ${this.parserPath}`,
			);
		}
	}

	async provideDocumentSymbols(
		document: vscode.TextDocument,
		token: vscode.CancellationToken,
	): Promise<vscode.DocumentSymbol[]> {
		// パーサーが存在しない場合は早期リターン
		if (!this.parserExists) {
			return [];
		}

		// キャンセレーションチェック
		if (token.isCancellationRequested) {
			return [];
		}

		// Go製解析ツールを実行
		return new Promise((resolve, reject) => {
			// execFileのほうがセキュリティ上好ましい
			const child = cp.execFile(
				this.parserPath,
				[document.fileName],
				{ timeout: 10000 }, // 10秒のタイムアウト
				(err, stdout, stderr) => {
					// キャンセレーションチェック
					if (token.isCancellationRequested) {
						return resolve([]);
					}

					if (err) {
						// エラーの種類に応じて適切なメッセージを表示
						if (err.code === "ETIMEDOUT") {
							vscode.window.showErrorMessage(
								"Go TDD Outline: パーサーがタイムアウトしました。ファイルが大きすぎる可能性があります。",
							);
						} else if (err.code === "ENOENT") {
							vscode.window.showErrorMessage(
								"Go TDD Outline: パーサーファイルが見つかりません。",
							);
							this.parserExists = false;
						} else {
							vscode.window.showErrorMessage(
								`Go TDD Outline: 解析中にエラーが発生しました: ${err.message}`,
							);
						}
						log(`Error executing go-outline-parser: ${err}`, "error");
						return resolve([]); // エラー時は空の配列を返す
					}
					if (stderr?.trim()) {
						log(`go-outline-parser stderr: ${stderr}`, "warn");
					}

					try {
						// 出力されたJSONをパース
						const goSymbols: GoSymbol[] = JSON.parse(stdout);
						if (!goSymbols) {
							return resolve([]);
						}

						// Goのシンボル情報をVS Codeのシンボル情報に変換
						const vsCodeSymbols = this.convertToVSCodeSymbols(goSymbols);
						resolve(vsCodeSymbols);
					} catch (e) {
						log(`Error parsing JSON from go-outline-parser: ${e}`, "error");
						log(`Raw output: ${stdout}`, "error");
						// JSONパースエラーの場合もユーザーに通知
						vscode.window.showErrorMessage(
							"Go TDD Outline: パーサーからの出力を解析できませんでした。Goファイルの構文を確認してください。",
						);
						return resolve([]); // エラー時は空の配列を返す
					}
				},
			);

			// キャンセレーション時にプロセスを終了
			token.onCancellationRequested(() => {
				if (!child?.killed) {
					child.kill();
				}
			});
		});
	}

	/**
	 * Go製ツールからのJSONオブジェクトをVSCodeのDocumentSymbol[]に再帰的に変換する
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
				s.kind as vscode.SymbolKind, // Go側でSymbolKindの番号を合わせている
				range,
				range, // selectionRangeも同じにしておく
			);

			if (s.children && s.children.length > 0) {
				symbol.children = this.convertToVSCodeSymbols(s.children);
			}

			return symbol;
		});
	}
}

// ログ出力用のヘルパー関数
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
