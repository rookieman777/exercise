package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"exercise/database"
	"exercise/models"
	"exercise/services"

	"gorm.io/gorm"
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
	userService := services.NewUserService()

	// 演示1：基本CRUD操作
	demoBasicCRUD(userService)

	// 演示2：关联关系和查询
	demoAssociations()

	// 演示3：事务管理
	demoTransactions()

	// 演示4：高级查询和统计
	demoAdvancedQueries(userService)

	// 演示5：性能优化技巧
	demoPerformanceTips()

	fmt.Println("\n🎉 演示完成！")
}

func demoBasicCRUD(service services.UserService) {
	fmt.Println("\n1️⃣ 基本CRUD操作演示")
	fmt.Println("----------------")

	// 1.1 创建用户
	user1 := &models.User{
		Username: "john_doe",
		Email:    "john@example.com",
		Password: "SecurePass123",
		Age:      25,
		IsActive: true,
	}

	fmt.Println("\n📝 创建用户:")
	if err := service.Register(user1); err != nil {
		log.Printf("创建用户失败: %v", err)
	} else {
		fmt.Printf("✅ 用户创建成功: %s (ID: %d)\n", user1.Username, user1.ID)
	}

	// 1.2 查询用户
	fmt.Println("\n🔍 查询用户:")
	fetchedUser, err := service.GetUserByID(user1.ID)
	if err != nil {
		log.Printf("查询用户失败: %v", err)
	} else {
		fmt.Printf("✅ 查询到用户: %s (邮箱: %s)\n", fetchedUser.Username, fetchedUser.Email)
	}

	// 1.3 更新用户 //这个功能没有使用接口，直接连接数据库了
	fmt.Println("\n✏️ 更新用户:")
	user1.Age = 26
	user1.Email = "john.updated@example.com"
	// 使用数据库直接更新
	db := database.GetDB()
	if err := db.Model(user1).Updates(map[string]interface{}{
		"age":   user1.Age,
		"email": user1.Email,
	}).Error; err != nil {
		log.Printf("更新用户失败: %v", err)
	} else {
		fmt.Printf("✅ 用户更新成功: 年龄更新为 %d\n", user1.Age)
	}

	// 1.4 软删除用户
	fmt.Println("\n🗑️ 软删除用户:")
	if err := service.DeactivateAccount(user1.ID); err != nil {
		log.Printf("删除用户失败: %v", err)
	} else {
		fmt.Println("✅ 用户已软删除（停用）")
	}

	// 1.5 分页查询
	fmt.Println("\n📄 分页查询演示:")
	users, total, err := service.SearchUsers("", 1, 10)
	if err != nil {
		log.Printf("分页查询失败: %v", err)
	} else {
		fmt.Printf("✅ 分页查询结果: 第1页，每页10条，共%d条记录\n", total)
		for _, u := range users {
			fmt.Printf("   - %s (%s)\n", u.Username, u.Email)
		}
	}

}

// demoAssociations 演示关联关系和查询
func demoAssociations() {
	fmt.Println("\n2️⃣ 关联关系演示")
	fmt.Println("--------------")

	db := database.GetDB()

	// 2.1 创建具有关联数据的用户
	fmt.Println("\n🤝 创建带关联数据的用户:")
	user := &models.User{
		Username: "alice_smith",
		Email:    "alice@example.com",
		Password: "AlicePass456",
		Age:      30,
		Profile: &models.Profile{
			FirstName: "Alice",
			LastName:  "Smith",
			Bio:       "Software Engineer",
			Location:  "San Francisco",
		},
		Posts: []models.Post{
			{
				Title:   "我的第一篇博客",
				Content: "这是Alice的第一篇博客内容...",
				Slug:    "my-first-post",
				Status:  "published",
			},
		},
	}

	// 使用事务创建关联数据
	err := db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(user).Error; err != nil {
			return fmt.Errorf("创建用户失败: %v", err)
		}
		fmt.Printf("✅ 用户创建成功: ID=%d\n", user.ID)
		return nil
	})

	if err != nil {
		log.Printf("创建关联数据失败: %v", err)
		return
	}

	// 2.2 预加载关联数据
	fmt.Println("\n🔍 预加载关联数据:")
	var loadedUser models.User
	err = db.Preload("Profile").Preload("Posts").Preload("Posts.Comments").First(&loadedUser, user.ID).Error
	if err != nil {
		log.Printf("预加载失败: %v", err)
	} else {
		fmt.Printf("✅ 用户: %s\n", loadedUser.Username)
		if loadedUser.Profile != nil {
			fmt.Printf("   📝 资料: %s %s - %s\n",
				loadedUser.Profile.FirstName, loadedUser.Profile.LastName,
				loadedUser.Profile.Location)
		}
		fmt.Printf("   📰 文章数: %d\n", len(loadedUser.Posts))
	}

	// 2.3 关联查询
	fmt.Println("\n🔗 关联查询:")
	type UserWithPostCount struct {
		ID        uint
		Username  string
		Email     string
		PostCount int
	}

	var usersWithPosts []UserWithPostCount
	err = db.Model(&models.User{}).
		Select("users.id, users.username, users.email, COUNT(posts.id) as post_count").
		Joins("LEFT JOIN posts ON posts.author_id = users.id").
		Group("users.id").
		Having("post_count > 0").
		Find(&usersWithPosts).Error

	if err != nil {
		log.Printf("关联查询失败: %v", err)
	} else {
		fmt.Println("✅ 用户及其文章数统计:")
		for _, u := range usersWithPosts {
			fmt.Printf("   - %s: %d 篇文章\n", u.Username, u.PostCount)
		}
	}
}

