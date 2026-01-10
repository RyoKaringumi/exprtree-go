package expr

import (
	"testing"
)

// TestNewRational tests the creation of rational numbers
// NewRationalは検証のみを行い、自動的な約分や正規化は行わない
func TestNewRational(t *testing.T) {
	tests := []struct {
		name        string
		numerator   int64
		denominator int64
		wantErr     bool
		description string
	}{
		{
			name:        "正の有理数",
			numerator:   3,
			denominator: 4,
			wantErr:     false,
			description: "3/4は有効な有理数",
		},
		{
			name:        "負の有理数（負の分子）",
			numerator:   -5,
			denominator: 7,
			wantErr:     false,
			description: "-5/7は有効な有理数",
		},
		{
			name:        "負の有理数（負の分母）",
			numerator:   5,
			denominator: -7,
			wantErr:     false,
			description: "5/-7は有効な有理数（正規化はせずそのまま保持）",
		},
		{
			name:        "負の有理数（両方負）",
			numerator:   -5,
			denominator: -7,
			wantErr:     false,
			description: "-5/-7は有効な有理数（正規化はせずそのまま保持）",
		},
		{
			name:        "整数（分母1）",
			numerator:   42,
			denominator: 1,
			wantErr:     false,
			description: "整数は分母1の有理数として表現可能",
		},
		{
			name:        "ゼロ",
			numerator:   0,
			denominator: 1,
			wantErr:     false,
			description: "0は0/1として表現可能",
		},
		{
			name:        "ゼロ除算エラー",
			numerator:   1,
			denominator: 0,
			wantErr:     true,
			description: "分母が0の場合はエラー",
		},
		{
			name:        "ゼロ除算エラー（分子も0）",
			numerator:   0,
			denominator: 0,
			wantErr:     true,
			description: "0/0は不定形でエラー",
		},
		{
			name:        "約分可能な数（未約分のまま保持）",
			numerator:   6,
			denominator: 9,
			wantErr:     false,
			description: "6/9は約分せずそのまま保持（Simplify()を呼ばない限り）",
		},
		{
			name:        "大きな数",
			numerator:   1000000,
			denominator: 3,
			wantErr:     false,
			description: "大きな整数も正確に表現可能",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 実装待ち
			r, err := NewRational(tt.numerator, tt.denominator)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewRational() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr {
				// 入力値がそのまま保持されていることを確認
				if r.Numerator != tt.numerator || r.Denominator != tt.denominator {
					t.Errorf("NewRational(%d,%d) = %d/%d, 入力値がそのまま保持されていない",
						tt.numerator, tt.denominator, r.Numerator, r.Denominator)
				}
			}
		})
	}
}

// TestRationalAdd tests addition of rational numbers
// 演算結果は自動的に簡約化されない。必要に応じてSimplify()を明示的に呼び出す
func TestRationalAdd(t *testing.T) {
	tests := []struct {
		name        string
		r1Num       int64
		r1Denom     int64
		r2Num       int64
		r2Denom     int64
		wantNum     int64
		wantDenom   int64
		description string
	}{
		{
			name:  "同じ分母の加算（未簡約）",
			r1Num: 1, r1Denom: 4,
			r2Num: 1, r2Denom: 4,
			wantNum: 2, wantDenom: 4,
			description: "1/4 + 1/4 = 2/4（自動簡約化なし）",
		},
		{
			name:  "異なる分母の加算（未簡約）",
			r1Num: 1, r1Denom: 3,
			r2Num: 1, r2Denom: 6,
			wantNum: 3, wantDenom: 6,
			description: "1/3 + 1/6 = 2/6 + 1/6 = 3/6（自動簡約化なし）",
		},
		{
			name:  "負の数の加算",
			r1Num: -1, r1Denom: 2,
			r2Num: 1, r2Denom: 4,
			wantNum: -1, wantDenom: 4,
			description: "-1/2 + 1/4 = -2/4 + 1/4 = -1/4",
		},
		{
			name:  "整数の加算",
			r1Num: 3, r1Denom: 1,
			r2Num: 5, r2Denom: 1,
			wantNum: 8, wantDenom: 1,
			description: "3 + 5 = 8",
		},
		{
			name:  "ゼロとの加算",
			r1Num: 3, r1Denom: 4,
			r2Num: 0, r2Denom: 1,
			wantNum: 3, wantDenom: 4,
			description: "3/4 + 0 = 3/4（加法の単位元）",
		},
		{
			name:  "逆数の加算",
			r1Num: 3, r1Denom: 4,
			r2Num: -3, r2Denom: 4,
			wantNum: 0, wantDenom: 4,
			description: "3/4 + (-3/4) = 0/4（加法の逆元、ゼロも分母を保持）",
		},
		{
			name:  "可換性のテスト（未簡約）",
			r1Num: 2, r1Denom: 5,
			r2Num: 3, r2Denom: 7,
			wantNum: 29, wantDenom: 35,
			description: "2/5 + 3/7 = 14/35 + 15/35 = 29/35（可換性: a+b=b+a）",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 実装待ち
			r1, _ := NewRational(tt.r1Num, tt.r1Denom)
			r2, _ := NewRational(tt.r2Num, tt.r2Denom)
			result := RationalAdd(r1, r2)
			if result.Numerator != tt.wantNum || result.Denominator != tt.wantDenom {
				t.Errorf("Add() = %d/%d, want %d/%d",
					result.Numerator, result.Denominator, tt.wantNum, tt.wantDenom)
			}
		})
	}
}

