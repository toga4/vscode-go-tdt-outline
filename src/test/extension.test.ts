import assert from "node:assert";
import * as fs from "node:fs";
import * as path from "node:path";
import * as vscode from "vscode";
import type { ExtensionApi, GoTddOutlineProvider } from "../extension";
import { expectMatchSnapshot } from "./snapshot-helper";

suite("Snapshot Tests", () => {
  let provider: GoTddOutlineProvider;
  const testFileDir = path.join(__dirname, "testdata");

  suiteSetup(async () => {
    const extension = vscode.extensions.getExtension<ExtensionApi>("toga4.go-tdt-outline");
    assert.ok(extension, "Extension not found");

    const api = await extension.activate();
    provider = api.documentSymbolProvider;
  });

  for (const file of fs.readdirSync(testFileDir)) {
    test(file, async () => {
      const testFilePath = path.join(testFileDir, file);
      const document = await vscode.workspace.openTextDocument(testFilePath);
      const tokenSource = new vscode.CancellationTokenSource();
      const token = tokenSource.token;

      const symbols = await provider.provideDocumentSymbols(document, token);

      expectMatchSnapshot(file, symbols);
    });
  }
});
