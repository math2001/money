package api

import (
	"encoding/json"
	"log"
	"net/http"
)

func respond(w http.ResponseWriter, r *http.Request, code int, kind string, parts ...interface{}) {
	if len(parts)%2 == 1 {
		log.Printf("%v cannot generate map from odd number of parts: %d", r, len(parts))
		return
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
				log.Printf("%v key 'kind' is reserved (currently set to %q)", r, parts[i+1])
				return
			}
			obj[key] = parts[i+1]
		}
	}

	enc := json.NewEncoder(w)
	if err := enc.Encode(obj); err != nil {
		log.Printf("%v writing json obj: %s", r, err)
	}
}
