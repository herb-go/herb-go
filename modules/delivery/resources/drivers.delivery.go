package drivers

import (
	_ "github.com/herb-go/notification-drivers/delivery/aliyunsmsdelivery"            //aliyunsmsdelivery delivery
	_ "github.com/herb-go/notification-drivers/delivery/emaildelivery"                //eamil delivery
	_ "github.com/herb-go/notification-drivers/delivery/loggerdelivery"               //logger delivery
	_ "github.com/herb-go/notification-drivers/delivery/mockdelivery"                 //mock delivery
	_ "github.com/herb-go/notification-drivers/delivery/noitificationapidelivery"     //noitification api delivery
	_ "github.com/herb-go/notification-drivers/delivery/tencentcloudsmsdelivery"      //tencentcloudsms delivery
	_ "github.com/herb-go/notification-drivers/delivery/tencentminiprogramumdelivery" //tencentminiprogramum delivery
	_ "github.com/herb-go/notification-drivers/delivery/wechattmdelivery"             //wechattm delivery
	_ "github.com/herb-go/notification-drivers/delivery/wechatworkdelivery"           //wechatwork delivery
)
