package test

import _ "github.com/herb-go/herb/cache/marshalers/msgpackmarshaler" // msgpack driver

// Cache drivers
// import _ "github.com/herb-go/herb/cache/drivers/freecache"    // freecache driver
import _ "github.com/herb-go/herb/cache/drivers/syncmapcache" // syncmapcachecache driver

// import _ "github.com/herb-go/herb/cache/drivers/versioncache" //versioncache driver
// import _ "github.com/herb-go/providers/redis/rediscache" //rediscache driver
// import _ "github.com/herb-go/providers/sql/sqlcache" //sqlcache driver

//Sql drivers

// import _ "github.com/herb-go/herb/model/sql/querybuilder/drivers/mysql" //mysql driver
// import _ "github.com/herb-go/herb/model/sql/querybuilder/drivers/sqlite" //sqlite driver
// import _ "github.com/herb-go/herb/model/sql/querybuilder/drivers/postgres" //postgres driver
// import _ "github.com/herb-go/herb/model/sql/querybuilder/drivers/mssql" //mssql driver

//UniqueID Drivers
//import _ "github.com/herb-go/uniqueid/uuid"
//import _ "github.com/herb-go/uniqueid/snowflake"
