package model

import (
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	Db          *gorm.DB
	Adapter     *gormadapter.Adapter
	sqlOrganize []tabler
)

type tabler interface {
	TableName() string
}

func InitDb() {
	a, err := gormadapter.NewAdapter("mysql", "root:012345678@tcp(127.0.0.1:3306)/") // Your driver and data source.
	if err != nil {
		panic(err.Error())
	}
	Adapter = a

	Db, err = gorm.Open(mysql.New(mysql.Config{
		DSN: "root:012345678@tcp(127.0.0.1:3306)/demo?charset=utf8mb4&parseTime=True&loc=Local", // DSN data source name
	}), &gorm.Config{})

	AutoMigrate()

}

func register(input ...tabler) {
	sqlOrganize = append(sqlOrganize, input...)
}

func AutoMigrate() {
	for _, v := range sqlOrganize {
		if !Db.Migrator().HasTable(v.TableName()) {
			Db.AutoMigrate(v)
		}
	}
}
