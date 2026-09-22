package access

import (
	"fmt"

	"github.com/samber/mo"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func newR2Client(accountID, accessKeyID, secretAccessKey string) mo.Result[*minio.Client] {
	client, err := minio.New(fmt.Sprintf("%s.r2.cloudflarestorage.com", accountID), &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: true,
	})
	if err != nil {
		return mo.Err[*minio.Client](err)
	}

	return mo.Ok(client)
}
