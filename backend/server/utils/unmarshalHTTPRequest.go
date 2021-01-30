package utils

import (
	"fmt"
	"errors"
	"net/http"

	"github.com/mholt/binding"
)

// UnmarshalRequest unmarshals request.FormValue into arbitrary structs based on mholt/binding
func UnmarshalRequest(r *http.Request, s interface{}) error {
	// binding uses pkg/errors so unwrap to return stdlib error
	if sBinding, ok := s.(binding.FieldMapper); ok {
		return errors.Unwrap(binding.Bind(r, sBinding))
	}

	return fmt.Errorf("%v does not implement binding.FieldMapper interface", s)
}