package latex

import "testing"

func TestLexer_Numbers(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
	}{
		{"123", 123.0},
		{"3.14", 3.14},
		{"0.5", 0.5},
		{"42", 42.0},
		{"0", 0.0},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			lexer := NewLexer(tt.input)
			tok := lexer.NextToken()

			if tok.Type != NUMBER {
				t.Errorf("expected NUMBER token, got %v", tok.Type)
			}
			if tok.Value != tt.expected {
				t.Errorf("expected value %f, got %f", tt.expected, tok.Value)
			}
			if tok.Literal != tt.input {
				t.Errorf("expected literal %s, got %s", tt.input, tok.Literal)
			}
		})
	}
}

func TestLexer_Operators(t *testing.T) {
	tests := []struct {
		input    string
		expected TokenType
	}{
		{"+", PLUS},
		{"-", MINUS},
		{"*", MULTIPLY},
		{"/", DIVIDE},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			lexer := NewLexer(tt.input)
			tok := lexer.NextToken()

			if tok.Type != tt.expected {
				t.Errorf("expected token type %v, got %v", tt.expected, tok.Type)
			}
			if tok.Literal != tt.input {
				t.Errorf("expected literal %s, got %s", tt.input, tok.Literal)
			}
		})
	}
}

func TestLexer_Parentheses(t *testing.T) {
	tests := []struct {
		input    string
		expected TokenType
	}{
		{"(", LPAREN},
		{")", RPAREN},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			lexer := NewLexer(tt.input)
			tok := lexer.NextToken()

			if tok.Type != tt.expected {
				t.Errorf("expected token type %v, got %v", tt.expected, tok.Type)
			}
		})
	}
}

func TestLexer_Expression(t *testing.T) {
	input := "2 + 3 * 4"
	lexer := NewLexer(input)

	expectedTokens := []struct {
		tokenType TokenType
		literal   string
		value     float64
	}{
		{NUMBER, "2", 2.0},
		{PLUS, "+", 0},
		{NUMBER, "3", 3.0},
		{MULTIPLY, "*", 0},
		{NUMBER, "4", 4.0},
		{EOF, "", 0},
	}

	for i, expected := range expectedTokens {
		tok := lexer.NextToken()

		if tok.Type != expected.tokenType {
			t.Errorf("token[%d] - expected type %v, got %v", i, expected.tokenType, tok.Type)
		}
		if tok.Literal != expected.literal {
			t.Errorf("token[%d] - expected literal %s, got %s", i, expected.literal, tok.Literal)
		}
		if expected.tokenType == NUMBER && tok.Value != expected.value {
			t.Errorf("token[%d] - expected value %f, got %f", i, expected.value, tok.Value)
		}
	}
}

func TestLexer_Whitespace(t *testing.T) {
	input := "  2  +  3  "
	lexer := NewLexer(input)

	expectedTokens := []TokenType{NUMBER, PLUS, NUMBER, EOF}

	for i, expected := range expectedTokens {
		tok := lexer.NextToken()
		if tok.Type != expected {
			t.Errorf("token[%d] - expected type %v, got %v", i, expected, tok.Type)
		}
	}
}

func TestLexer_ComplexExpression(t *testing.T) {
	input := "(2 + 3) * 4 / 2"
	lexer := NewLexer(input)

	expectedTokens := []TokenType{
		LPAREN, NUMBER, PLUS, NUMBER, RPAREN,
		MULTIPLY, NUMBER, DIVIDE, NUMBER, EOF,
	}

	for i, expected := range expectedTokens {
		tok := lexer.NextToken()
		if tok.Type != expected {
			t.Errorf("token[%d] - expected type %v, got %v", i, expected, tok.Type)
		}
	}
}

func TestLexer_EmptyInput(t *testing.T) {
	lexer := NewLexer("")
	tok := lexer.NextToken()

	if tok.Type != EOF {
		t.Errorf("expected EOF for empty input, got %v", tok.Type)
	}
}

func TestLexer_IllegalCharacter(t *testing.T) {
	input := "2 @ 3"
	lexer := NewLexer(input)

	tokens := []Token{}
	for {
		tok := lexer.NextToken()
		tokens = append(tokens, tok)
		if tok.Type == EOF {
			break
		}
	}

	if tokens[1].Type != ILLEGAL {
		t.Errorf("expected ILLEGAL token for '@', got %v", tokens[1].Type)
	}
}

func TestLexer_DecimalNumbers(t *testing.T) {
	input := "2.5 + 1.5"
	lexer := NewLexer(input)

	tok1 := lexer.NextToken()
	if tok1.Type != NUMBER || tok1.Value != 2.5 {
		t.Errorf("expected NUMBER 2.5, got %v with value %f", tok1.Type, tok1.Value)
	}

	tok2 := lexer.NextToken()
	if tok2.Type != PLUS {
		t.Errorf("expected PLUS, got %v", tok2.Type)
	}

	tok3 := lexer.NextToken()
	if tok3.Type != NUMBER || tok3.Value != 1.5 {
		t.Errorf("expected NUMBER 1.5, got %v with value %f", tok3.Type, tok3.Value)
	}
}

