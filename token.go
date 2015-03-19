package goose

import "fmt"

type Token struct {
	Start int
	Stop  int
	Type  TokenType
	Raw   []rune
	Value interface{}
}

func (t *Token) String() string {
	switch t.Type {
	case TokenErr:
		if asserted, ok := t.Value.(error); ok {
			return asserted.Error()
		}

		return "error value is not an error"
	default:
		return fmt.Sprintf("[%d,%d) %d %q", t.Start, t.Stop, t.Type, t.Value)
	}
}

type TokenType uint

var keywords = map[string]TokenType{
	"add":         TokenAdd,
	"asc":         TokenAsc,
	"binary":      TokenBinary,
	"boolean":     TokenBoolean,
	"change":      TokenChange,
	"column":      TokenColumn,
	"create":      TokenCreate,
	"date":        TokenDate,
	"datetime":    TokenDatetime,
	"decimal":     TokenDecimal,
	"default":     TokenDefault,
	"desc":        TokenDesc,
	"down":        TokenDown,
	"end":         TokenEnd,
	"false":       TokenFalse,
	"float":       TokenFloat,
	"index":       TokenIndex,
	"integer":     TokenInteger,
	"name":        TokenName,
	"null":        TokenNull,
	"order":       TokenOrder,
	"precision":   TokenPrecision,
	"primary_key": TokenPrimaryKey,
	"raw":         TokenRaw,
	"references":  TokenReferences,
	"remove":      TokenRemove,
	"rename":      TokenRename,
	"scale":       TokenScale,
	"size":        TokenSize,
	"string":      TokenString,
	"table":       TokenTable,
	"text":        TokenText,
	"time":        TokenTime,
	"timestamp":   TokenTimestamp,
	"timestamps":  TokenTimestamps,
	"true":        TokenTrue,
	"unique":      TokenUnique,
	"up":          TokenUp,
}

const (
	TokenErr TokenType = iota
	TokenEOF

	TokenAdd
	TokenAsc
	TokenBinary
	TokenBoolean
	TokenChange
	TokenColon
	TokenColumn
	TokenComma
	TokenCreate
	TokenDate
	TokenDatetime
	TokenDecimal
	TokenDefault
	TokenDesc
	TokenDown
	TokenEnd
	TokenFalse
	TokenFloat
	TokenFloatLiteral
	TokenIndex
	TokenInteger
	TokenIntegerLiteral
	TokenLeftBrace
	TokenName
	TokenNull
	TokenOrder
	TokenPrecision
	TokenPrimaryKey
	TokenRaw
	TokenReferences
	TokenRemove
	TokenRename
	TokenRightBrace
	TokenScale
	TokenSize
	TokenString
	TokenStringLiteral
	TokenTable
	TokenText
	TokenTime
	TokenTimestamp
	TokenTimestamps
	TokenTrue
	TokenUnique
	TokenUp
)