// TestRationalSubtract tests subtraction of rational numbers
func TestRationalSubtract(t *testing.T) {
	tests := []struct {
		name        string
		r1Num       int64
		r1Denom     int64
		r2Num       int64
		r2Denom     int64
		wantNum     int64
		wantDenom   int64
		description string
	}{
		{
			name:  "同じ分母の減算（未簡約）",
			r1Num: 3, r1Denom: 4,
			r2Num: 1, r2Denom: 4,
			wantNum: 2, wantDenom: 4,
			description: "3/4 - 1/4 = 2/4（自動簡約化なし）",
		},
		{
			name:  "異なる分母の減算（未簡約）",
			r1Num: 1, r1Denom: 2,
			r2Num: 1, r2Denom: 3,
			wantNum: 1, wantDenom: 6,
			description: "1/2 - 1/3 = 3/6 - 2/6 = 1/6",
		},
		{
			name:  "負の結果",
			r1Num: 1, r1Denom: 4,
			r2Num: 1, r2Denom: 2,
			wantNum: -1, wantDenom: 4,
			description: "1/4 - 1/2 = 1/4 - 2/4 = -1/4",
		},
		{
			name:  "ゼロからの減算",
			r1Num: 0, r1Denom: 1,
			r2Num: 3, r2Denom: 4,
			wantNum: -3, wantDenom: 4,
			description: "0 - 3/4 = -3/4",
		},
		{
			name:  "自分自身の減算",
			r1Num: 5, r1Denom: 7,
			r2Num: 5, r2Denom: 7,
			wantNum: 0, wantDenom: 7,
			description: "5/7 - 5/7 = 0/7（分母を保持）",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 実装待ち
		})
	}
}

// TestRationalMultiply tests multiplication of rational numbers
func TestRationalMultiply(t *testing.T) {
	tests := []struct {
		name        string
		r1Num       int64
		r1Denom     int64
		r2Num       int64
		r2Denom     int64
		wantNum     int64
		wantDenom   int64
		description string
	}{
		{
			name:  "基本的な乗算（未簡約）",
			r1Num: 2, r1Denom: 3,
			r2Num: 3, r2Denom: 4,
			wantNum: 6, wantDenom: 12,
			description: "2/3 * 3/4 = 6/12（自動簡約化なし）",
		},
		{
			name:  "整数との乗算",
			r1Num: 3, r1Denom: 5,
			r2Num: 2, r2Denom: 1,
			wantNum: 6, wantDenom: 5,
			description: "3/5 * 2 = 6/5",
		},
		{
			name:  "負の数の乗算（未簡約）",
			r1Num: -2, r1Denom: 3,
			r2Num: 3, r2Denom: 5,
			wantNum: -6, wantDenom: 15,
			description: "-2/3 * 3/5 = -6/15（自動簡約化なし）",
		},
		{
			name:  "両方負の数の乗算（未簡約）",
			r1Num: -2, r1Denom: 3,
			r2Num: -3, r2Denom: 5,
			wantNum: 6, wantDenom: 15,
			description: "-2/3 * -3/5 = 6/15（負×負=正、自動簡約化なし）",
		},
		{
			name:  "ゼロとの乗算",
			r1Num: 5, r1Denom: 7,
			r2Num: 0, r2Denom: 1,
			wantNum: 0, wantDenom: 7,
			description: "5/7 * 0 = 0/7（乗法の零元、分母を保持）",
		},
		{
			name:  "1との乗算",
			r1Num: 5, r1Denom: 7,
			r2Num: 1, r2Denom: 1,
			wantNum: 5, wantDenom: 7,
			description: "5/7 * 1 = 5/7（乗法の単位元）",
		},
		{
			name:  "逆数との乗算（未簡約）",
			r1Num: 3, r1Denom: 4,
			r2Num: 4, r2Denom: 3,
			wantNum: 12, wantDenom: 12,
			description: "3/4 * 4/3 = 12/12（乗法の逆元、自動簡約化なし）",
		},
		{
			name:  "可換性のテスト",
			r1Num: 2, r1Denom: 5,
			r2Num: 3, r2Denom: 7,
			wantNum: 6, wantDenom: 35,
			description: "2/5 * 3/7 = 6/35（可換性: a*b=b*a）",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 実装待ち
		})
	}
}

// TestRationalDivide tests division of rational numbers
func TestRationalDivide(t *testing.T) {
	tests := []struct {
		name        string
		r1Num       int64
		r1Denom     int64
		r2Num       int64
		r2Denom     int64
		wantNum     int64
		wantDenom   int64
		wantErr     bool
		description string
	}{
		{
			name:  "基本的な除算（未簡約）",
			r1Num: 1, r1Denom: 2,
			r2Num: 1, r2Denom: 3,
			wantNum: 3, wantDenom: 2,
			wantErr:     false,
			description: "1/2 ÷ 1/3 = 1/2 * 3/1 = 3/2",
		},
		{
			name:  "整数による除算（未簡約）",
			r1Num: 6, r1Denom: 5,
			r2Num: 2, r2Denom: 1,
			wantNum: 6, wantDenom: 10,
			wantErr:     false,
			description: "6/5 ÷ 2 = 6/5 * 1/2 = 6/10（自動簡約化なし）",
		},
		{
			name:  "ゼロによる除算",
			r1Num: 5, r1Denom: 7,
			r2Num: 0, r2Denom: 1,
			wantNum: 0, wantDenom: 0,
			wantErr:     true,
			description: "5/7 ÷ 0 はエラー（ゼロ除算）",
		},
		{
			name:  "ゼロを除算",
			r1Num: 0, r1Denom: 1,
			r2Num: 5, r2Denom: 7,
			wantNum: 0, wantDenom: 5,
			wantErr:     false,
			description: "0 ÷ 5/7 = 0 * 7/5 = 0/5",
		},
		{
			name:  "1による除算",
			r1Num: 5, r1Denom: 7,
			r2Num: 1, r2Denom: 1,
			wantNum: 5, wantDenom: 7,
			wantErr:     false,
			description: "5/7 ÷ 1 = 5/7",
		},
		{
			name:  "自分自身による除算（未簡約）",
			r1Num: 5, r1Denom: 7,
			r2Num: 5, r2Denom: 7,
			wantNum: 35, wantDenom: 35,
			wantErr:     false,
			description: "5/7 ÷ 5/7 = 5/7 * 7/5 = 35/35（自動簡約化なし）",
		},
		{
			name:  "負の数の除算（未簡約）",
			r1Num: -3, r1Denom: 4,
			r2Num: 2, r2Denom: 5,
			wantNum: -15, wantDenom: 8,
			wantErr:     false,
			description: "-3/4 ÷ 2/5 = -3/4 * 5/2 = -15/8",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 実装待ち
		})
	}
}

