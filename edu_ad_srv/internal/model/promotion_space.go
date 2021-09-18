package model

//CREATE TABLE `promotion_space` (
//`id` int(11) NOT NULL AUTO_INCREMENT,
//`name` varchar(255) DEFAULT NULL COMMENT '名称',
//`spaceKey` varchar(255) DEFAULT NULL COMMENT '广告位key',
//`createTime` datetime DEFAULT NULL,
//`updateTime` datetime DEFAULT NULL,
//`isDel` int(2) DEFAULT '0',
//PRIMARY KEY (`id`) USING BTREE,
//KEY `promotion_space_key_isDel` (`spaceKey`,`isDel`) USING BTREE
//) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8 ROW_FORMAT=COMPACT;

type PromotionSpace struct {
	BaseModel
	Name     string `json:"name" gorm:"type:varchar(255);default:null;comment:名称"`
	SpaceKey string `json:"space_key" gorm:"type:varchar(255);default:null;comment:广告位key"`
}

func (PromotionSpace) TableName() string {
	return "promotion_space"
}
