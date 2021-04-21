package utils

func InternalError() {
	println("Oops, internal error")
}

type RuntimeError interface {
	ErrorData() interface{}
	Message() string
}

type CantOpenFileError struct {
	Filename string
	Message  string
}

func (error *CantOpenFileError) ErrorData() interface{} {
	return error
}

func (error *CantOpenFileError) GetMessage() string {
	return "Can't open file '" + error.Filename + "'"
}

func (error *CantOpenFileError) Error() string {
	return "Can't open file '" + error.Filename + "'"
}
