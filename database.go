package main

import (
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	dsn := "root:123456@tcp(127.0.0.1:3306)/gorm_exercise?charset=utf8mb4&parseTime=True&loc=Local"

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database: " + err.Error())
	}

	DB = db

	fmt.Println("数据库连接成功！")

	err = DB.AutoMigrate(&User{})
	if err != nil {
		panic("数据库迁移失败: " + err.Error())
	}
	fmt.Println("数据库迁移成功，User 表已创建！")

	sqlDB, err := DB.DB() // 获取底层 *sql.DB
	if err != nil {
		panic("获取底层 DB 失败: " + err.Error())
	}

	sqlDB.SetMaxOpenConns(100)                // 最大连接数 100
	sqlDB.SetMaxIdleConns(10)                 // 空闲连接数 10
	sqlDB.SetConnMaxLifetime(5 * time.Minute) // 连接存活时间 5 分钟
}
