package server

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"testing"

	"cloud.google.com/go/storage"
	"github.com/szpp-dev-team/szpp-judge-judge/model"
	"google.golang.org/api/option"
)

func TestSubmitAC(t *testing.T) {
	if err := testEasily("0"); err != nil {
		log.Fatal(err)
	}
}

func TestSubmitAllOdd(t *testing.T) {
	if err := testEasily("1"); err != nil {
		log.Fatal(err)
	}
}

func TestSubmitCE(t *testing.T) {
	if err := testEasily("2"); err != nil {
		log.Fatal(err)
	}
}

func TestSubmitMLE2(t *testing.T) {
	if err := testEasily("3"); err != nil {
		log.Fatal(err)
	}
}

func TestSubmitMLE1(t *testing.T) {
	if err := testEasily("4"); err != nil {
		log.Fatal(err)
	}
}

func TestSubmitAllOdd2(t *testing.T) {
	if err := testEasily("5"); err != nil {
		log.Fatal(err)
	}
}

func TestSubmit1(t *testing.T) {
	if err := testEasily("6"); err != nil {
		log.Fatal(err)
	}
}

func TestSubmit2(t *testing.T) {
	if err := testEasily("7"); err != nil {
		log.Fatal(err)
	}
}

func testEasily(submitID string) error {
	var testcases []model.Testcase
	testcases = append(testcases, model.Testcase{ID: 0, Name: "sample01.txt"})
	return testRequest(makeSrv(), submitID, "2", "cpp", testcases)
}

func makeSrv() *Server {
	ctx := context.Background()
	gcs, err := storage.NewClient(ctx, option.WithCredentialsFile("../credentials.json"))
	if err != nil {
		log.Fatal(err)
	}

	return New(gcs)
}

func testRequest(srv *Server, submitID string, taskID string, langID string, testCaseIDs []model.Testcase) error {
	fmt.Println("######## " + submitID + " ########")
	submitIDNum, err := strconv.Atoi(submitID)
	if err != nil {
		return err
	}
	taskIDNum, err := strconv.Atoi(taskID)
	if err != nil {
		return err
	}
	req := &model.JudgeRequest{SubmitID: submitIDNum, TaskID: taskIDNum, LanguageID: langID, Testcases: testCaseIDs}
	res, err := srv.HandleJudgeRequest(req)
	if err != nil {
		return err
	}
	fmt.Println(res)
	fmt.Println("######## " + submitID + " ########")
	return nil
}
