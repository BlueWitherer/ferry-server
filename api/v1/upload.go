package v1

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"net/http"

	"github.com/BlueWitherer/GDDataSyncServer/log"
	"github.com/BlueWitherer/GDDataSyncServer/utils"
	"github.com/samber/mo"
)

func parseGameVarBody(data []byte) mo.Result[map[string]bool] {
	buf := bytes.NewReader(data)

	var size uint64
	if err := binary.Read(buf, binary.LittleEndian, &size); err != nil {
		log.Error("failed to read size: %v", err)
		return mo.Err[map[string]bool](err)
	}

	res := make(map[string]bool)

	for buf.Len() > 0 {
		var strLen uint8
		if err := binary.Read(buf, binary.LittleEndian, &strLen); err != nil {
			return mo.Err[map[string]bool](err)
		}

		strBytes := make([]byte, strLen)
		if _, err := buf.Read(strBytes); err != nil {
			return mo.Err[map[string]bool](err)
		}

		var b bool
		if err := binary.Read(buf, binary.LittleEndian, &b); err != nil {
			return mo.Err[map[string]bool](err)
		}

		res[string(strBytes)] = b
	}

	fmt.Printf("Size: %d, Map: %+v\n", size, res)
	return mo.Ok(res)
}

func init() {
	http.HandleFunc("/v1/upload", func(w http.ResponseWriter, r *http.Request) { // geometry dash gamevars
		header := w.Header()
		utils.WriteHeaders(&header, http.MethodPost)

		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "failed to read body", http.StatusInternalServerError)
			return
		}
		defer r.Body.Close()

		parseGameVarBody(body)
	})

	http.HandleFunc("/v1/upload-mods", func(w http.ResponseWriter, r *http.Request) { // all mod settings
		header := w.Header()
		utils.WriteHeaders(&header, http.MethodPost)
	})

	http.HandleFunc("/v1/upload-mods-saves", func(w http.ResponseWriter, r *http.Request) { // all mod save data (supporter only probably)
		header := w.Header()
		utils.WriteHeaders(&header, http.MethodPost)
	})
}
