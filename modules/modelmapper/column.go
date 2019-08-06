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
	Raw         string
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

func getLoaderFormDB(conn db.Database) (columns.Loader, error) {
	drivername := conn.Driver()
	driver, err := columns.Driver(drivername)
	if err != nil {
		return nil, err
	}
	return driver, nil
}

func NewModelCulumns(conn db.Database, database string, table string, field_prefix string) (*ModelColumns, error) {
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