func TestLexer_Position(t *testing.T) {
	input := "2+3"
	lexer := NewLexer(input)

	tok1 := lexer.NextToken()
	if tok1.Pos != 0 {
		t.Errorf("expected position 0 for first token, got %d", tok1.Pos)
	}

	tok2 := lexer.NextToken()
	if tok2.Pos != 1 {
		t.Errorf("expected position 1 for second token, got %d", tok2.Pos)
	}

	tok3 := lexer.NextToken()
	if tok3.Pos != 2 {
		t.Errorf("expected position 2 for third token, got %d", tok3.Pos)
	}
}

func TestLexer_Variables(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"x", "x"},
		{"y", "y"},
		{"a", "a"},
		{"Z", "Z"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			lexer := NewLexer(tt.input)
			tok := lexer.NextToken()

			if tok.Type != VARIABLE {
				t.Errorf("expected VARIABLE token, got %v", tok.Type)
			}
			if tok.Literal != tt.expected {
				t.Errorf("expected literal %s, got %s", tt.expected, tok.Literal)
			}
		})
	}
}

func TestLexer_VariableExpression(t *testing.T) {
	input := "x + 2"
	lexer := NewLexer(input)

	expectedTokens := []struct {
		tokenType TokenType
		literal   string
	}{
		{VARIABLE, "x"},
		{PLUS, "+"},
		{NUMBER, "2"},
		{EOF, ""},
	}

	for i, expected := range expectedTokens {
		tok := lexer.NextToken()

		if tok.Type != expected.tokenType {
			t.Errorf("token[%d] - expected type %v, got %v", i, expected.tokenType, tok.Type)
		}
		if tok.Literal != expected.literal {
			t.Errorf("token[%d] - expected literal %s, got %s", i, expected.literal, tok.Literal)
		}
	}
}

func TestLexer_MultiCharacterIdentifier(t *testing.T) {
	// Multi-character sequences are tokenized as separate single-character variables
	// (implicit multiplication is handled by the parser)
	input := "abc"
	lexer := NewLexer(input)

	expectedTokens := []struct {
		tokenType TokenType
		literal   string
	}{
		{VARIABLE, "a"},
		{VARIABLE, "b"},
		{VARIABLE, "c"},
		{EOF, ""},
	}

	for i, expected := range expectedTokens {
		tok := lexer.NextToken()
		if tok.Type != expected.tokenType {
			t.Errorf("token[%d] - expected type %v, got %v", i, expected.tokenType, tok.Type)
		}
		if tok.Literal != expected.literal {
			t.Errorf("token[%d] - expected literal %q, got %q", i, expected.literal, tok.Literal)
		}
	}
}

func TestLexer_Caret(t *testing.T) {
	input := "2^3"
	lexer := NewLexer(input)

	expectedTokens := []TokenType{NUMBER, CARET, NUMBER, EOF}

	for i, expected := range expectedTokens {
		tok := lexer.NextToken()
		if tok.Type != expected {
			t.Errorf("token[%d] - expected type %v, got %v", i, expected, tok.Type)
		}
	}
}

func TestLexer_Braces(t *testing.T) {
	input := "{2+3}"
	lexer := NewLexer(input)

	expectedTokens := []TokenType{LBRACE, NUMBER, PLUS, NUMBER, RBRACE, EOF}

	for i, expected := range expectedTokens {
		tok := lexer.NextToken()
		if tok.Type != expected {
			t.Errorf("token[%d] - expected type %v, got %v", i, expected, tok.Type)
		}
	}
}

func TestLexer_Brackets(t *testing.T) {
	input := "[3]"
	lexer := NewLexer(input)

	expectedTokens := []TokenType{LBRACKET, NUMBER, RBRACKET, EOF}

	for i, expected := range expectedTokens {
		tok := lexer.NextToken()
		if tok.Type != expected {
			t.Errorf("token[%d] - expected type %v, got %v", i, expected, tok.Type)
		}
	}
}

func TestLexer_SqrtCommand(t *testing.T) {
	input := "\\sqrt{4}"
	lexer := NewLexer(input)

	tok1 := lexer.NextToken()
	if tok1.Type != COMMAND {
		t.Errorf("expected COMMAND token, got %v", tok1.Type)
	}
	if tok1.Literal != "sqrt" {
		t.Errorf("expected literal 'sqrt', got '%s'", tok1.Literal)
	}

	tok2 := lexer.NextToken()
	if tok2.Type != LBRACE {
		t.Errorf("expected LBRACE, got %v", tok2.Type)
	}

	tok3 := lexer.NextToken()
	if tok3.Type != NUMBER || tok3.Value != 4.0 {
		t.Errorf("expected NUMBER 4, got %v with value %f", tok3.Type, tok3.Value)
	}

	tok4 := lexer.NextToken()
	if tok4.Type != RBRACE {
		t.Errorf("expected RBRACE, got %v", tok4.Type)
	}
}

