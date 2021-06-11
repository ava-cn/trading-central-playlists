package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"net/http"
	"strings"
)

// CorsMiddleware 跨域中间件
func CorsMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var (
			currentOrigin       string
			allowOriginsFormYml = viper.GetString("server.corsWebsites") // 从配置中读取域名
			allowOrigin         string
		)
		if strings.Contains(allowOriginsFormYml, "*") {
			allowOrigin = "*"
		} else {
			currentOrigin = ctx.GetHeader("Origin") // 当前请求的域
			if strings.Contains(allowOriginsFormYml, currentOrigin) {
				allowOrigin = currentOrigin
			}
		}
		ctx.Writer.Header().Set("Access-Control-Allow-Origin", allowOrigin)    // 设置允许跨域的域名
		ctx.Writer.Header().Set("Access-Control-Max-Age", "86400")             // 缓存时间
		ctx.Writer.Header().Set("Access-Control-Allow-Methods", "GET,OPTIONS") // 允许请求的方法
		ctx.Writer.Header().Set("Access-Control-Allow-Headers", "*")           // 请求的头
		ctx.Writer.Header().Set("Access-Control-Allow-Credentials", "true")    // 是否允许https

		if ctx.Request.Method == http.MethodOptions {
			ctx.AbortWithStatus(http.StatusOK)
		} else {
			ctx.Next()
		}
	}
}
