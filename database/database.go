package database

import (
	"fmt"
	"udabayar-go-api-di/config"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func ConnectDatabase(dbconfig *config.DBConfig) (*gorm.DB, error) {
	dbUri := fmt.Sprintf(
		"%s:%s@/%s?charset=utf8&parseTime=True&loc=Local",
		dbconfig.User,
		dbconfig.Password,
		dbconfig.Name,
	)

	return gorm.Open(dbconfig.Client, dbUri)
}
