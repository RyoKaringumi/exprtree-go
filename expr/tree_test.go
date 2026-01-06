package expr

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
	add := NewAddExpression(left, right)
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
	sub := NewSubtractExpression(left, right)
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
	mul := NewMultiplyExpression(left, right)
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
	div := NewDivideExpression(left, right)
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
	div := NewDivideExpression(left, right)
	_, ok := div.Eval()
	if ok {
		t.Errorf("Expected evaluation to fail due to division by zero")
	}
}

func TestAddExpressionChildren(t *testing.T) {
	left := &Constant{Value: NumberValue{Value: 3.0}}
	right := &Constant{Value: NumberValue{Value: 4.0}}
	add := NewAddExpression(left, right)
	children := add.Children()
	if len(children) != 2 {
		t.Errorf("Expected 2 children, got %d", len(children))
	}
	if children[0] != left || children[1] != right {
		t.Errorf("Children do not match expected left and right")
	}
}

func TestConstantChildren(t *testing.T) {
	constant := &Constant{Value: NumberValue{Value: 5.0}}
	children := constant.Children()
	if len(children) != 0 {
		t.Errorf("Expected 0 children for Constant, got %d", len(children))
	}
}

func TestVariableEval(t *testing.T) {
	variable := &Variable{Name: "x"}
	_, ok := variable.Eval()
	if ok {
		t.Errorf("Expected evaluation to fail for variable without context")
	}
}

func TestVariableChildren(t *testing.T) {
	variable := &Variable{Name: "x"}
	children := variable.Children()
	if len(children) != 0 {
		t.Errorf("Expected 0 children for Variable, got %d", len(children))
	}
}

func TestExpressionWithVariable(t *testing.T) {
	// x + 2 should fail to evaluate because x has no value
	left := &Variable{Name: "x"}
	right := &Constant{Value: NumberValue{Value: 2.0}}
	add := NewAddExpression(left, right)
	_, ok := add.Eval()
	if ok {
		t.Errorf("Expected evaluation to fail due to variable")
	}
}
