package api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func respond(w http.ResponseWriter, code int, kind string, parts ...interface{}) error {
	if len(parts)%2 == 1 {
		return fmt.Errorf("cannot generate map from odd number of parts: %d", len(parts))
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	obj := make(map[string]interface{}, len(parts)/2)
	obj["kind"] = kind
	for i, part := range parts {
		if i%2 == 1 {
			continue
		}
		if key, ok := part.(string); ok {
			if key == "kind" {
				return fmt.Errorf("key 'kind' is reserved (currently set to %q)", parts[i+1])
			}
			obj[key] = parts[i+1]
		}
	}

	enc := json.NewEncoder(w)
	if err := enc.Encode(obj); err != nil {
		return fmt.Errorf("writing json obj: %s", err)
	}

	return nil
}
