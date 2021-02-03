package postservice

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/mholt/binding"
)

//go:generate go run ../requestBinding/bindingGenerator.go -f post.go -out APIBindings.go -package postservice

func parseRequest(r *http.Request) (*PostRequest, error) {
	var request PostRequest

	err := UnmarshalRequest(r, &request)
	if err != nil {
		return nil, err
	}

	request.Slug = chi.URLParam(r, "slug")

	return &request, nil
}

// UnmarshalRequest unmarshals request.FormValue into arbitrary structs based on mholt/binding
func UnmarshalRequest(r *http.Request, s interface{}) error {
	// binding uses pkg/errors so unwrap to return stdlib error
	if sBinding, ok := s.(binding.FieldMapper); ok {
		return errors.Unwrap(binding.Bind(r, sBinding))
	}

	return fmt.Errorf("%v does not implement binding.FieldMapper interface", s)
}