package main

import (
	//"encoding/json"
	//"errors"
	"exercise/database"
	"fmt"
	"log"
	//"os"
	//"strings"
	//"time"
	//"exercise/databse"
	//"exercise/models"
	//"exercise/services"
	//"gorm.io/gorm"
)

func main() {
	fmt.Println("GORM项目练习")
	fmt.Println("==================")

	if err := database.InitDatabase(); err != nil {
		log.Fatalf("初始化数据库失败: %v", err)
	}
	defer database.CloseDatabase() //main结束执行，关闭连接

}
