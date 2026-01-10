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

### 明示的変換の原則

**重要な設計方針**: このパッケージでは、式の形を変える操作は**常に明示的**に行う必要があります。

#### 暗黙的な操作の禁止

以下のような暗黙的な変換は行いません：

- **自動約分**: `6/9` を自動的に `2/3` に簡約化しない
- **自動符号正規化**: `5/-7` を自動的に `-5/7` に変換しない
- **自動正規形**: `0/5` を自動的に `0/1` に変換しない

#### 明示的な変換メソッド

式を変形する場合は、明示的なメソッド・関数を呼び出す必要があります：

```go
// ❌ 悪い例: コンストラクタで暗黙的に約分
r := NewRational(6, 9)  // 内部で自動的に 2/3 に変換される

// ✅ 良い例: 明示的に約分メソッドを呼び出す
r := NewRational(6, 9)       // 6/9 のまま
r = Simplify(r)              // 明示的に 2/3 に変換
```

#### この原則の理由

1. **透明性**: ユーザーが式の変化を追跡しやすい
2. **予測可能性**: 意図しない変換によるバグを防ぐ
3. **数学的厳密性**: 変換のタイミングと条件を明確にする
4. **デバッグ容易性**: どこで式が変わったかを把握しやすい

#### 有理数の設計

`expr/real.go` では、有理数を以下の方針で実装します：

- **内部表現**: 分子（numerator）と分母（denominator）のペア
- **生成時の制約**: 分母がゼロでないことのみをチェック
- **変換の明示性**: 約分、符号正規化などは明示的なメソッドで提供
- **不変性**: 一度生成された有理数オブジェクトは変更しない（イミュータブル）

例：
```go
// 基本的な生成（検証のみ、変換なし）
r, err := NewRational(6, 9)  // 6/9 のまま保持

// 明示的な簡約化
r = Simplify(r)              // 2/3 に変換

// 明示的な正規化（符号を分子に移動）
r = Normalize(NewRational(5, -7))  // -5/7 に変換
```

## ファイル構成

```
exprtree-go/
├── main.go                    # メインプログラム（使用例）
├── integration_test.go        # 統合テスト
├── go.mod                     # Goモジュール定義
├── CLAUDE.md                  # このファイル
├── expr/                      # Expression treeパッケージ
│   ├── tree.go               # すべての式型と評価ロジック
│   ├── tree_test.go          # Expressionのユニットテスト
│   ├── monomial.go           # 単項式関連の関数
│   ├── monomial_test.go      # 単項式のユニットテスト
│   ├── polynomial.go         # 多項式関連の関数
│   ├── polynomial_test.go    # 多項式のユニットテスト
│   ├── real.go               # 有理数の型と演算
│   └── real_test.go          # 有理数のユニットテスト
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
