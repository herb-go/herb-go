package drivers

import (
	_ "github.com/herb-go/notification-drivers/queue/cronqueue/embeddedqueue" //embeddedqueue notification queue
	_ "github.com/herb-go/notification-drivers/queue/passthroughqueue"        //passthroughqueue notification queue

	_ "github.com/herb-go/notification-drivers/store/embeddedstore/embeddeddraftbox" //embeddeddraftbox notification draft box
)
