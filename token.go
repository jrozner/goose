package goose

import "fmt"

type Token struct {
	Start int
	Stop  int
	Raw   []rune
	Value interface{}
	Type  TokenType
}

func (t *Token) String() string {
	switch t.Type {
	case Err:
		if asserted, ok := t.Value.(error); ok {
			return asserted.Error()
		}

		return "error value is not an error"
	case EOF:
		return "EOF"
	default:
		return fmt.Sprintf("%q", t.Value)
	}
}

type TokenType int

var keywords = map[string]TokenType{
	"add":         Add,
	"asc":         Asc,
	"binary":      Binary,
	"boolean":     Boolean,
	"change":      Change,
	"column":      Column,
	"create":      Create,
	"date":        Date,
	"datetime":    Datetime,
	"decimal":     Decimal,
	"default":     Default,
	"desc":        Desc,
	"down":        Down,
	"end":         End,
	"false":       False,
	"float":       Float,
	"index":       Index,
	"integer":     Integer,
	"name":        Name,
	"null":        Null,
	"order":       Order,
	"precision":   Precision,
	"primary_key": PrimaryKey,
	"raw":         Raw,
	"references":  References,
	"remove":      Remove,
	"rename":      Rename,
	"scale":       Scale,
	"size":        Size,
	"string":      String,
	"table":       Table,
	"text":        Text,
	"time":        Time,
	"timestamp":   Timestamp,
	"timestamps":  Timestamps,
	"true":        True,
	"unique":      Unique,
	"up":          Up,
}

const (
	Err TokenType = iota
	EOF

	Add
	Asc
	Binary
	Boolean
	Change
	Colon
	Column
	Comma
	Create
	Date
	Datetime
	Decimal
	Default
	Desc
	Down
	End
	False
	Float
	FloatLiteral
	Index
	Integer
	IntegerLiteral
	LeftBrace
	Name
	Null
	Order
	Precision
	PrimaryKey
	Raw
	References
	Remove
	Rename
	RightBrace
	Scale
	Size
	String
	StringLiteral
	Table
	Text
	Time
	Timestamp
	Timestamps
	True
	Unique
	Up
)
