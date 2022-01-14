package errors

type AccessDenied struct {
	Code string
	Msg  string
}

func (e AccessDenied) GetCode() string {
	if e.Code == "" {
		e.Code = ApiNotAllowed
	}

	return e.Code
}

func (e AccessDenied) Error() string {
	if e.Msg == "" {
		e.Msg = "you don't have permission to access this resource"
	}

	return e.Msg
}
