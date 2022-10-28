package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"cloud.google.com/go/storage"
	"github.com/szpp-dev-team/szpp-judge-judge/server"
)

func main() {
	ctx := context.Background()
	gcs, err := storage.NewClient(ctx)
	if err != nil {
		log.Fatal(err)
	}

	port := "8001"
	if os.Getenv("PORT") != "" {
		port = os.Getenv("PORT")
	}

	srv := server.New(gcs)
	if err := http.ListenAndServe(":"+port, srv); err != nil {
		log.Fatal(err)
	}
}
