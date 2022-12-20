package lexer

type TokenType int

type Token interface {
	Type() TokenType
	Value() string
}

type Lexer interface {
	Next() Token
}