// demoTransactions 演示事务管理
func demoTransactions() {
	fmt.Println("\n3️⃣ 事务管理演示")
	fmt.Println("--------------")

	db := database.GetDB()

	// 3.1 简单事务示例
	fmt.Println("\n🔁 简单事务:")
	err := db.Transaction(func(tx *gorm.DB) error {
		// 操作1：创建用户
		user := &models.User{
			Username: "bob_jones",
			Email:    "bob@example.com",
			Password: "BobPass789",
			Age:      35,
		}
		if err := tx.Create(user).Error; err != nil {
			return fmt.Errorf("创建用户失败: %v", err)
		}
		fmt.Printf("✅ 步骤1: 用户创建成功 (ID: %d)\n", user.ID)

		// 操作2：创建用户资料
		profile := &models.Profile{
			UserID:    user.ID,
			FirstName: "Bob",
			LastName:  "Jones",
			Bio:       "Database Administrator",
		}
		if err := tx.Create(profile).Error; err != nil {
			return fmt.Errorf("创建资料失败: %v", err)
		}
		fmt.Printf("✅ 步骤2: 用户资料创建成功\n")

		// 操作3：创建文章
		post := &models.Post{
			AuthorID: user.ID,
			Title:    "数据库优化技巧",
			Content:  "分享一些数据库性能优化的实践经验...",
			Slug:     "database-optimization",
			Status:   "published",
		}
		if err := tx.Create(post).Error; err != nil {
			return fmt.Errorf("创建文章失败: %v", err)
		}
		fmt.Printf("✅ 步骤3: 文章创建成功\n")

		return nil // 提交事务
	})

	if err != nil {
		log.Printf("事务执行失败: %v", err)
	} else {
		fmt.Println("🎉 所有操作已成功提交")
	}

	// 3.2 嵌套事务示例，内层回滚不影响外层
	fmt.Println("\n🔁 嵌套事务:")
	err = db.Transaction(func(tx *gorm.DB) error {
		// 外层事务
		user := &models.User{
			Username: "carol_wilson",
			Email:    "carol@example.com",
			Password: "CarolPass101",
			Age:      28,
		}
		if err := tx.Create(user).Error; err != nil {
			return err
		}
		fmt.Printf("✅ 外层事务: 用户创建成功\n")

		// 嵌套事务（保存点）
		nestedErr := tx.Transaction(func(tx2 *gorm.DB) error {
			// 内层事务操作
			profile := &models.Profile{
				UserID:    user.ID,
				FirstName: "Carol",
				LastName:  "Wilson",
			}
			if err := tx2.Create(profile).Error; err != nil {
				return err
			}
			fmt.Printf("✅ 内层事务: 资料创建成功\n")

			// 模拟一个可能失败的操作
			var count int64
			if err := tx2.Model(&models.User{}).Where("email = ?", "nonexistent@example.com").Count(&count).Error; err != nil {
				fmt.Println("⚠️ 内层事务: 查询失败（预期行为）")
				return err // 这将回滚内层事务但不影响外层
			}

			return nil
		})

		if nestedErr != nil {
			fmt.Printf("⚠️ 内层事务已回滚，但外层事务继续执行\n")
		}

		// 外层事务继续执行其他操作
		post := &models.Post{
			AuthorID: user.ID,
			Title:    "嵌套事务示例",
			Content:  "这是一个嵌套事务的演示...",
			Status:   "draft",
		}
		if err := tx.Create(post).Error; err != nil {
			return err
		}
		fmt.Printf("✅ 外层事务: 文章创建成功\n")

		return nil
	})

	if err != nil {
		log.Printf("嵌套事务失败: %v", err)
	} else {
		fmt.Println("🎉 嵌套事务执行完成")
	}

}

