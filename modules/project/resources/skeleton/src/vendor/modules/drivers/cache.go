package drivers

// Cache drivers
// import _ "github.com/herb-go/herb/cache/drivers/freecache"    // freecache driver
import _ "github.com/herb-go/herb/cache/drivers/syncmapcache" // syncmapcachecache driver

// import _ "github.com/herb-go/herb/cache/drivers/versioncache" //versioncache driver
// import _ "github.com/herb-go/herb-drivers/cache/rediscache" //rediscache driver
// import _ "github.com/herb-go/providers/sql/sqlcache" //sqlcache driver
import _ "github.com/herb-go/herb-drivers/cache/hiredcache" //hired worker cache driver
