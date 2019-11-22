package azureblob

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"os"

	"github.com/Azure/azure-storage-blob-go/azblob"
)

// Storing implementation for azure blob
type Storing struct {
	ContainerName    string
	StorageAccount   string
	StorageAccessKey string
}

// New create struct to deal with Azure blobs
func New() (*Storing, error) {
	st := &Storing{
		StorageAccount:   os.Getenv("AZURE_STORAGE_ACCOUNT"),
		ContainerName:    os.Getenv("AZURE_STORAGE_CONTAINER"),
		StorageAccessKey: os.Getenv("AZURE_STORAGE_ACESS_KEY"),
	}

	err := allStoringFieldsCorrect(st)

	if err != nil {
		return nil, err
	}

	return st, nil

}

// Upload upload a file to blob
func (s *Storing) Upload(name string, contentType string, content []byte) (string, error) {
	credential, err := azblob.NewSharedKeyCredential(s.StorageAccount, s.StorageAccessKey)
	p := azblob.NewPipeline(credential, azblob.PipelineOptions{})
	URL, _ := url.Parse(
		fmt.Sprintf("https://%s.blob.core.windows.net/%s", s.StorageAccount, s.ContainerName))
	containerURL := azblob.NewContainerURL(*URL, p)
	blobURL := containerURL.NewBlockBlobURL(name)

	ctx := context.Background()
	_, err = azblob.UploadBufferToBlockBlob(ctx, content, blobURL, azblob.UploadToBlockBlobOptions{
		BlockSize:   4 * 1024 * 1024,
		Parallelism: 16})

	if err != nil {
		return "", err
	}

	return name, err
}

func (s *Storing) Download(name string) ([]byte, error) {
	return nil, nil
}

func Provider() string {
	return "blob"
}

func (s *Storing) Delete(name string) error {
	return nil
}

func allStoringFieldsCorrect(st *Storing) error {
	if st.ContainerName == "" || st.StorageAccount == "" || st.StorageAccessKey == "" {
		return errors.New("You need to set AZURE_STORAGE_ACCOUNT, AZURE_STORAGE_CONTAINER and AZURE_STORAGE_ACESS_KEY as env variables")
	}

	return nil
}
