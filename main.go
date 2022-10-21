package main

import (
	"context"
	"log"
	"net/http"

	"cloud.google.com/go/storage"
	"github.com/szpp-dev-team/szpp-judge-judge/server"
)

func main() {
	ctx := context.Background()
	gcs, err := storage.NewClient(ctx)
	if err != nil {
		log.Fatal(err)
	}

	srv := server.New(gcs)
	if err := http.ListenAndServe(":8080", srv); err != nil {
		log.Fatal(err)
	}
}
