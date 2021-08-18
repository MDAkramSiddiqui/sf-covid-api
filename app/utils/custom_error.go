package utils

// custom error struct having status code as well
type CustomErr struct {
	Err        error
	StatusCode int
}

func (err *CustomErr) Message() string {
	return err.Err.Error()
}
