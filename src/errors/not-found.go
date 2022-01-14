package errors

type NotFound struct {
	Code string
	Msg  string
}

func (e NotFound) GetCode() string {
	if e.Code == "" {
		e.Code = ApiNotFound
	}

	return e.Code
}

func (e NotFound) Error() string {
	if e.Msg == "" {
		e.Msg = "the requested resource was not found"
	}

	return e.Msg
}
