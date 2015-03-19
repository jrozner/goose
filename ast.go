package goose

type Node struct {
	Type     NodeType
	Children []interface{}
}

type NodeType uint

const (
	NodeAddColumn NodeType = iota
	NodeAddIndex
	NodeBoolean
	NodeChangeColumn
	NodeColumnDefinition
	NodeColumnName
	NodeCreateTable
	NodeDataType
	NodeDefaultValue
	NodeDownStatement
	NodeIndexOptions
	NodeIndexOptionsBlock
	NodeNewName
	NodeOption
	NodeOptionsBlock
	NodeRanameColumn
	NodeRaw
	NodeRawBody
	NodeRemoveColumn
	NodeRemoveIndex
	NodeRemoveTable
	NodeRenameTable
	NodeRoot
	NodeStatement
	NodeTableName
	NodeUpStatement
)
