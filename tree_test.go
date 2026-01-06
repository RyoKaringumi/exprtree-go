package main

import "testing"

func TestConstantEval(t *testing.T) {
	constant := &Constant{Value: NumberValue{Value: 5.0}}
	result, ok := constant.Eval()
	if !ok {
		t.Errorf("Expected evaluation to succeed")
	}
	if num, ok := result.(*NumberValue); ok {
		if num.Value != 5.0 {
			t.Errorf("Expected 5.0, got %f", num.Value)
		}
	} else {
		t.Errorf("Expected NumberValue")
	}
}

func TestAddExpression(t *testing.T) {
	left := &Constant{Value: NumberValue{Value: 3.0}}
	right := &Constant{Value: NumberValue{Value: 4.0}}
	add := &AddExpression{Left: left, Right: right}
	result, ok := add.Eval()
	if !ok {
		t.Errorf("Expected evaluation to succeed")
	}
	if num, ok := result.(*NumberValue); ok {
		if num.Value != 7.0 {
			t.Errorf("Expected 7.0, got %f", num.Value)
		}
	} else {
		t.Errorf("Expected NumberValue")
	}
}

func TestSubtractExpression(t *testing.T) {
	left := &Constant{Value: NumberValue{Value: 10.0}}
	right := &Constant{Value: NumberValue{Value: 3.0}}
	sub := &SubtractExpression{Left: left, Right: right}
	result, ok := sub.Eval()
	if !ok {
		t.Errorf("Expected evaluation to succeed")
	}
	if num, ok := result.(*NumberValue); ok {
		if num.Value != 7.0 {
			t.Errorf("Expected 7.0, got %f", num.Value)
		}
	} else {
		t.Errorf("Expected NumberValue")
	}
}

func TestMultiplyExpression(t *testing.T) {
	left := &Constant{Value: NumberValue{Value: 6.0}}
	right := &Constant{Value: NumberValue{Value: 7.0}}
	mul := &MultiplyExpression{Left: left, Right: right}
	result, ok := mul.Eval()
	if !ok {
		t.Errorf("Expected evaluation to succeed")
	}
	if num, ok := result.(*NumberValue); ok {
		if num.Value != 42.0 {
			t.Errorf("Expected 42.0, got %f", num.Value)
		}
	} else {
		t.Errorf("Expected NumberValue")
	}
}

func TestDivideExpression(t *testing.T) {
	left := &Constant{Value: NumberValue{Value: 15.0}}
	right := &Constant{Value: NumberValue{Value: 3.0}}
	div := &DivideExpression{Left: left, Right: right}
	result, ok := div.Eval()
	if !ok {
		t.Errorf("Expected evaluation to succeed")
	}
	if num, ok := result.(*NumberValue); ok {
		if num.Value != 5.0 {
			t.Errorf("Expected 5.0, got %f", num.Value)
		}
	} else {
		t.Errorf("Expected NumberValue")
	}
}

func TestDivideByZero(t *testing.T) {
	left := &Constant{Value: NumberValue{Value: 10.0}}
	right := &Constant{Value: NumberValue{Value: 0.0}}
	div := &DivideExpression{Left: left, Right: right}
	_, ok := div.Eval()
	if ok {
		t.Errorf("Expected evaluation to fail due to division by zero")
	}
}
