package errors

const (
	// ApiOK is the code for success
	ApiOK = "E0000"

	// ApiError is the code for unknown errors
	ApiError = "E0001"

	// ApiBadRequest is the code for malformed requests
	ApiBadRequest = "E0400"

	// ApiNotAllowed is the code for requests that are not allowed
	ApiNotAllowed = "E0403"

	// ApiNotFound is the code for requests that are not found
	ApiNotFound = "E0404"

	// ApiConflict is the code for requests that are in conflict on creation
	ApiConflict = "E0409"

	// ApiServerError is the code for internal server errors
	ApiServerError = "E0500"
)
