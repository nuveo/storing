package mock

import (
	"sync"
)

type item struct {
	ContentType string
	Content     []byte
}

// Storing mock
type Storing struct {
	mtx         sync.RWMutex
	UploadErr   error
	DownloadErr error
	DeleteErr   error
	Items       map[string]item
}

// Provider returns the name of the provider of the current adapter.
func (s *Storing) Provider() string {
	return "Mock"
}

// Upload upload file to mock
func (s *Storing) Upload(name string, contentType string, content []byte) (path string, err error) {
	err = s.UploadErr
	if err == nil {
		s.mtx.Lock()
		s.Items[name] = item{ContentType: contentType, Content: content}
		s.mtx.Unlock()
		path = name
	}
	return
}

// Download file from mock
func (s *Storing) Download(path string) (b []byte, err error) {
	err = s.DownloadErr
	if err == nil {
		s.mtx.Lock()
		b = s.Items[path].Content
		s.mtx.Unlock()
	}
	return
}

// Delete from mock
func (s *Storing) Delete(key string) (err error) {
	err = s.DeleteErr
	if err == nil {
		s.mtx.Lock()
		delete(s.Items, key)
		s.mtx.Unlock()
	}
	return
}
