package token

const (
	ILIEGAL   = "ILIEGAL"
	EOF       = "EOF"
	IDENT     = "IDENT"
	INT       = "INT"
	COMMA     = ","
	SEMICOLON = ";"
	LPAREN    = "("
	RPAREN    = ")"
	LBRACE    = "{"
	RBRACE    = "}"

	// keyword
	FUNCTION = "FUNCTION"
	LET      = "LET"
	TRUE     = "TRUE"
	FALSE    = "FALSE"
	IF       = "IF"
	ELSE     = "ELSE"
	RETURN   = "RETURN"
	STRING   = "STRING"
	IMPORT   = "IMPORT"
	WHILE    = "WHILE"
	BREAK    = "BREAK"
	PACKAGE  = "PACKAGE"
	FOR      = "FOR"
	DEFER    = "DEFER"

	// oop keyword
	CLASS     = "CLASS"
	NEW       = "NEW"
	EXTENDS   = "EXTENDS"
	IMPLEMENT = "IMPLEMENT"
	INTERFACE = "INTERFACE"
	// oop object install value access
	OBJET_GET = "->"
	CLASS_GET = "::"

	ASSIGN   = "="
	PLUS     = "+"
	MINUS    = "-"
	BANG     = "!"
	ASTERISK = "*"
	SLASH    = "/"

	PLUSASSIGN     = "+="
	MINUSASSIGN    = "-="
	ASTERISKASSIGN = "*="
	SLASHASSIGN    = "/="
	PLUSPLUS       = "++"
	MINUSMINUS     = "--"

	LT = "<"
	LE = "<="
	GT = ">"
	GE = ">="

	EQ     = "=="
	NOT_EQ = "!="

	LBRACKET = "["
	RBRACKET = "]"

	COLON = ":"

	AND = "&&"
	OR  = "||"

	FLOAT = "FLOAT"
)

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

var keywords = map[string]TokenType{
	"fn":        FUNCTION,
	"let":       LET,
	"true":      TRUE,
	"false":     FALSE,
	"if":        IF,
	"else":      ELSE,
	"return":    RETURN,
	"import":    IMPORT,
	"while":     WHILE,
	"package":   PACKAGE,
	"break":     BREAK,
	"for":       FOR,
	"class":     CLASS,
	"new":       NEW,
	"extends":   EXTENDS,
	"implement": IMPLEMENT,
	"interface": INTERFACE,
	"defer":     DEFER,
}

func LookIndent(indent string) TokenType {
	token, ok := keywords[indent]
	if ok {
		return token
	}
	return IDENT
}
