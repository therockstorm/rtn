package rtn

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http/cookiejar"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
)

type fedClient interface {
	getFile() (string, error)
}

// Updater updates routing information.
type Updater struct {
	fed       fedClient
	s3        s3iface.S3API
	bucket    string
	marshaler func(v interface{}) ([]byte, error)
}

// NewUpdater constructor for Updater.
func NewUpdater() *Updater {
	return &Updater{
		fed:       newClient(cookiejar.New),
		s3:        s3.New(session.Must(session.NewSession(&aws.Config{Region: aws.String("us-west-2")}))),
		bucket:    "rtn-lookup-info",
		marshaler: json.Marshal,
	}
}

// Update updates an S3 file containing bank name and routing information.
func (i *Updater) Update() error {
	res, err := i.fed.getFile()
	if err != nil {
		return err
	}

	lines := strings.Split(strings.TrimSpace(res), "\n")
	entries := make([]achEntry, len(lines))
	for i, l := range lines {
		entries[i] = makeFrom(l)
	}

	e, err := i.marshaler(entries)
	if err != nil {
		return err
	}

	return i.upload(e)
}

func (i *Updater) upload(b []byte) error {
	_, herr := i.s3.HeadBucket(&s3.HeadBucketInput{Bucket: aws.String(i.bucket)})
	if herr != nil {
		return herr
	}

	log.Println("Uploading to " + i.bucket)

	_, err := i.s3.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(i.bucket),
		Key:    aws.String("info.json"),
		Body:   aws.ReadSeekCloser(bytes.NewReader(b)),
	})

	return err
}
