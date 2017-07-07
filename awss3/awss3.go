package awss3

import (
	"bytes"
	"io/ioutil"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

// Storing implementation for aws s3
type Storing struct {
	session *session.Session
}

// Provider returns the name of the provider of the current adapter.
func (s *Storing) Provider() string {
	return "s3"
}

func (s *Storing) startS3Session() (err error) {
	if s.session == nil {
		s.session, err = session.NewSession()
		return
	}
	return
}

// Upload upload file to S3
func (s *Storing) Upload(name string, contentType string, content []byte) (path string, err error) {
	err = s.startS3Session()
	if err != nil {
		return
	}

	uploader := s3manager.NewUploader(s.session)

	var file *s3manager.UploadOutput
	file, err = uploader.Upload(&s3manager.UploadInput{
		Bucket:      aws.String("nuveofs"),
		ACL:         aws.String("private"),
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
	err = s.startS3Session()
	if err != nil {
		return
	}

	downloader := s3manager.NewDownloader(s.session)
	if err != nil {
		return
	}

	tmpfile, err := ioutil.TempFile("", "nuveo")
	if err != nil {
		return
	}
	defer os.Remove(tmpfile.Name())

	_, err = downloader.Download(tmpfile, &s3.GetObjectInput{
		Bucket: aws.String("nuveofs"),
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
	err = s.startS3Session()
	if err != nil {
		return
	}

	svc := s3.New(s.session)
	obj := &s3.DeleteObjectInput{
		Bucket: aws.String("nuveofs"),
		Key:    aws.String(key),
	}
	_, err = svc.DeleteObject(obj)
	return
}
