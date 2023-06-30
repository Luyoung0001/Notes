package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"net/http"
)

// 小清单的小项目
// 前\后端分离
// 本项目只写后端,主要对数据库进行创建,初始化,然后进行一些增删查改

var (
	DB *gorm.DB
)

type Todo struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Status bool   `json:"status"`
}

func initMySQL() (err error) {
	// 连接数据库
	dsn := "root:791975457@qq.com@tcp(127.0.0.1:3306)/bubble?charset=utf8mb4&parseTime=True&loc=Local"
	DB, err = gorm.Open("mysql", dsn)
	if err != nil {
		return err
	}
	// 判断是否连通
	return DB.DB().Ping()
}
func main() {
	// 连接数据库
	err := initMySQL()
	if err != nil {
		panic(err)
	}
	// 模型绑定
	DB.AutoMigrate(&Todo{}) // todos
	defer func(DB *gorm.DB) {
		err := DB.Close()
		if err != nil {
			panic(err)
		}
	}(DB)

	r := gin.Default()
	r.Static("/static", "/Users/luliang/GoLand/gin_practice/chap18/static")
	r.LoadHTMLGlob("/Users/luliang/GoLand/gin_practice/chap18/templates/*")
	r.GET("/bubble", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})
	// v1
	v1Group := r.Group("v1")
	{
		// 待办事项
		// 添加
		v1Group.POST("/todo", func(c *gin.Context) {
			// 前端页面提交待办事项,接受请求
			// 返回响应
			var todo Todo
			c.BindJSON(&todo)
			// 从请求中把数据捞出来,存储到数据库
			if err := DB.Create(&todo).Error; err != nil {
				// 失败时的响应
				c.JSON(http.StatusOK, gin.H{
					"err": err.Error(),
				})
				return
			} else {
				// 成功的响应
				c.JSON(http.StatusOK, todo)
				//c.JSON(http.StatusOK, gin.H{
				//	"code": 2000,
				//	"msg":  "Ok",
				//	"data": todo,
				//})
			}

		})
		// 查看
		// 查看所有的待办事项
		v1Group.GET("/todo", func(c *gin.Context) {
			// 查询todo数据库中的数据所有的数据
			var todoList []Todo
			if err = DB.Find(&todoList).Error; err != nil {
				// 如果出错
				c.JSON(http.StatusOK, gin.H{
					"err": err.Error(),
				})
				return

			} else {
				// 如果成功
				c.JSON(http.StatusOK, todoList)

			}

		})
		// 查看某一个待办事项,前端用不到
		v1Group.GET("/todo/:id", func(c *gin.Context) {

		})

		// 修改(更新) 某一个事项
		v1Group.PUT("/todo/:id", func(c *gin.Context) {
			// 拿到 id,然后查询,修改
			id, ok := c.Params.Get("id")
			if !ok {
				c.JSON(http.StatusOK, gin.H{
					"error": "id不存在",
				})
				return
			}
			var todo Todo
			if err = DB.Where("id=?", id).First(&todo).Error; err != nil {
				// 如果出错
				c.JSON(http.StatusOK, gin.H{
					"err": err.Error(),
				})
				return

			} // 更新
			c.BindJSON(&todo)
			if err = DB.Save(&todo).Error; err != nil {
				c.JSON(http.StatusOK, gin.H{
					"error": err.Error(),
				})
				return

			}
			// 成功
			c.JSON(http.StatusOK, todo)
		})
		// 删除
		v1Group.DELETE("/todo/:id", func(c *gin.Context) {
			// 拿到 id,然后查询,修改
			id, ok := c.Params.Get("id")
			if !ok {
				c.JSON(http.StatusOK, gin.H{
					"error": "id不存在",
				})
				return
			}
			var todo Todo
			if err = DB.Where("id=?", id).Delete(&todo).Error; err != nil {
				// 如果出错
				c.JSON(http.StatusOK, gin.H{
					"err": err.Error(),
				})
				return
			}
			// 成功
			c.JSON(http.StatusOK, todo)
		})

	}
	r.Run(":9001")

}
