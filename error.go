package nrpc

type userError struct {
	err error
}

func (ue *userError) IsUserError() {}

func (ue *userError) Unwrap() error {
	return ue.err
}

func (ue *userError) Error() string {
	return ue.err.Error()
}

func IsUserError(err error) bool {
	_, ok := err.(interface{ IsUserError() })
	return ok
}

func UserError(err error) error {
	if err == nil {
		return nil
	}
	return &userError{err: err}
}
