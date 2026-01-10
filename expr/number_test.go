package expr

import (
	"reflect"
	"testing"
)

// 注: TestGCDはreal_test.goにあります（number演算と有理数演算で共有されるため）

// TestLCM は最小公倍数関数のテスト
func TestLCM(t *testing.T) {
	tests := []struct {
		name     string
		a, b     int64
		expected int64
	}{
		{"両方正の数", 12, 18, 36},
		{"一方が他方の倍数", 25, 100, 100},
		{"互いに素", 7, 11, 77},
		{"同じ数", 15, 15, 15},
		{"一方がゼロ", 5, 0, 0},
		{"小さな素数", 3, 5, 15},
		{"素数のべき乗", 8, 12, 24},
		{"大きな互いに素な数", 13, 17, 221},
		{"1との計算", 1, 42, 42},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := LCM(tt.a, tt.b)
			if result != tt.expected {
				t.Errorf("LCM(%d, %d) = %d; expected %d", tt.a, tt.b, result, tt.expected)
			}
		})
	}
}

// TestIsPrime は素数判定関数のテスト
func TestIsPrime(t *testing.T) {
	tests := []struct {
		name     string
		n        int64
		expected bool
	}{
		{"2は素数", 2, true},
		{"3は素数", 3, true},
		{"4は素数でない", 4, false},
		{"5は素数", 5, true},
		{"1は素数でない", 1, false},
		{"0は素数でない", 0, false},
		{"負の数は素数でない", -5, false},
		{"17は素数", 17, true},
		{"100は素数でない", 100, false},
		{"97は素数", 97, true},
		{"大きな素数", 7919, true},
		{"大きな合成数", 7920, false},
		{"完全平方数", 49, false},
		{"メルセンヌ素数", 31, true},
		{"2より大きい偶数", 100, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsPrime(tt.n)
			if result != tt.expected {
				t.Errorf("IsPrime(%d) = %v; expected %v", tt.n, result, tt.expected)
			}
		})
	}
}

// TestPrimeFactorization は素因数分解関数のテスト
func TestPrimeFactorization(t *testing.T) {
	tests := []struct {
		name     string
		n        int64
		expected []int64
	}{
		{"素数", 7, []int64{7}},
		{"素数のべき乗", 8, []int64{2, 2, 2}},
		{"2つの異なる素数", 6, []int64{2, 3}},
		{"複数の素因数", 60, []int64{2, 2, 3, 5}},
		{"大きな素数", 97, []int64{97}},
		{"完全平方数", 36, []int64{2, 2, 3, 3}},
		{"3つの素数の積", 30, []int64{2, 3, 5}},
		{"2のべき乗", 64, []int64{2, 2, 2, 2, 2, 2}},
		{"大きな合成数", 1001, []int64{7, 11, 13}},
		{"数1", 1, []int64{}},
		{"小さな合成数", 12, []int64{2, 2, 3}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := PrimeFactorization(tt.n)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("PrimeFactorization(%d) = %v; expected %v", tt.n, result, tt.expected)
			}
		})
	}
}

// TestDivisors は約数列挙関数のテスト
func TestDivisors(t *testing.T) {
	tests := []struct {
		name     string
		n        int64
		expected []int64
	}{
		{"素数", 7, []int64{1, 7}},
		{"完全平方数", 16, []int64{1, 2, 4, 8, 16}},
		{"小さな合成数", 12, []int64{1, 2, 3, 4, 6, 12}},
		{"数1", 1, []int64{1}},
		{"素数のべき乗", 27, []int64{1, 3, 9, 27}},
		{"2つの素数の積", 15, []int64{1, 3, 5, 15}},
		{"高度合成数", 24, []int64{1, 2, 3, 4, 6, 8, 12, 24}},
		{"大きな素数", 97, []int64{1, 97}},
		{"完全数", 6, []int64{1, 2, 3, 6}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Divisors(tt.n)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Divisors(%d) = %v; expected %v", tt.n, result, tt.expected)
			}
		})
	}
}

// TestNumDivisors は約数の個数を返す関数のテスト
func TestNumDivisors(t *testing.T) {
	tests := []struct {
		name     string
		n        int64
		expected int
	}{
		{"素数", 7, 2},
		{"完全平方数", 16, 5},
		{"小さな合成数", 12, 6},
		{"数1", 1, 1},
		{"素数のべき乗", 27, 4},
		{"2つの素数の積", 15, 4},
		{"高度合成数", 24, 8},
		{"大きな素数", 97, 2},
		{"完全数", 28, 6},
		{"2のべき乗", 64, 7},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NumDivisors(tt.n)
			if result != tt.expected {
				t.Errorf("NumDivisors(%d) = %d; expected %d", tt.n, result, tt.expected)
			}
		})
	}
}

// TestIsPerfectNumber は完全数判定関数のテスト
func TestIsPerfectNumber(t *testing.T) {
	tests := []struct {
		name     string
		n        int64
		expected bool
	}{
		{"6は完全数", 6, true},
		{"28は完全数", 28, true},
		{"496は完全数", 496, true},
		{"8128は完全数", 8128, true},
		{"12は完全数でない", 12, false},
		{"1は完全数でない", 1, false},
		{"素数は完全数でない", 7, false},
		{"2のべき乗は完全数でない", 16, false},
		{"大きな非完全数", 1000, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsPerfectNumber(tt.n)
			if result != tt.expected {
				t.Errorf("IsPerfectNumber(%d) = %v; expected %v", tt.n, result, tt.expected)
			}
		})
	}
}

