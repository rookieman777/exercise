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

	// 运行数据库迁移（创建表）
	fmt.Println("\n🔧 开始数据库迁移...")
	if err := database.Migrate(); err != nil {
		log.Fatalf("数据库迁移失败: %v", err)
	}

	// 清理旧数据（避免重复运行时的冲突）
	fmt.Println("\n🧹 清理演示数据...")
	db := database.GetDB()
	db.Exec("DELETE FROM comments")
	db.Exec("DELETE FROM posts")
	db.Exec("DELETE FROM profiles")
	db.Exec("DELETE FROM users")
	fmt.Println("✅ 旧数据已清理")

	fmt.Println("\n✅ 数据库连接成功，开始演示...")

	// 创建服务实例
	//userService := services.NewUserService()

	// 演示1：基本CRUD操作
	//demoBasicCRUD(userService)

	// 演示2：关联关系和查询
	//demoAssociations()

	// 演示3：事务管理
	//demoTransactions()

	// 演示4：高级查询和统计
	//demoAdvancedQueries(userService)

	// 演示5：性能优化技巧
	//demoPerformanceTips()

	fmt.Println("\n🎉 演示完成！")
}
