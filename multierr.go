package stealthpool

import "strings"

// MultiError holds multiple errors and concatenates them
type MultiError struct {
	errs []error
}

func newMultiErr() *MultiError {
	return &MultiError{}
}

// Add adds an error to the existing list
func (e *MultiError) Add(err error) {
	if err != nil {
		e.errs = append(e.errs, err)
	}
}

// Return is the return value of the multierror. If it doesn't contain any errors, nil will be returned
func (e *MultiError) Return() error {
	if len(e.errs) == 0 {
		return nil
	}
	return e
}

// Error returns the concatenation of all errors stored inside the multierror
func (e *MultiError) Error() string {
	result := strings.Builder{}
	for _, e := range e.errs {
		result.WriteString(e.Error() + "\n")
	}
	return result.String()
}
