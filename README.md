## 简介

gin-x是采用golang编写的web框架，其实现处处体现了大名鼎鼎的Gin框架的设计思想。

## 快速入门

本示例将展示如何快速在本地启动一个gin-x服务。

要安装gin-x，需要先安装Go1.23或以上版本，并配置GOPATH工作空间。

1. 创建项目目录及main文件：

   ```shell
   $ mkdir myproject && cd myproject
   $ touch main.go
   ```

2. 初始化go mod：

   ```shell
   $ go mod init myproject
   ```

3. 执行以下命令安装gin-x：

   ```shell
   $ go get github.com/cauliflower-beep/gin-x@latest
   ```

4. 打开喜欢的文本编辑器，键入以下代码到main.go中：

   ```go
   package main
   
   import ginX "github.com/cauliflower-beep/gin-x"
   
   func main() {
   	r := ginX.New()
   	r.GET("/hello", func(ctx *ginX.Context) {
   		ctx.String(200, "hello gin-x!\n")
   	})
   	_ = r.Run(":9999")
   }
   ```

运行这段代码，就可以启动一个http服务了：

```shell
$ go run main.go
```

测试：

```shell
$ curl http://localhost:9999/hello
hello gin-x!
```

## API

### 创建应用

gin-x应用中，所有客户端请求都由一个`engine`实例接管。创建一个engine实例：

```go
package main

import ginX "github.com/cauliflower-beep/gin-x"

func main() {
    // 创建engine实例
	r := ginX.New()
	// 路由注册...
    // 启动服务
	_ = r.Run(":9999")
}
```

### 路由

gin-x定义了一组api，可基于`engine`实例创建路由映射规则：

```go
package main

import ginX "github.com/cauliflower-beep/gin-x"

func main() {
	r := ginX.New()
    
    // 定义路由及对应的handler
	r.GET("/hello", hello)
    r.POST("/upload",upload)
    
	_ = r.Run(":9999")
}
```

gin-x也支持定义分组路由。分组路由通常具备相同的路径前缀：

```go
package main

import ginX "github.com/cauliflower-beep/gin-x"

func main() {
	r := ginX.New()
    
    // group:v1
    v1 := r.Group("/v1")
    {
        v1.POST("/login", loginEndpoint)
		v1.POST("/submit", submitEndpoint)
    }
    
    // group:v2
    v2 := r.Group("/v2")
    {
        v2.POST("/login", loginEndpoint)
		v2.POST("/submit", submitEndpoint)
    }
    
	_ = r.Run(":9999")
}
```

### 路径参数

url路径中的参数可以通过`(*context).Param`获取：

```go
package main

import ginX "github.com/cauliflower-beep/gin-x"

func main() {
	r := ginX.New()
    
    // 定义路由及对应的handler
    r.GET("/hello/:name", func(ctx *ginX.Context) {
        name := ctx.Param("name")
		ctx.String(200, "hello, %s\n", name)
	})
    
    // 通配符匹配
    r.GET("/user/:name/*action", func(c *ginX.Context) {
		name := c.Param("name")
		action := c.Param("action")
		msg := name + " is " + action + "\n"
		c.String(200, msg)
	})
    
	_ = r.Run(":9999")
}
```

测试：

```shell
$ curl http://localhost:9999/hello/goku
hello, goku

$ curl http://localhost:9999/user/goku/fly
goku is fly
```

### 中间件

gin-x预置了两个中间件：`logger`和`recovery`，分别提供请求日志记录及故障恢复的能力，可作用于全局`engine`实例：

```go
package main

import ginX "github.com/cauliflower-beep/gin-x"

func main() {
	r := ginX.New()
    r.Use(Logger(), Recovery())
	...
	_ = r.Run(":9999")
}
```

同时，gin-x也支持自定义中间件并作用于特定的分组路由：

```go
package main

import (
	"fmt"
	ginX "github.com/cauliflower-beep/gin-x"
	"net/http"
)

// 只作用于v1分组的中间件
func print4V1() ginX.HandlerFunc {
	return func(c *ginX.Context) {
		fmt.Println("This is v1's middleware.")
	}
}

func main() {
	r := ginX.New()

	v1 := r.Group("/v1")
	// 将中间件作用于v1分组
	v1.Use(print4V1())
	{
		v1.GET("/hello", func(c *ginX.Context) {
			c.String(http.StatusOK, "hello, you're at %s\n", c.Request.URL.Path)
		})
	}

	v2 := r.Group("/v2")
	{
		v2.GET("/hello", func(c *ginX.Context) {
			c.String(http.StatusOK, "hello, you're at %s\n", c.Request.URL.Path)
		})
	}

	_ = r.Run(":9999")
}

```

