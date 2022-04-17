package LogService

import "github.com/sirupsen/logrus"

func GlobalLog() logrus.Fields {
	return logrus.Fields{"test": 2, "test_s": "dsada"}
}
