package nrpc

type SolidError interface {
	error
	IsSolid() bool
	Unwrap() error
}

type solidError struct {
	err error
}

func (ue *solidError) IsSolid() bool {
	return true
}

func (ue *solidError) Unwrap() error {
	return ue.err
}

func (ue *solidError) Error() string {
	return ue.err.Error()
}

func IsSolid(err error) bool {
	if err == nil {
		return false
	}
	if se, ok := err.(SolidError); ok {
		return se.IsSolid()
	}
	return false
}

func Solid(err error) error {
	if err == nil {
		return nil
	}
	if IsSolid(err) {
		return err
	}
	return &solidError{err: err}
}
