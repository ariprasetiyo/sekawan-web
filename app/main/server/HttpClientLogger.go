package server

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"time"

	"github.com/sirupsen/logrus"
)

func NewHttpClientLogger(isEnable bool, isEnableAllRequestHeaderAndBody bool) *HttpClientLogger {
	return &HttpClientLogger{isEnableLog: isEnable, isEnableAllRequestHeaderAndBody: isEnableAllRequestHeaderAndBody}
}

type HttpClientLogger struct {
	isEnableLog                     bool
	isEnableAllRequestHeaderAndBody bool
}

/*
isEnableLog value true meaning print http client request & response
*/
func (t HttpClientLogger) RoundTrip(req *http.Request) (*http.Response, error) {
	if t.isEnableLog && t.isEnableAllRequestHeaderAndBody {
		return isLoggingHeaderAndBody(t, req)
	} else if t.isEnableLog {
		return isLoggingBody(t, req)
	} else {
		resp, err := http.DefaultTransport.RoundTrip(req)
		return resp, err
	}
}

func isLoggingBody(t HttpClientLogger, req *http.Request) (*http.Response, error) {
	requestTime := time.Now().UnixMilli()
	requestBody, err := ioutil.ReadAll(req.Body)
	if err != nil {
		logrus.Fatal(err)
	}
	req.Body = ioutil.NopCloser(bytes.NewBuffer(requestBody))
	logrus.Infoln("http request", req.Method, req.URL, "body=", string(requestBody))

	resp, err := http.DefaultTransport.RoundTrip(req)
	if err != nil {
		return resp, err
	}
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.Fatal(err)
	}
	resp.Body = ioutil.NopCloser(bytes.NewBuffer(responseBody))
	latency := time.Now().UnixMilli() - requestTime
	logrus.Infoln("http response status", resp.StatusCode, ", duration", latency, "ms,", resp.Request.URL, " , body=", string(responseBody))

	return resp, err
}

func isLoggingHeaderAndBody(t HttpClientLogger, req *http.Request) (*http.Response, error) {

	reqDump, err := httputil.DumpRequestOut(req, true)
	if err != nil {
		logrus.Fatal(err)
	}
	logrus.Infoln("http request", req.Method, req.URL, "body=", string(reqDump))

	resp, err := http.DefaultTransport.RoundTrip(req)
	if err != nil {
		return resp, err
	}
	resDump, errDump := httputil.DumpResponse(resp, true)
	if errDump != nil {
		logrus.Fatal(err)
	}
	logrus.Infoln("http response", resp.StatusCode, resp.Request.URL, "body=", string(resDump))

	return resp, err
}
