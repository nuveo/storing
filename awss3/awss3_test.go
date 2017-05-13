package awss3

import (
	"bytes"
	"io/ioutil"
	"testing"
)

func TestS3(t *testing.T) {

	bUp, err := ioutil.ReadFile("../testdata/testfile.txt")
	if err != nil {
		t.Fatal(err)
	}

	// Upload
	_, err = Upload("testfile.txt", "text/plain", bUp)
	if err != nil {
		t.Fatal(err)
	}

	// Download
	var bDown []byte
	bDown, err = Download("testfile.txt")
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(bUp, bDown) {
		t.Fatal("Uploaded data are different from the data received in the download.")
	}

	// Delete
	err = Delete("testfile.txt")
	if err != nil {
		t.Fatal(err)
	}

	_, err = Download("testfile.txt")
	if err == nil {
		t.Fatal("An error was expected")
	}
}
