package input

import (
	"io/ioutil"
	"log"
	"os"
	"testing"
)

func TestScan(t *testing.T) {
	content := []byte("xyz@xyz.com")
	tmpfile, err := ioutil.TempFile("", "test")
	if err != nil {
		log.Fatal(err)
	}

	defer os.Remove(tmpfile.Name()) // clean up

	if _, err := tmpfile.Write(content); err != nil {
		log.Fatal(err)
	}

	if _, err := tmpfile.Seek(0, 0); err != nil {
		log.Fatal(err)
	}

	originalStdin := os.Stdin
	defer func() { os.Stdin = originalStdin }() // Restore original Stdin

	os.Stdin = tmpfile
	if input := Scan("Email: "); input != "xyz@xyz.com" {
		t.Errorf("Expectations mismatched: \n\t\t expected: xyz@xyz.com \n\t\t actual: %s", input)
	}

	if err := tmpfile.Close(); err != nil {
		log.Fatal(err)
	}

}