// TestRationalSimplify tests simplification of rational numbers
// Simplify()は明示的に呼び出される変換関数
func TestRationalSimplify(t *testing.T) {
	tests := []struct {
		name        string
		num         int64
		denom       int64
		wantNum     int64
		wantDenom   int64
		description string
	}{
		{
			name: "既に既約",
			num:  3, denom: 5,
			wantNum: 3, wantDenom: 5,
			description: "3/5は既に既約分数",
		},
		{
			name: "約分可能",
			num:  6, denom: 9,
			wantNum: 2, wantDenom: 3,
			description: "6/9 = 2/3（GCD=3で約分）",
		},
		{
			name: "大きな公約数",
			num:  100, denom: 150,
			wantNum: 2, wantDenom: 3,
			description: "100/150 = 2/3（GCD=50で約分）",
		},
		{
			name: "負の分子（符号はそのまま）",
			num:  -4, denom: 6,
			wantNum: -2, wantDenom: 3,
			description: "-4/6 = -2/3（約分のみ、符号正規化はしない）",
		},
		{
			name: "負の分母（符号はそのまま）",
			num:  4, denom: -6,
			wantNum: 2, wantDenom: -3,
			description: "4/-6 = 2/-3（約分のみ、符号正規化はしない）",
		},
		{
			name: "両方負（符号はそのまま）",
			num:  -4, denom: -6,
			wantNum: -2, wantDenom: -3,
			description: "-4/-6 = -2/-3（約分のみ、符号正規化はしない）",
		},
		{
			name: "ゼロ",
			num:  0, denom: 5,
			wantNum: 0, wantDenom: 1,
			description: "0/5 = 0/1（ゼロの場合は分母を1に）",
		},
		{
			name: "整数",
			num:  15, denom: 5,
			wantNum: 3, wantDenom: 1,
			description: "15/5 = 3/1（整数形式）",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 実装待ち
			// r, _ := NewRational(tt.num, tt.denom)
			// simplified := Simplify(r)
			// if simplified.Numerator() != tt.wantNum || simplified.Denominator() != tt.wantDenom {
			//     t.Errorf("Simplify(%d/%d) = %d/%d, want %d/%d",
			//         tt.num, tt.denom,
			//         simplified.Numerator(), simplified.Denominator(),
			//         tt.wantNum, tt.wantDenom)
			// }
		})
	}
}

// TestRationalNormalize tests sign normalization
// Normalize()は明示的に呼び出される変換関数（符号を分子に移動）
func TestRationalNormalize(t *testing.T) {
	tests := []struct {
		name        string
		num         int64
		denom       int64
		wantNum     int64
		wantDenom   int64
		description string
	}{
		{
			name: "既に正規化済み（正）",
			num:  3, denom: 4,
			wantNum: 3, wantDenom: 4,
			description: "3/4は既に正規化済み",
		},
		{
			name: "既に正規化済み（負の分子）",
			num:  -3, denom: 4,
			wantNum: -3, wantDenom: 4,
			description: "-3/4は既に正規化済み",
		},
		{
			name: "負の分母を正規化",
			num:  3, denom: -4,
			wantNum: -3, wantDenom: 4,
			description: "3/-4 → -3/4（符号を分子に移動）",
		},
		{
			name: "両方負を正規化",
			num:  -3, denom: -4,
			wantNum: 3, wantDenom: 4,
			description: "-3/-4 → 3/4（両方負は正）",
		},
		{
			name: "ゼロの正規化",
			num:  0, denom: -5,
			wantNum: 0, wantDenom: 5,
			description: "0/-5 → 0/5（ゼロは常に正の分母）",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 実装待ち
			// r, _ := NewRational(tt.num, tt.denom)
			// normalized := Normalize(r)
			// if normalized.Numerator() != tt.wantNum || normalized.Denominator() != tt.wantDenom {
			//     t.Errorf("Normalize(%d/%d) = %d/%d, want %d/%d",
			//         tt.num, tt.denom,
			//         normalized.Numerator(), normalized.Denominator(),
			//         tt.wantNum, tt.wantDenom)
			// }
		})
	}
}

// TestRationalCanonical tests canonical form (simplify + normalize)
// Canonical()は簡約化と正規化を同時に行う明示的な変換関数
func TestRationalCanonical(t *testing.T) {
	tests := []struct {
		name        string
		num         int64
		denom       int64
		wantNum     int64
		wantDenom   int64
		description string
	}{
		{
			name: "正規形への変換",
			num:  6, denom: -9,
			wantNum: -2, wantDenom: 3,
			description: "6/-9 → -2/3（約分 + 符号正規化）",
		},
		{
			name: "両方負の正規形",
			num:  -6, denom: -9,
			wantNum: 2, wantDenom: 3,
			description: "-6/-9 → 2/3（約分 + 符号正規化）",
		},
		{
			name: "ゼロの正規形",
			num:  0, denom: -15,
			wantNum: 0, wantDenom: 1,
			description: "0/-15 → 0/1（ゼロの正規形）",
		},
		{
			name: "整数の正規形",
			num:  -15, denom: -5,
			wantNum: 3, wantDenom: 1,
			description: "-15/-5 → 3/1（整数の正規形）",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 実装待ち
			// r, _ := NewRational(tt.num, tt.denom)
			// canonical := Canonical(r)
			// if canonical.Numerator() != tt.wantNum || canonical.Denominator() != tt.wantDenom {
			//     t.Errorf("Canonical(%d/%d) = %d/%d, want %d/%d",
			//         tt.num, tt.denom,
			//         canonical.Numerator(), canonical.Denominator(),
			//         tt.wantNum, tt.wantDenom)
			// }
		})
	}
}

