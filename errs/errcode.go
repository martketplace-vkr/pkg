package errs

type errCode int

const (
	ErrCodeBadRequest     errCode = 400
	ErrCodeUnauthorized   errCode = 401
	ErrCodeForbidden      errCode = 403
	ErrCodeNotFound       errCode = 404
	ErrCodeNotAllowed     errCode = 405
	ErrCodeNotAcceptable  errCode = 406
	ErrCodeRequestTimeout errCode = 408
	ErrCodeLocked         errCode = 423
	ErrCodeOutOfLimit     errCode = 430

	ErrCodeInternal       errCode = 500
	ErrCodeNotImplemented errCode = 501
	ErrCodeBadGateway     errCode = 502
	ErrCodeUnknown        errCode = 520
)
