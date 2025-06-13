// src/extension.ts
import * as vscode from "vscode";
import * as cp from "node:child_process";
import * as path from "node:path";

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

export function activate(context: vscode.ExtensionContext) {
	console.log(
		'Congratulations, your extension "go-test-outline" is now active!',
	);

	const goTestOutlineProvider = new GoTestOutlineProvider(context);

	// Go言語ファイルに対してDocumentSymbolProviderを登録する
	context.subscriptions.push(
		vscode.languages.registerDocumentSymbolProvider(
			{ language: "go", scheme: "file" },
			goTestOutlineProvider,
		),
	);
}

class GoTestOutlineProvider implements vscode.DocumentSymbolProvider {
	private parserPath: string;

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
	}

	async provideDocumentSymbols(
		document: vscode.TextDocument,
		token: vscode.CancellationToken,
	): Promise<vscode.DocumentSymbol[]> {
		// Go製解析ツールを実行
		return new Promise((resolve, reject) => {
			// execFileのほうがセキュリティ上好ましい
			cp.execFile(
				this.parserPath,
				[document.fileName],
				(err, stdout, stderr) => {
					if (err) {
						console.error(`Error executing go-outline-parser: ${err}`);
						return reject(err);
					}
					if (stderr) {
						console.error(`go-outline-parser stderr: ${stderr}`);
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
						console.error(`Error parsing JSON from go-outline-parser: ${e}`);
						console.error(`Raw output: ${stdout}`);
						return reject(e);
					}
				},
			);
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

export function deactivate() {}