测试：

```shell
$ curl http://localhost:9999/v1/hello
hello, you're at /v1/hello

// 服务器终端能看到中间件输出的内容：
$ go run main.go 
2025/01/08 17:00:35 add route:  GET - /v1/hello
2025/01/08 17:00:35 add route:  GET - /v2/hello
This is v1's middleware.

$ curl http://localhost:9999/v2/hello
hello, you're at /v2/hello
// 服务器终端没有任何内容输出
```

### 静态文件

gin-x支持对外提供静态文件资源：

```go
package main

import ginX "github.com/cauliflower-beep/gin-x"

func main() {
	r := ginX.New()

	// 设置静态资源路径
	r.Static("/assets", "./static")

	_ = r.Run(":9999")
}
```

测试：在当前目录创建`/static/chunxiao.txt`文件，启动服务执行如下请求：

```shell
$ curl http://localhost:9999/assets/chunxiao.txt
春眠不觉晓，处处闻啼鸟。
夜来风雨声，花落知多少。
```

### 模板Template

gin-x支持模板渲染：

```go
package main

import (
	"fmt"
	ginX "github.com/cauliflower-beep/gin-x"
	"html/template"
	"time"
)

// 自定义模板渲染函数，可在模板中调用
func FormatAsDate(t time.Time) string {
	year, month, day := t.Date()
	return fmt.Sprintf("%d-%02d-%02d", year, month, day)
}

func main() {
	r := ginX.New()

	// 注册模板渲染函数
	r.SetFuncMap(template.FuncMap{
		"FormatAsDate": FormatAsDate,
	})
	// 加载模板集合
	r.LoadHTMLGlob("templates/*")

	r.GET("/date", func(c *ginX.Context) {
		c.HTML(200, "date.tmpl", ginX.H{
			"title": "gin-x",
			"now":   time.Now(),
		})
	})

	_ = r.Run(":9999")
}
```

测试：在当前目录下创建`/templates/date.tmpl`

```html
<html>
<body>
    <p>hello, {{.title}}</p>
    <p>Date: {{.now | FormatAsDate}}</p>
</body>
</html>
```

发起请求：

```shell
curl http://localhost:9999/date
```

模板渲染响应：

```html
<html>
<body>
    <p>hello, gin-x</p>
    <p>Date: 2025-01-08</p>
</body>
</html>
```

### 错误恢复

框架提供了错误恢复的能力，避免因为一些可能的bug导致服务异常宕机。可以使用框架预置的`Recovery`中间件来引入这种能力：

```go
package main

import ginX "github.com/cauliflower-beep/gin-x"

func main() {
	r := ginX.New()

	// 应用Recovery中间件
	r.Use(ginX.Recovery())

	r.GET("/count", func(ctx *ginX.Context) {
		arr := []int{1, 2, 3}
		// 数组越界 panic
		ctx.String(200, "Count %d\n", arr[3])
	})

	_ = r.Run(":9999")
}
```

测试：

```shell
$ curl http://localhost:9999/count
{"message":"Internal Server Error"}
```

服务器终端日志：

```shell
$ go run main.go 
2025/01/08 17:09:25 runtime error: index out of range [3] with length 3
Traceback:
        /usr/local/go/src/runtime/panic.go:781
        /usr/local/go/src/runtime/panic.go:115
        /home/lsx01/goSpace/src/myproject/main.go:14
        /home/lsx01/goSpace/pkg/mod/github.com/cauliflower-beep/gin-x@v1.0.0/context.go:42
        /home/lsx01/goSpace/pkg/mod/github.com/cauliflower-beep/gin-x@v1.0.0/recovery.go:38
        /home/lsx01/goSpace/pkg/mod/github.com/cauliflower-beep/gin-x@v1.0.0/context.go:42
        /home/lsx01/goSpace/pkg/mod/github.com/cauliflower-beep/gin-x@v1.0.0/router.go:104
        /home/lsx01/goSpace/pkg/mod/github.com/cauliflower-beep/gin-x@v1.0.0/ginX.go:74
        /usr/local/go/src/net/http/server.go:3211
        /usr/local/go/src/net/http/server.go:2093
        /usr/local/go/src/runtime/asm_amd64.s:1701
```

应用不会宕机，可继续对外提供服务。