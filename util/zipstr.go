package util

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"io/ioutil"
)

func CompressString(s string) string {
	var buf bytes.Buffer
	zw, _ := gzip.NewWriterLevel(&buf, gzip.BestCompression)
	_, _ = zw.Write([]byte(s))
	zw.Close()
	zipResult := buf.Bytes()
	return base64.URLEncoding.EncodeToString(zipResult)
}

func DecompressString(s string) (string, error) {
	zipData, err := base64.URLEncoding.DecodeString(s)
	if err != nil {
		return "", err
	}
	zipReader, err := gzip.NewReader(bytes.NewBuffer(zipData))
	if err != nil {
		return "", err
	}
	result, err := ioutil.ReadAll(zipReader)
	if err != nil {
		return "", err
	}
	return string(result), nil
}
