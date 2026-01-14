package value

type Value interface {
	Kind() ValueKind
	Equals(other any) bool
}

type ValueKind int

const (
	RealKind ValueKind = iota

	// 真偽値
	BoolKind
)

type RealValue struct {
	Value
	v float64
}

func NewRealValue(value float64) *RealValue {
	return &RealValue{
		v: value,
	}
}

func (r *RealValue) Kind() ValueKind {
	return RealKind
}

func (r *RealValue) Float64() float64 {
	return r.v
}

func (r *RealValue) Eval() (Value, bool) {
	return r, true
}

func (r *RealValue) Equals(other any) bool {
	otherReal, ok := other.(*RealValue)
	if !ok {
		return false
	}
	return r.v == otherReal.v
}

type BoolValue struct {
	Value
	v bool
}

func NewBoolValue(value bool) *BoolValue {
	return &BoolValue{
		v: value,
	}
}

func (b *BoolValue) Kind() ValueKind {
	return BoolKind
}

func (b *BoolValue) Bool() bool {
	return b.v
}

func (b *BoolValue) Eval() (Value, bool) {
	return b, true
}

func (b *BoolValue) Equals(other any) bool {
	otherBool, ok := other.(*BoolValue)
	if !ok {
		return false
	}
	return b.v == otherBool.v
}
