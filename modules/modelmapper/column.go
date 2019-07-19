package modelmapper

import (
	"github.com/herb-go/herb/model/sql/db"
	"github.com/herb-go/herb/model/sql/db/columns"
	_ "github.com/herb-go/herb/model/sql/db/columns/mysqlcolumns"  //mysql driver
	_ "github.com/herb-go/herb/model/sql/db/columns/sqlitecolumns" //sqlite driver
	"github.com/herb-go/util/cli/name"
)

type Column struct {
	*columns.Column
}

func (c *Column) Name() string {
	return MustGetColumnName(c.Column)
}

type ModelColumns struct {
	Columns     []*Column
	Name        *name.Name
	Database    string
	PrimaryKeys []*Column
	HasTime     bool
}

func MustGetColumnName(c *columns.Column) string {
	n, err := name.New(false, c.Field)
	if err != nil {
		panic(err)
	}
	return n.Pascal
}
func (m *ModelColumns) FirstPrimayKey() *Column {
	return m.Columns[0]
}

func (m *ModelColumns) HasPrimayKey() bool {
	return len(m.PrimaryKeys) > 0
}

func (m *ModelColumns) IsSinglePrimayKey() bool {
	return len(m.PrimaryKeys) == 1
}

func (m *ModelColumns) IsMultiPrimayKey() bool {
	return len(m.PrimaryKeys) > 1
}

func (m *ModelColumns) PrimayKeyField() string {
	if m.IsSinglePrimayKey() {
		return "interface{}"
	}
	return "*" + m.Name.Pascal + "PrimaryKey"
}

// func (m *ModelColumns) PrimaryKeyType() string {
// 	output := "//" + m.Name.Pascal + "PrimaryKey : table " + m.Name.Raw + " primary key type\n"
// 	switch len(m.PrimaryKeys) {
// 	case 0:
// 		output = output + "type " + m.Name.Pascal + "PrimaryKey map[string]interface{}\n"
// 	case 1:
// 		output = output + "type " + m.Name.Pascal + "PrimaryKey "
// 		if !m.PrimaryKeys[0].NotNull {
// 			output = output + "*"
// 		}
// 		output = output + m.Columns[0].ColumnType + "\n"
// 	default:
// 		output = output + "type " + m.Name.Pascal + "PrimaryKey struct{\n"
// 		for _, v := range m.PrimaryKeys {
// 			output = output + "    " + v.Name() + " "
// 			if !v.NotNull {
// 				output = output + "*"
// 			}
// 			output = output + v.ColumnType + "\n"
// 		}
// 		output = output + "}\n"
// 	}
// 	return output
// }

// func (m *ModelColumns) BuildByPKQuery() string {
// 	output := "//BuildByPKQuery : build by pk query for table " + m.Name.Raw + "\n"
// 	output = output + "func (mapper *" + m.Name.Pascal + "Mapper) BuildByPKQuery(pk interface{}) *querybuilder.PlainQuery {\n"
// 	output = output + "var query= mapper.QueryBuilder()\n"
// 	output = output + "    var q = query.New(\"\")\n"
// 	switch len(m.PrimaryKeys) {
// 	case 0:
// 		output = output + "    for k,v :=range pk {\n"
// 		output = output + " q.And(query.Equal( " + m.Name.Pascal + ".FieldAlias(k),v))\n"
// 		output = output + "    }"
// 	case 1:
// 		output = output + "    q.And(query.Equal(" + m.Name.Pascal + "FieldAlias" + m.Columns[0].Name() + ",pk))\n"
// 	default:
// 		for _, v := range m.PrimaryKeys {
// 			output = output + "    q.And(query.Equal(" + m.Name.Pascal + "FieldAlias" + v.Name() + ",pk." + v.Name() + "))\n"
// 		}
// 	}
// 	output = output + "    return q\n}\n"
// 	return output
// }
// func (m *ModelColumns) ModelPrimaryKey() string {
// 	output := "//ModelPrimaryKey :  get primary key from model.\n"
// 	output = output + "func (mapper *" + m.Name.Pascal + "Mapper) ModelPrimaryKey( model *" + m.Name.Pascal + "Model ) *" + m.Name.Pascal + "PrimaryKey {\n"
// 	switch len(m.PrimaryKeys) {
// 	case 0:
// 		output = output + "    return nil\n"
// 	case 1:
// 		output = output + "    var pk " + m.Name.Pascal + "PrimaryKey\n"
// 		output = output + "    pk=" + m.Name.Pascal + "PrimaryKey(model." + m.Columns[0].Name() + ")\n"
// 		output = output + "    return &pk\n"
// 	default:
// 		output = output + "    pk:=" + m.Name.Pascal + "PrimaryKey{}\n"
// 		for _, v := range m.PrimaryKeys {
// 			output = output + "    pk." + v.Name() + " = model." + v.Name() + "\n"
// 		}
// 		output = output + "    return &pk\n"
// 	}
// 	output = output + "}\n"
// 	return output
// }

// func (m *ModelColumns) ColumnsToModelStruct() string {
// 	output := "//" + m.Name.Pascal + "Model :" + m.Name.Raw + " model.\n"
// 	output = output + "type " + m.Name.Pascal + "Model struct{\n"
// 	for _, v := range m.Columns {
// 		output = output + "    " + v.Name() + " "
// 		if !v.NotNull {
// 			output = output + "*"
// 		}
// 		output = output + v.ColumnType + "\n"
// 	}
// 	output = output + "}\n"
// 	return output
// }

func getLoaderFormDB(conn db.Database) (columns.Loader, error) {
	drivername := conn.Driver()
	driver, err := columns.Driver(drivername)
	if err != nil {
		return nil, err
	}
	return driver, nil
}

func NewModelCulumns(conn db.Database, database string, table string) (*ModelColumns, error) {
	loader, err := getLoaderFormDB(conn)
	if err != nil {
		return nil, err
	}
	err = loader.Load(conn, table)
	if err != nil {
		return nil, err
	}
	columns, err := loader.Columns()
	if err != nil {
		return nil, err
	}
	c := make([]*Column, len(columns))
	for k := range columns {
		c[k] = &Column{columns[k]}
	}
	pks := []*Column{}
	var hasTime bool
	for _, v := range c {
		if v.PrimayKey {
			pks = append(pks, v)
		}
		if v.ColumnType == "time.Time" {
			hasTime = true
		}
	}
	tablename, err := name.New(false, table)
	if err != nil {
		return nil, err
	}
	mc := &ModelColumns{
		Columns:     c,
		Name:        tablename,
		Database:    database,
		PrimaryKeys: pks,
		HasTime:     hasTime,
	}
	return mc, nil
}
