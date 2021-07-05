package main

import (
	"github.com/gogf/gf/frame/g"
	_ "onvif-gf-demos/boot"
	_ "onvif-gf-demos/router"
)

// @title       `gf-demo`示例服务API
// @version     1.0
// @description `GoFrame`基础开发框架示例服务API接口文档。
// @schemes     http
func main() {
	g.Server().Run()
}
