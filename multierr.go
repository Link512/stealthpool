package stealthpool

import "strings"

type multiErr struct {
	errs []error
}

func newMultiErr() *multiErr {
	return &multiErr{}
}

func (e *multiErr) Add(err error) {
	if err != nil {
		e.errs = append(e.errs, err)
	}
}

func (e *multiErr) Return() error {
	if len(e.errs) == 0 {
		return nil
	}
	return e
}

func (e *multiErr) Error() string {
	result := strings.Builder{}
	for _, e := range e.errs {
		result.WriteString(e.Error() + "\n")
	}
	return strings.TrimRight(result.String(), "\n")
}
