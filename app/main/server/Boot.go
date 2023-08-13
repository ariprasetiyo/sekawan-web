package server

import (
	"bytes"
	"context"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"runtime"
	"sekawan-web/app/main/util"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"google.golang.org/grpc/credentials"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	APP_NAME              string
	OTEL_COLLECTOR_URL    string
	IS_OTEL_INSECURE_MODE bool
	HTTP_SERVER_PORT      string
	GIN_MODE              string
)

const (
	logrusTimestampFormat string = "2006-01-02T15:04:05.999999999Z07:00"
)

func InitConfig() {
	err := godotenv.Load()
	util.IsErrorDoPanic(err)

	if os.Getenv(util.CONFIG_APP_ENV) == "development" {
		logrus.Infoln("load devel environment variable ")
	} else if os.Getenv(util.CONFIG_APP_ENV) == "staging" {
		logrus.Infoln("load staging environment variable ")
	} else if os.Getenv(util.CONFIG_APP_ENV) == "production" {
		logrus.Infoln("load production environment variable ")
	}

	IS_OTEL_INSECURE_MODE, err = strconv.ParseBool(os.Getenv(util.CONFIG_OTEL_INSECURE_MODE))
	util.IsErrorDoPanic(err)
	APP_NAME = os.Getenv(util.CONFIG_APP_NAME)
	OTEL_COLLECTOR_URL = os.Getenv(util.CONFIG_OTEL_EXPORTER_OTLP_ENDPOINT)
	HTTP_SERVER_PORT = os.Getenv(util.CONFIG_HTTP_SERVER_PORT)
	GIN_MODE = os.Getenv(util.CONFIG_GIN_MODE)
}

func InitLogrusFormat() {

	logrus.SetReportCaller(true)
	logrus.SetFormatter(&logrus.JSONFormatter{
		CallerPrettyfier: func(frame *runtime.Frame) (function string, file string) {
			fileName := path.Base(frame.File) + util.COLON + strconv.Itoa(frame.Line)
			//return frame.Function, fileName
			return util.EMPTY_STRING, fileName
		},
		TimestampFormat: logrusTimestampFormat,
	})
}

func HttpClient() *http.Client {
	httpClientLogger := HttpClientLogger{}
	client := &http.Client{Timeout: 10 * time.Second, Transport: httpClientLogger}
	return client
}

// running open telemetry
func InitTracer() func(context.Context) error {

	secureOption := otlptracegrpc.WithTLSCredentials(credentials.NewClientTLSFromCert(nil, ""))
	if IS_OTEL_INSECURE_MODE {
		secureOption = otlptracegrpc.WithInsecure()
	}

	exporter, err := otlptrace.New(
		context.Background(),
		otlptracegrpc.NewClient(
			secureOption,
			otlptracegrpc.WithEndpoint(OTEL_COLLECTOR_URL),
		),
	)

	if err != nil {
		logrus.Fatalln(err)
	}
	resources, err := resource.New(
		context.Background(),
		resource.WithAttributes(
			attribute.String("service.name", APP_NAME),
			attribute.String("library.language", "go"),
		),
	)
	if err != nil {
		logrus.Printf("Could not set resources: ", err)
	}

	otel.SetTracerProvider(
		sdktrace.NewTracerProvider(
			sdktrace.WithSampler(sdktrace.AlwaysSample()),
			sdktrace.WithBatcher(exporter),
			sdktrace.WithResource(resources),
		),
	)
	return exporter.Shutdown
}

func InitPostgreSQL() *PostgreSQLClientRepository {
	host := os.Getenv(util.CONFIG_DB_HOST)
	port, err := strconv.Atoi(os.Getenv(util.CONFIG_DB_PORT))
	util.IsErrorDoPanic(err)
	dbname := os.Getenv(util.CONFIG_DB_NAME)
	uname := os.Getenv(util.CONFIG_DB_USERNAME)
	pass := os.Getenv(util.CONFIG_DB_PASSWORD)

	var gConfig *gorm.Config
	isShowQuery, err := strconv.ParseBool(os.Getenv(util.CONFIG_DB_IS_SHOW_QUERY))
	util.IsErrorDoPanic(err)
	if isShowQuery {
		showQuery := logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags),
			logger.Config{
				LogLevel: logger.Info,
			})

		gConfig = &gorm.Config{Logger: showQuery, PrepareStmt: true}
	} else {
		gConfig = &gorm.Config{}
	}

	newDbConn, error := NewPostgreSQLRepository(host, uname, pass, dbname, port, gConfig)
	util.IsErrorDoPanic(error)
	return newDbConn

}

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func readBody(reader io.Reader) string {
	buf := new(bytes.Buffer)
	buf.ReadFrom(reader)

	s := buf.String()
	return s
}
