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
	if err == nil {
		return false
	}
	_, ok := err.(interface{ IsUserError() })
	return ok
}

func UserError(err error) error {
	if err == nil {
		return nil
	}
	if IsUserError(err) {
		return err
	}
	return &userError{err: err}
}
