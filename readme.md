## 关于 goship
goship是一个根据表结构自动生成 restful api 项目框架.  


## 依赖
- golang   
- gin-swagger  

### golang 安装
golang 官网: https://golang.org/  
如果没有开启 go module,则需要开启,可以通过 `go env` 查看`GO111MODULE`开启情况  

开启命令
```shell script
export GO111MODULE=on
```

### 安装 goship
`goship`官网: https://github.com/gohouse/goship    

简洁安装
```shell script
go get -u github.com/gohouse/goship/cmd/goship
```
若 `$GOPATH/bin` 没有加入$PATH中，你需要执行将其可执行文件移动到$GOBIN下
```shell script
mv $GOPATH/bin/goship /usr/local/go/bin
```

### 安装 swagger
linux/unix 会自动尝试安装, Windows下请手动安装,手动安装请参考  
`gin-swagger`官网: https://github.com/swaggo/gin-swagger  

linux/unix下手动简洁安装
```shell script
go get -u github.com/swaggo/swag/cmd/swag
```
若 `$GOPATH/bin` 没有加入$PATH中，你需要执行将其可执行文件移动到$GOBIN下
```shell script
mv $GOPATH/bin/swag /usr/local/go/bin
```

## 运行
### 1. 生成项目
```shell script
# 如果没有配置文件,可以使用如下命令导出配置模板,config.toml为导出文件名,可任意指定
goship -e config.toml

# 修改导出的 config.toml 配置,然后使用 -f 参数指定配置文件运行  
goship -f config.toml
# 或者,如果不指定 -f,则默认读取当前目录的 config.toml
goship
```
> 也可以克隆源码,手动在根目录 `/path/to/goship`执行 `go run main.go`来生成项目

### 2. 运行项目
```shell script
# 进入生成的项目,如默认执行goship后在当前目录生成一个 goship-demo 的项目
cd goship-demo
# 可以编辑当前目录的配置文件,然后运行,如果不做修改,默认是8088端口
go run mian.go
# 或者指定配置文件
go run main.go -f config.toml
```

