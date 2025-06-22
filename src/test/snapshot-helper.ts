import * as fs from "node:fs";
import * as path from "node:path";
import type * as vscode from "vscode";

/**
 * スナップショットファイルのパスを生成
 * srcディレクトリ以下のsnapshotsディレクトリに保存
 */
export function getSnapshotPath(testName: string): string {
  const snapshotsDir = path.resolve("src/test/snapshots");
  if (!fs.existsSync(snapshotsDir)) {
    fs.mkdirSync(snapshotsDir, { recursive: true });
  }
  return path.join(snapshotsDir, `${testName}.json`);
}

/**
 * スナップショットテストを実行
 * updateSnapshots が true の場合、期待値を更新する
 */
export function expectMatchSnapshot(testName: string, actualSymbols: vscode.DocumentSymbol[]): void {
  const actualJson = JSON.stringify(actualSymbols, null, 2);
  const snapshotPath = getSnapshotPath(testName);

  if (process.env.UPDATE_SNAPSHOTS === "true") {
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