// TestGCD tests the greatest common divisor function
func TestGCD(t *testing.T) {
	tests := []struct {
		name        string
		a           int64
		b           int64
		want        int64
		description string
	}{
		{
			name: "互いに素",
			a:    7, b: 11,
			want:        1,
			description: "GCD(7,11) = 1（素数同士）",
		},
		{
			name: "公約数あり",
			a:    12, b: 18,
			want:        6,
			description: "GCD(12,18) = 6",
		},
		{
			name: "一方が0",
			a:    5, b: 0,
			want:        5,
			description: "GCD(5,0) = 5（定義により）",
		},
		{
			name: "両方0",
			a:    0, b: 0,
			want:        0,
			description: "GCD(0,0) = 0（定義により）",
		},
		{
			name: "同じ数",
			a:    15, b: 15,
			want:        15,
			description: "GCD(15,15) = 15",
		},
		{
			name: "負の数（絶対値で計算）",
			a:    -12, b: 18,
			want:        6,
			description: "GCD(-12,18) = 6（絶対値で計算）",
		},
		{
			name: "大きな数",
			a:    1071, b: 462,
			want:        21,
			description: "GCD(1071,462) = 21（ユークリッドの互除法の例）",
		},
		{
			name: "2のべき乗",
			a:    64, b: 48,
			want:        16,
			description: "GCD(64,48) = 16",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 実装待ち
			// result := GCD(tt.a, tt.b)
			// if result != tt.want {
			//     t.Errorf("GCD(%d,%d) = %d, want %d", tt.a, tt.b, result, tt.want)
			// }
		})
	}
}

// TestRationalEqual tests equality comparison
// Equal()は正規形で比較するが、元のオブジェクトは変更しない
func TestRationalEqual(t *testing.T) {
	tests := []struct {
		name        string
		r1Num       int64
		r1Denom     int64
		r2Num       int64
		r2Denom     int64
		want        bool
		description string
	}{
		{
			name:  "等しい（同じ表現）",
			r1Num: 1, r1Denom: 2,
			r2Num: 1, r2Denom: 2,
			want:        true,
			description: "1/2 = 1/2",
		},
		{
			name:  "等しい（異なる表現）",
			r1Num: 1, r1Denom: 2,
			r2Num: 2, r2Denom: 4,
			want:        true,
			description: "1/2 = 2/4（内部的に正規化して比較）",
		},
		{
			name:  "等しい（符号が異なる表現）",
			r1Num: -1, r1Denom: 2,
			r2Num: 1, r2Denom: -2,
			want:        true,
			description: "-1/2 = 1/-2（内部的に正規化して比較）",
		},
		{
			name:  "等しい（複雑な表現）",
			r1Num: -6, r1Denom: -9,
			r2Num: 2, r2Denom: 3,
			want:        true,
			description: "-6/-9 = 2/3（内部的に正規化して比較）",
		},
		{
			name:  "等しくない",
			r1Num: 1, r1Denom: 2,
			r2Num: 1, r2Denom: 3,
			want:        false,
			description: "1/2 ≠ 1/3",
		},
		{
			name:  "ゼロの等価性",
			r1Num: 0, r1Denom: 1,
			r2Num: 0, r2Denom: 5,
			want:        true,
			description: "0/1 = 0/5 = 0",
		},
		{
			name:  "符号の違い",
			r1Num: 1, r1Denom: 2,
			r2Num: -1, r2Denom: 2,
			want:        false,
			description: "1/2 ≠ -1/2",
		},
		{
			name:  "反射律",
			r1Num: 3, r1Denom: 5,
			r2Num: 3, r2Denom: 5,
			want:        true,
			description: "a = a（反射律）",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 実装待ち
			// r1, _ := NewRational(tt.r1Num, tt.r1Denom)
			// r2, _ := NewRational(tt.r2Num, tt.r2Denom)
			// result := Equal(r1, r2)
			// if result != tt.want {
			//     t.Errorf("Equal(%d/%d, %d/%d) = %v, want %v",
			//         tt.r1Num, tt.r1Denom, tt.r2Num, tt.r2Denom, result, tt.want)
			// }
			// // 元のオブジェクトが変更されていないことを確認
			// if r1.Numerator() != tt.r1Num || r1.Denominator() != tt.r1Denom {
			//     t.Errorf("Equal()がr1を変更した: %d/%d → %d/%d",
			//         tt.r1Num, tt.r1Denom, r1.Numerator(), r1.Denominator())
			// }
			// if r2.Numerator() != tt.r2Num || r2.Denominator() != tt.r2Denom {
			//     t.Errorf("Equal()がr2を変更した: %d/%d → %d/%d",
			//         tt.r2Num, tt.r2Denom, r2.Numerator(), r2.Denominator())
			// }
		})
	}
}

// TestRationalCompare tests comparison of rational numbers
// Compare()は正規形で比較するが、元のオブジェクトは変更しない
func TestRationalCompare(t *testing.T) {
	tests := []struct {
		name        string
		r1Num       int64
		r1Denom     int64
		r2Num       int64
		r2Denom     int64
		want        int // -1: less, 0: equal, 1: greater
		description string
	}{
		{
			name:  "より小さい",
			r1Num: 1, r1Denom: 3,
			r2Num: 1, r2Denom: 2,
			want:        -1,
			description: "1/3 < 1/2",
		},
		{
			name:  "等しい",
			r1Num: 1, r1Denom: 2,
			r2Num: 2, r2Denom: 4,
			want:        0,
			description: "1/2 = 2/4",
		},
		{
			name:  "より大きい",
			r1Num: 3, r1Denom: 4,
			r2Num: 1, r2Denom: 2,
			want:        1,
			description: "3/4 > 1/2",
		},
		{
			name:  "負の数との比較",
			r1Num: -1, r1Denom: 2,
			r2Num: 1, r2Denom: 2,
			want:        -1,
			description: "-1/2 < 1/2",
		},
		{
			name:  "負の数同士の比較",
			r1Num: -1, r1Denom: 2,
			r2Num: -1, r2Denom: 3,
			want:        -1,
			description: "-1/2 < -1/3",
		},
		{
			name:  "ゼロとの比較（正）",
			r1Num: 1, r1Denom: 2,
			r2Num: 0, r2Denom: 1,
			want:        1,
			description: "1/2 > 0",
		},
		{
			name:  "ゼロとの比較（負）",
			r1Num: -1, r1Denom: 2,
			r2Num: 0, r2Denom: 1,
			want:        -1,
			description: "-1/2 < 0",
		},
		{
			name:  "異なる表現での比較",
			r1Num: 2, r1Denom: -4,
			r2Num: -3, r2Denom: 6,
			want:        0,
			description: "2/-4 = -3/6 = -1/2（正規化して比較）",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 実装待ち
		})
	}
}

