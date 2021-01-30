package main

import (
	"bytes"
	"encoding/base32"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	uuid "github.com/satori/go.uuid"
)

var uuidNamespace uuid.UUID = uuid.NewV5(uuid.Nil, "test namespace")

func deterministicUUID(r io.Reader) uuid.UUID {
	content, err := ioutil.ReadAll(r)
	if err != nil {
		return uuid.Nil
	}

	encoded := base32.StdEncoding.EncodeToString(content)

	return uuid.NewV5(uuidNamespace, encoded)
}

// WriteFileToUUIDPath saves each file deterministically and returns the path 
func WriteFileToUUIDPath(r io.Reader, rootDir string) string {
	r1, r2 := splitReader(r)

	u := deterministicUUID(r1)
	path := filepath.Join(rootDir, pathFromUUID(u))

	f, err := createFile(path)
	if err != nil {
		return ""
	}
	defer f.Close()

	_, err = io.Copy(f, r2)
	if err != nil {
		return ""
	}

	return path
}

func createFile(path string) (*os.File, error) {
	os.MkdirAll(filepath.Dir(path), os.ModePerm)
	return os.Create(path)
}

func pathFromUUID(u uuid.UUID) string {
	id := u.String()
	return filepath.Join(id[0:2], id[2:4], id[4:])
}

// splitReader creates a buffered TeeReader allowing an io.Reader to be read twice
// Must use first io.Reader first otherwise second buffer will be empty
func splitReader(r io.Reader) (io.Reader, io.Reader) {
	var buf bytes.Buffer
	return io.TeeReader(r, &buf), &buf
}

func deleteFile(path string) {
	os.Remove(path)
}