func TestLexer_SqrtWithOptional(t *testing.T) {
	input := "\\sqrt[3]{8}"
	lexer := NewLexer(input)

	expectedTokens := []struct {
		tokenType TokenType
		literal   string
		value     float64
	}{
		{COMMAND, "sqrt", 0},
		{LBRACKET, "[", 0},
		{NUMBER, "3", 3.0},
		{RBRACKET, "]", 0},
		{LBRACE, "{", 0},
		{NUMBER, "8", 8.0},
		{RBRACE, "}", 0},
		{EOF, "", 0},
	}

	for i, expected := range expectedTokens {
		tok := lexer.NextToken()
		if tok.Type != expected.tokenType {
			t.Errorf("token[%d] - expected type %v, got %v", i, expected.tokenType, tok.Type)
		}
		if tok.Literal != expected.literal && expected.literal != "" {
			t.Errorf("token[%d] - expected literal %s, got %s", i, expected.literal, tok.Literal)
		}
	}
}

func TestLexer_UnknownCommand(t *testing.T) {
	input := "\\unknown"
	lexer := NewLexer(input)

	tok := lexer.NextToken()
	if tok.Type != ILLEGAL {
		t.Errorf("expected ILLEGAL token for unknown command, got %v", tok.Type)
	}
}

func TestLexer_ComplexPowerExpression(t *testing.T) {
	input := "2^3^4"
	lexer := NewLexer(input)

	expectedTokens := []TokenType{NUMBER, CARET, NUMBER, CARET, NUMBER, EOF}

	for i, expected := range expectedTokens {
		tok := lexer.NextToken()
		if tok.Type != expected {
			t.Errorf("token[%d] - expected type %v, got %v", i, expected, tok.Type)
		}
	}
}

func TestLexer_Equal(t *testing.T) {
	input := "="
	lexer := NewLexer(input)

	tok := lexer.NextToken()
	if tok.Type != EQUAL {
		t.Errorf("expected EQUAL token, got %v", tok.Type)
	}
	if tok.Literal != "=" {
		t.Errorf("expected literal '=', got '%s'", tok.Literal)
	}
}

func TestLexer_EqualInExpression(t *testing.T) {
	input := "2 + 3 = 5"
	lexer := NewLexer(input)

	expectedTokens := []struct {
		tokenType TokenType
		literal   string
		value     float64
	}{
		{NUMBER, "2", 2.0},
		{PLUS, "+", 0},
		{NUMBER, "3", 3.0},
		{EQUAL, "=", 0},
		{NUMBER, "5", 5.0},
		{EOF, "", 0},
	}

	for i, expected := range expectedTokens {
		tok := lexer.NextToken()

		if tok.Type != expected.tokenType {
			t.Errorf("token[%d] - expected type %v, got %v", i, expected.tokenType, tok.Type)
		}
		if tok.Literal != expected.literal {
			t.Errorf("token[%d] - expected literal %s, got %s", i, expected.literal, tok.Literal)
		}
		if expected.tokenType == NUMBER && tok.Value != expected.value {
			t.Errorf("token[%d] - expected value %f, got %f", i, expected.value, tok.Value)
		}
	}
}

func TestLexer_MultipleEquals(t *testing.T) {
	input := "a = b = c"
	lexer := NewLexer(input)

	expectedTokens := []TokenType{
		VARIABLE, EQUAL, VARIABLE, EQUAL, VARIABLE, EOF,
	}

	for i, expected := range expectedTokens {
		tok := lexer.NextToken()
		if tok.Type != expected {
			t.Errorf("token[%d] - expected type %v, got %v", i, expected, tok.Type)
		}
	}
}

func TestLexer_EqualWithWhitespace(t *testing.T) {
	input := "  2  =  3  "
	lexer := NewLexer(input)

	expectedTokens := []TokenType{NUMBER, EQUAL, NUMBER, EOF}

	for i, expected := range expectedTokens {
		tok := lexer.NextToken()
		if tok.Type != expected {
			t.Errorf("token[%d] - expected type %v, got %v", i, expected, tok.Type)
		}
	}
}

func TestLexer_ComplexEqualExpression(t *testing.T) {
	input := "(2 + 3) = (1 + 4)"
	lexer := NewLexer(input)

	expectedTokens := []TokenType{
		LPAREN, NUMBER, PLUS, NUMBER, RPAREN,
		EQUAL,
		LPAREN, NUMBER, PLUS, NUMBER, RPAREN,
		EOF,
	}

	for i, expected := range expectedTokens {
		tok := lexer.NextToken()
		if tok.Type != expected {
			t.Errorf("token[%d] - expected type %v, got %v", i, expected, tok.Type)
		}
	}
}
