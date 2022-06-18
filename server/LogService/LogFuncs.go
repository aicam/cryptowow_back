package LogService

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http/httputil"
)

func LogPotentialCyberAttack(c *gin.Context, happendPart string) {
	reqDump, err := httputil.DumpRequest(c.Request, true)
	if err != nil {
		reqDump = []byte("Problem in dumping request")
	}
	logrus.WithFields(logrus.Fields{
		"ip_address":   c.ClientIP(),
		"full_request": string(reqDump),
		"part":         happendPart,
	}).
		Info("Cyber Attack")
}

func LogSucceedJoinOperation(inviter, invited string) {
	logrus.WithFields(logrus.Fields{
		"inviter": inviter,
		"invited": invited,
	}).Info("Join Operation")
}

func LogCrash(service, part string) {
	logrus.WithFields(logrus.Fields{
		"Service": service,
		"Part":    part,
	}).Error("Crash in service")
}

func LogCrashinShop(service, part string) {
	logrus.WithFields(logrus.Fields{
		"Service": service,
		"Part":    part,
	})
}
