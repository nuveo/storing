package azureblob

import (
	"io/ioutil"
	"testing"
)

func TestBlob(t *testing.T) {
	s, err := New()
	if err != nil {
		t.Fatal(err)
	}

	bUp, err := ioutil.ReadFile("../testdata/testfile.txt")
	if err != nil {
		t.Fatal(err)
	}

	// Upload
	_, err = s.Upload("testfile.txt", "text/plain", bUp)
	if err != nil {
		t.Fatal(err)
	}

}
