package model

// CommonModel 公共模型
type CommonModel struct {
	ID         int64 `gorm:"column:id;type:bigint(20);not null;primary_key"`
	Status     int8  `gorm:"column:status;type:tinyint(4);default:0;not null"`
	Createtime int64 `gorm:"column:createtime" json:"createtime" form:"createtime"`
	Updatetime int64 `gorm:"column:updatetime" json:"updatetime" form:"updatetime"`
}
