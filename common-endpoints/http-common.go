package ce

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func writeJson(w http.ResponseWriter, code int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	switch data.(type) {
	case []byte:
		w.Write(data.([]byte))
	default:
		enc := json.NewEncoder(w)
		enc.Encode(data)
	}
}

func writeError(w http.ResponseWriter, code int, msg string) {
	writeJson(w, code, map[string]interface{}{"code": code, "msg": msg})
}

func readJson(r *http.Request, res interface{}) (status int, err error) {
	if r.Body == nil {
		return http.StatusBadRequest, fmt.Errorf("bad request")
	}
	defer r.Body.Close()

	if err = json.NewDecoder(r.Body).Decode(res); err != nil {
		return http.StatusBadRequest, err
	}
	return http.StatusOK, nil
}
