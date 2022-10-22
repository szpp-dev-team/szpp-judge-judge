package server

import (
	"context"
	"fmt"
	"log"

	"testing"

	"cloud.google.com/go/storage"
	"github.com/szpp-dev-team/szpp-judge-judge/model"
	"google.golang.org/api/option"
)

func TestHandleSubmit(t *testing.T) {
	ctx := context.Background()
	gcs, err := storage.NewClient(ctx, option.WithCredentialsFile("../credentials.json"))
	if err != nil {
		log.Fatal(err)
	}

	srv := New(gcs)

	req := &model.JudgeRequest{SubmitID: "test", TaskID: "test", LanguageID: "cpp", TestcaseIDs: []string{"sample1.txt", "sample2.txt", "sample3.txt"}}

	res, err := srv.HandleJudgeRequest(req)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(res)
}