package iptables

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"

	"github.com/pkg/errors"
)

func writeFile(filePath string, contents []byte) error {
	fileObj, err := os.Create(filePath)
	if err != nil {
		return errors.Wrapf(err, "could not open file to write")
	}
	defer fileObj.Close()

	if _, err = io.Copy(fileObj, bytes.NewReader(contents)); err != nil {
		return errors.Wrapf(err, "could not write to file")
	}
	return nil
}

func getTempFile() (*os.File, error) {
	file, err := ioutil.TempFile(os.TempDir(), "shield")
	if err != nil {
		return nil, err
	}
	return file, nil
}
