import assert from "node:assert";
import * as fs from "node:fs";
import * as path from "node:path";
import * as vscode from "vscode";
import type { ExtensionApi, GoTddOutlineProvider } from "../extension";
import { expectMatchSnapshot } from "./snapshot-helper";

suite("Snapshot Tests", () => {
  let provider: GoTddOutlineProvider;
  const fixturesDir = path.join(__dirname, "fixtures");

  suiteSetup(async () => {
    const extension = vscode.extensions.getExtension<ExtensionApi>("toga4.go-tdt-outline");
    assert.ok(extension, "Extension not found");

    const api = await extension.activate();
    provider = api.documentSymbolProvider;
  });

  test("basic_table_test.go - full structure snapshot", async () => {
    const testFilePath = path.join(fixturesDir, "basic_table_test.go");
    const document = await vscode.workspace.openTextDocument(testFilePath);
    const tokenSource = new vscode.CancellationTokenSource();
    const token = tokenSource.token;

    const symbols = await provider.provideDocumentSymbols(document, token);

    expectMatchSnapshot("basic_table_test", symbols);
  });

  test("multiple_functions_test.go - full structure snapshot", async () => {
    const testFilePath = path.join(fixturesDir, "multiple_functions_test.go");
    const document = await vscode.workspace.openTextDocument(testFilePath);
    const tokenSource = new vscode.CancellationTokenSource();
    const token = tokenSource.token;

    const symbols = await provider.provideDocumentSymbols(document, token);

    expectMatchSnapshot("multiple_functions_test", symbols);
  });

  test("case_insensitive_test.go - full structure snapshot", async () => {
    const testFilePath = path.join(fixturesDir, "case_insensitive_test.go");
    const document = await vscode.workspace.openTextDocument(testFilePath);
    const tokenSource = new vscode.CancellationTokenSource();
    const token = tokenSource.token;

    const symbols = await provider.provideDocumentSymbols(document, token);

    expectMatchSnapshot("case_insensitive_test", symbols);
  });

  test("various_fields_test.go - full structure snapshot", async () => {
    const testFilePath = path.join(fixturesDir, "various_fields_test.go");

    if (!fs.existsSync(testFilePath)) {
      return; // Skip
    }

    const document = await vscode.workspace.openTextDocument(testFilePath);
    const tokenSource = new vscode.CancellationTokenSource();
    const token = tokenSource.token;

    const symbols = await provider.provideDocumentSymbols(document, token);

    expectMatchSnapshot("various_fields_test", symbols);
  });

  test("typed_test_cases.go - full structure snapshot", async () => {
    const testFilePath = path.join(fixturesDir, "typed_test_cases.go");

    if (!fs.existsSync(testFilePath)) {
      return; // Skip
    }

    const document = await vscode.workspace.openTextDocument(testFilePath);
    const tokenSource = new vscode.CancellationTokenSource();
    const token = tokenSource.token;

    const symbols = await provider.provideDocumentSymbols(document, token);

    expectMatchSnapshot("typed_test_cases", symbols);
  });

  test("map_test_cases.go - full structure snapshot", async () => {
    const testFilePath = path.join(fixturesDir, "map_test_cases.go");

    if (!fs.existsSync(testFilePath)) {
      return; // Skip
    }

    const document = await vscode.workspace.openTextDocument(testFilePath);
    const tokenSource = new vscode.CancellationTokenSource();
    const token = tokenSource.token;

    const symbols = await provider.provideDocumentSymbols(document, token);

    expectMatchSnapshot("map_test_cases", symbols);
  });

  test("backtick_strings_test.go - full structure snapshot", async () => {
    const testFilePath = path.join(fixturesDir, "backtick_strings_test.go");

    if (!fs.existsSync(testFilePath)) {
      return; // Skip
    }

    const document = await vscode.workspace.openTextDocument(testFilePath);
    const tokenSource = new vscode.CancellationTokenSource();
    const token = tokenSource.token;

    const symbols = await provider.provideDocumentSymbols(document, token);

    expectMatchSnapshot("backtick_strings_test", symbols);
  });

  test("multiple_tables_test.go - full structure snapshot", async () => {
    const testFilePath = path.join(fixturesDir, "multiple_tables_test.go");

    if (!fs.existsSync(testFilePath)) {
      return; // Skip
    }

    const document = await vscode.workspace.openTextDocument(testFilePath);
    const tokenSource = new vscode.CancellationTokenSource();
    const token = tokenSource.token;

    const symbols = await provider.provideDocumentSymbols(document, token);

    expectMatchSnapshot("multiple_tables_test", symbols);
  });

  test("no_name_field_test.go - full structure snapshot", async () => {
    const testFilePath = path.join(fixturesDir, "no_name_field_test.go");

    if (!fs.existsSync(testFilePath)) {
      return; // Skip
    }

    const document = await vscode.workspace.openTextDocument(testFilePath);
    const tokenSource = new vscode.CancellationTokenSource();
    const token = tokenSource.token;

    const symbols = await provider.provideDocumentSymbols(document, token);

    expectMatchSnapshot("no_name_field_test", symbols);
  });

  test("non_test_functions.go - full structure snapshot", async () => {
    const testFilePath = path.join(fixturesDir, "non_test_functions.go");

    if (!fs.existsSync(testFilePath)) {
      return; // Skip
    }

    const document = await vscode.workspace.openTextDocument(testFilePath);
    const tokenSource = new vscode.CancellationTokenSource();
    const token = tokenSource.token;

    const symbols = await provider.provideDocumentSymbols(document, token);

    expectMatchSnapshot("non_test_functions", symbols);
  });
});
