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

- `tree.go`: すべての式型と評価ロジック
- `tree_test.go`: 各演算子と評価メカニズムのユニットテスト
- `main.go`: ライブラリの使用例
