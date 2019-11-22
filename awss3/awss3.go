package awss3

import (
	"bytes"
	"io/ioutil"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/nuveo/log"
)

// Storing implementation for aws s3
type Storing struct {
	Session *session.Session
	Bucket  string
	ACL     string
}

// New s3 storing
func New(opts ...func(*Storing) error) (st *Storing, err error) {
	ss, err := session.NewSession()
	if err != nil {
		return
	}
	st = &Storing{
		Session: ss,
		ACL:     os.Getenv("AWS_ACL"),
		Bucket:  os.Getenv("AWS_BUCKET"),
	}
	for _, opt := range opts {
		if err = opt(st); err != nil {
			st = nil
			return
		}
	}
	return
}

// ACL is an option to set when storing is created
func ACL(acl string) func(*Storing) error {
	return func(st *Storing) error {
		st.ACL = acl
		return nil
	}
}

// Bucket is an option to set when storing is created
func Bucket(bucket string) func(*Storing) error {
	return func(st *Storing) error {
		st.Bucket = bucket
		return nil
	}
}

// CustomSession is an option to set when storing is created
func CustomSession(s *session.Session) func(*Storing) error {
	return func(st *Storing) error {
		st.Session = s
		return nil
	}
}

// Provider returns the name of the provider of the current adapter.
func (s *Storing) Provider() string {
	return "s3"
}

// Upload upload file to S3
func (s *Storing) Upload(name string, contentType string, content []byte) (path string, err error) {
	log.Printf("upload file to s3 %v\n", name)
	uploader := s3manager.NewUploader(s.Session)

	var file *s3manager.UploadOutput
	file, err = uploader.Upload(&s3manager.UploadInput{
		Bucket:      aws.String(s.Bucket),
		ACL:         aws.String(s.ACL),
		Key:         aws.String(name),
		ContentType: aws.String(contentType),
		Body:        bytes.NewReader(content),
	})
	if err != nil {
		return
	}
	path = file.Location
	return
}

// Download file from s3
func (s *Storing) Download(path string) (b []byte, err error) {
	log.Printf("download file from s3 %v\n", path)
	downloader := s3manager.NewDownloader(s.Session)
	if err != nil {
		return
	}

	tmpfile, err := ioutil.TempFile("", "nuveo")
	if err != nil {
		return
	}
	defer func() {
		rmErr := os.Remove(tmpfile.Name())
		if rmErr != nil {
			log.Errorln(rmErr)
		}
	}()

	_, err = downloader.Download(tmpfile, &s3.GetObjectInput{
		Bucket: aws.String(s.Bucket),
		Key:    aws.String(path),
	})
	if err != nil {
		return
	}

	b, err = ioutil.ReadAll(tmpfile)
	return
}

// Delete from s3
func (s *Storing) Delete(key string) (err error) {
	log.Printf("delete file from s3 %v\n", key)
	svc := s3.New(s.Session)
	obj := &s3.DeleteObjectInput{
		Bucket: aws.String(s.Bucket),
		Key:    aws.String(key),
	}
	_, err = svc.DeleteObject(obj)
	return
}
