package ast

import (
	"testing"
	"z/token"
)

func TestString(t *testing.T) {
	program := &Program{
		Statements: []Statement{
			&LetStatement{
				Token: token.Token{Type: token.LET, Literal: "let"},
				Name: &Identifier{
					Token: token.Token{Type: token.IDENT, Literal: "myVar"},
					Value: "myVar",
				},
				Value: &Identifier{
					Token: token.Token{Type: token.IDENT, Literal: "anotherVar"},
					Value: "anotherVar",
				},
			},
		},
	}

	if program.String() != "let myVar = anotherVar;" {
		t.Errorf("program.String() wrong. got=%q", program.String())
	}
}

func TestForStatement(t *testing.T) {
	program := &Program{
		Statements: []Statement{
			&ExpressionStatement{
				Expression: &ForExpression{
					Token: token.Token{Type: token.FOR, Literal: "for"},
					Initor: &LetStatement{
						Token: token.Token{Type: token.LET, Literal: "let"},
						Name: &Identifier{
							Token: token.Token{Type: token.IDENT, Literal: "i"},
							Value: "i",
						},
						Value: &IntegerLiteral{Token: token.Token{Literal: "0"}, Value: 0},
					},
					Condition: &InfixExpression{
						Left: &Identifier{
							Token: token.Token{Type: token.IDENT, Literal: "i"},
							Value: "i",
						},
						Operator: "<",
						Right: &IntegerLiteral{
							Token: token.Token{Literal: "10"},
							Value: 10,
						},
					},
					After: &InfixExpression{
						Left: &Identifier{
							Token: token.Token{Type: token.IDENT, Literal: "i"},
							Value: "i",
						},
						Operator: "++",
						Right: &IntegerLiteral{
							Token: token.Token{Literal: "1"},
							Value: 1,
						},
					},
				},
			},
		},
	}
	if program.String() != "for(let i = 0;;(i < 10);(i ++ 1)){}" {
		t.Errorf("program.String() wrong. got=%q", program.String())
	}
}
