# Casbin 的两个重要文件
casbin 有两个配置文件, `model.conf` 和 `policy.csv`。其中, `model.conf` 存储了我们的访问模型, 而 `policy.csv` 存储的是我们具体的用户权限配置。

## 模型文件
Casbin 中最基本，最简单的模型是 ACL。ACL 的 CONF 模型为：
```conf
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
```
声明了：
1.  `r = sub, obj, act` 定义有限请求将由 3 部分组成：*sub*ject - 用户，*obj*ect - URL 或更一般的资源，最后是 *act*ion - 操作。
2.  `p = sub, obj, act` 定义策略的格式。例如，`admin, data, write` 意味着 `All admins can write data.`
3.  `e = some(where (p.eft == allow))`意味着用户可以做某事，只要有定义的策略允许他这样做。
4.  `g = _, _` 定义用户角色定义的格式。例如，`Alice, admin` 表示 Alice 是管理员。
5.  `m = g(r.sub, p.sub) && r.obj == p.obj && r.act == p.act` 定义授权的工作流程：检查用户的角色 -> 检查用户试图访问的资源 -> 检查用户想做什么。

## 权限配置文件
以下是权限配置：
```
p, user, data, read
p, admin, data, read
p, admin, data, write
g, Alice, admin
g, Bob, user
```
定义了 3 个权限和两个角色：
1. user 角色可以读数据
2. admin 角色可以读数据
3. admin 角色可以写数据

我们也可以使用数据库来作为权限配置的文件：
```go
// 使用postgres数据库初始化一个Gorm适配器
a, err := gormadapter.NewAdapter(
"postgres", 
"host=127.0.0.1 user=postgres dbname=go_blog port=5432 sslmode=disable TimeZone=Asia/Shanghai password=root", 
true)
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
```


# 示例代码
```go
func main() {
	// ...省略部分代码

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
```



访问本地 9000 端口的 `/api/v1/hello` ，控制台会打印:

```
很遗憾,权限验证没有通过
```

访问本地 9000 端口的 `/api/v1/add` ，添加 `admin, /api/v1/hello, GET` 策略

之后再次访问 `/api/v1/hello` ，控制台会打印:

```
恭喜您,权限验证通过
Hello 接收到GET请求..
```


> 参考:
>
> https://casbin.org/zh/docs/get-started
>
> https://dev.to/maxwellhertz/tutorial-integrate-gin-with-cabsin-56m0
>
> https://www.topgoer.com/gin%E6%A1%86%E6%9E%B6/%E5%85%B6%E4%BB%96/%E6%9D%83%E9%99%90%E7%AE%A1%E7%90%86.html