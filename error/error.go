package error

import (
	"net/http"
	"net/url"
	"runtime"
	"strings"

	"github.com/labstack/echo"
	log "github.com/sirupsen/logrus"
)

type AppError struct {
	Status       string `json:"status"`
	FieldName    string `json:"field_name"`
	ErrorCode    string `json:"error_code"`
	ErrorMessage string `json:"error_message"`
}

func (e AppError) Error() string {
	return e.ErrorMessage
}

func New(message string) error {
	return AppError{
		ErrorMessage: message,
	}
}

func NewWithArgs(message string, args ...string) error {
	return AppError{
		FieldName: strings.Join(args, ","),
		ErrorCode: message,
	}
}

// JSON is used to send http error response in JSON.
// This method assumed that it uses echo.Context.
// err params is the error that we want to show to user.
// originalErrs is the error that we want to log to the backend log system. Sometimes
// what we want to show to user and what we want to know for debugging purposes is different.
// So to accomodate that case, we give optional originalErrs params so we can give a user friendly,
// translated message, and still have an original error for debugging.
func JSON(c echo.Context, httpStatus int, appErr error, originalErrs ...error) error {
	errorResponse := AppError{
		ErrorMessage: appErr.Error(),
	}

	req := c.Request()
	res := c.Response()

	qs := ""
	var fv url.Values
	if !strings.HasPrefix(req.Header.Get(echo.HeaderContentType), echo.MIMEMultipartForm) {
		qs = c.QueryString()

		var err error
		fv, err = c.FormParams()
		if err != nil {
			c.Error(err)
		}
	}

	file, line := getFileAndLineNumber()

	logFields := map[string]interface{}{
		"user_id":      c.Get("USER_UID"),
		"request_id":   res.Header().Get(echo.HeaderXRequestID),
		"ip":           c.RealIP(),
		"host":         req.Host,
		"uri":          req.RequestURI,
		"method":       req.Method,
		"user_agent":   req.UserAgent(),
		"status":       httpStatus,
		"query_string": qs,
		"form_values":  fv,
		"file":         file,
		"line":         line,
	}

	if ae, ok := appErr.(AppError); ok {
		errorResponse.ErrorCode = ae.ErrorCode
		errorResponse.ErrorMessage = ae.ErrorMessage
		errorResponse.FieldName = ae.FieldName

		logFields["error_code"] = ae.ErrorCode
		logFields["error_message"] = ae.ErrorMessage
		logFields["field_name"] = ae.FieldName

		if len(originalErrs) > 0 {
			logFields["original_error_message"] = originalErrs[0]
		}

		logData := log.WithFields(logFields)

		if httpStatus == http.StatusInternalServerError {
			logData.Error()
		} else {
			logData.Info()
		}
	} else {
		errorResponse.ErrorMessage = appErr.Error()

		logFields["error_message"] = appErr.Error()

		log.WithFields(logFields).Error()

	}

	return c.JSON(httpStatus, errorResponse)
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
