import * as fs from "node:fs";
import * as path from "node:path";
import * as vscode from "vscode";
import { GoTddOutlineProvider } from "../extension";
import { expectMatchSnapshot, shouldUpdateSnapshots } from "./snapshot-helper";

suite("Snapshot Tests", () => {
  let provider: GoTddOutlineProvider;
  let extensionContext: Partial<vscode.ExtensionContext>;
  let outputChannel: vscode.OutputChannel;
  const fixturesDir = path.join(__dirname, "fixtures");
  const updateSnapshots = shouldUpdateSnapshots();

  suiteSetup(() => {
    // 実際のExtensionContextを模擬
    extensionContext = {
      extensionPath: path.join(__dirname, "../.."),
      subscriptions: [],
      workspaceState: {
        get: () => undefined,
        update: () => Promise.resolve(),
        keys: () => [],
      },
      globalState: {
        get: () => undefined,
        update: () => Promise.resolve(),
        keys: () => [],
        setKeysForSync: () => {},
      },
    };

    outputChannel = vscode.window.createOutputChannel("Snapshot Test");

    // 実際のGoTddOutlineProviderを作成
    provider = new GoTddOutlineProvider(extensionContext as vscode.ExtensionContext, outputChannel);
  });

  suiteTeardown(() => {
    outputChannel.dispose();
  });

  test("basic_table_test.go - full structure snapshot", async () => {
    const testFilePath = path.join(fixturesDir, "basic_table_test.go");
    const document = await vscode.workspace.openTextDocument(testFilePath);
    const tokenSource = new vscode.CancellationTokenSource();
    const token = tokenSource.token;

    const symbols = await provider.provideDocumentSymbols(document, token);

    expectMatchSnapshot("basic_table_test", symbols, updateSnapshots);
  });

  test("multiple_functions_test.go - full structure snapshot", async () => {
    const testFilePath = path.join(fixturesDir, "multiple_functions_test.go");
    const document = await vscode.workspace.openTextDocument(testFilePath);
    const tokenSource = new vscode.CancellationTokenSource();
    const token = tokenSource.token;

    const symbols = await provider.provideDocumentSymbols(document, token);

    expectMatchSnapshot("multiple_functions_test", symbols, updateSnapshots);
  });

  test("case_insensitive_test.go - full structure snapshot", async () => {
    const testFilePath = path.join(fixturesDir, "case_insensitive_test.go");
    const document = await vscode.workspace.openTextDocument(testFilePath);
    const tokenSource = new vscode.CancellationTokenSource();
    const token = tokenSource.token;

    const symbols = await provider.provideDocumentSymbols(document, token);

    expectMatchSnapshot("case_insensitive_test", symbols, updateSnapshots);
  });

  test("various_fields_test.go - full structure snapshot", async () => {
    const testFilePath = path.join(fixturesDir, "various_fields_test.go");

    if (!fs.existsSync(testFilePath)) {
      return; // スキップ
    }

    const document = await vscode.workspace.openTextDocument(testFilePath);
    const tokenSource = new vscode.CancellationTokenSource();
    const token = tokenSource.token;

    const symbols = await provider.provideDocumentSymbols(document, token);

    expectMatchSnapshot("various_fields_test", symbols, updateSnapshots);
  });

  test("typed_test_cases.go - full structure snapshot", async () => {
    const testFilePath = path.join(fixturesDir, "typed_test_cases.go");

    if (!fs.existsSync(testFilePath)) {
      return; // スキップ
    }

    const document = await vscode.workspace.openTextDocument(testFilePath);
    const tokenSource = new vscode.CancellationTokenSource();
    const token = tokenSource.token;

    const symbols = await provider.provideDocumentSymbols(document, token);

    expectMatchSnapshot("typed_test_cases", symbols, updateSnapshots);
  });

  test("map_test_cases.go - full structure snapshot", async () => {
    const testFilePath = path.join(fixturesDir, "map_test_cases.go");

    if (!fs.existsSync(testFilePath)) {
      return; // スキップ
    }

    const document = await vscode.workspace.openTextDocument(testFilePath);
    const tokenSource = new vscode.CancellationTokenSource();
    const token = tokenSource.token;

    const symbols = await provider.provideDocumentSymbols(document, token);

    expectMatchSnapshot("map_test_cases", symbols, updateSnapshots);
  });

  test("backtick_strings_test.go - full structure snapshot", async () => {
    const testFilePath = path.join(fixturesDir, "backtick_strings_test.go");

    if (!fs.existsSync(testFilePath)) {
      return; // スキップ
    }

    const document = await vscode.workspace.openTextDocument(testFilePath);
    const tokenSource = new vscode.CancellationTokenSource();
    const token = tokenSource.token;

    const symbols = await provider.provideDocumentSymbols(document, token);

    expectMatchSnapshot("backtick_strings_test", symbols, updateSnapshots);
  });

  test("multiple_tables_test.go - full structure snapshot", async () => {
    const testFilePath = path.join(fixturesDir, "multiple_tables_test.go");

    if (!fs.existsSync(testFilePath)) {
      return; // スキップ
    }

    const document = await vscode.workspace.openTextDocument(testFilePath);
    const tokenSource = new vscode.CancellationTokenSource();
    const token = tokenSource.token;

    const symbols = await provider.provideDocumentSymbols(document, token);

    expectMatchSnapshot("multiple_tables_test", symbols, updateSnapshots);
  });

  test("no_name_field_test.go - full structure snapshot", async () => {
    const testFilePath = path.join(fixturesDir, "no_name_field_test.go");

    if (!fs.existsSync(testFilePath)) {
      return; // スキップ
    }

    const document = await vscode.workspace.openTextDocument(testFilePath);
    const tokenSource = new vscode.CancellationTokenSource();
    const token = tokenSource.token;

    const symbols = await provider.provideDocumentSymbols(document, token);

    expectMatchSnapshot("no_name_field_test", symbols, updateSnapshots);
  });

  test("non_test_functions.go - full structure snapshot", async () => {
    const testFilePath = path.join(fixturesDir, "non_test_functions.go");

    if (!fs.existsSync(testFilePath)) {
      return; // スキップ
    }

    const document = await vscode.workspace.openTextDocument(testFilePath);
    const tokenSource = new vscode.CancellationTokenSource();
    const token = tokenSource.token;

    const symbols = await provider.provideDocumentSymbols(document, token);

    expectMatchSnapshot("non_test_functions", symbols, updateSnapshots);
  });
});
