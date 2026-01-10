package expr

import "math"

// GCD はユークリッドの互除法を用いて2つのint64値の最大公約数を計算します。
// 負の入力に対しては、その絶対値の最大公約数を返します。
// GCD(0, 0) は慣例により0を返します。
// この関数は有理数演算で内部的に使用されますが、単独でも使用できます。
func GCD(a, b int64) int64 {
	// TODO: ユークリッドの互除法を実装
	if a < 0 {
		a = -a
	}
	if b < 0 {
		b = -b
	}
	if a == 0 && b == 0 {
		return 0
	}
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

// LCM は2つのint64値の最小公倍数を計算します。
// 関係式 LCM(a, b) * GCD(a, b) = |a * b| を使用します。
// LCM(a, 0) または LCM(0, b) は0を返します。
func LCM(a, b int64) int64 {
	// TODO: GCDを使用して実装
	if a == 0 || b == 0 {
		return 0
	}
	gcd := GCD(a, b)
	return int64(math.Abs(float64(a*b)) / float64(gcd))
}

// IsPrime はnが素数かどうかを判定します。
// n <= 1 の場合はfalseを返します。
// sqrt(n)までの試し割りを使用します。
func IsPrime(n int64) bool {
	// TODO: 試し割り法を実装
	for i := int64(2); i*i <= n; i++ {
		if n%i == 0 {
			return false
		}
	}
	return n > 1
}

// PrimeFactorization はnの素因数分解を素因数のスライスとして返します。
// 素因数は昇順で重複を含めて返されます。
// n = 1 の場合は空のスライスを返します。
// n <= 0 の場合の動作は未定義です（panicまたはnilを返す）。
//
// 例: PrimeFactorization(60) は [2, 2, 3, 5] を返します（60 = 2^2 * 3 * 5）
func PrimeFactorization(n int64) []int64 {
	// TODO: 試し割り法による素因数分解を実装
	if n <= 0 {
		panic("PrimeFactorization: n must be positive")
	}
	if n == 1 {
		return []int64{}
	}
	var factors []int64
	for i := int64(2); i*i <= n; i++ {
		for n%i == 0 {
			factors = append(factors, i)
			n /= i
		}
	}
	if n > 1 {
		factors = append(factors, n)
	}
	return factors
}

// Divisors はnのすべての正の約数を昇順で返します。
// n = 1 の場合は [1] を返します。
// n <= 0 の場合の動作は未定義です。
//
// 例: Divisors(12) は [1, 2, 3, 4, 6, 12] を返します
func Divisors(n int64) []int64 {
	if n <= 0 {
		panic("Divisors: n must be positive")
	}
	var divisors []int64
	for i := int64(1); i <= n; i++ {
		if n%i == 0 {
			divisors = append(divisors, i)
		}
	}
	return divisors
}

// NumDivisors はnの正の約数の個数を返します。
// 素因数分解を使用して効率的に計算できます:
// n = p1^a1 * p2^a2 * ... * pk^ak のとき、
// NumDivisors(n) = (a1+1) * (a2+1) * ... * (ak+1)
func NumDivisors(n int64) int {
	// TODO: 素因数分解または直接カウントを使用して実装
	panic("not implemented")
}

// IsPerfectNumber はnが完全数かどうかを判定します。
// 完全数は自分自身を除く約数の和に等しい数です。
// 例: 6 = 1+2+3, 28 = 1+2+4+7+14
func IsPerfectNumber(n int64) bool {
	// TODO: 約数の和を計算して実装
	panic("not implemented")
}

// EulerPhi はオイラーのトーシェント関数φ(n)を計算します。
// 1からnまでの整数のうち、nと互いに素なものの個数を返します。
// φ(n) = n * ∏(1 - 1/p) （すべての素因数pに対して）を使用して計算できます。
//
// 性質:
// - φ(1) = 1
// - φ(p) = p-1 （pが素数のとき）
// - φ(mn) = φ(m)φ(n) （gcd(m,n) = 1のとき、乗法的）
func EulerPhi(n int64) int64 {
	// TODO: 素因数分解を使用して実装
	panic("not implemented")
}

// IsCoprime は2つのint64値が互いに素（relatively prime）かどうかを判定します。
// 2つの数が互いに素であるとは、GCD(a, b) = 1 であることです。
func IsCoprime(a, b int64) bool {
	// TODO: GCDを使用して実装
	panic("not implemented")
}
