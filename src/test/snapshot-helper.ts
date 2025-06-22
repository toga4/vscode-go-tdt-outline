import * as fs from "node:fs";
import * as path from "node:path";
import { diff_match_patch } from "diff-match-patch";
import * as pc from "picocolors";
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
    throw new Error(`Snapshot mismatch for: ${testName}.${pc.reset("")}

${diff(expectedJson, actualJson)}
Run with UPDATE_SNAPSHOTS=true to update snapshots.`);
  }
}

function diff(expectedJson: string, actualJson: string): string {
  const dmp = new diff_match_patch();
  const a = dmp.diff_linesToChars_(expectedJson, actualJson);
  const diffs = dmp.diff_main(a.chars1, a.chars2);
  dmp.diff_charsToLines_(diffs, a.lineArray);
  dmp.diff_cleanupSemantic(diffs);

  const lines = [];
  for (const [n, text] of diffs) {
    for (const line of text.split("\n").slice(0, -1)) {
      if (n < 0) {
        lines.push(pc.red(`- ${line}`));
      } else if (n > 0) {
        lines.push(pc.green(`+ ${line}`));
      } else {
        lines.push(`  ${line}`);
      }
    }
  }
  return lines.join("\n");
}
