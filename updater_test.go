package rtn

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/stretchr/testify/mock"
)

const fedFile = `011000015O0110000150122415000000000FEDERAL RESERVE BANK                1000 PEACHTREE ST N.E.              ATLANTA             GA303094470877372245711     
011000028O0110000151072811000000000STATE STREET BANK AND TRUST COMPANY JAB2NW                              N. QUINCY           MA021710000617664240011     
invalid line length adds empty entry instead of panicing
`

var entries = []achEntry{achEntry{City: "ATLANTA", State: "GA", RoutingNumber: "011000015", Name: "FEDERAL RESERVE BANK", LastModified: "12-24-15"}, achEntry{City: "N. QUINCY", State: "MA", RoutingNumber: "011000028", Name: "STATE STREET BANK AND TRUST COMPANY", LastModified: "07-28-11"}, achEntry{City: "", State: "", RoutingNumber: "", Name: "", LastModified: ""}}

func init() {
	log.SetOutput(ioutil.Discard)
}

func TestNewUpdater(t *testing.T) {
	if NewUpdater() == nil {
		t.Errorf("Nil updater.")
	}
}

type tc struct {
	getFileErr error
	marshalErr error
	headErr    error
	putErr     error
	expected   error
}

func TestUpdate(t *testing.T) {
	e := errors.New("error")
	awsErr := awserr.New("code", "message", nil)
	cases := map[string]tc{
		"getFile error":    {e, nil, nil, nil, e},
		"marshal error":    {nil, e, nil, nil, e},
		"HeadBucket error": {nil, nil, awsErr, nil, awsErr},
		"PutObject error":  {nil, nil, nil, awsErr, awsErr},
		"success":          {nil, nil, nil, nil, nil},
	}

	for k, c := range cases {
		u := setupUpdater(c)

		err := u.Update()

		if err != c.expected {
			t.Errorf("Update() == %v, expected %v. %v", err, c.expected, k)
		}
	}
}

// Ignored, hits the Fed site and uploads to S3
func _TestFunctional(t *testing.T) {
	err := NewUpdater().Update()
	if err != nil {
		t.Error("err=" + err.Error())
	}
}

func setupUpdater(c tc) Updater {
	const bucket = "bucket-name"
	mm := new(marshalerMock)
	fcm := new(fedClientMock)
	s3m := new(s3Mock)
	str, _ := json.Marshal(entries)

	fcm.On("getFile").Return(string(fedFile), c.getFileErr)
	mm.On("marshal", entries).Return(str, c.marshalErr)
	s3m.On("HeadBucket", &s3.HeadBucketInput{Bucket: aws.String(bucket)}).
		Return(&s3.HeadBucketOutput{}, c.headErr)
	s3m.On("PutObject", &s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String("info.json"),
		Body:   aws.ReadSeekCloser(bytes.NewReader(str)),
	}).Return(&s3.PutObjectOutput{}, c.putErr)

	return Updater{fcm, s3m, bucket, mm.marshal}
}

type fedClientMock struct {
	mock.Mock
}

func (m *fedClientMock) getFile() (string, error) {
	args := m.Called()
	return args.String(0), args.Error(1)
}

type s3Mock struct {
	s3iface.S3API
	mock.Mock
}

func (m *s3Mock) HeadBucket(i *s3.HeadBucketInput) (*s3.HeadBucketOutput, error) {
	args := m.Called(i)
	return args.Get(0).(*s3.HeadBucketOutput), args.Error(1)
}

func (m *s3Mock) PutObject(i *s3.PutObjectInput) (*s3.PutObjectOutput, error) {
	args := m.Called(i)
	return args.Get(0).(*s3.PutObjectOutput), args.Error(1)
}

type marshalerMock struct {
	mock.Mock
}

func (m *marshalerMock) marshal(i interface{}) ([]byte, error) {
	args := m.Called(i)
	return args.Get(0).([]byte), args.Error(1)
}
