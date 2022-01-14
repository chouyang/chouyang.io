package errors

type Throwable interface {
	GetCode() string
	error
}
