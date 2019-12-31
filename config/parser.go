package config

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/gohouse/gorose/v2"
	"io/ioutil"
)

type DirConf struct {
	Routing string `toml:"routing"`
	Model   string `toml:"model"`
	Api     string `toml:"api"`
}
type JwtConf struct {
	Switch bool   `toml:"switch"`
	Secret string `toml:"secret"`
	Expire int64  `toml:"expire"`
}
type ErrorConf struct {
	Switch          bool  `toml:"switch"`
	ErrorStackLayer int64 `toml:"error_stack_layer"`
}
type Http struct {
	Port       string `toml:"port"`
	ApiVersion string `toml:"api_version"`
}
type SiteInfo struct {
	RootDir     string `toml:"root_dir"`
	ProjectName string `toml:"project_name"`
	TestToken   string `toml:"test_token"`
	GoModule    string `toml:"go_module"`
}
type Config struct {
	SiteInfo SiteInfo `toml:"site_info"`
	//Http     Http                  `toml:"http"`
	Join     map[string][][]string `toml:"join"`
	Database gorose.Config         `toml:"database"`
	//Jwt      JwtConf               `toml:"jwt"`
	//Error    ErrorConf             `toml:"error"`
	Dir DirConf `toml:"dir"`
}

func ParseConfig(file string, c *Config) {
	// 解析配置
	bytes, err := ioutil.ReadFile(file)
	if err != nil {
		panic(err.Error())
	}
	err = toml.Unmarshal(bytes, c)
	// 拼接项目根目录
	c.SiteInfo.RootDir = fmt.Sprintf("%s/%s", c.SiteInfo.RootDir, c.SiteInfo.ProjectName)
	if err != nil {
		panic(err.Error())
	}
	// 默认路径
	c.Dir = DirConf{
		Routing: "routing",
		Model:   "model",
		Api:     "api",
	}
}
