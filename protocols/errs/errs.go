package errs

import "fmt"

type OctlError struct {
	code int
	emsg string
}

func (e *OctlError) Error() string {
	if e == nil {
		return "OK"
	}
	return e.emsg
}

func (e *OctlError) Code() int {
	if e == nil {
		return 0
	}
	return e.code
}

func (e *OctlError) String() string {
	if e == nil {
		return "code=0: OK"
	}
	return fmt.Sprintf("code=%d: %s", e.code, e.emsg)
}

func New(code int, emsg string) *OctlError {
	return &OctlError{
		code: code,
		emsg: emsg,
	}
}

const (
	OctlReadConfigError = 1 + iota
	OctlWriteConfigError
	OctlInitClientError
	OctlWorkgroupAuthError
	OctlHttpRequestError
	OctlHttpStatusError
	OctlMessageParseError
	OctlNodeParseError
	OctlFileOperationError
	OctlGitOperationError
	OctlTaskWaitingError
	OctlArgumentError
	OctlSdkNotInitializedError
	OctlSdkPanicRecoverError
	OctlSdkBufferError
	OctlContextCancelError
)
