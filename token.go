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
	case TokenComma:
		fallthrough
	case TokenColon:
		fallthrough
	case TokenLeftBrace:
		fallthrough
	case TokenRightBrace:
		if asserted, ok := t.Value.(rune); ok {
			return fmt.Sprintf("%s", string(asserted))
		}

		return "error asserting value"
	case TokenStringLiteral:
		if asserted, ok := t.Value.(string); ok {
			runes := []rune(asserted)
			return fmt.Sprintf("%s", string(runes[1:len(asserted)-1]))
		}

		return "error asserting value"
	case TokenFloatLiteral:
		if asserted, ok := t.Value.(float64); ok {
			return fmt.Sprintf("%.02f", asserted)
		}

		return "error asserting value"
	case TokenIntegerLiteral:
		if asserted, ok := t.Value.(int64); ok {
			return fmt.Sprintf("%d", asserted)
		}

		return "error asserting value"
	default:
		return fmt.Sprintf("%q", t.Value)
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