// TestRationalAbs tests absolute value
func TestRationalAbs(t *testing.T) {
	tests := []struct {
		name        string
		num         int64
		denom       int64
		wantNum     int64
		wantDenom   int64
		description string
	}{
		{
			name: "正の数",
			num:  3, denom: 4,
			wantNum: 3, wantDenom: 4,
			description: "|3/4| = 3/4",
		},
		{
			name: "負の数（負の分子）",
			num:  -3, denom: 4,
			wantNum: 3, wantDenom: 4,
			description: "|-3/4| = 3/4",
		},
		{
			name: "負の数（負の分母）",
			num:  3, denom: -4,
			wantNum: 3, wantDenom: 4,
			description: "|3/-4| = 3/4",
		},
		{
			name: "ゼロ",
			num:  0, denom: 1,
			wantNum: 0, wantDenom: 1,
			description: "|0| = 0",
		},
		{
			name: "負の整数",
			num:  -5, denom: 1,
			wantNum: 5, wantDenom: 1,
			description: "|-5| = 5",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 実装待ち
		})
	}
}

// TestRationalSign tests sign determination
func TestRationalSign(t *testing.T) {
	tests := []struct {
		name        string
		num         int64
		denom       int64
		want        int // -1: negative, 0: zero, 1: positive
		description string
	}{
		{
			name: "正の数",
			num:  3, denom: 4,
			want:        1,
			description: "sign(3/4) = 1",
		},
		{
			name: "負の数（負の分子）",
			num:  -3, denom: 4,
			want:        -1,
			description: "sign(-3/4) = -1",
		},
		{
			name: "負の数（負の分母）",
			num:  3, denom: -4,
			want:        -1,
			description: "sign(3/-4) = -1",
		},
		{
			name: "両方負（正）",
			num:  -3, denom: -4,
			want:        1,
			description: "sign(-3/-4) = 1",
		},
		{
			name: "ゼロ",
			num:  0, denom: 1,
			want:        0,
			description: "sign(0) = 0",
		},
		{
			name: "負の整数",
			num:  -7, denom: 1,
			want:        -1,
			description: "sign(-7) = -1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 実装待ち
		})
	}
}

// TestRationalNegate tests negation
func TestRationalNegate(t *testing.T) {
	tests := []struct {
		name        string
		num         int64
		denom       int64
		wantNum     int64
		wantDenom   int64
		description string
	}{
		{
			name: "正の数の符号反転",
			num:  3, denom: 4,
			wantNum: -3, wantDenom: 4,
			description: "-(3/4) = -3/4",
		},
		{
			name: "負の数の符号反転",
			num:  -3, denom: 4,
			wantNum: 3, wantDenom: 4,
			description: "-(-3/4) = 3/4",
		},
		{
			name: "負の分母の符号反転",
			num:  3, denom: -4,
			wantNum: -3, wantDenom: -4,
			description: "-(3/-4) = -3/-4（分子の符号を反転）",
		},
		{
			name: "ゼロの符号反転",
			num:  0, denom: 1,
			wantNum: 0, wantDenom: 1,
			description: "-0 = 0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 実装待ち
			// r, _ := NewRational(tt.num, tt.denom)
			// negated := Negate(r)
			// if negated.Numerator() != tt.wantNum || negated.Denominator() != tt.wantDenom {
			//     t.Errorf("Negate(%d/%d) = %d/%d, want %d/%d",
			//         tt.num, tt.denom,
			//         negated.Numerator(), negated.Denominator(),
			//         tt.wantNum, tt.wantDenom)
			// }
			// // 二重否定のテスト
			// doubleNegated := Negate(negated)
			// if doubleNegated.Numerator() != tt.num || doubleNegated.Denominator() != tt.denom {
			//     t.Errorf("Negate(Negate(%d/%d)) = %d/%d, want %d/%d",
			//         tt.num, tt.denom,
			//         doubleNegated.Numerator(), doubleNegated.Denominator(),
			//         tt.num, tt.denom)
			// }
		})
	}
}

// TestRationalReciprocal tests reciprocal (multiplicative inverse)
func TestRationalReciprocal(t *testing.T) {
	tests := []struct {
		name        string
		num         int64
		denom       int64
		wantNum     int64
		wantDenom   int64
		wantErr     bool
		description string
	}{
		{
			name: "正の有理数の逆数",
			num:  3, denom: 4,
			wantNum: 4, wantDenom: 3,
			wantErr:     false,
			description: "1/(3/4) = 4/3（分子と分母を入れ替え）",
		},
		{
			name: "負の有理数の逆数（負の分子）",
			num:  -3, denom: 4,
			wantNum: -4, wantDenom: 3,
			wantErr:     false,
			description: "1/(-3/4) = -4/3",
		},
		{
			name: "負の有理数の逆数（負の分母）",
			num:  3, denom: -4,
			wantNum: -4, wantDenom: 3,
			wantErr:     false,
			description: "1/(3/-4) = -4/3",
		},
		{
			name: "整数の逆数",
			num:  5, denom: 1,
			wantNum: 1, wantDenom: 5,
			wantErr:     false,
			description: "1/5 = 1/5",
		},
		{
			name: "ゼロの逆数",
			num:  0, denom: 1,
			wantNum: 0, wantDenom: 0,
			wantErr:     true,
			description: "1/0は未定義",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 実装待ち
			// r, _ := NewRational(tt.num, tt.denom)
			// reciprocal, err := Reciprocal(r)
			// if (err != nil) != tt.wantErr {
			//     t.Errorf("Reciprocal() error = %v, wantErr %v", err, tt.wantErr)
			// }
			// if !tt.wantErr {
			//     if reciprocal.Numerator() != tt.wantNum || reciprocal.Denominator() != tt.wantDenom {
			//         t.Errorf("Reciprocal(%d/%d) = %d/%d, want %d/%d",
			//             tt.num, tt.denom,
			//             reciprocal.Numerator(), reciprocal.Denominator(),
			//             tt.wantNum, tt.wantDenom)
			//     }
			//     // 二重逆数のテスト（正規化して比較）
			//     doubleReciprocal, _ := Reciprocal(reciprocal)
			//     if !Equal(doubleReciprocal, r) {
			//         t.Errorf("Reciprocal(Reciprocal(%d/%d)) ≠ %d/%d",
			//             tt.num, tt.denom, tt.num, tt.denom)
			//     }
			// }
		})
	}
}

