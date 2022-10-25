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

	testCaseIDs := []string{"sample1.txt", "sample2.txt", "sample3.txt", "sample4.txt", "sample5.txt", "sample6.txt",
							"sample7.txt", "sample8.txt", "sample9.txt","sample10.txt", "sample11.txt", "sample12.txt"}

	fmt.Println("######## AC ########")	
	req := &model.JudgeRequest{SubmitID: "test", TaskID: "test", LanguageID: "cpp", TestcaseIDs: testCaseIDs}
	res, err := srv.HandleJudgeRequest(req)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(res)
	fmt.Println("######## AC ########")

	fmt.Println("######## WA ########")
	req = &model.JudgeRequest{SubmitID: "test-wa", TaskID: "test", LanguageID: "cpp", TestcaseIDs: testCaseIDs}
	res, err = srv.HandleJudgeRequest(req)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(res)
	fmt.Println("######## WA ########")

	// これ失敗する
	fmt.Println("######## MLE ########")
	req = &model.JudgeRequest{SubmitID: "test-mle", TaskID: "test", LanguageID: "cpp", TestcaseIDs: testCaseIDs}
	res, err = srv.HandleJudgeRequest(req)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(res)
	fmt.Println("######## MLE ########")

	// これ失敗する
	fmt.Println("######## TLE ########")
	req = &model.JudgeRequest{SubmitID: "test-tle", TaskID: "test", LanguageID: "cpp", TestcaseIDs: testCaseIDs}
	res, err = srv.HandleJudgeRequest(req)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(res)
	fmt.Println("######## TLE ########")

	// これ失敗する
	fmt.Println("######## CE ########")
	req = &model.JudgeRequest{SubmitID: "test-ce", TaskID: "test", LanguageID: "cpp", TestcaseIDs: testCaseIDs}
	res, err = srv.HandleJudgeRequest(req)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(res)
	fmt.Println("######## CE ########")
}