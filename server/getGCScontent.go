package server

import (
	"context"
	"io"
	"io/ioutil"
	"os"

	"cloud.google.com/go/storage"
)

func saveGCSContentAsFile(ctx context.Context, bkt *storage.BucketHandle, gcsPath string, filePath string) error {
	obj := bkt.Object(gcsPath)

	r, err := obj.NewReader(ctx)
	if err != nil {
		return err
	}
	defer r.Close()

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, r)
	if err != nil {
		return err
	}

	return nil
}

func getGCSContentAsBytes(ctx context.Context, bkt *storage.BucketHandle, gcsPath string) ([]byte, error) {
	obj := bkt.Object(gcsPath)

	r, err := obj.NewReader(ctx)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	ans, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	return ans, nil
}
