package test

import _ "github.com/herb-go/herb/cache/marshalers/msgpackmarshaler" // msgpack driver

// Cache drivers
// import _ "github.com/herb-go/herb/cache/drivers/freecache"    // freecache driver
import _ "github.com/herb-go/herb/cache/drivers/syncmapcache" // syncmapcachecache driver

// import _ "github.com/herb-go/herb/cache/drivers/versioncache" //versioncache driver
// import _ "github.com/herb-go/providers/redis/rediscache" //rediscache driver
// import _ "github.com/herb-go/providers/sql/sqlcache" //sqlcache driver

//Sql drivers
// import _ "github.com/go-sql-driver/mysql" //mysql driver

// import _ "github.com/mattn/go-sqlite3" //sqlite driver
// import _ "github.com/lib/pq" //postgresql driver
// import _ "github.com/denisenkom/go-mssqldb" //mssql drvier

//UniqueID Drivers
//import _ "github.com/herb-go/uniqueid/uuid"
//import _ "github.com/herb-go/uniqueid/snowflake"
