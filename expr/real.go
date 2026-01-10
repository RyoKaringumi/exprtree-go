package expr

import "fmt"

// real.go - 実数（有理数）の厳密な表現と演算を提供する
//
// このファイルは有理数を分数形式（分子/分母）で表現し、
// 浮動小数点数の精度問題を回避した数学的に厳密な演算を実現する。

type RationalValue struct {
	Numerator   int64 // 分子
	Denominator int64 // 分母（常に正の整数）
}

func NewRational(numerator int64, denominator int64) (*RationalValue, error) {
	if denominator == 0 {
		return nil, fmt.Errorf("denominator cannot be zero")
	}
	return &RationalValue{
		Numerator:   numerator,
		Denominator: denominator,
	}, nil
}
