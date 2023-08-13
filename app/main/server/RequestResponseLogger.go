package server

import (
	"bytes"
	"io/ioutil"
	"sekawan-web/app/main/util"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func RequestResponseLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.URL.Path == util.URI_HEALTH_CHECK {
			return
		}

		startTime := time.Now().UnixMilli()
		buf, _ := ioutil.ReadAll(c.Request.Body)
		rdr1 := ioutil.NopCloser(bytes.NewBuffer(buf))
		rdr2 := ioutil.NopCloser(bytes.NewBuffer(buf)) //We have to create a new Buffer, because rdr1 will be read.
		logrus.Infoln("request: http method", c.Request.Method, ", client ip", c.ClientIP(), ", user-agent", c.Request.UserAgent(), c.Request.URL, ", body=", (rdr1))
		c.Request.Body = rdr2

		blw := &bodyLogWriter{body: bytes.NewBufferString(util.EMPTY_STRING), ResponseWriter: c.Writer}
		c.Writer = blw
		c.Next()

		statusCode := c.Writer.Status()
		duration := time.Now().UnixMilli() - startTime
		logrus.Infoln("Response: http status", statusCode, ", duration:", duration, "ms ", c.Request.URL, " , body=", blw.body.String())
	}
}
