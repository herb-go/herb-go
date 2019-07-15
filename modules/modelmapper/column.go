package modelmapper

import (
	"github.com/herb-go/herb/model/sql/db"
	"github.com/herb-go/herb/model/sql/db/columns"
	_ "github.com/herb-go/herb/model/sql/db/columns/mysqlcolumns"  //mysql driver
	_ "github.com/herb-go/herb/model/sql/db/columns/sqlitecolumns" //sqlite driver
	"github.com/herb-go/util/cli/name"
)

type ModelColumns struct {
	Columns     []*columns.Column
	Name        *name.Name
	Database    string
	PrimaryKeys []*columns.Column
	HasTime     bool
}

func MustGetColumnName(c *columns.Column) string {
	n, err := name.New(false, c.Field)
	if err != nil {
		panic(err)
	}
	return n.Pascal
}
func (m *ModelColumns) FirstPrimayKey() *columns.Column {
	return m.Columns[0]
}

func (m *ModelColumns) CanCreate() bool {
	if len(m.PrimaryKeys) == 1 {
		if m.PrimaryKeys[0].ColumnType == "string" {
			return true
		}
		if m.PrimaryKeys[0].AutoValue {
			return true
		}
	}
	return false
}

func (m *ModelColumns) HasPrimayKey() bool {
	return len(m.Columns) > 0
}

func (m *ModelColumns) PrimaryKeyType() string {
	output := "//" + m.Name.Pascal + "PrimaryKey : table " + m.Name.Raw + " primary key type\n"
	switch len(m.PrimaryKeys) {
	case 0:
		output = output + "type " + m.Name.Pascal + "PrimaryKey map[string]interface{}\n"
	case 1:
		output = output + "type " + m.Name.Pascal + "PrimaryKey "
		if !m.PrimaryKeys[0].NotNull {
			output = output + "*"
		}
		output = output + m.Columns[0].ColumnType + "\n"
	default:
		output = output + "type " + m.Name.Pascal + "PrimaryKey struct{\n"
		for _, v := range m.PrimaryKeys {
			output = output + "    " + MustGetColumnName(v) + " "
			if !v.NotNull {
				output = output + "*"
			}
			output = output + v.ColumnType + "\n"
		}
		output = output + "}\n"
	}
	return output
}

func (m *ModelColumns) BuildByPKQuery() string {
	output := "//BuildByPKQuery : build by pk query for table " + m.Name.Raw + "\n"
	output = output + "func (mapper *" + m.Name.Pascal + "Mapper) BuildByPKQuery(pk interface{}) *querybuilder.PlainQuery {\n"
	output = output + "var query= mapper.QueryBuilder()\n"
	output = output + "    var q = query.New(\"\")\n"
	switch len(m.PrimaryKeys) {
	case 0:
		output = output + "    for k,v :=range pk {\n"
		output = output + " q.And(query.Equal( " + m.Name.Pascal + ".FieldAlias(k),v))\n"
		output = output + "    }"
	case 1:
		output = output + "    q.And(query.Equal(" + m.Name.Pascal + "FieldAlias" + MustGetColumnName(m.Columns[0]) + ",pk))\n"
	default:
		for _, v := range m.PrimaryKeys {
			output = output + "    q.And(query.Equal(" + m.Name.Pascal + "FieldAlias" + MustGetColumnName(v) + ",pk." + MustGetColumnName(v) + "))\n"
		}
	}
	output = output + "    return q\n}\n"
	return output
}
func (m *ModelColumns) ModelPrimaryKey() string {
	output := "//ModelPrimaryKey :  get primary key from model.\n"
	output = output + "func (mapper *" + m.Name.Pascal + "Mapper) ModelPrimaryKey( model *" + m.Name.Pascal + "Model ) *" + m.Name.Pascal + "PrimaryKey {\n"
	switch len(m.PrimaryKeys) {
	case 0:
		output = output + "    return nil\n"
	case 1:
		output = output + "    var pk " + m.Name.Pascal + "PrimaryKey\n"
		output = output + "    pk=" + m.Name.Pascal + "PrimaryKey(model." + MustGetColumnName(m.Columns[0]) + ")\n"
		output = output + "    return &pk\n"
	default:
		output = output + "    pk:=" + m.Name.Pascal + "PrimaryKey{}\n"
		for _, v := range m.PrimaryKeys {
			output = output + "    pk." + MustGetColumnName(v) + " = model." + MustGetColumnName(v) + "\n"
		}
		output = output + "    return &pk\n"
	}
	output = output + "}\n"
	return output
}
func (m *ModelColumns) ColumnsToModelStruct() string {
	output := "//" + m.Name.Pascal + "Model :" + m.Name.Raw + " model.\n"
	output = output + "type " + m.Name.Pascal + "Model struct{\n"
	for _, v := range m.Columns {
		output = output + "    " + MustGetColumnName(v) + " "
		if !v.NotNull {
			output = output + "*"
		}
		output = output + v.ColumnType + "\n"
	}
	output = output + "}\n"
	return output
}

