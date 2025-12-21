package database

import (
	"exercise/models"
	"fmt"
	"log"
	"strings"

	"gorm.io/gorm"
)

// Migrate 运行数据库迁移
func Migrate() error {
	db := GetDB()
	if db == nil {
		return fmt.Errorf("数据库连接未初始化")
	}

	log.Println("开始数据库迁移...")

	// 禁用外键约束（避免顺序问题）
	if err := db.Exec("SET FOREIGN_KEY_CHECKS = 0").Error; err != nil {
		return fmt.Errorf("禁用外键约束失败: %v", err)
	}

	// 按依赖顺序创建表
	tables := []interface{}{
		&models.User{},
		&models.Profile{},
		&models.Post{},
		&models.Comment{},
		&models.Course{},
	}

	for _, table := range tables {
		if err := db.AutoMigrate(table); err != nil {
			return fmt.Errorf("迁移表 %T 失败: %v", table, err)
		}
		log.Printf("✅ 表已迁移: %T", table)
	}

	// 重新启用外键约束
	if err := db.Exec("SET FOREIGN_KEY_CHECKS = 1").Error; err != nil {
		return fmt.Errorf("启用外键约束失败: %v", err)
	}

	// 创建索引
	if err := createIndexes(db); err != nil {
		return fmt.Errorf("创建索引失败: %v", err)
	}

	log.Println("✅ 数据库迁移完成")
	return nil
}

// createIndexes 创建额外的索引
func createIndexes(db *gorm.DB) error {
	indexes := []struct {
		table   string
		columns []string
	}{
		{"users", []string{"email"}},
		{"users", []string{"username"}},
		{"posts", []string{"author_id"}},
		{"posts", []string{"slug"}},
		{"comments", []string{"user_id", "post_id"}},
	}

	for _, idx := range indexes {
		indexName := fmt.Sprintf("idx_%s_%s", idx.table, strings.Join(idx.columns, "_"))

		// 先检查索引是否存在（兼容MySQL 5.x）
		var count int64
		checkQuery := fmt.Sprintf(`
			SELECT COUNT(*) 
			FROM information_schema.statistics 
			WHERE table_schema = DATABASE() 
			AND table_name = '%s' 
			AND index_name = '%s'
		`, idx.table, indexName)

		if err := db.Raw(checkQuery).Scan(&count).Error; err != nil {
			log.Printf("⚠️ 检查索引 %s 失败: %v", indexName, err)
			continue
		}

		// 如果索引不存在，则创建
		if count == 0 {
			query := fmt.Sprintf("CREATE INDEX %s ON %s (%s)",
				indexName, idx.table, strings.Join(idx.columns, ", "))
			if err := db.Exec(query).Error; err != nil {
				log.Printf("⚠️ 创建索引 %s 失败: %v", indexName, err)
			} else {
				log.Printf("✅ 索引已创建: %s", indexName)
			}
		} else {
			log.Printf("ℹ️ 索引已存在: %s", indexName)
		}
	}

	return nil
}

// DropAll 删除所有表（仅用于开发）
func DropAll() error {
	db := GetDB()
	if db == nil {
		return fmt.Errorf("数据库连接未初始化")
	}

	log.Println("删除所有表...")

	tables := []interface{}{
		&models.Comment{},
		&models.Post{},
		&models.Course{},
		&models.Profile{},
		&models.User{},
	}

	// 禁用外键约束
	if err := db.Exec("SET FOREIGN_KEY_CHECKS = 0").Error; err != nil {
		return fmt.Errorf("禁用外键约束失败: %v", err)
	}

	for _, table := range tables {
		if err := db.Migrator().DropTable(table); err != nil {
			log.Printf("删除表 %T 失败: %v", table, err)
		} else {
			log.Printf("🗑️ 表已删除: %T", table)
		}
	}

	// 启用外键约束
	if err := db.Exec("SET FOREIGN_KEY_CHECKS = 1").Error; err != nil {
		return fmt.Errorf("启用外键约束失败: %v", err)
	}

	log.Println("✅ 所有表已删除")
	return nil
}

// Reset 重置数据库（删除并重新创建）
func Reset() error {
	if err := DropAll(); err != nil {
		return err
	}
	return Migrate()
}

// CheckStatus 检查数据库状态
func CheckStatus() error {
	db := GetDB()
	if db == nil {
		return fmt.Errorf("数据库连接未初始化")
	}

	// 检查连接
	var result int
	if err := db.Raw("SELECT 1").Scan(&result).Error; err != nil {
		return fmt.Errorf("数据库连接检查失败: %v", err)
	}

	// 检查表状态
	tables := []string{"users", "profiles", "posts", "comments", "courses"}
	for _, table := range tables {
		var exists bool
		query := fmt.Sprintf("SELECT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_schema = DATABASE() AND table_name = '%s')", table)
		if err := db.Raw(query).Scan(&exists).Error; err != nil {
			log.Printf("检查表 %s 失败: %v", table, err)
			continue
		}
		if exists {
			log.Printf("✅ 表存在: %s", table)
		} else {
			log.Printf("⚠️ 表不存在: %s", table)
		}
	}

	log.Println("✅ 数据库状态检查完成")
	return nil
}
