package main

import (
	"fmt"
)

func main() {
	//初始化数据库链接
	InitDB()

	//创建对象
	user := User{
		Name: "Alice",
		Age:  20,
	}

	//插入数据库
	result := DB.Create(&user)
	if result.Error != nil {
		fmt.Println("新增用户失败:", result.Error)
		return
	}

	fmt.Println("新增用户成功，ID:", user.ID)

	//查询用户
	var users []User
	result = DB.Find(&users) // result 是 *gorm.DB
	if result.Error != nil {
		fmt.Println("查询失败:", result.Error)
	} else {
		fmt.Println("查询成功，记录数:", result.RowsAffected)
	}

	//修改用户
	result = DB.Model(&user).Updates(User{Name: "Bob", Age: 25})
	if result.Error != nil {
		fmt.Println("更新用户失败:", result.Error)
	} else {
		fmt.Println("更新用户成功")
	}

	var updatedUser User
	result = DB.First(&updatedUser, user.ID)
	if result.Error != nil {
		fmt.Println("查询修改后的用户失败:", result.Error)
	} else {
		fmt.Printf("修改后: ID=%d, Name=%s, Age=%d\n", updatedUser.ID, updatedUser.Name, updatedUser.Age)
	}

	//删除用户
	result = DB.Delete(&updatedUser)
	if result.Error != nil {
		fmt.Println("删除用户失败:", result.Error)
	} else {
		fmt.Println("删除用户成功")
	}

	var checkUser User
	result = DB.First(&checkUser, user.ID)
	if result.Error != nil {
		fmt.Println("查询已删除用户:", result.Error) // 一般会显示 record not found
	} else {
		fmt.Printf("用户仍存在: ID=%d, Name=%s, Age=%d\n", checkUser.ID, checkUser.Name, checkUser.Age)
	}
}
