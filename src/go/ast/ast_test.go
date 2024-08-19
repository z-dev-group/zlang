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
		t.Errorf("program.String() wrong. got=%s", program.String())
	}
}

func TestInterfaceExpression(t *testing.T) {
	program := Program{
		Statements: []Statement{
			&ExpressionStatement{
				Expression: &InterfaceExpress{
					Token: token.Token{Literal: "interface"},
					Name: Identifier{
						Token: token.Token{Literal: "intf"},
						Value: "intf",
					},
					Functions: []FunctionLiteral{
						{
							Token:      token.Token{Literal: "fn"},
							Name:       "hello",
							Body:       &BlockStatement{},
							Parameters: []*Identifier{},
						},
					},
				},
			},
		},
	}
	if program.String() != "interface intf{fn hello() {} }" {
		t.Errorf("program.String() wrong. got=%s", program.String())
	}
}

func TestClassExpression(t *testing.T) {
	program := Program{
		Statements: []Statement{
			&ExpressionStatement{
				Expression: &ClassExpress{
					Token: token.Token{Literal: "class", Type: token.CLASS},
					Name: Identifier{
						Token: token.Token{Literal: "Person", Type: token.IDENT},
						Value: "Person",
					},
					Parents: []ClassExpress{
						{
							Token: token.Token{Literal: "class", Type: token.CLASS},
							Name: Identifier{
								Token: token.Token{Literal: "TwoLegs", Type: token.IDENT},
								Value: "TwoLegs",
							},
						},
						{
							Token: token.Token{Literal: "class", Type: token.CLASS},
							Name: Identifier{
								Token: token.Token{Literal: "Animal", Type: token.IDENT},
								Value: "Animal",
							},
						},
					},
					LetStatements: []LetStatement{
						{
							Token: token.Token{Type: token.LET, Literal: "let"},
							Name: &Identifier{
								Token: token.Token{Literal: "name", Type: token.IDENT},
								Value: "name",
							},
							Value: &StringLiteral{
								Token: token.Token{Literal: "seven"},
								Value: "seven",
							},
						},
						{
							Token: token.Token{Literal: "let", Type: token.LET},
							Name: &Identifier{
								Token: token.Token{Literal: "age", Type: token.IDENT},
								Value: "age",
							},
							Value: &IntegerLiteral{
								Token: token.Token{Literal: "12"},
								Value: 12,
							},
						},
					},
					Functions: []FunctionLiteral{
						{
							Token: token.Token{Literal: "fn"},
							Name:  "work",
							Body:  &BlockStatement{},
							Parameters: []*Identifier{
								{
									Token: token.Token{Literal: "hours"},
									Value: "hours",
								},
							},
						},
					},
				},
			},
		},
	}
	if program.String() != "class Person extends TwoLegs,Animal {let name = seven;let age = 12;fn work(hours) {} }" {
		t.Fatalf("program.String() wrong. got=%s", program.String())
	}
}

func TestObjectExpress(t *testing.T) {
	program := Program{
		Statements: []Statement{
			&ExpressionStatement{
				Expression: &ObjectExpress{
					Token: token.Token{Literal: "new", Type: token.NEW},
					Class: ClassExpress{
						Name: Identifier{
							Token: token.Token{Literal: "Hello"},
						},
					},
				},
			},
		},
	}
	if program.String() != "new Hello()" {
		t.Fatalf("program.String() wrong. got=%s", program.String())
	}
}
