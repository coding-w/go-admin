package middleware

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"go-admin/global"
	"go-admin/model/system"
	"go-admin/service"
	"go-admin/utils"
	"go.uber.org/zap"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"
)

var operationRecordService = service.ServiceGroup.OperationRecordService

var respPool sync.Pool
var bufferSize = 1024

func init() {
	respPool.New = func() interface{} {
		return make([]byte, bufferSize)
	}
}

func OperationRecord() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 临时存储请求体
		var bodyBuffer bytes.Buffer
		// 通过 io.TeeReader 避免重复读取请求体
		bodyReader := io.TeeReader(c.Request.Body, &bodyBuffer)
		// 确保请求体可以被多次读取
		c.Request.Body = io.NopCloser(bodyReader)
		// 自定义函数，用于从 bodyBuffer 中解析请求体内容
		body, err := parseRequestBody(c, &bodyBuffer)
		if err != nil {
			global.GA_LOG.Error("read body from request error:", zap.Error(err))
		}
		userId := getUserID(c)
		record := system.SysOperationRecord{
			Ip:     c.ClientIP(),
			Method: c.Request.Method,
			Path:   c.Request.URL.Path,
			Agent:  c.Request.UserAgent(),
			Body:   maskSensitiveData(determineRequestBody(c, body)),
			UserID: userId,
		}
		// 拦截响应并记录响应体
		writer := responseBodyWriter{
			ResponseWriter: c.Writer,
			body:           &bytes.Buffer{},
		}
		c.Writer = writer
		// 当前时间，作为请求开始的时间
		now := time.Now()
		// 调用 c.Next() 继续处理后续中间件和最终的处理函数
		c.Next()
		// 计算请求处理的总耗时，记录延迟
		latency := time.Since(now)
		record.ErrorMessage = c.Errors.ByType(gin.ErrorTypePrivate).String()
		record.Status = c.Writer.Status()
		record.Latency = latency
		record.Resp = determineResponseBody(c, writer.body.String())

		if err := operationRecordService.CreateSysOperationRecord(record); err != nil {
			global.GA_LOG.Error("create operation record error:", zap.Error(err))
		}
	}
}

// 获取用户ID的逻辑
func getUserID(c *gin.Context) int {
	claims, _ := utils.GetClaims(c)
	if claims != nil && claims.BaseClaims.ID != 0 {
		return int(claims.BaseClaims.ID)
	}
	id, err := strconv.Atoi(c.Request.Header.Get("x-user-id"))
	if err != nil {
		return 0
	}
	return id
}

// 解析请求体内容
func parseRequestBody(c *gin.Context, bodyBuffer *bytes.Buffer) ([]byte, error) {
	if c.Request.Method == http.MethodGet {
		query := c.Request.URL.RawQuery
		query, _ = url.QueryUnescape(query)
		m := parseQueryToMap(query)
		return json.Marshal(&m)
	}
	return io.ReadAll(bodyBuffer)
}

// 将查询字符串解析为键值对
func parseQueryToMap(query string) map[string]string {
	m := make(map[string]string)
	split := strings.Split(query, "&")
	for _, v := range split {
		kv := strings.Split(v, "=")
		if len(kv) == 2 {
			m[kv[0]] = kv[1]
		}
	}
	return m
}

// 判断请求体内容，不记录较大的请求体和上传的文件
func determineRequestBody(c *gin.Context, body []byte) string {
	if strings.Contains(c.GetHeader("Content-Type"), "multipart/form-data") {
		return "[文件]"
	}
	if len(body) > bufferSize {
		return "[超出记录长度]"
	}
	return string(body)
}

// 处理敏感数据屏蔽
func maskSensitiveData(data string) string {
	// 示例：屏蔽可能的密码字段内容
	return strings.ReplaceAll(data, "password", "****")
}

// 判断响应内容，不记录较大的请求体和下载的文件
func determineResponseBody(c *gin.Context, responseBody string) string {
	if isLargeFileResponse(c) {
		return "[文件下载]"
	}
	if len(responseBody) > bufferSize {
		return "[超出记录长度]"
	}
	return responseBody
}

// 判断是否为下载文件
func isLargeFileResponse(c *gin.Context) bool {
	contentType := c.Writer.Header().Get("Content-Type")
	return strings.Contains(contentType, "application/force-download") ||
		strings.Contains(contentType, "application/octet-stream") ||
		strings.Contains(contentType, "application/vnd.ms-excel") ||
		strings.Contains(contentType, "application/download") ||
		strings.Contains(contentType, "attachment")
}

// 自定义响应写入器，用于捕获响应内容
type responseBodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (r responseBodyWriter) Write(b []byte) (int, error) {
	r.body.Write(b)
	return r.ResponseWriter.Write(b)
}
