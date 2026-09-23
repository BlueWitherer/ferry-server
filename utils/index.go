package utils

import (
	"encoding/json"
	"net/http"

	"github.com/BlueWitherer/GDDataSyncServer/log"
	"github.com/samber/mo"
)

type WebRes[T any] struct {
	Payload *T     `json:"payload"`
	Error   string `json:"error"`
}

func WriteHeaders(header *http.Header, method string) {
	header.Set("Access-Control-Allow-Origin", "*")
	header.Set("Access-Control-Allow-Methods", method)
	header.Set("Access-Control-Allow-Headers", "Content-Type")
	header.Set("Content-Type", "application/json")
}

func WriteWebRes[T any](w http.ResponseWriter, payload mo.Option[T], code int) mo.Result[bool] {
	var out WebRes[T]

	if v, ok := payload.Get(); ok {
		out.Payload = &v
	} else {
		out.Error = "Failed to encode response"

		w.WriteHeader(code)
		if err := json.NewEncoder(w).Encode(out); err != nil {
			log.Error("Failed to encode response: %s", err.Error())
			http.Error(w, "Failed to encode response", code)

			return mo.Err[bool](err)
		}

		return mo.Errf[bool]("%s", out.Error)
	}

	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(out); err != nil {
		log.Error("Failed to encode response: %s", err.Error())
		http.Error(w, "Failed to encode response", code)

		return mo.Err[bool](err)
	}

	return mo.Ok(true)
}

func WriteWebErr(w http.ResponseWriter, message string, code int) mo.Result[bool] {
	var out WebRes[bool]
	out.Payload = nil
	out.Error = message

	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(out); err != nil {
		log.Error("Failed to encode error response: %s", err.Error())
		http.Error(w, "Failed to encode error response", code)

		return mo.Err[bool](err)
	}

	return mo.Ok(true)
}
