package response

import (
	"runtime"
	"strings"

	"github.com/labstack/echo"
	log "github.com/sirupsen/logrus"
)

type Response struct {
	Result interface{}            `json:"result"`
	Errors map[string]interface{} `json:"errors"`
}

func JSON(c echo.Context, httpStatus int, response Response) error {
	res := c.Response()

	file, line := getFileAndLineNumber()

	logFields := map[string]interface{}{
		"request_id": res.Header().Get(echo.HeaderXRequestID),
		"file":       file,
		"line":       line,
		"errors":     response.Errors,
	}

	log.WithFields(logFields).Error()

	return c.JSON(httpStatus, response)
}

func getFileAndLineNumber() (string, int) {
	_, file, line, ok := runtime.Caller(2)
	if !ok {
		file = "<???>"
		line = 1
	} else {
		slash := strings.LastIndex(file, "/")
		if slash >= 0 {
			file = file[slash+1:]
		}
	}

	return file, line
}
