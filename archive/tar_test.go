package archive

import (
	"archive/tar"
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"testing"
)

const TAR_FILE = "a.tar"

func TestCompress(t *testing.T) {
	// Create a buffer to write our archive to.
	buf := new(bytes.Buffer)

	//Create a new tar archive
	tw := tar.NewWriter(buf)

	var files = []struct {
		Name, Body string
	}{
		{"readme.txt", "This archive contains some text files."},
		{"gopher.txt", "Gopher names:\nGeorge\nGeoffrey\nGonzo"},
	}

	for _, file := range files {
		hdr := &tar.Header{
			Name: file.Name,
			Mode: 0600,
			Size: int64(len(file.Body)),
		}
		if err := tw.WriteHeader(hdr); err != nil {
			t.Fatal(err)
		}

		if _, err := tw.Write([]byte(file.Body)); err != nil {
			t.Fatal(err)
		}
	}

	if err := tw.Close(); err != nil {
		t.Fatal(err)
	}

	t.Log(len(buf.Bytes()))

	f, err := os.Create(TAR_FILE)
	if err != nil {
		t.Fatal(err)
	}

	f.Write(buf.Bytes())
	f.Close()
}

func TestUncompress(t *testing.T) {
	f, err := os.Open(TAR_FILE)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	buf, err := ioutil.ReadAll(f)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(len(buf))
	r := bytes.NewReader(buf)
	tr := tar.NewReader(r)
	for {
		hdr, err := tr.Next()
		if err == io.EOF { //end of file
			break
		}
		if err != nil {
			t.Fatal(err)
		}
		t.Logf("file name : %s, len : %d \n", hdr.Name, hdr.Size)
		b := make([]byte, hdr.Size)
		_, err = tr.Read(b)
		if err != nil {
			t.Fatal(err)
		}
		t.Log(string(b))
	}
}

func TestFileInfoHeader(t *testing.T) {
	fi, err := os.Stat(TAR_FILE)
	if err != nil {
		t.Fatal(err)
	}

	h, err := tar.FileInfoHeader(fi, "")
	if err != nil {
		t.Fatal(err)
	}

	t.Log(h.Name)
}
