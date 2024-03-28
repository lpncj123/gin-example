package main

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)
import "github.com/thinkerou/favicon"

// 自定义中间件
func myHandler() gin.HandlerFunc {
	return func(context *gin.Context) {
		context.Set("userSession", "userid_1")
		context.Next() //放行
		//context.Abort() //拦截
	}
}
func main() {
	//创建一个服务
	ginServer := gin.Default()
	//设置访问左上角小图标
	ginServer.Use(favicon.New("./favicon.ico"))
	//访问地址，处理我们请求 Request Response
	//gin restful风格是十分简单的
	ginServer.GET("/hello", func(context *gin.Context) {
		context.JSON(200, gin.H{
			"message": "helloWorld",
		})
	})
	//加载静态页面
	ginServer.LoadHTMLGlob("templates/*")
	//加载资源文件
	ginServer.Static("/static", "./static")
	//响应一个页面给前端
	ginServer.GET("/index", func(context *gin.Context) {
		context.HTML(200, "index.html", gin.H{
			"msg": "gin第一次使用",
		})
	})
	//接收前端传递过来参数
	//usl?userid=xxx&username=lp
	ginServer.GET("/user/info", myHandler(), func(context *gin.Context) {
		userSession := context.MustGet("userSession").(string)
		log.Println("+======>", userSession)
		userid := context.Query("userid")
		username := context.Query("username")
		context.JSON(http.StatusOK, gin.H{"userid": userid, "username": username})
	})

	//user/info/1/lp
	ginServer.GET("/user/info/:userid/:username", func(context *gin.Context) {
		userid := context.Param("userid")
		username := context.Param("username")
		context.JSON(http.StatusOK, gin.H{"userid": userid, "username": username})
	})

	//前端给后端传递JSON
	ginServer.POST("/json", func(context *gin.Context) {
		//request body
		data, _ := context.GetRawData()
		var m map[string]interface{}
		//包装为JSON对象[]byte
		_ = json.Unmarshal(data, &m)
		context.JSON(http.StatusOK, m)
	})
	//支持函数式编程=》
	ginServer.POST("/user/add", func(context *gin.Context) {
		username := context.PostForm("username")
		password := context.PostForm("password")
		context.JSON(http.StatusOK, gin.H{
			"username": username,
			"password": password,
			"msg":      "ok",
		})
	})
	//重定向
	ginServer.GET("/test", func(context *gin.Context) {
		context.Redirect(http.StatusMovedPermanently, "https://www.baidu.com")
	})

	//404 Router
	ginServer.NoRoute(func(context *gin.Context) {
		context.HTML(http.StatusNotFound, "404.html", nil)
	})
	//路由组 /user/add
	userGroup := ginServer.Group("/user")
	{
		userGroup.GET("/add")
		userGroup.POST("/login")
		userGroup.POST("/logout")
	}
	orderGroup := ginServer.Group("/order")
	{
		orderGroup.GET("/add")
		orderGroup.DELETE("/delete")
	}
	//设置服务器端口
	ginServer.Run(":8082")
}
