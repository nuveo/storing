package storing

// Storing interface
type Storing interface {
	Upload(string, string, []byte) (string, error)
	Download(string) ([]byte, error)
	Provider() string
	Delete(string) error
}
