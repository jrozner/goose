package goose

type Node struct {
	Type     NodeType
	Children []interface{}
}

type NodeType uint

var Nodes = map[NodeType]string{
	NodeAddColumn:         "AddColumn",
	NodeAddIndex:          "AddIndex",
	NodeBoolean:           "Boolean",
	NodeChangeColumn:      "ChangeColumn",
	NodeColumnDefinition:  "ColumnDefinition",
	NodeColumnName:        "ColumnName",
	NodeCreateTable:       "CreateTable",
	NodeDataType:          "DataType",
	NodeDefaultValue:      "DefaultValue",
	NodeDownStatement:     "DownStatement",
	NodeFloat:             "Float",
	NodeIndexOptions:      "IndexOptions",
	NodeIndexOptionsBlock: "IndexOptionsBlock",
	NodeInteger:           "Integer",
	NodeNewName:           "NewName",
	NodeNull:              "Null",
	NodeOption:            "Option",
	NodeOptionsBlock:      "OptionsBlock",
	NodeRanameColumn:      "RenameColumn",
	NodeRaw:               "Raw",
	NodeRawBody:           "RawBody",
	NodeRemoveColumn:      "RemoveColumn",
	NodeRemoveIndex:       "RemoveIndex",
	NodeRemoveTable:       "RemoveTable",
	NodeRenameTable:       "RenameTable",
	NodeRoot:              "Root",
	NodeStatement:         "Statement",
	NodeString:            "String",
	NodeTableName:         "TableName",
	NodeUpStatement:       "UpStatement",
}

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
	NodeFloat
	NodeIndexOptions
	NodeIndexOptionsBlock
	NodeInteger
	NodeNewName
	NodeNull
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
	NodeString
	NodeTableName
	NodeUpStatement
)
