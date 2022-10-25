package server

import (
	"context"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/szpp-dev-team/szpp-judge-judge/lib/exec"
	"github.com/szpp-dev-team/szpp-judge-judge/model"
	"github.com/szpp-dev-team/szpp-judge-judge/proglang"
)

func (srv *Server) HandleJudgeRequest(judgeReq *model.JudgeRequest) (*model.JudgeResponse, error) {
	// GCSの準備
	ctx := context.Background()
	bkt := srv.gcs.Bucket("szpp-judge")

	// tmp directory 作成
	submitsDir := filepath.Join("tmp", "submits", judgeReq.SubmitID)
	os.Chmod(submitsDir, os.ModePerm)

	err := os.MkdirAll(submitsDir, os.ModePerm)
	if err != nil {
		return nil, err
	}

	testCasesDir := filepath.Join("tmp", "test-cases", judgeReq.SubmitID)
	os.Chmod(testCasesDir, os.ModePerm)

	err = os.MkdirAll(filepath.Join(testCasesDir, "in"), os.ModePerm)
	if err != nil {
		return nil, err
	}
	err = os.MkdirAll(filepath.Join(testCasesDir, "out"), os.ModePerm)
	if err != nil {
		return nil, err
	}

	// ソースコードをGCPから取得
	err = saveGCSContentAsFile(ctx, bkt, filepath.Join("submits", judgeReq.SubmitID), filepath.Join(submitsDir, "Main.cpp"))
	if err != nil {
		return nil, err
	}

	// テストケースをGCSから取得
	testCaseOut := [][]byte{}
	for _, testCaseID := range judgeReq.TestcaseIDs {
		err = saveGCSContentAsFile(ctx, bkt, filepath.Join("testcases", judgeReq.SubmitID, "in", testCaseID), filepath.Join(testCasesDir, "in", testCaseID))
		if err != nil {
			return nil, err
		}

		tmp, err := getGCSContentAsBytes(ctx, bkt, filepath.Join("testcases", judgeReq.SubmitID, "out", testCaseID))
		if err != nil {
			return nil, err
		}
		testCaseOut = append(testCaseOut, tmp)
	}

	// ソースコードをコンパイルする
	cmd := proglang.NewCommand(judgeReq.LanguageID, submitsDir)
	result, err := exec.RunCommand(cmd.CompileCommand, submitsDir)
	if err != nil {
		return nil, err
	}

	// ソースコードを全てのテストケースに対して実行する
	var execResult []*exec.Result
	for _, testCaseID := range judgeReq.TestcaseIDs {
		execCmd := cmd.ExecuteCommand + "  <" + filepath.Join(testCasesDir, "in", testCaseID)
		result, err = exec.RunCommand(execCmd, submitsDir)
		execResult = append(execResult, result)
	}

	// 判定してレスポンスを返す。
	resp := *makeResp(judgeReq.TestcaseIDs, execResult, testCaseOut)

	return &resp, nil
}

func makeResp(testCaseIDs []string, execResult []*exec.Result, correctAns [][]byte) *model.JudgeResponse {
	var ans model.JudgeResponse
	ans.TestcaseResults = make([]model.TestcaseResult, len(execResult))

	ans.Status = model.StatusAC

	for i, r := range execResult {
		var tcr model.TestcaseResult
		tcr.ID = testCaseIDs[i]
		tcr.ExecutionMemory = int64(r.ExecutionMemory)
		tcr.ExecutionTime = r.ExecutionTime.Milliseconds()

		if !(r.Success) {
			tcr.Status = model.StatusCE
			ans.Status = model.StatusCE
			ans.CompileMessage = &r.Stderr
		} else if r.Stderr != "" {
			tcr.Status = model.StatusRE
			ans.Status = model.StatusCE
			ans.ErrorMessage = &r.Stderr
		} else if r.ExecutionTime.Milliseconds() > 2000 {
			tcr.Status = model.StatusTLE
			ans.Status = model.StatusTLE
		} else if r.ExecutionMemory > 1024*100 {
			tcr.Status = model.StatusMLE
			ans.Status = model.StatusMLE
		} else if false {
			tcr.Status = model.StatusOLE
			ans.Status = model.StatusOLE
		} else {
			userAns := strings.Fields(r.Stdout)
			correct := strings.Fields(string(correctAns[i]))
			if reflect.DeepEqual(userAns, correct) {
				tcr.Status = model.StatusAC
			} else {
				tcr.Status = model.StatusWA
				ans.Status = model.StatusWA
			}
		}

		ans.TestcaseResults[i] = tcr
	}

	return &ans
}
