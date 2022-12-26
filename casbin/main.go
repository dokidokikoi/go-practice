package main

import (
	"fmt"
	"log"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/gin-gonic/gin"
)

func main() {
	// 使用postgres数据库初始化一个Gorm适配器
	a, err := gormadapter.NewAdapter("postgres", "host=127.0.0.1 user=postgres dbname=go_blog port=5432 sslmode=disable TimeZone=Asia/Shanghai password=root", true)
	if err != nil {
		log.Fatalf("error: adapter: %s", err)
	}
	m, err := model.NewModelFromString(`
		[request_definition]
		r = sub, obj, act

		[policy_definition]
		p = sub, obj, act

		[role_definition]
		g = _, _

		[policy_effect]
		e = some(where (p.eft == allow))

		[matchers]
		m = r.sub == p.sub && r.obj == p.obj && r.act == p.act
	`)
	if err != nil {
		log.Fatalf("error: model: %s", err)
	}

	e, err := casbin.NewEnforcer(m, a)
	if err != nil {
		log.Fatalf("error: enforcer: %s", err)
	}

	//从DB加载策略
	e.LoadPolicy()

	//获取router路由对象
	r := gin.New()

	r.POST("/api/v1/add", func(c *gin.Context) {
		fmt.Println("增加Policy")
		if ok, _ := e.AddPolicy("admin", "/api/v1/hello", "GET"); !ok {
			fmt.Println("Policy已经存在")
		} else {
			fmt.Println("增加成功")
		}
	})
	//删除policy
	r.DELETE("/api/v1/delete", func(c *gin.Context) {
		fmt.Println("删除Policy")
		if ok, _ := e.RemovePolicy("admin", "/api/v1/hello", "GET"); !ok {
			fmt.Println("Policy不存在")
		} else {
			fmt.Println("删除成功")
		}
	})
	//获取policy
	r.GET("/api/v1/get", func(c *gin.Context) {
		fmt.Println("查看policy")
		list := e.GetPolicy()
		for _, vlist := range list {
			for _, v := range vlist {
				fmt.Printf("value: %s, ", v)
			}
		}
	})
	//使用自定义拦截器中间件
	r.Use(Authorize(e))
	//创建请求
	r.GET("/api/v1/hello", func(c *gin.Context) {
		fmt.Println("Hello 接收到GET请求..")
	})

	r.Run(":9000") //参数为空 默认监听8080端口
}

// 拦截器
func Authorize(e *casbin.Enforcer) gin.HandlerFunc {

	return func(c *gin.Context) {

		//获取请求的URI
		obj := c.Request.URL.RequestURI()
		//获取请求方法
		act := c.Request.Method
		//获取用户的角色
		sub := "admin"

		//判断策略中是否存在
		if ok, _ := e.Enforce(sub, obj, act); ok {
			fmt.Println("恭喜您,权限验证通过")
			c.Next()
		} else {
			fmt.Println("很遗憾,权限验证没有通过")
			c.Abort()
		}
	}
}
