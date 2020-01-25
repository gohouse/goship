package templates

import "io/ioutil"

// getXXXXFromFile 意思是读取模板文件, 方便开发调试用
// getXXXXFromRaw 把模板内容放入字符串, 方便打包可执行文件使用

func GetRoutingTemplate() string {
	return getRoutingTemplateFromRaw()
}

func GetApiTemplate() string {
	return getApiTemplateFromRaw()
}

func GetConfigTemplate() []byte {
	return getConfigTemplateFromRaw()
}



func getRoutingTemplateFromFile() string {
	file, err := ioutil.ReadFile("templates/routing.tmpl")
	if err != nil {
		panic(err.Error())
	}
	return string(file)
}
func getRoutingTemplateFromRaw() string {
	return `package routing

import (
	"github.com/gin-gonic/gin"
	"{{.GoModule}}/api"
)

func ApiRun(route *gin.RouterGroup)  {
{{range .Tabs}}
    // {{.TableComment}} - 列表
    route.GET("/{{.TableName}}",v1(r(api.{{.TableName}}List)))
    // {{.TableComment}} - 详情
    route.GET("/{{.TableName}}/:pkid",v1(r(api.{{.TableName}}Info)))
    // {{.TableComment}} - 删除
    route.DELETE("/{{.TableName}}/:pkid",v1(r(api.{{.TableName}}Delete)))
    // {{.TableComment}} - 修改
    route.PUT("/{{.TableName}}/:pkid",v1(r(api.{{.TableName}}Edit)))
    // {{.TableComment}} - 新增
    route.POST("/{{.TableName}}",v1(r(api.{{.TableName}}Add)))
{{end}}
}`
}


