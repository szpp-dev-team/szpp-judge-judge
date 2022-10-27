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

func TestSubmitAC(t *testing.T) {
	if err := testEasily("test"); err != nil {
		log.Fatal(err)
	}
}

func TestSubmitWA(t *testing.T) {
	if err := testEasily("test-wa"); err != nil {
		log.Fatal(err)
	}
}

func TestSubmitRE(t *testing.T) {
	if err := testEasily("test-mle"); err != nil {
		log.Fatal(err)
	}
}

func TestSubmitTLE(t *testing.T) {
	if err := testEasily("test-tle"); err != nil {
		log.Fatal(err)
	}
}

func TestSubmitCE(t *testing.T) {
	if err := testEasily("test-ce"); err != nil {
		log.Fatal(err)
	}
}

func TestSubmitAllOdd(t *testing.T) {
	if err := testEasily("test-all-odd"); err != nil {
		log.Fatal(err)
	}
}

func testEasily(submitID string) error {
	testCaseIDs := []string{"sample1.txt", "sample2.txt", "sample3.txt", "sample4.txt", "sample5.txt", "sample6.txt",
		"sample7.txt", "sample8.txt", "sample9.txt", "sample10.txt", "sample11.txt", "sample12.txt"}
	return testRequest(makeSrv(), submitID, "test", "cpp", testCaseIDs)
}

func makeSrv() *Server {
	ctx := context.Background()
	gcs, err := storage.NewClient(ctx, option.WithCredentialsFile("../credentials.json"))
	if err != nil {
		log.Fatal(err)
	}

	return New(gcs)
}

func testRequest(srv *Server, submitID string, taskID string, langID string, testCaseIDs []string) error {
	fmt.Println("######## " + submitID + " ########")
	req := &model.JudgeRequest{SubmitID: submitID, TaskID: taskID, LanguageID: langID, TestcaseIDs: testCaseIDs}
	res, err := srv.HandleJudgeRequest(req)
	if err != nil {
		return err
	}
	fmt.Println(res)
	fmt.Println("######## " + req.SubmitID + " ########")
	return nil
}
