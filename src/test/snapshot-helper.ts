import * as fs from "node:fs";
import * as path from "node:path";
import type * as vscode from "vscode";

/**
 * DocumentSymbolをスナップショット用にシリアライズ可能な形式に変換
 */
export interface SerializableSymbol {
  name: string;
  detail: string;
  kind: number;
  range: {
    start: { line: number; character: number };
    end: { line: number; character: number };
  };
  selectionRange: {
    start: { line: number; character: number };
    end: { line: number; character: number };
  };
  children: SerializableSymbol[];
}

/**
 * VSCodeのDocumentSymbolをシリアライズ可能な形式に変換
 */
export function serializeSymbols(symbols: vscode.DocumentSymbol[]): SerializableSymbol[] {
  return symbols.map((symbol) => ({
    name: symbol.name,
    detail: symbol.detail,
    kind: symbol.kind,
    range: {
      start: {
        line: symbol.range.start.line,
        character: symbol.range.start.character,
      },
      end: {
        line: symbol.range.end.line,
        character: symbol.range.end.character,
      },
    },
    selectionRange: {
      start: {
        line: symbol.selectionRange.start.line,
        character: symbol.selectionRange.start.character,
      },
      end: {
        line: symbol.selectionRange.end.line,
        character: symbol.selectionRange.end.character,
      },
    },
    children: serializeSymbols(symbol.children),
  }));
}

/**
 * スナップショットファイルのパスを生成
 * srcディレクトリ以下のsnapshotsディレクトリに保存
 */
export function getSnapshotPath(testName: string): string {
  // __dirnameは out/test を指すので、src/test に変換する
  const srcTestDir = __dirname.replace(/^(.*)\/out\/test/, "$1/src/test");
  const snapshotsDir = path.join(srcTestDir, "snapshots");
  if (!fs.existsSync(snapshotsDir)) {
    fs.mkdirSync(snapshotsDir, { recursive: true });
  }
  return path.join(snapshotsDir, `${testName}.json`);
}

/**
 * スナップショットテストを実行
 * updateSnapshots が true の場合、期待値を更新する
 */
export function expectMatchSnapshot(
  testName: string,
  actualSymbols: vscode.DocumentSymbol[],
  updateSnapshots = false,
): void {
  const actualJson = JSON.stringify(serializeSymbols(actualSymbols), null, 2);
  const snapshotPath = getSnapshotPath(testName);

  if (updateSnapshots) {
    fs.writeFileSync(snapshotPath, actualJson, "utf8");
    console.log(`📸 Updated snapshot for: ${testName}`);
    return;
  }

  if (!fs.existsSync(snapshotPath)) {
    // スナップショットが存在しない場合は新規作成
    fs.writeFileSync(snapshotPath, actualJson, "utf8");
    console.log(`📸 Created new snapshot for: ${testName}`);
    return;
  }

  const expectedJson = fs.readFileSync(snapshotPath, "utf8");

  // 単純な文字列比較
  if (actualJson !== expectedJson) {
    throw new Error(`Snapshot mismatch for: ${testName}

Expected:
${expectedJson}

Actual:
${actualJson}

Run with UPDATE_SNAPSHOTS=true to update snapshots.`);
  }
}

/**
 * 環境変数からスナップショット更新フラグを取得
 */
export function shouldUpdateSnapshots(): boolean {
  return process.env.UPDATE_SNAPSHOTS === "true";
}
