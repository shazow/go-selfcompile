package selfcompile

import (
	"fmt"
	"strings"
)

type combinedErrors []error

func (errs combinedErrors) Error() string {
	s := fmt.Sprintf("%d errors: ", len(errs))
	msgs := make([]string, 0, len(errs))
	for _, err := range errs {
		msgs = append(msgs, err.Error())
	}
	return s + strings.Join(msgs, "; ")
}

// combineErrors merges multiple non-nil errors into one.
func combineErrors(errs ...error) error {
	validErrors := []error{}

	for _, err := range errs {
		if err != nil {
			validErrors = append(validErrors, err)
		}
	}
	if len(validErrors) == 0 {
		return nil
	}
	if len(validErrors) == 1 {
		return validErrors[0]
	}
	return combinedErrors(validErrors)
}
