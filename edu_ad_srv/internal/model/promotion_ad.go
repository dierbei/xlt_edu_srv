package model

import "time"

//CREATE TABLE `promotion_ad` (
//`id` int(11) NOT NULL AUTO_INCREMENT,
//`name` varchar(255) DEFAULT NULL COMMENT '广告名',
//`spaceId` int(11) DEFAULT NULL COMMENT '广告位id',
//`keyword` varchar(255) DEFAULT NULL COMMENT '精确搜索关键词',
//`htmlContent` text COMMENT '静态广告的内容',
//`text` varchar(255) DEFAULT NULL COMMENT '文字一',
//`link` varchar(255) DEFAULT NULL COMMENT '链接一',
//`startTime` datetime DEFAULT NULL COMMENT '开始时间',
//`endTime` datetime DEFAULT NULL COMMENT '结束时间',
//`createTime` datetime DEFAULT NULL,
//`updateTime` datetime DEFAULT NULL,
//`status` int(2) NOT NULL DEFAULT '0',
//`priority` int(4) DEFAULT '0' COMMENT '优先级',
//`img` varchar(255) DEFAULT NULL,
//PRIMARY KEY (`id`) USING BTREE,
//KEY `promotion_ad_SEG` (`spaceId`,`startTime`,`endTime`,`status`) USING BTREE
//) ENGINE=InnoDB AUTO_INCREMENT=1095 DEFAULT CHARSET=utf8 ROW_FORMAT=COMPACT;

type PromotionAd struct {
	BaseModel
	Name        string    `json:"name" gorm:"type:varchar(255);not null;comment:广告名"`
	SpaceID     int       `json:"space_id" gorm:"type:int(11);default:null;comment:广告位id"`
	Keyword     string    `json:"keyword" gorm:"type:varchar(255);default:null;comment:精确搜索关键词"`
	HtmlContent string    `json:"html_content" gorm:"type:text;comment:静态广告的内容"`
	Text        string    `json:"text" gorm:"type:varchar(255);default:null;comment:文字"`
	Link        string    `json:"link" gorm:"type:varchar(255);default:null;comment:链接"`
	StartTime   time.Time `json:"start_time" gorm:"default:null;comment:开始时间"`
	EndTime     time.Time `json:"end_time" gorm:"default:null;comment:结束时间"`
	Status      int       `json:"status" gorm:"type:int(2);default:0"`
	Priority    int       `json:"priority" gorm:"type:int(4);default:0;comment:优先级"`
	Img         string    `json:"img" gorm:"type:varchar(255);default:null"`
}

func (PromotionAd) TableName() string {
	return "promotion_ad"
}
