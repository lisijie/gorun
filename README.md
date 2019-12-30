# gorun

gorun 是一个用于提高开发效率的Go程序热编译工具。

## 使用

你可以直接使用编译后的 gorun 命令或者以库的形式嵌入到项目中。

### 使用 gorun 命令

1. 编译本项目并将 gorun 程序放到你的 `PATH` 环境变量相关目录下。

2. 在项目下执行 gorun init 命令，修改生成 gorun.toml 配置文件，将 `build_cmd` 改成编译项目的命令，将 `run_cmd` 改成运行项目的命令。

3. 在项目目录下使用 gorun 命令运行项目。

### 嵌入到项目中

新建一个入口文件，例如名为 `cmd/dev/main.go`，内容如下：

```
package main

import (
	"github.com/lisijie/gorun/gorun"
	"log"
)

func main() {
	cfg := &gorun.Config{
		AppName:          "app",
		BuildCommand:     "go build -o app cmd/app/main.go",
		RunCommand:       "./app",
		WatchExtensions:  ".go",
		WatchExcludeDirs: "vendor,resource",
	}
	app := gorun.New(cfg)
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
```

执行：

```
go run cmd/dev/main.go 
```