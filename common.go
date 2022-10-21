package zont

import (
	"github.com/sirupsen/logrus"
)

var (
	ContextLogger = logrus.WithFields(logrus.Fields{
		"application_name": "go-zont",
	})
)
