package entity

type User struct {
	ID       int64  `gorm:"column:id;type:bigint(20);not null;primary_key"`
	Username string `gorm:"column:username;type:varchar(255);not null;index:username,unique"`
	Field1   int64  `gorm:"column:field1;type:int(11);index:unq,unique"`
	Field2   string `gorm:"column:field2;type:varchar(255);not null;index:unq,unique;index:index"`
	Field3   int32  `gorm:"column:field3;type:tinyint(4);index:index12"`
	Field4   string `gorm:"column:field4;type:varchar(255)"`
}
