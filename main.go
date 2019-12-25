package main

import (
	"flag"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/lisijie/gorun/gorun"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var configFile = "gorun.toml"

var configTpl = `app_name = "{app_name}"
app_path = "./"
watch_exclude_dirs = ""
watch_extensions = ".go,.toml,.ini,.yml"
build_cmd = "go build -o gorun_app"
run_cmd = "./gorun_app"
`

var debug bool

func main() {
	flag.Parse()

	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "init":
			var appName string
			if len(os.Args) > 2 {
				appName = os.Args[2]
			} else {
				wd, _ := os.Getwd()
				appName = filepath.Base(wd)
			}
			createConfigFile(appName)
			return
		case "debug":
			debug = true
		default:
			fmt.Println("unknown command")
			return
		}
	}

	conf := &gorun.Config{}
	if _, err := toml.DecodeFile(configFile, conf); err != nil {
		log.Fatalln(err)
	}
	app := gorun.New(conf)
	if debug {
		app.SetLogger(gorun.StdLogger)
	}
	if err := app.Run(); err != nil {
		log.Fatalln(err)
	}
}

func createConfigFile(name string) {
	if isFile(configFile) {
		fmt.Println(configFile, "already exists")
		return
	}
	data := strings.Replace(configTpl, "{app_name}", name, 1)
	if err := ioutil.WriteFile(configFile, []byte(data), 0644); err != nil {
		fmt.Println(err)
	}
}

func isFile(filename string) bool {
	info, err := os.Stat(filename)
	if err == nil && !info.IsDir() {
		return true
	}
	return false
}
