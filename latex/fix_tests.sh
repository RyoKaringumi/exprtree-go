#!/bin/bash
# Fix converter tests to handle interface{} return type

# Pattern 1: Replace "expression, err := converter.Convert" with "result, err := converter.Convert"
# and add type assertion for Expression types that need Eval()
sed -i '
/expression, err := converter.Convert(node)/,/^}$/ {
    s/expression, err := converter.Convert(node)/result, err := converter.Convert(node)/
    /if err != nil {/a\
\
	expression, ok := result.(expr.Expression)\
	if !ok {\
		t.Fatalf("expected Expression, got %T", result)\
	}
}
' converter_test.go