func (m *ModelColumns) ColumnsToFieldsMethod(query string) string {
	qn, err := name.New(false, query)
	if err != nil {
		panic(err)
	}
	output := "//" + qn.Pascal + "Fields : map model " + qn.Raw + " fields to database column.\n"
	output = output + "func (mapper *" + m.Name.Pascal + "Mapper) " + qn.Pascal + "Fields(model *" + m.Name.Pascal + "Model)* querybuilder.Fields {\n"
	output = output + "	   if model == nil {\n"
	output = output + "        model = New" + m.Name.Pascal + "Model()\n"
	output = output + "	   }\n"
	output = output + "    return model.BuildFields(true,\n"
	for _, v := range m.Columns {

		output = output + "	//Field \"" + m.Name.Raw + "." + v.Field + "\"\n	" + m.Name.Pascal + "Field" + MustGetColumnName(v) + ","
		output = output + "\n"
	}
	output = output + "    )\n"
	output = output + "}\n"
	return output
}

func (m *ModelColumns) ColumnsToFieldsInsertMethod(query string) string {
	qn, err := name.New(false, query)
	if err != nil {
		panic(err)
	}
	output := "//" + qn.Pascal + "FieldsInsert : map model " + qn.Raw + " fields to database column used in insert query.\n"
	output = output + "func (mapper *" + m.Name.Pascal + "Mapper) " + qn.Pascal + "FieldsInsert(model *" + m.Name.Pascal + "Model)* querybuilder.Fields {\n"
	output = output + "    return model.BuildFields(false,"
	skiped := ""
	fields := ""
	for _, v := range m.Columns {
		if v.AutoValue {
			skiped = skiped + "\n		//Skip field \"" + v.Field + "\" which should be set by database"
			skiped = skiped + "\n		 //Field \"" + m.Name.Raw + "." + v.Field + "\"\n		//" + m.Name.Pascal + "Field" + MustGetColumnName(v) + ","
		} else {
			fields = fields + "\n		 //Field \"" + m.Name.Raw + "." + v.Field + "\"\n		" + m.Name.Pascal + "Field" + MustGetColumnName(v) + ","
		}
	}

	output = output + skiped + fields + "\n    )\n}\n"
	return output
}
func (m *ModelColumns) ColumnsToFieldsUpdateMethod(query string) string {
	qn, err := name.New(false, query)
	if err != nil {
		panic(err)
	}
	output := "//" + qn.Pascal + "FieldsUpdate : map model " + qn.Raw + " fields to database column used in update query.\n"
	output = output + "func (mapper *" + m.Name.Pascal + "Mapper) " + qn.Pascal + "FieldsUpdate(model *" + m.Name.Pascal + "Model)* querybuilder.Fields {\n"
	output = output + "    return model.BuildFields(false,"
	skiped := ""
	primaryKey := ""
	fields := ""
	for _, v := range m.Columns {
		if v.AutoValue {
			skiped = skiped + "\n		//Skip field \"" + v.Field + "\" which should be set by database"
			skiped = skiped + "\n		//Field \"" + m.Name.Raw + "." + v.Field + "\"\n		//" + m.Name.Pascal + "Field" + MustGetColumnName(v) + ","
		} else if v.PrimayKey {
			primaryKey = primaryKey + "\n		//Skip primary key field \"" + v.Field + "\""
			primaryKey = primaryKey + "\n		//Field \"" + m.Name.Raw + "." + v.Field + "\"\n		//" + m.Name.Pascal + "Field" + MustGetColumnName(v) + ","

		} else {
			fields = fields + "\n		//Field \"" + m.Name.Raw + "." + v.Field + "\"\n		" + m.Name.Pascal + "Field" + MustGetColumnName(v) + ","
		}
	}

	output = output + skiped + fields + primaryKey + "\n    )\n}\n"
	return output
}
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
	c, err := loader.Columns()
	if err != nil {
		return nil, err
	}
	pks := []*columns.Column{}
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
