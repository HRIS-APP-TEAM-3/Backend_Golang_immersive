package database

import (
	"fmt"
	"hris-app-golang/app/config"
	divisionsData "hris-app-golang/feature/divisions/data"
	rolesData "hris-app-golang/feature/roles/data"
	usersData "hris-app-golang/feature/users/data"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitMysql(cfg *config.AppConfig) *gorm.DB {
	connectionString := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
		cfg.DBUsername, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)

	DB, err := gorm.Open(mysql.Open(connectionString), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	return DB
}

func InittialMigration(db *gorm.DB) {
	db.AutoMigrate(&usersData.User{})
	db.AutoMigrate(&usersData.UserImportant{})
	db.AutoMigrate(&usersData.UserEducation{})
	db.AutoMigrate(&rolesData.Role{})
	db.AutoMigrate(&divisionsData.Division{})
}