// demoAdvancedQueries 演示高级查询和统计
func demoAdvancedQueries(service services.UserService) {
	fmt.Println("\n4️⃣ 高级查询和统计")
	fmt.Println("----------------")

	// 4.1 复杂条件查询
	fmt.Println("\n🔍 复杂条件查询:")
	db := database.GetDB()

	type UserStats struct {
		Status   string
		AgeGroup string
		Count    int
	}

	var stats []UserStats
	err := db.Model(&models.User{}).
		Select(`
			CASE 
				WHEN is_active = 1 THEN '活跃'
				ELSE '非活跃'
			END as status,
			CASE 
				WHEN age < 20 THEN '青少年'
				WHEN age BETWEEN 20 AND 40 THEN '青年'
				WHEN age BETWEEN 41 AND 60 THEN '中年'
				ELSE '老年'
			END as age_group,
			COUNT(*) as count
		`).
		Group("status, age_group").
		Order("count DESC").
		Find(&stats).Error

	if err != nil {
		log.Printf("统计查询失败: %v", err)
	} else {
		fmt.Println("✅ 用户统计分组:")
		for _, s := range stats {
			fmt.Printf("   - %s | %s: %d人\n", s.Status, s.AgeGroup, s.Count)
		}
	}

	// 4.2 原生SQL查询
	fmt.Println("\n🔍 原生SQL查询:")
	var activeUserCount int64
	err = db.Raw(`
		SELECT COUNT(*) 
		FROM users 
		WHERE is_active = 1 
		AND created_at > ?
	`, time.Now().AddDate(0, -1, 0)).
		Scan(&activeUserCount).Error

	if err != nil {
		log.Printf("原生SQL查询失败: %v", err)
	} else {
		fmt.Printf("✅ 最近一个月活跃用户: %d人\n", activeUserCount)
	}

	// 4.3 统计功能
	fmt.Println("\n📊 统计功能演示:")
	statsData, err := service.GetUserStats()
	if err != nil {
		log.Printf("获取统计失败: %v", err)
	} else {
		fmt.Println("✅ 用户统计数据:")
		fmt.Printf("   总用户数: %d\n", statsData.TotalUsers)
		fmt.Printf("   活跃用户: %d\n", statsData.ActiveUsers)
		fmt.Printf("   今日注册: %d\n", statsData.TodayRegisters)
		fmt.Printf("   平均年龄: %.2f\n", statsData.AvgAge)
		if len(statsData.TopDomains) > 0 {
			fmt.Printf("   常用邮箱域名: %v\n", statsData.TopDomains)
		}
	}
}

// demoPerformanceTips 演示性能优化技巧
func demoPerformanceTips() {
	fmt.Println("\n5️⃣ 性能优化技巧")
	fmt.Println("----------------")

	db := database.GetDB()

	// 5.1 查询优化：只选择需要的字段
	fmt.Println("\n🔍 查询优化 - 只选择需要的字段:")
	type MinimalUser struct {
		ID       uint
		Username string
		Email    string
	}

	var minimalUsers []MinimalUser
	err := db.Model(&models.User{}).
		Select("id, username, email").
		Limit(5).
		Find(&minimalUsers).Error

	if err != nil {
		log.Printf("查询优化失败: %v", err)
	} else {
		fmt.Println("✅ 只查询必要字段，减少数据传输")
		for _, u := range minimalUsers {
			fmt.Printf("   - %s (%s)\n", u.Username, u.Email)
		}
	}

	// 5.2 批量操作
	fmt.Println("\n📦 批量操作优化:")
	users := []models.User{}
	for i := 1; i <= 5; i++ {
		users = append(users, models.User{
			Username: fmt.Sprintf("bulk_user_%d", i),
			Email:    fmt.Sprintf("bulk%d@example.com", i),
			Password: "BulkPass123",
			Age:      20 + i,
		})
	}

	// 批量创建
	if err := db.CreateInBatches(users, 3).Error; err != nil {
		log.Printf("批量创建失败: %v", err)
	} else {
		fmt.Println("✅ 批量创建成功（分批处理，每批3条）")
	}

	// 5.3 连接池监控
	fmt.Println("\n📊 连接池状态:")

	sqlDB, err := db.DB()
	if err != nil {
		log.Printf("获取连接池失败: %v", err)
	} else {
		stats := sqlDB.Stats()
		fmt.Printf("✅ 连接池状态:\n")
		fmt.Printf("   最大连接数: %d\n", stats.MaxOpenConnections)
		fmt.Printf("   打开连接: %d\n", stats.OpenConnections)
		fmt.Printf("   正在使用: %d\n", stats.InUse)
		fmt.Printf("   空闲连接: %d\n", stats.Idle)
	}

	// 5.4 索引使用演示
	fmt.Println("\n🔍 索引使用演示:")
	// 查询使用索引的字段
	var indexedUsers []models.User
	startTime := time.Now()
	err = db.Where("username LIKE ?", "bulk%").
		Find(&indexedUsers).Error
	if err != nil {
		log.Printf("索引查询失败: %v", err)
	} else {
		duration := time.Since(startTime)
		fmt.Printf("✅ 索引查询完成，耗时: %v，查询到 %d 条记录\n", duration, len(indexedUsers))
	}
}

