package access

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"os"

	"github.com/BlueWitherer/ferry-server/log"
	"github.com/samber/mo"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var client *minio.Client
var bucket string

func getR2Client() mo.Result[*minio.Client] {
	if client != nil {
		return mo.Ok(client)
	}

	c, err := minio.New(fmt.Sprintf("%s.r2.cloudflarestorage.com", os.Getenv("R2_ACCOUNT_ID")), &minio.Options{
		Creds:  credentials.NewStaticV4(os.Getenv("R2_ACCESS_KEY_ID"), os.Getenv("R2_SECRET_ACCESS_KEY"), ""),
		Secure: true,
	})
	if err != nil {
		return mo.Err[*minio.Client](err)
	}

	client = c

	return mo.Ok(c)
}

func ParseGameVarBody(data []byte) mo.Result[map[string]bool] {
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

		if strLen != 4 {
			return mo.Errf[map[string]bool]("Key string length is invalid")
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

	log.Trace("Size: %v, Map: %+v\n", size, res)
	return mo.Ok(res)
}

func r2Read(key, name string) mo.Result[[]byte] {
	cRes := getR2Client()
	if cRes.IsError() {
		return mo.Err[[]byte](cRes.Error())
	}
	client := cRes.MustGet()

	objectName := fmt.Sprintf("%s/%s", key, name)

	obj, err := client.GetObject(context.Background(), bucket, objectName, minio.GetObjectOptions{})
	if err != nil {
		return mo.Err[[]byte](err)
	}
	defer obj.Close()

	b, err := io.ReadAll(obj)
	if err != nil {
		return mo.Err[[]byte](err)
	}

	return mo.Ok(b)
}

func r2Write(key, name string, data []byte) mo.Result[bool] {
	cRes := getR2Client()
	if cRes.IsError() {
		return mo.Err[bool](cRes.Error())
	}
	client := cRes.MustGet()

	objectName := fmt.Sprintf("%s/%s", key, name)

	_, err := client.PutObject(context.Background(), bucket, objectName,
		bytes.NewReader(data), int64(len(data)),
		minio.PutObjectOptions{ContentType: "application/octet-stream"},
	)
	if err != nil {
		return mo.Err[bool](err)
	}

	return mo.Ok(true)
}

func R2ReadGameVarsRaw(accountID int) mo.Result[[]byte] {
	dataRes := r2Read(fmt.Sprintf("%v", accountID), "gv")
	if dataRes.IsError() {
		return mo.Err[[]byte](dataRes.Error())
	}

	return dataRes
}

func R2ReadGameVars(accountID int) mo.Result[map[string]bool] {
	dataRes := R2ReadGameVarsRaw(accountID)
	if dataRes.IsError() {
		return mo.Err[map[string]bool](dataRes.Error())
	}

	return ParseGameVarBody(dataRes.MustGet())
}

func R2WriteGameVars(accountID int, data []byte) mo.Result[bool] {
	return r2Write(fmt.Sprintf("%v", accountID), "gv", data)
}

func init() {
	bucket = os.Getenv("R2_BUCKET")
}