func getApiTemplateFromFile() string {
	file, err := ioutil.ReadFile("templates/api.tmpl")
	if err != nil {
		panic(err.Error())
	}
	return string(file)
}
func getApiTemplateFromRaw() string {
	return `package api

import (
	"github.com/gin-gonic/gin"
	"github.com/gohouse/gorose/v2"
	"github.com/gohouse/t"
	"{{.GoModule}}/helper"
	"{{.GoModule}}/model"
	"strings"
)

// {{.TableName}}List	godoc
// @Summary 	{{.TableComment}} - 获取列表
// @Description {{.TableComment}} - 获取列表
// @Tags 		{{.TableName}}
// @Accept  	json
// @Produce		json
// @Param		Authorization header string true "header中放入token" default({{.Token}})
{{- range .Fields}}
// @Param		{{.ColumnName}} query string false "{{.ColumnName}}: {{.ColumnComment}}" default()
{{- end}}
// @Success 	200 {object} helper.GinBack 成功时自定义返回
// @Failure 	400 {object} helper.GinBack 失败时自定义返回
// @Failure 	401 {object} helper.GinBack 认证失败
// @Failure 	404 {string} string http返回，找不到资源
// @Failure 	500 {string} string http返回，服务器错误
// @Router    /{{.TableName}}List [get]
func {{.TableName}}List(c *gin.Context) helper.ApiReturn {
	// build where
	where, _ := helper.BuildWhereMulti(c, []string{
	    {{range .Fields}}
		    "{{.ColumnName}}", {{end}}
	})

	// build paginate
	limit, page := helper.BuildPageParams(c)

	// query data
	res, err := DB().Table(model.{{.TableName}}{}).Where(where).Limit(limit).Page(page).Paginate()

	return helper.QueryReturn(res, err)
}

// {{.TableName}}Info	godoc
// @Summary 	{{.TableComment}} - 获取详情
// @Description {{.TableComment}} - 获取详情
// @Tags 		{{.TableName}}
// @Accept  	json
// @Produce		json
// @Param		Authorization header string true "header中放入token" default({{.Token}})
// @Param		pkid path string true "主键" default()
// @Success 	200 {object} helper.GinBack 成功时自定义返回
// @Failure 	400 {object} helper.GinBack 失败时自定义返回
// @Failure 	401 {object} helper.GinBack 认证失败
// @Failure 	404 {string} string http返回，找不到资源
// @Failure 	500 {string} string http返回，服务器错误
// @Router    /{{.TableName}}Info [get]
func {{.TableName}}Info(c *gin.Context) helper.ApiReturn {
	// query row
	res, err := DB().Table(model.{{.TableName}}{}).Where(gorose.Data{"{{.PkidName}}": c.Param("pkid")}).First()

	return helper.QueryReturn(res, err)
}

// {{.TableName}}Delete	godoc
// @Summary 	{{.TableComment}} - 删除
// @Description {{.TableComment}} - 删除
// @Tags 		{{.TableName}}
// @Accept  	json
// @Produce		json
// @Param		Authorization header string true "header中放入token" default({{.Token}})
// @Param		pkid path string true "主键" default()
// @Success 	200 {object} helper.GinBack 成功时自定义返回
// @Failure 	400 {object} helper.GinBack 失败时自定义返回
// @Failure 	401 {object} helper.GinBack 认证失败
// @Failure 	404 {string} string http返回，找不到资源
// @Failure 	500 {string} string http返回，服务器错误
// @Router    /{{.TableName}}Delete [delete]
func {{.TableName}}Delete(c *gin.Context) helper.ApiReturn {
	// parse pkid
	var pkids = c.Param("pkid")
	var pkid = t.New(strings.Split(pkids, ",")).SliceInterface()

	// delete
	aff, err := DB().Table(model.{{.TableName}}{}).WhereIn("{{.PkidName}}", pkid).Delete()

	return helper.ExecReturn(aff, err)
}

// {{.TableName}}Edit	godoc
// @Summary 	{{.TableComment}} - 编辑
// @Description {{.TableComment}} - 编辑
// @Tags 		{{.TableName}}
// @Accept  	multipart/form-data
// @Produce		json
// @Param		Authorization header string true "header中放入token" default({{.Token}})
{{- range .Fields}}
// @Param		{{.ColumnName}} formData string false "{{.ColumnName}}: {{.ColumnComment}}" default()
{{- end}}
// @Success 	200 {object} helper.GinBack 成功时自定义返回
// @Failure 	400 {object} helper.GinBack 失败时自定义返回
// @Failure 	401 {object} helper.GinBack 认证失败
// @Failure 	404 {string} string http返回，找不到资源
// @Failure 	500 {string} string http返回，服务器错误
// @Router    /{{.TableName}}Edit [put]
func {{.TableName}}Edit(c *gin.Context) helper.ApiReturn {
	// build data
	data, i := helper.BuildWhere(c, []string{
		{{- range .Fields}}
			"{{.ColumnName}}",
        {{- end}}
	})
	if i == 0 {
		return helper.FailReturn("params needed")
	}

	// build where
	where := gorose.Data{"{{.PkidName}}": c.Param("pkid")}

	// insert
	aff, err := DB().Table(model.{{.TableName}}{}).Where(where).Update(data)

	return helper.ExecReturn(aff, err)
}

// {{.TableName}}Edit	godoc
// @Summary 	{{.TableComment}} - 新增
// @Description {{.TableComment}} - 新增
// @Tags 		{{.TableName}}
// @Accept  	multipart/form-data
// @Produce		json
// @Param		Authorization header string true "header中放入token" default({{.Token}})
{{- range .Fields}}
// @Param		{{.ColumnName}} formData string false "{{.ColumnName}}: {{.ColumnComment}}" default()
{{- end}}
// @Success 	200 {object} helper.GinBack 成功时自定义返回
// @Failure 	400 {object} helper.GinBack 失败时自定义返回
// @Failure 	401 {object} helper.GinBack 认证失败
// @Failure 	404 {string} string http返回，找不到资源
// @Failure 	500 {string} string http返回，服务器错误
// @Router    /{{.TableName}}Edit [post]
func {{.TableName}}Add(c *gin.Context) helper.ApiReturn {
	// build data
	data, i := helper.BuildWhere(c, []string{
		{{- range .Fields}}
			"{{.ColumnName}}",
        {{- end}}
	})
	if i == 0 {
		return helper.FailReturn("params needed")
	}

	// insert
	aff, err := DB().Table(model.{{.TableName}}{}).Insert(data)

	return helper.ExecReturn(aff, err)
}
`
}

func getConfigTemplateFromFile() []byte {
	file, err := ioutil.ReadFile("templates/config.toml.tmpl")
	if err != nil {
		panic(err.Error())
	}
	return file
}
func getConfigTemplateFromRaw() []byte {
	return []byte(`############## 基本信息 #############
[site_info]
    root_dir = "./"   # 工作目录
    project_name = "goship-demo"    # 项目目录名字
    go_module = "goshipdemo"    # module名字
    test_token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.xxxxxx.xxx"    # api 测试使用的 jwt token
    goship_template = "https://github.com/gohouse/goship-template.git" # 框架模板

############## 数据库配置 #############
[database]
    driver = "mysql" # 数据库驱动
    dsn = "root:123456@tcp(localhost:3306)/goship?charset=utf8mb4" # dsn链接
    setMaxOpenConns = 300 # 连接池 - 最大打开连接数
    setMaxIdleConns = 50  # 连接池 - 最大空闲连接数
    prefix = "pre_"  # 表前缀

############## join ############# todo
[join]
    # 获取用户角色信息,默认 inner join
    user_role = [
        ["user as a"],
        ["role as b", "a.role_name", "b.role_name"]
    ]
    # 获取用户角色信息,指定第二个参数"left", 则为 left join
    user_role_userinfo = [
        ["user as a"],
        ["role as b",       "a.role_name",  "b.role_name"],
        ["userinfo as c",   "a.id",         "c.user_id"]
    ]
`)
}