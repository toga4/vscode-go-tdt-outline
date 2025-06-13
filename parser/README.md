# Go TDD Outline Parser

GoのテストファイルからVS Code用のアウトラインシンボルを生成するパーサーです。テーブル駆動テストのテストケースを個別のシンボルとして認識し、VS Codeのアウトラインビューで階層的に表示できます。

## 使用方法

```bash
# パーサーの実行
go run ./parser <test_file.go>

# JSON形式で出力される結果を整形して表示
go run ./parser <test_file.go> | jq '.'
```

## 対応するテストパターン

### 1. 匿名構造体のスライス
```go
tests := []struct {
    name string
    // ...
}{
    {name: "test1"},
    // ...
}
```

### 2. 型定義された構造体のスライス
```go
type Test struct {
    name string
    // ...
}

data := []Test{
    {name: "test1"},
    // ...
}
```

### 3. 型エイリアスを使用
```go
type Tests []Test

scenarios := Tests{
    {name: "test1"},
    // ...
}
```

### 4. map型のテストケース
```go
// 構造体をvalueに持つmap
tests := map[string]struct {
    input int
    want  int
}{
    "正常系": {input: 1, want: 2},
    "異常系": {input: -1, want: -1},
}

// シンプルなmap
cases := map[string]int{
    "one": 1,
    "two": 2,
}
```

**注意**: 
- スライス型の場合、変数名に制限はありません。`name`フィールド（または互換性のあるフィールド）を持つすべての構造体スライスが検出されます。
- map型の場合、文字列キーを持つmapのキーがテストケース名として使用されます。

## 実装された改善点

### 1. エラーハンドリングの改善
- `log.Fatal`の代わりに適切なエラー返却を実装
- ファイルパスの検証を追加（空文字列、非Goファイルのチェック）

### 2. パニック防止
- 文字列のインデックスアクセスを`strings.HasPrefix`に変更
- より安全な実装に改善

### 3. より柔軟なテストケース認識
- 複数のフィールド名をサポート: `name`, `testName`, `desc`, `description`, `title`, `scenario`
- 大文字小文字を区別しない比較（`strings.EqualFold`）
- 複数のテストテーブルに対応
- 型定義されたテストケースに対応

### 4. コード品質の向上
- 包括的なコメントを追加
- 定数の定義（`SymbolKindFunction`, `SymbolKindStruct`）
- 関数の責務を明確に分離

### 5. テストの追加
- 基本的なテーブル駆動テストの解析
- 複数のテスト関数とテストテーブル
- エラーケースの検証
- 様々なフィールド名のサポート確認
- 型定義されたテストケースの検証

## 出力形式

```json
[
  {
    "name": "TestExample",
    "detail": "test function",
    "kind": 11,
    "range": {...},
    "children": [
      {
        "name": "正常系テストケース",
        "detail": "test case",
        "kind": 12,
        "range": {...}
      }
    ]
  }
]
```

## 今後の拡張案

### 1. より複雑な型定義への対応
現在対応していない複雑なパターン：
```go
// パッケージ外の型を使用
tests := []somepackage.TestCase{
    {Name: "test1"},
}

// インターフェースを含む複雑な型定義
type TestCase interface {
    Name() string
}
```

### 2. サブテストへの対応
ネストしたt.Runによるサブテストの階層構造の認識

### 3. パフォーマンス最適化
大規模なテストファイルに対する処理速度の改善

### 4. カスタマイズ可能な設定
- 認識するフィールド名の設定
- 出力フォーマットの選択
- フィルタリング機能

## 開発

```bash
# テストの実行
cd parser
go test ./...

# カバレッジの確認
go test -cover ./...
``` 