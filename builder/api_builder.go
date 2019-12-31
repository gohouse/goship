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

type Api struct {
	sm   *schema.Schema
	ge   *gorose.Engin
	conf *config.Config
}

func NewApi(ge *gorose.Engin, c *config.Config) *Api {
	return &Api{sm: schema.NewSchema(ge), ge: ge, conf: c}
}

type ApiResult struct {
	Fields       []schema.TableColumn
	GoModule     string
	TableName    string
	Token        string
	PkidName     string
	TableComment string
}

func (r *Api) Build() {
	// 获取表信息
	tableList := r.sm.TableList()
	for k, comment := range tableList {
		tabName := util.CamelCase(k)
		var apiRes = ApiResult{
			Fields:       r.sm.TableColumnList(k),
			GoModule:     r.conf.SiteInfo.GoModule,
			TableName:    tabName,
			Token:        r.conf.SiteInfo.TestToken,
			TableComment: comment,
			PkidName:     r.sm.TablePkidName(k),
		}
		// 读取模板
		//tmpl, err := template.ParseFiles("templates/api.tmpl")
		tmpl := template.New("api_tmpl")
		_, err := tmpl.Parse(templates.GetApiTemplate())
		if err != nil {
			panic(err.Error())
		}

		var apiFile = fmt.Sprintf("%s/%s/%s.go", r.conf.SiteInfo.RootDir, r.conf.Dir.Api, k)
		//f := util.NewFile(apiFile)
		f := file.NewFile(apiFile).OpenFile()
		defer f.Close()

		err = tmpl.Execute(f, apiRes)

		if err != nil {
			log.Fatal(err.Error())
		}

		util.GofmtW(apiFile)
	}
}