// demoErrorHandling 演示错误处理
func demoErrorHandling() {
	fmt.Println("\n⚠️ 错误处理演示")
	fmt.Println("---------------")

	db := database.GetDB()

	// 6.1 GORM错误类型
	fmt.Println("\n🔍 GORM错误类型:")
	var user models.User
	err := db.Where("id = ?", 999999).First(&user).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			fmt.Println("✅ 正确处理: 记录不存在")
		} else if strings.Contains(err.Error(), "connection") {
			fmt.Println("⚠️ 连接错误")
		} else {
			fmt.Printf("❌ 其他错误: %v\n", err)
		}
	}

	// 6.2 事务错误处理
	fmt.Println("\n🔍 事务错误处理:")
	err = db.Transaction(func(tx *gorm.DB) error {
		// 尝试创建重复数据
		dupUser := &models.User{
			Username: "john_doe", // 已存在
			Email:    "duplicate@example.com",
			Password: "Pass123",
		}
		if err := tx.Create(dupUser).Error; err != nil {
			fmt.Printf("✅ 事务捕获错误: %v\n", err)
			return err // 回滚事务
		}
		return nil
	})

	if err != nil {
		fmt.Printf("✅ 事务已回滚: %v\n", err)
	}
}

// demoHooks 演示钩子和回调
func demoHooks() {
	fmt.Println("\n🔔 钩子(Hooks)演示")
	fmt.Println("---------------")

	// 注意: 钩子在models/user.go中定义
	fmt.Println("1. BeforeCreate - 创建前自动设置年龄默认值")
	fmt.Println("2. AfterCreate - 创建后自动创建Profile")
	fmt.Println("3. BeforeUpdate - 更新前记录日志")
	fmt.Println("\n✅ 钩子已定义在models/user.go中")
}

// 导出演示报告
func generateReport(){
	//创建报告结构
	report := map[string]interface{}{
		"timestamp": time.Now().Format(time.RFC3339),
		"demonstrations": []string{
			"基本CRUD操作",
			"关联关系查询",
			"事务管理",
			"高级查询统计",
			"性能优化",
			"错误处理",
			"钩子(Hooks)",
		},
		"best_practices": []string{
			"使用连接池提高性能",
			"批量操作减少数据库请求",
			"正确使用事务确保数据一致性",
			"预加载关联数据避免N+1查询",
			"合理使用索引优化查询速度",
			"及时关闭数据库连接",
			"编写可读的查询语句",
		},
		"learning_outcomes": []string{
			"掌握了 GORM 的基本 CRUD 操作",
			"理解了关联关系的定义和使用",
			"学会了使用事务确保数据完整性",
			"掌握了查询优化技巧",
			"了解了错误处理的最佳实践",
			"学会了如何使用钩子扩展功能",
		},
	}

	// 导出为 JSON
	jsonData, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		log.Printf("生成报告失败: %v", err)
		return
	}

	filename := fmt.Sprintf("gorm_demo_report_%s.json", time.Now().Format("20060102_150405"))
	if err := os.WriteFile(filename, jsonData, 0644); err != nil {
		log.Printf("保存报告失败: %v", err)
	} else {
		fmt.Printf("\n📋 演示报告已生成: %s\n", filename)
	}
}