package pkg

type CollectError struct {
	Err  error
	Type string
}

func (e *CollectError) Error() string {
	return e.Err.Error()
}
