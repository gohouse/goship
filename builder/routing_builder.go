package builder

import (
	"fmt"
	"github.com/gohouse/file"
	"github.com/gohouse/gorose/v2"
	"github.com/gohouse/goship/config"
	"github.com/gohouse/goship/templates"
	"github.com/gohouse/goship/util"
	"github.com/gohouse/schema"
	"html/template"
	"log"
)

type Routing struct {
	sm   *schema.Schema
	ge   *gorose.Engin
	conf *config.Config
}

func NewRouting(ge *gorose.Engin, c *config.Config) *Routing {
	return &Routing{sm: schema.NewSchema(ge), ge: ge, conf: c}
}

type TableList struct {
	TableName    string
	TableComment string
}
type Result struct {
	Tabs     []TableList
	GoModule string
}

func (r *Routing) Build() {
	// 获取表信息
	tableList := r.sm.TableList()

	var tbs []TableList
	for k, v := range tableList {
		tbs = append(tbs, TableList{
			TableName:    util.CamelCase(k),
			TableComment: v,
		})
	}

	// 读取模板
	//tmpl, err := r.ReadTemplateFile()
	tmpl := template.New("api_tmpl")
	_, err := tmpl.Parse(templates.GetRoutingTemplate())
	if err != nil {
		panic(err.Error())
	}

	var apiFile = fmt.Sprintf("%s/%s/api.go", r.conf.SiteInfo.RootDir, r.conf.Dir.Routing)
	//f := util.NewFile(apiFile)
	f := file.NewFile(apiFile).OpenFile()
	defer f.Close()

	var res = Result{
		Tabs:     tbs,
		GoModule: r.conf.SiteInfo.GoModule,
	}

	err = tmpl.Execute(f, res)

	if err != nil {
		log.Fatal(err.Error())
	}

	util.GofmtW(apiFile)
}