// TestEulerPhi はオイラーのトーシェント関数のテスト
func TestEulerPhi(t *testing.T) {
	tests := []struct {
		name     string
		n        int64
		expected int64
	}{
		{"素数", 7, 6},
		{"素数のべき乗", 9, 6},
		{"2つの素数の積", 15, 8},
		{"数1", 1, 1},
		{"2のべき乗", 16, 8},
		{"小さな合成数", 12, 4},
		{"大きな素数", 97, 96},
		{"完全平方数", 36, 12},
		{"3つの素数の積", 30, 8},
		{"すべての下位数と互いに素", 2, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := EulerPhi(tt.n)
			if result != tt.expected {
				t.Errorf("EulerPhi(%d) = %d; expected %d", tt.n, result, tt.expected)
			}
		})
	}
}

// TestIsCoprime は互いに素判定関数のテスト
func TestIsCoprime(t *testing.T) {
	tests := []struct {
		name     string
		a, b     int64
		expected bool
	}{
		{"互いに素な素数", 7, 11, true},
		{"互いに素でない", 12, 18, false},
		{"連続する数", 14, 15, true},
		{"同じ数", 7, 7, false},
		{"一方が1", 1, 42, true},
		{"べき乗と素数", 8, 9, true},
		{"同じ素数の倍数", 6, 9, false},
		{"大きな互いに素な数", 97, 101, true},
		{"フィボナッチ数列の連続項", 89, 55, true},
		{"両方偶数", 4, 6, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsCoprime(tt.a, tt.b)
			if result != tt.expected {
				t.Errorf("IsCoprime(%d, %d) = %v; expected %v", tt.a, tt.b, result, tt.expected)
			}
		})
	}
}

// TestGCDProperties はGCDの数学的性質のテスト
func TestGCDProperties(t *testing.T) {
	t.Run("可換性", func(t *testing.T) {
		var a, b int64 = 48, 18
		if GCD(a, b) != GCD(b, a) {
			t.Error("GCDは可換であるべき")
		}
	})

	t.Run("結合性", func(t *testing.T) {
		var a, b, c int64 = 48, 18, 30
		if GCD(GCD(a, b), c) != GCD(a, GCD(b, c)) {
			t.Error("GCDは結合的であるべき")
		}
	})

	t.Run("0との単位元", func(t *testing.T) {
		var a int64 = 42
		if GCD(a, 0) != a {
			t.Error("GCD(a, 0) は a と等しいべき")
		}
	})
}

// TestLCMProperties はLCMの数学的性質のテスト
func TestLCMProperties(t *testing.T) {
	t.Run("可換性", func(t *testing.T) {
		var a, b int64 = 12, 18
		if LCM(a, b) != LCM(b, a) {
			t.Error("LCMは可換であるべき")
		}
	})

	t.Run("GCD-LCM関係式", func(t *testing.T) {
		var a, b int64 = 12, 18
		if a*b != GCD(a, b)*LCM(a, b) {
			t.Error("a*b は GCD(a,b)*LCM(a,b) と等しいべき")
		}
	})
}

// TestPrimeFactorizationProduct は素因数の積が元の数に等しいことをテスト
func TestPrimeFactorizationProduct(t *testing.T) {
	tests := []int64{12, 60, 97, 1001, 8, 36}

	for _, n := range tests {
		t.Run(string(rune(n)), func(t *testing.T) {
			factors := PrimeFactorization(n)
			var product int64 = 1
			for _, f := range factors {
				product *= f
			}
			if product != n {
				t.Errorf("%dの素因数の積は%d; 期待値 %d", n, product, n)
			}
		})
	}
}

// TestDivisorsSorted は約数がソート順で返されることをテスト
func TestDivisorsSorted(t *testing.T) {
	tests := []int64{12, 24, 36, 60}

	for _, n := range tests {
		t.Run(string(rune(n)), func(t *testing.T) {
			divisors := Divisors(n)
			for i := 1; i < len(divisors); i++ {
				if divisors[i-1] >= divisors[i] {
					t.Errorf("%dの約数がソートされていない: %v", n, divisors)
					break
				}
			}
		})
	}
}

// TestEulerPhiMultiplicative は互いに素な数に対してphiが乗法的であることをテスト
func TestEulerPhiMultiplicative(t *testing.T) {
	tests := []struct {
		a, b int64
	}{
		{3, 5},
		{7, 11},
		{4, 9},
	}

	for _, tt := range tests {
		if !IsCoprime(tt.a, tt.b) {
			t.Fatalf("%dと%dは互いに素でない", tt.a, tt.b)
		}
		phiProduct := EulerPhi(tt.a) * EulerPhi(tt.b)
		phiAB := EulerPhi(tt.a * tt.b)
		if phiProduct != phiAB {
			t.Errorf("phi(%d)*phi(%d) = %d; phi(%d) = %d; 互いに素な数では等しいべき",
				tt.a, tt.b, phiProduct, tt.a*tt.b, phiAB)
		}
	}
}
