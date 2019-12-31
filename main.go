package main

import (
	"flag"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gohouse/converter"
	"github.com/gohouse/file"
	"github.com/gohouse/gorose/v2"
	"github.com/gohouse/goship/builder"
	"github.com/gohouse/goship/config"
	"github.com/gohouse/goship/templates"
	"github.com/gohouse/goship/util"
	"html/template"
	"log"
	"os"
)

var f string
var e string

func init() {
	flag.StringVar(&f, "f", "", "file:运行需要指定的配置文件")
	flag.StringVar(&e, "e", "", "export:导出配置文件模板")
	flag.Parse()
	if f == "" {
		f = "./config.toml"
	}
}
func main() {
	// 检查导出配置
	if e != "" {
		checkExport(e)
		return
	}

	// 检查依赖
	if !util.CheckCommandExists("go") {
		log.Println("未检测到go,请先安装golang")
		return
	}
	if !util.CheckCommandExists("swag") {
		log.Println("未检测到swag,请先安装gin-swagger")
		return
	}

	// 解析配置
	var c config.Config
	config.ParseConfig(f, &c)

	// 驱动数据库
	engin, err := gorose.Open(&c.Database)
	if err != nil {
		panic(err.Error())
	}

	// 初始化项目
	projectInit(&c)

	// 替换 module
	replaceModulePlaceholder(&c)

	// 生成路由
	builder.NewRouting(engin, &c).Build()

	// 生成api
	builder.NewApi(engin, &c).Build()

	// 生成 model
	genModel(engin, &c)

	// 写入数据库配置到配置文件
	addDbConf(&c)

	// 生成 swagger api 文档
	if err = genSwag(&c); err != nil {
		log.Println("生成swagger api失败:", err.Error())
	}

	// 生成 go.mod
	if err = genGoMod(&c); err != nil {
		log.Println("生成go.mod失败:", err.Error())
	}
	log.Println("finish")
}

func checkExport(e string) {
	_, err := file.NewFile(e).Write(templates.GetConfigTemplate())
	if err != nil {
		panic(err.Error())
	}
}

func projectInit(c *config.Config) {
	util.RunCmd(fmt.Sprintf("rm -rf %s", c.SiteInfo.RootDir))
	//util.RunCmd(fmt.Sprintf("cp -r ~/go/src/github.com/gohouse/goship-template %s", c.SiteInfo.RootDir))
	util.RunCmd(fmt.Sprintf("git clone https://github.com/gohouse/goship-template.git %s", c.SiteInfo.RootDir))
}
func replaceModulePlaceholder(c *config.Config) {
	files, err := file.GetAllFiles(c.SiteInfo.RootDir)
	if err != nil {
		panic(err.Error())
	}
	for _, f := range files {
		file.ReplaceFileContent(f, "goshipdemo", c.SiteInfo.GoModule)
	}
}
func genGoMod(c *config.Config) error {
	if c.SiteInfo.GoModule == "goshipdemo" {
		return util.RunCmds([]string{
			fmt.Sprintf("cd %s", c.SiteInfo.RootDir),
			"go mod tidy",
		})
	}
	var gomodfile = fmt.Sprintf("%s/go.mod", c.SiteInfo.RootDir)
	if err := file.ReplaceFileContent(gomodfile, "goshipdemo", c.SiteInfo.GoModule); err != nil {
		return err
	}

	return util.RunCmds([]string{
		fmt.Sprintf("cd %s", c.SiteInfo.RootDir),
		"go mod tidy",
	})
}
func genModel(engin *gorose.Engin, c *config.Config) {
	err := converter.NewTable2Struct().
		Config(&converter.T2tConfig{StructNameToHump: true}).
		SavePath(fmt.Sprintf("%s/%s/%s", c.SiteInfo.RootDir, c.Dir.Model, "model.go")).
		DB(engin.GetQueryDB()).
		TagKey("gorose").
		RealNameMethod("TableName").
		EnableJsonTag(true).
		Prefix(engin.GetPrefix()).
		Run()
	if err != nil {
		panic(err.Error())
	}
}

func addDbConf(c *config.Config) {
	// 初始化配置文件
	var conffile_example = fmt.Sprintf("%s/config.toml.example", c.SiteInfo.RootDir)
	var conffile = fmt.Sprintf("%s/config.toml", c.SiteInfo.RootDir)
	err := util.RunCmd(fmt.Sprintf("cp %s %s", conffile_example, conffile))
	if err != nil {
		panic(err)
	}

	// 替换数据库配置
	tmpl := template.New("conf")
	_, err = tmpl.Parse(`
# 数据库配置
[database]
driver = "{{.Driver}}" # 数据库驱动
dsn = "{{.Dsn}}" # dsn链接
setMaxOpenConns = {{.SetMaxOpenConns}} # 连接池 - 最大打开连接数
setMaxIdleConns = {{.SetMaxIdleConns}}  # 连接池 - 最大空闲连接数
prefix = "{{.Prefix}}"  # 表前缀`)
	if err != nil {
		panic(err.Error())
	}
	err = tmpl.Execute(util.NewFileWithMod(conffile, os.O_APPEND|os.O_CREATE|os.O_WRONLY), c.Database)
	if err != nil {
		panic(err.Error())
	}
}
func genSwag(c *config.Config) error {
	return util.RunCmds([]string{
		fmt.Sprintf("cd %s", c.SiteInfo.RootDir),
		"swag init",
	})
}