// TestRationalIsInteger tests integer detection
func TestRationalIsInteger(t *testing.T) {
	tests := []struct {
		name        string
		num         int64
		denom       int64
		want        bool
		description string
	}{
		{
			name: "整数",
			num:  5, denom: 1,
			want:        true,
			description: "5/1は整数",
		},
		{
			name: "簡約化すると整数",
			num:  10, denom: 5,
			want:        true,
			description: "10/5 = 2は整数",
		},
		{
			name: "真分数",
			num:  3, denom: 4,
			want:        false,
			description: "3/4は整数でない",
		},
		{
			name: "ゼロ",
			num:  0, denom: 1,
			want:        true,
			description: "0は整数",
		},
		{
			name: "負の整数",
			num:  -5, denom: 1,
			want:        true,
			description: "-5は整数",
		},
		{
			name: "負の分母だが整数",
			num:  -10, denom: -5,
			want:        true,
			description: "-10/-5 = 2は整数",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 実装待ち
		})
	}
}

// TestRationalToFloat64 tests conversion to float64
func TestRationalToFloat64(t *testing.T) {
	tests := []struct {
		name        string
		num         int64
		denom       int64
		want        float64
		tolerance   float64
		description string
	}{
		{
			name: "単純な分数",
			num:  1, denom: 2,
			want:        0.5,
			tolerance:   1e-10,
			description: "1/2 = 0.5",
		},
		{
			name: "1/3（無限小数）",
			num:  1, denom: 3,
			want:        0.333333333333333,
			tolerance:   1e-10,
			description: "1/3 ≈ 0.333...",
		},
		{
			name: "整数",
			num:  5, denom: 1,
			want:        5.0,
			tolerance:   1e-10,
			description: "5/1 = 5.0",
		},
		{
			name: "ゼロ",
			num:  0, denom: 1,
			want:        0.0,
			tolerance:   1e-10,
			description: "0/1 = 0.0",
		},
		{
			name: "負の数",
			num:  -3, denom: 4,
			want:        -0.75,
			tolerance:   1e-10,
			description: "-3/4 = -0.75",
		},
		{
			name: "負の分母",
			num:  3, denom: -4,
			want:        -0.75,
			tolerance:   1e-10,
			description: "3/-4 = -0.75",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 実装待ち
			// r, _ := NewRational(tt.num, tt.denom)
			// result := ToFloat64(r)
			// if math.Abs(result-tt.want) > tt.tolerance {
			//     t.Errorf("ToFloat64(%d/%d) = %v, want %v", tt.num, tt.denom, result, tt.want)
			// }
		})
	}
}

// TestRationalString tests string representation
func TestRationalString(t *testing.T) {
	tests := []struct {
		name        string
		num         int64
		denom       int64
		want        string
		description string
	}{
		{
			name: "通常の分数",
			num:  3, denom: 4,
			want:        "3/4",
			description: "3/4の文字列表現",
		},
		{
			name: "整数",
			num:  5, denom: 1,
			want:        "5",
			description: "整数は分母を省略",
		},
		{
			name: "ゼロ",
			num:  0, denom: 1,
			want:        "0",
			description: "ゼロの文字列表現",
		},
		{
			name: "負の分数",
			num:  -3, denom: 4,
			want:        "-3/4",
			description: "負の分数",
		},
		{
			name: "未簡約の分数（そのまま表示）",
			num:  6, denom: 9,
			want:        "6/9",
			description: "6/9は簡約化せずそのまま表示",
		},
		{
			name: "負の分母（そのまま表示）",
			num:  3, denom: -4,
			want:        "3/-4",
			description: "3/-4は正規化せずそのまま表示",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 実装待ち
		})
	}
}

// TestRationalPower tests integer exponentiation
func TestRationalPower(t *testing.T) {
	tests := []struct {
		name        string
		num         int64
		denom       int64
		exponent    int
		wantNum     int64
		wantDenom   int64
		description string
	}{
		{
			name: "正のべき乗（未簡約）",
			num:  2, denom: 3,
			exponent: 2,
			wantNum:  4, wantDenom: 9,
			description: "(2/3)^2 = 4/9",
		},
		{
			name: "ゼロ乗",
			num:  2, denom: 3,
			exponent: 0,
			wantNum:  1, wantDenom: 1,
			description: "(2/3)^0 = 1（ゼロ乗の定義）",
		},
		{
			name: "1乗",
			num:  2, denom: 3,
			exponent: 1,
			wantNum:  2, wantDenom: 3,
			description: "(2/3)^1 = 2/3",
		},
		{
			name: "負のべき乗",
			num:  2, denom: 3,
			exponent: -1,
			wantNum:  3, wantDenom: 2,
			description: "(2/3)^-1 = 3/2（逆数）",
		},
		{
			name: "負のべき乗（高次、未簡約）",
			num:  2, denom: 3,
			exponent: -2,
			wantNum:  9, wantDenom: 4,
			description: "(2/3)^-2 = (3/2)^2 = 9/4",
		},
		{
			name: "3乗",
			num:  1, denom: 2,
			exponent: 3,
			wantNum:  1, wantDenom: 8,
			description: "(1/2)^3 = 1/8",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 実装待ち
		})
	}
}

// TestRationalAssociativity tests associative property
func TestRationalAssociativity(t *testing.T) {
	tests := []struct {
		name        string
		op          string // "add" or "multiply"
		a           struct{ num, denom int64 }
		b           struct{ num, denom int64 }
		c           struct{ num, denom int64 }
		description string
	}{
		{
			name:        "加法の結合律",
			op:          "add",
			a:           struct{ num, denom int64 }{1, 2},
			b:           struct{ num, denom int64 }{1, 3},
			c:           struct{ num, denom int64 }{1, 4},
			description: "(a+b)+c = a+(b+c)",
		},
		{
			name:        "乗法の結合律",
			op:          "multiply",
			a:           struct{ num, denom int64 }{2, 3},
			b:           struct{ num, denom int64 }{3, 4},
			c:           struct{ num, denom int64 }{4, 5},
			description: "(a*b)*c = a*(b*c)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 実装待ち
			// a, _ := NewRational(tt.a.num, tt.a.denom)
			// b, _ := NewRational(tt.b.num, tt.b.denom)
			// c, _ := NewRational(tt.c.num, tt.c.denom)
			//
			// var left, right *Rational
			// if tt.op == "add" {
			//     left = Add(Add(a, b), c)   // (a+b)+c
			//     right = Add(a, Add(b, c))  // a+(b+c)
			// } else { // multiply
			//     left = Multiply(Multiply(a, b), c)   // (a*b)*c
			//     right = Multiply(a, Multiply(b, c))  // a*(b*c)
			// }
			//
			// if !Equal(left, right) {
			//     t.Errorf("結合律が成り立たない: (%s %s %s) %s %s ≠ %s %s (%s %s %s)",
			//         formatRat(a), tt.op, formatRat(b), tt.op, formatRat(c),
			//         formatRat(a), tt.op, formatRat(b), tt.op, formatRat(c))
			// }
		})
	}
}

