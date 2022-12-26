package lexer

type TokenType int
type Tokens []Token

type Token interface {
	Type() TokenType
	Value() string
}

type Lexer interface {
	Next() Token
	Tokens() Tokens
}
