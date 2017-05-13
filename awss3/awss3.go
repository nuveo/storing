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

var awsSession *session.Session

func Provider() string {
	return "s3"
}

func startS3Session() (err error) {
	if awsSession == nil {
		awsSession, err = session.NewSession()
		return
	}
	return
}

// Upload upload file to S3
func Upload(name string, contentType string, content []byte) (path string, err error) {
	err = startS3Session()
	if err != nil {
		return nil, err
	}

	uploader := s3manager.NewUploader(awsSession)

	var file *s3manager.UploadOutput
	file, err = uploader.Upload(&s3manager.UploadInput{
		Bucket:      aws.String("nuveofs"),
		ACL:         aws.String("public-read"),
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
func Download(path string) (b []byte, err error) {
	err = startS3Session()
	if err != nil {
		return
	}

	downloader := s3manager.NewDownloader(awsSession)
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

	b, err = ioutil.ReadAll(f)
	return
}

// Delete from s3
func Delete(key string) (err error) {
	err = startS3Session()
	if err != nil {
		return nil, err
	}

	svc := s3.New(awsSession)
	obj := &s3.DeleteObjectInput{
		Bucket: aws.String("nuveofs"),
		Key:    aws.String(key),
	}
	_, err = svc.DeleteObject(obj)
	return
}
