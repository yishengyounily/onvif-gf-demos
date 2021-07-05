package router

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"onvif-gf-demos/app/api"
)

// 你可以将路由注册放到一个文件中管理，
// 也可以按照模块拆分到不同的文件中管理，
// 但统一都放到router目录下。
func init() {
	s := g.Server()
	// 分组路由注册方式
	s.Group("/", func(group *ghttp.RouterGroup) {
		group.Group("/", func(group *ghttp.RouterGroup) {
			group.GET("/api/discovery", api.ONVIF.Discovery)
			group.POST("/api/:service/:method", api.ONVIF.PostMethod)
		})
	})
}