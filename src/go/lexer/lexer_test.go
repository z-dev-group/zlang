// lexer/lexer_test.go

package lexer

import (
	"testing"

	"z/token"
)

func TestNextToken(t *testing.T) {
	input := `let five = 5
let ten = 10

let add = fn(x, y) {
  x + y
}

let result = add(five, ten)
!-*/5
5 < 10 > 5

if (5  < 10) {
  return true
} else {
  return false
};

10 == 10;
10 != 9
"foobar"
"foo bar"
[1, 2]
{"foo": "bar"}
/* annotation */
// this is annotation
import "package/include_file.z"
>=
<=
while (2 > 1) {x}
a = b + 1
return []
let money = 100.98
package string
name.age
a += 1
b -= 2
c *= 3
d /= 4
i++
i--
break
for
class
new
extends
implement
interface
->
_name
__age
defer {}
&& 
||
?
:
`
	input += "`hello`"

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.LET, "let"},
		{token.IDENT, "five"},
		{token.ASSIGN, "="},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"},
		{token.IDENT, "ten"},
		{token.ASSIGN, "="},
		{token.INT, "10"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"},
		{token.IDENT, "add"},
		{token.ASSIGN, "="},
		{token.FUNCTION, "fn"},
		{token.LPAREN, "("},
		{token.IDENT, "x"},
		{token.COMMA, ","},
		{token.IDENT, "y"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.IDENT, "x"},
		{token.PLUS, "+"},
		{token.IDENT, "y"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"},
		{token.IDENT, "result"},
		{token.ASSIGN, "="},
		{token.IDENT, "add"},
		{token.LPAREN, "("},
		{token.IDENT, "five"},
		{token.COMMA, ","},
		{token.IDENT, "ten"},
		{token.RPAREN, ")"},
		{token.SEMICOLON, ";"},
		{token.BANG, "!"},
		{token.MINUS, "-"},
		{token.ASTERISK, "*"},
		{token.SLASH, "/"},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},
		{token.INT, "5"},
		{token.LT, "<"},
		{token.INT, "10"},
		{token.GT, ">"},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},
		{token.IF, "if"},
		{token.LPAREN, "("},
		{token.INT, "5"},
		{token.LT, "<"},
		{token.INT, "10"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.RETURN, "return"},
		{token.TRUE, "true"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.ELSE, "else"},
		{token.LBRACE, "{"},
		{token.RETURN, "return"},
		{token.FALSE, "false"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.SEMICOLON, ";"},
		{token.INT, "10"},
		{token.EQ, "=="},
		{token.INT, "10"},
		{token.SEMICOLON, ";"},
		{token.INT, "10"},
		{token.NOT_EQ, "!="},
		{token.INT, "9"},
		{token.SEMICOLON, ";"},
		{token.STRING, "foobar"},
		{token.SEMICOLON, ";"},
		{token.STRING, "foo bar"},
		{token.SEMICOLON, ";"},
		{token.LBRACKET, "["},
		{token.INT, "1"},
		{token.COMMA, ","},
		{token.INT, "2"},
		{token.RBRACKET, "]"},
		{token.SEMICOLON, ";"},
		{token.LBRACE, "{"},
		{token.STRING, "foo"},
		{token.COLON, ":"},
		{token.STRING, "bar"},
		{token.RBRACE, "}"},
		{token.SEMICOLON, ";"},
		{token.IMPORT, "import"},
		{token.STRING, "package/include_file.z"},
		{token.SEMICOLON, ";"},
		{token.GE, ">="},
		{token.SEMICOLON, ";"},
		{token.LE, "<="},
		{token.SEMICOLON, ";"},
		{token.WHILE, "while"},
		{token.LPAREN, "("},
		{token.INT, "2"},
		{token.GT, ">"},
		{token.INT, "1"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.IDENT, "x"},
		{token.RBRACE, "}"},
		{token.SEMICOLON, ";"},
		{token.IDENT, "a"},
		{token.ASSIGN, "="},
		{token.IDENT, "b"},
		{token.PLUS, "+"},
		{token.INT, "1"},
		{token.SEMICOLON, ";"},
		{token.RETURN, "return"},
		{token.LBRACKET, "["},
		{token.RBRACKET, "]"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"},
		{token.IDENT, "money"},
		{token.ASSIGN, "="},
		{token.FLOAT, "100.98"},
		{token.SEMICOLON, ";"},
		{token.PACKAGE, "package"},
		{token.IDENT, "string"},
		{token.SEMICOLON, ";"},
		{token.IDENT, "name.age"},
		{token.SEMICOLON, ";"},
		{token.IDENT, "a"},
		{token.PLUSASSIGN, "+="},
		{token.INT, "1"},
		{token.SEMICOLON, ";"},
		{token.IDENT, "b"},
		{token.MINUSASSIGN, "-="},
		{token.INT, "2"},
		{token.SEMICOLON, ";"},
		{token.IDENT, "c"},
		{token.ASTERISKASSIGN, "*="},
		{token.INT, "3"},
		{token.SEMICOLON, ";"},
		{token.IDENT, "d"},
		{token.SLASHASSIGN, "/="},
		{token.INT, "4"},
		{token.SEMICOLON, ";"},
		{token.IDENT, "i"},
		{token.PLUSPLUS, "++"},
		{token.SEMICOLON, ";"},
		{token.IDENT, "i"},
		{token.MINUSMINUS, "--"},
		{token.SEMICOLON, ";"},
		{token.BREAK, "break"},
		{token.SEMICOLON, ";"},
		{token.FOR, "for"},
		{token.SEMICOLON, ";"},
		{token.CLASS, "class"},
		{token.SEMICOLON, ";"},
		{token.NEW, "new"},
		{token.SEMICOLON, ";"},
		{token.EXTENDS, "extends"},
		{token.SEMICOLON, ";"},
		{token.IMPLEMENT, "implement"},
		{token.SEMICOLON, ";"},
		{token.INTERFACE, "interface"},
		{token.SEMICOLON, ";"},
		{token.OBJET_GET, "->"},
		{token.SEMICOLON, ";"},
		{token.IDENT, "_name"},
		{token.SEMICOLON, ";"},
		{token.IDENT, "__age"},
		{token.SEMICOLON, ";"},
		{token.DEFER, "defer"},
		{token.LBRACE, "{"},
		{token.RBRACE, "}"},
		{token.SEMICOLON, ";"},
		{token.AND, "&&"},
		{token.SEMICOLON, ";"},
		{token.OR, "||"},
		{token.SEMICOLON, ";"},
		{token.QUESTION, "?"},
		{token.SEMICOLON, ";"},
		{token.COLON, ":"},
		{token.SEMICOLON, ";"},
		// special string
		{token.STRING, "hello"},
		{token.SEMICOLON, ";"},
		{token.EOF, ""},
	}
	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}