// TestRationalDistributivity tests distributive property
func TestRationalDistributivity(t *testing.T) {
	tests := []struct {
		name        string
		a           struct{ num, denom int64 }
		b           struct{ num, denom int64 }
		c           struct{ num, denom int64 }
		description string
	}{
		{
			name:        "分配律",
			a:           struct{ num, denom int64 }{2, 3},
			b:           struct{ num, denom int64 }{3, 4},
			c:           struct{ num, denom int64 }{4, 5},
			description: "a*(b+c) = a*b + a*c",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 実装待ち
			// a, _ := NewRational(tt.a.num, tt.a.denom)
			// b, _ := NewRational(tt.b.num, tt.b.denom)
			// c, _ := NewRational(tt.c.num, tt.c.denom)
			//
			// left := Multiply(a, Add(b, c))        // a*(b+c)
			// right := Add(Multiply(a, b), Multiply(a, c))  // a*b + a*c
			//
			// if !Equal(left, right) {
			//     t.Errorf("分配律が成り立たない: a*(b+c) ≠ a*b + a*c")
			// }
		})
	}
}

// TestRationalCommutativity tests commutative property
func TestRationalCommutativity(t *testing.T) {
	tests := []struct {
		name        string
		op          string // "add" or "multiply"
		a           struct{ num, denom int64 }
		b           struct{ num, denom int64 }
		description string
	}{
		{
			name:        "加法の可換性",
			op:          "add",
			a:           struct{ num, denom int64 }{2, 5},
			b:           struct{ num, denom int64 }{3, 7},
			description: "a+b = b+a",
		},
		{
			name:        "乗法の可換性",
			op:          "multiply",
			a:           struct{ num, denom int64 }{2, 5},
			b:           struct{ num, denom int64 }{3, 7},
			description: "a*b = b*a",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 実装待ち
			// a, _ := NewRational(tt.a.num, tt.a.denom)
			// b, _ := NewRational(tt.b.num, tt.b.denom)
			//
			// var left, right *Rational
			// if tt.op == "add" {
			//     left = Add(a, b)   // a+b
			//     right = Add(b, a)  // b+a
			// } else { // multiply
			//     left = Multiply(a, b)   // a*b
			//     right = Multiply(b, a)  // b*a
			// }
			//
			// if !Equal(left, right) {
			//     t.Errorf("可換性が成り立たない")
			// }
		})
	}
}

// TestRationalEdgeCases tests edge cases and boundary conditions
func TestRationalEdgeCases(t *testing.T) {
	tests := []struct {
		name        string
		description string
	}{
		{
			name:        "最大値近辺の計算",
			description: "int64の最大値に近い数での演算オーバーフローテスト",
		},
		{
			name:        "最小値近辺の計算",
			description: "int64の最小値に近い数での演算アンダーフローテスト",
		},
		{
			name:        "繰り返し演算の精度",
			description: "多数回の演算を繰り返しても精度が保たれることの確認",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 実装待ち
		})
	}
}

// TestRationalImmutability tests that rational numbers are immutable
func TestRationalImmutability(t *testing.T) {
	tests := []struct {
		name        string
		description string
	}{
		{
			name:        "加算後の不変性",
			description: "Add()を呼び出しても元のオブジェクトが変更されないことを確認",
		},
		{
			name:        "Simplify後の不変性",
			description: "Simplify()を呼び出しても元のオブジェクトが変更されないことを確認",
		},
		{
			name:        "Normalize後の不変性",
			description: "Normalize()を呼び出しても元のオブジェクトが変更されないことを確認",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 実装待ち
			// r, _ := NewRational(6, 9)
			// originalNum := r.Numerator()
			// originalDenom := r.Denominator()
			//
			// _ = Simplify(r)  // 簡約化を実行
			//
			// // 元のオブジェクトが変更されていないことを確認
			// if r.Numerator() != originalNum || r.Denominator() != originalDenom {
			//     t.Errorf("Simplify()が元のオブジェクトを変更した: %d/%d → %d/%d",
			//         originalNum, originalDenom, r.Numerator(), r.Denominator())
			// }
		})
	}
}

