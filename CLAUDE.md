# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## プロジェクト概要

exprtree-goは、数式の構文木(expression tree)を構築・評価するためのGoライブラリです。四則演算(加算、減算、乗算、除算)をサポートし、定数と変数を含む数式を評価できます。

## 開発コマンド

```bash
# テストの実行
go test

# 特定のテストの実行
go test -run TestAddExpression

# 詳細出力でテストを実行
go test -v

# プログラムの実行
go run .
```

## アーキテクチャ

### コア設計パターン

このプロジェクトはCompositeパターンを使用して式ツリーを実装しています:

- **Expression インターフェース**: すべての式タイプが実装する必要があるインターフェース
  - `Children() []Expression`: 子ノードを返す
  - `Eval() (ExpressValue, bool)`: 式を評価し、結果と成功フラグを返す

- **ExpressValue インターフェース**: 評価結果を表す
  - `NumberValue`: float64値をラップする具象型

### 式の階層構造

```
Expression (インターフェース)
├── BinaryExpression (基底構造体)
│   ├── AddExpression (加算)
│   ├── SubtractExpression (減算)
│   ├── MultiplyExpression (乗算)
│   └── DivideExpression (除算)
├── Constant (定数リーフノード)
└── Variable (変数リーフノード - 未実装)
```

**重要な実装パターン:**
- すべての二項演算子は`BinaryExpression`を埋め込み、`Left`と`Right`フィールドを継承
- コンストラクタ関数(例: `NewAddExpression`)を使用して式を構築
- `Eval()`メソッドは`(ExpressValue, bool)`を返し、評価の成功/失敗を示す
- リーフノード(`Constant`, `Variable`)は空のスライスを返す`Children()`を実装

### 評価メカニズム

評価は再帰的に行われます:
1. 二項演算子は両方の子を評価
2. 両方の評価が成功し、`NumberValue`を返す場合のみ演算を実行
3. ゼロ除算などのエラーケースは`(nil, false)`を返す

## ファイル構成

```
exprtree-go/
├── main.go                    # メインプログラム（使用例）
├── integration_test.go        # 統合テスト
├── go.mod                     # Goモジュール定義
├── CLAUDE.md                  # このファイル
├── expr/                      # Expression treeパッケージ
│   ├── tree.go               # すべての式型と評価ロジック
│   └── tree_test.go          # Expressionのユニットテスト
└── latex/                     # LaTeXパーサーパッケージ
    ├── lexer.go              # 字句解析器
    ├── lexer_test.go         # Lexerのユニットテスト
    ├── parser.go             # 構文解析器（Pratt Parsing）
    ├── parser_test.go        # Parserのユニットテスト
    ├── converter.go          # LaTeX AST → Expression tree変換
    └── converter_test.go     # Converterのユニットテスト
```

### パッケージ構成

- **exprtree (ルート)**: メインパッケージ
- **expr**: Expression tree関連の型と評価ロジック
- **latex**: LaTeX文字列のパース機能

## LaTeXパーサーの使用方法

### 基本的な使い方

```go
import (
    "exprtree/expr"
    "exprtree/latex"
)

// 数式文字列をパースして評価
result, err := latex.ParseAndEval("2 + 3 * 4")
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Result: %.2f\n", result.Value) // Result: 14.00
```

### Expression treeを取得

```go
import (
    "exprtree/expr"
    "exprtree/latex"
)

// 数式文字列からExpression treeを構築
expression, err := latex.ParseLatex("(2 + 3) * 4")
if err != nil {
    log.Fatal(err)
}

// ツリーを走査
children := expression.Children()
fmt.Printf("Number of children: %d\n", len(children))

// 評価
result, ok := expression.Eval()
if ok {
    if num, ok := result.(*expr.NumberValue); ok {
        fmt.Printf("Result: %.2f\n", num.Value)
    }
}
```

### サポートされる構文

現在サポートされている構文:
- **数値**: 整数 (`42`) および小数 (`3.14`)
- **四則演算子**: `+`, `-`, `*`, `/`
- **括弧**: `(`, `)`
- **演算子優先順位**: `*`, `/` > `+`, `-`
- **左結合性**: `10 - 3 - 2` → `(10 - 3) - 2`

例:
```
"2 + 3"              → 5
"2 + 3 * 4"          → 14 (優先順位)
"(2 + 3) * 4"        → 20 (括弧)
"10 - 2 - 3"         → 5 (左結合)
"2.5 + 1.5"          → 4 (小数)
"(1 + 2) * (3 + 4)"  → 21 (複雑な式)
```

### LaTeXパーサーのアーキテクチャ

3段階のパイプライン構成:

```
LaTeX文字列 → [Lexer] → Tokens → [Parser] → LaTeX AST → [Converter] → Expression Tree
```

1. **Lexer (字句解析)**: 入力文字列をトークンに分割
2. **Parser (構文解析)**: トークンからLaTeX ASTを構築（Pratt Parsingで優先順位を処理）
3. **Converter (変換)**: LaTeX ASTを既存のExpression treeに変換

この分離により、各レイヤーを独立してテスト可能で、将来の拡張が容易です。

### 将来の拡張予定

- **LaTeXコマンド**: `\frac{分子}{分母}`, `\sqrt{x}` 等
- **変数**: `x`, `y` 等（評価コンテキスト付き）
- **暗黙の乗算**: `2x`, `3(x+1)`
- **累乗**: `x^2`
- **その他のLaTeX構文**
