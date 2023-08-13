package util

import (
	"strings"

	"github.com/sirupsen/logrus"
)

func IsEmptyString(value string) bool {
	return len(strings.TrimSpace(value)) == 0
}

func IsErrorDoPanic(e error) {
	if e != nil {
		logrus.Panicln(e)
	}
}

func IsErrorDoPanicWithMessage(customMessage string, e error) {
	if e != nil {
		logrus.Panicln(customMessage, e)
	}
}

func IsErrorDoPrint(e error) {
	if e != nil {
		logrus.Errorln(e)
	}
}

func IsErrorDoPrintWithMessage(customMessage string, e error) {
	if e != nil {
		logrus.Errorln(customMessage, e)
	}
}

func GetAggregatorId(clientId, clientUser string) string {
	if IsEmptyString(clientId) {
		clientId = strings.Split(clientUser, SEMICOLON)[0]
	}
	return clientId
}
