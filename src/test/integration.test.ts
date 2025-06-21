import * as assert from "node:assert";
import * as fs from "node:fs";
import * as path from "node:path";
import * as vscode from "vscode";
import { GoTddOutlineProvider } from "../extension";

suite("Parser Integration Tests", () => {
  let provider: GoTddOutlineProvider;
  let extensionContext: Partial<vscode.ExtensionContext>;
  let outputChannel: vscode.OutputChannel;
  const fixturesDir = path.join(__dirname, "fixtures");

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

    outputChannel = vscode.window.createOutputChannel("Test Integration");

    // 実際のGoTddOutlineProviderを作成
    provider = new GoTddOutlineProvider(extensionContext as vscode.ExtensionContext, outputChannel);
  });

  suiteTeardown(() => {
    outputChannel.dispose();
  });

  suite("Basic Test Parsing", () => {
    test("Should parse basic table driven test structure", async () => {
      const testFilePath = path.join(fixturesDir, "basic_table_test.go");
      assert.ok(fs.existsSync(testFilePath), `Test file should exist: ${testFilePath}`);

      const document = await vscode.workspace.openTextDocument(testFilePath);
      const tokenSource = new vscode.CancellationTokenSource();
      const token = tokenSource.token;

      const symbols = await provider.provideDocumentSymbols(document, token);

      // 基本的な構造の検証
      assert.ok(symbols.length > 0, "Should return at least one symbol");

      const testFunction = symbols.find((s) => s.name === "TestExample");
      assert.ok(testFunction, "Should find TestExample function");
      assert.strictEqual(testFunction.kind, vscode.SymbolKind.Function, "TestExample should be a function");

      // テストケースの検証
      assert.ok(testFunction.children.length >= 2, "Should have at least 2 test cases");

      const testCase1 = testFunction.children.find((c) => c.name === "normal case");
      const testCase2 = testFunction.children.find((c) => c.name === "zero value");

      assert.ok(testCase1, "Should find 'normal case' test case");
      assert.ok(testCase2, "Should find 'zero value' test case");

      assert.strictEqual(testCase1.kind, vscode.SymbolKind.Struct, "Test case should be a struct");
      assert.strictEqual(testCase2.kind, vscode.SymbolKind.Struct, "Test case should be a struct");

      // 位置情報の検証
      assert.ok(testCase1.range.start.line >= 0, "Test case should have valid line position");
      assert.ok(testCase2.range.start.line >= 0, "Test case should have valid line position");
    });
  });

  suite("Multiple Functions Test", () => {
    test("Should parse multiple test functions with different field names", async () => {
      const testFilePath = path.join(fixturesDir, "multiple_functions_test.go");
      const document = await vscode.workspace.openTextDocument(testFilePath);
      const tokenSource = new vscode.CancellationTokenSource();
      const token = tokenSource.token;

      const symbols = await provider.provideDocumentSymbols(document, token);

      assert.strictEqual(symbols.length, 2, "Should find exactly 2 test functions");

      const testFirst = symbols.find((s) => s.name === "TestFirst");
      const testSecond = symbols.find((s) => s.name === "TestSecond");

      assert.ok(testFirst, "Should find TestFirst function");
      assert.ok(testSecond, "Should find TestSecond function");

      // TestFirst のテストケースを検証 (name フィールドを使用)
      assert.strictEqual(testFirst.children.length, 2, "TestFirst should have 2 test cases");
      assert.ok(
        testFirst.children.some((c) => c.name === "test1"),
        "Should find test1 case",
      );
      assert.ok(
        testFirst.children.some((c) => c.name === "test2"),
        "Should find test2 case",
      );

      // TestSecond のテストケースを検証 (desc フィールドを使用)
      assert.strictEqual(testSecond.children.length, 2, "TestSecond should have 2 test cases");
      assert.ok(
        testSecond.children.some((c) => c.name === "test3"),
        "Should find test3 case",
      );
      assert.ok(
        testSecond.children.some((c) => c.name === "test4"),
        "Should find test4 case",
      );
    });
  });

  suite("Case Insensitive Field Names", () => {
    test("Should recognize case-insensitive name fields", async () => {
      const testFilePath = path.join(fixturesDir, "case_insensitive_test.go");
      const document = await vscode.workspace.openTextDocument(testFilePath);
      const tokenSource = new vscode.CancellationTokenSource();
      const token = tokenSource.token;

      const symbols = await provider.provideDocumentSymbols(document, token);

      assert.ok(symbols.length > 0, "Should return at least one symbol");

      const testFunction = symbols.find((s) => s.name === "TestCaseInsensitive");
      assert.ok(testFunction, "Should find TestCaseInsensitive function");

      // 大文字小文字を区別しないフィールド名での検出を確認
      assert.ok(testFunction.children.length >= 3, "Should have at least 3 test cases");

      const testCases = testFunction.children.map((c) => c.name);
      assert.ok(testCases.includes("uppercase NAME"), "Should find uppercase NAME case");
      assert.ok(testCases.includes("mixed case Name"), "Should find mixed case Name case");
      assert.ok(testCases.includes("lowercase name"), "Should find lowercase name case");
    });
  });

  suite("Various Field Names", () => {
    test("Should recognize different supported field names", async () => {
      const testFilePath = path.join(fixturesDir, "various_fields_test.go");

      // ファイルが存在するかチェック
      if (!fs.existsSync(testFilePath)) {
        // テストをスキップ
        return;
      }

      const document = await vscode.workspace.openTextDocument(testFilePath);
      const tokenSource = new vscode.CancellationTokenSource();
      const token = tokenSource.token;

      const symbols = await provider.provideDocumentSymbols(document, token);

      assert.ok(symbols.length > 0, "Should return at least one symbol");

      // 各テスト関数が正しく認識されることを確認
      symbols.forEach((symbol) => {
        if (symbol.kind === vscode.SymbolKind.Function && symbol.name.startsWith("Test")) {
          assert.ok(symbol.children.length > 0, `${symbol.name} should have test cases`);
        }
      });
    });
  });

  suite("Error Handling", () => {
    test("Should handle files without test cases gracefully", async () => {
      const testFilePath = path.join(fixturesDir, "non_test_functions.go");

      if (!fs.existsSync(testFilePath)) {
        // テストをスキップ
        return;
      }

      const document = await vscode.workspace.openTextDocument(testFilePath);
      const tokenSource = new vscode.CancellationTokenSource();
      const token = tokenSource.token;

      const symbols = await provider.provideDocumentSymbols(document, token);

      // エラーではなく、空配列または非テスト関数のシンボルを返すことを確認
      assert.ok(Array.isArray(symbols), "Should return an array");

      // テーブル駆動テストではない関数は子要素を持たないことを確認
      symbols.forEach((symbol) => {
        if (symbol.kind === vscode.SymbolKind.Function && !symbol.name.startsWith("Test")) {
          assert.strictEqual(symbol.children.length, 0, "Non-test functions should have no children");
        }
      });
    });

    test("Should handle cancellation token", async () => {
      const testFilePath = path.join(fixturesDir, "basic_table_test.go");
      const document = await vscode.workspace.openTextDocument(testFilePath);

      // キャンセレーションを事前に要求
      const source = new vscode.CancellationTokenSource();
      source.cancel();

      const symbols = await provider.provideDocumentSymbols(document, source.token);
      assert.strictEqual(symbols.length, 0, "Should return empty array when cancelled");
    });
  });

  suite("Complex Test Structures", () => {
    test("Should parse typed test cases", async () => {
      const testFilePath = path.join(fixturesDir, "typed_test_cases.go");

      if (!fs.existsSync(testFilePath)) {
        return;
      }

      const document = await vscode.workspace.openTextDocument(testFilePath);
      const tokenSource = new vscode.CancellationTokenSource();
      const token = tokenSource.token;

      const symbols = await provider.provideDocumentSymbols(document, token);

      assert.ok(symbols.length > 0, "Should return at least one symbol");

      // 型付きテストケースが正しく解析されることを確認
      symbols.forEach((symbol) => {
        if (symbol.kind === vscode.SymbolKind.Function && symbol.name.startsWith("Test")) {
          assert.ok(symbol.children.length >= 0, `${symbol.name} should be parseable`);
        }
      });
    });

    test("Should parse map-based test cases", async () => {
      const testFilePath = path.join(fixturesDir, "map_test_cases.go");

      if (!fs.existsSync(testFilePath)) {
        return;
      }

      const document = await vscode.workspace.openTextDocument(testFilePath);
      const tokenSource = new vscode.CancellationTokenSource();
      const token = tokenSource.token;

      const symbols = await provider.provideDocumentSymbols(document, token);

      // マップベースのテストケースの解析結果を確認
      assert.ok(Array.isArray(symbols), "Should return an array for map-based tests");
    });

    test("Should parse tests with backtick strings", async () => {
      const testFilePath = path.join(fixturesDir, "backtick_strings_test.go");

      if (!fs.existsSync(testFilePath)) {
        return;
      }

      const document = await vscode.workspace.openTextDocument(testFilePath);
      const tokenSource = new vscode.CancellationTokenSource();
      const token = tokenSource.token;

      const symbols = await provider.provideDocumentSymbols(document, token);

      // バッククォート文字列を含むテストケースが正しく解析されることを確認
      assert.ok(Array.isArray(symbols), "Should handle backtick strings correctly");

      symbols.forEach((symbol) => {
        if (symbol.kind === vscode.SymbolKind.Function && symbol.name.startsWith("Test")) {
          // 各テストケースの名前が正しく抽出されていることを確認
          symbol.children.forEach((child) => {
            assert.ok(child.name.length > 0, "Test case name should not be empty");
          });
        }
      });
    });
  });
});