// TestCommonDenominator tests converting two rationals to have a common denominator
// CommonDenominator()は2つの有理数を受け取り、共通の分母(最小公倍数)を持つように変換する
func TestCommonDenominator(t *testing.T) {
	tests := []struct {
		name        string
		r1Num       int64
		r1Denom     int64
		r2Num       int64
		r2Denom     int64
		wantR1Num   int64
		wantR1Denom int64
		wantR2Num   int64
		wantR2Denom int64
		description string
	}{
		{
			name:  "基本的な通分",
			r1Num: 1, r1Denom: 2,
			r2Num: 1, r2Denom: 3,
			wantR1Num: 3, wantR1Denom: 6,
			wantR2Num: 2, wantR2Denom: 6,
			description: "1/2 と 1/3 → 3/6 と 2/6 (LCM=6)",
		},
		{
			name:  "既に同じ分母",
			r1Num: 1, r1Denom: 4,
			r2Num: 3, r2Denom: 4,
			wantR1Num: 1, wantR1Denom: 4,
			wantR2Num: 3, wantR2Denom: 4,
			description: "1/4 と 3/4 → そのまま（既に分母が同じ）",
		},
		{
			name:  "分母が倍数関係",
			r1Num: 1, r1Denom: 3,
			r2Num: 1, r2Denom: 6,
			wantR1Num: 2, wantR1Denom: 6,
			wantR2Num: 1, wantR2Denom: 6,
			description: "1/3 と 1/6 → 2/6 と 1/6 (LCM=6)",
		},
		{
			name:  "分母が互いに素",
			r1Num: 2, r1Denom: 5,
			r2Num: 3, r2Denom: 7,
			wantR1Num: 14, wantR1Denom: 35,
			wantR2Num: 15, wantR2Denom: 35,
			description: "2/5 と 3/7 → 14/35 と 15/35 (LCM=35)",
		},
		{
			name:  "整数との通分",
			r1Num: 1, r1Denom: 2,
			r2Num: 3, r2Denom: 1,
			wantR1Num: 1, wantR1Denom: 2,
			wantR2Num: 6, wantR2Denom: 2,
			description: "1/2 と 3/1 → 1/2 と 6/2 (LCM=2)",
		},
		{
			name:  "負の数を含む通分",
			r1Num: -1, r1Denom: 2,
			r2Num: 1, r2Denom: 3,
			wantR1Num: -3, wantR1Denom: 6,
			wantR2Num: 2, wantR2Denom: 6,
			description: "-1/2 と 1/3 → -3/6 と 2/6 (LCM=6)",
		},
		{
			name:  "両方負の数",
			r1Num: -2, r1Denom: 3,
			r2Num: -3, r2Denom: 4,
			wantR1Num: -8, wantR1Denom: 12,
			wantR2Num: -9, wantR2Denom: 12,
			description: "-2/3 と -3/4 → -8/12 と -9/12 (LCM=12)",
		},
		{
			name:  "ゼロを含む通分",
			r1Num: 0, r1Denom: 1,
			r2Num: 1, r2Denom: 3,
			wantR1Num: 0, wantR1Denom: 3,
			wantR2Num: 1, wantR2Denom: 3,
			description: "0/1 と 1/3 → 0/3 と 1/3 (LCM=3)",
		},
		{
			name:  "複雑な分母のGCD",
			r1Num: 1, r1Denom: 12,
			r2Num: 1, r2Denom: 18,
			wantR1Num: 3, wantR1Denom: 36,
			wantR2Num: 2, wantR2Denom: 36,
			description: "1/12 と 1/18 → 3/36 と 2/36 (LCM=36, GCD=6)",
		},
		{
			name:  "負の分母を含む通分（符号を分子に）",
			r1Num: 1, r1Denom: -2,
			r2Num: 1, r2Denom: 3,
			wantR1Num: -3, wantR1Denom: 6,
			wantR2Num: 2, wantR2Denom: 6,
			description: "1/-2 と 1/3 → -3/6 と 2/6 (符号を分子に移動、分母は正)",
		},
		{
			name:  "大きな分母の通分",
			r1Num: 1, r1Denom: 100,
			r2Num: 1, r2Denom: 150,
			wantR1Num: 3, wantR1Denom: 300,
			wantR2Num: 2, wantR2Denom: 300,
			description: "1/100 と 1/150 → 3/300 と 2/300 (LCM=300)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 実装待ち
			r1, _ := NewRational(tt.r1Num, tt.r1Denom)
			r2, _ := NewRational(tt.r2Num, tt.r2Denom)

			result1, result2 := CommonDenominator(r1, r2)

			if result1.Numerator != tt.wantR1Num || result1.Denominator != tt.wantR1Denom {
				t.Errorf("CommonDenominator() result1 = %d/%d, want %d/%d",
					result1.Numerator, result1.Denominator, tt.wantR1Num, tt.wantR1Denom)
			}
			if result2.Numerator != tt.wantR2Num || result2.Denominator != tt.wantR2Denom {
				t.Errorf("CommonDenominator() result2 = %d/%d, want %d/%d",
					result2.Numerator, result2.Denominator, tt.wantR2Num, tt.wantR2Denom)
			}

			// 元のオブジェクトが変更されていないことを確認（イミュータブル性）
			if r1.Numerator != tt.r1Num || r1.Denominator != tt.r1Denom {
				t.Errorf("CommonDenominator()がr1を変更した: %d/%d → %d/%d",
					tt.r1Num, tt.r1Denom, r1.Numerator, r1.Denominator)
			}
			if r2.Numerator != tt.r2Num || r2.Denominator != tt.r2Denom {
				t.Errorf("CommonDenominator()がr2を変更した: %d/%d → %d/%d",
					tt.r2Num, tt.r2Denom, r2.Numerator, r2.Denominator)
			}
		})
	}
}

// TestCommonDenominatorEquivalence tests that common denominator preserves value
// 通分後も元の値と等価であることを確認
func TestCommonDenominatorEquivalence(t *testing.T) {
	tests := []struct {
		name        string
		r1Num       int64
		r1Denom     int64
		r2Num       int64
		r2Denom     int64
		description string
	}{
		{
			name:  "通分後の値の等価性",
			r1Num: 1, r1Denom: 2,
			r2Num: 1, r2Denom: 3,
			description: "1/2と1/3を通分しても、それぞれの値は変わらない",
		},
		{
			name:  "複雑な分数での等価性",
			r1Num: 5, r1Denom: 12,
			r2Num: 7, r2Denom: 18,
			description: "5/12と7/18を通分しても、それぞれの値は変わらない",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 実装待ち
			// r1, _ := NewRational(tt.r1Num, tt.r1Denom)
			// r2, _ := NewRational(tt.r2Num, tt.r2Denom)

			// result1, result2 := CommonDenominator(r1, r2)

			// // 通分前後で値が等しいことを確認
			// if !Equal(r1, result1) {
			// 	t.Errorf("通分後の値が変わった: %d/%d ≠ %d/%d",
			// 		r1.Numerator, r1.Denominator, result1.Numerator, result1.Denominator)
			// }
			// if !Equal(r2, result2) {
			// 	t.Errorf("通分後の値が変わった: %d/%d ≠ %d/%d",
			// 		r2.Numerator, r2.Denominator, result2.Numerator, result2.Denominator)
			// }

			// // 通分後の分母が同じことを確認
			// if result1.Denominator != result2.Denominator {
			// 	t.Errorf("通分後の分母が異なる: %d ≠ %d",
			// 		result1.Denominator, result2.Denominator)
			// }
		})
	}
}
