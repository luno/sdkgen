package clientgen

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"

	openapi "github.com/go-openapi/spec"
)

// LoadSpec loads the OpenAPI spec.
// fn can be a URL or a filename.
func LoadSpec(fn string) (*openapi.Swagger, error) {
	if strings.HasPrefix(fn, "http://") || strings.HasPrefix(fn, "https://") {
		f, err := http.Get(fn)
		if err != nil {
			return nil, err
		}
		defer f.Body.Close()

		var spec openapi.Swagger
		if err := json.NewDecoder(f.Body).Decode(&spec); err != nil {
			return nil, err
		}

		return &spec, nil
	}

	b, err := os.ReadFile(fn)
	if err != nil {
		return nil, err
	}

	var spec openapi.Swagger
	if err := json.Unmarshal(b, &spec); err != nil {
		return nil, err
	}

	return &spec, nil
}
