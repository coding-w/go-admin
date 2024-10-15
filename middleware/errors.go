package middleware

import (
	"github.com/gin-gonic/gin"
	"go-admin/global"
	"go-admin/model/common/response"
	"go.uber.org/zap"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime/debug"
	"strings"
)

func GinRecovery(stack bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		defer func() {
			if err := recover(); err != nil {
				var brokenPie bool
				// 判断是否为网络错误
				if ne, ok := err.(*net.OpError); ok {
					// 判断是否为系统调用错误
					if se, ok := ne.Err.(*os.SyscallError); ok {
						// 检查错误信息中是否包含 "broken pipe" 或 "connection reset by peer"
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") ||
							strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPie = true
						}
					}
				}
				// 请求的详细信息
				request, _ := httputil.DumpRequest(c.Request, true)
				// 判断是否为网络错误
				if brokenPie {
					// 记录错误日志
					global.GA_LOG.Error("broken pipe or connection reset by peer",
						zap.Any("error", err),
						zap.String("request", string(request)),
					)
					_ = c.Error(err.(error))
					c.Abort()
					return
				}
				// 判断是否开启堆栈跟踪
				if stack {
					global.GA_LOG.Error("recover from panic with stack trace",
						zap.Any("error", err),
						zap.String("request", string(request)),
						zap.String("stack", string(debug.Stack())),
					)
				} else {
					global.GA_LOG.Error("recover from panic",
						zap.Any("error", err),
						zap.String("request", string(request)),
					)
				}
				response.Fail(c)
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()
	}
}
