package server

import (
	"context"
	"io"
	"io/ioutil"
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
	obj := bkt.Object(filepath.Join("submits", judgeReq.SubmitID))

	r, err := obj.NewReader(ctx)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	file, err := os.Create(filepath.Join(submitsDir, "Main.cpp"))
	if err != nil {
		return nil, err
	}
	defer file.Close()

	_, err = io.Copy(file, r)
	if err != nil {
		return nil, err
	}

	// テストケースをGCSから取得
	testCaseOut := [][]byte{}
	for i, testCaseID := range judgeReq.TestcaseIDs {
		obj = bkt.Object(filepath.Join("testcases", judgeReq.SubmitID, "in", testCaseID))
		r, err = obj.NewReader(ctx)
		if err != nil {
			return nil, err
		}

		file, err = os.Create(filepath.Join(testCasesDir, "in", testCaseID))
		if err != nil {
			return nil, err
		}

		_, err = io.Copy(file, r)
		if err != nil {
			return nil, err
		}
		file.Close()

		obj = bkt.Object(filepath.Join("testcases", judgeReq.SubmitID, "out", testCaseID))
		r, err = obj.NewReader(ctx)
		if err != nil {
			return nil, err
		}

		out, err := ioutil.ReadAll(r)
		if err != nil {
			return nil, err
		}
		testCaseOut = append(testCaseOut, []byte(""))
		testCaseOut[i] = out
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
	var resp model.JudgeResponse
	resp.TestcaseResults = make([]model.TestcaseResult, len(execResult))
	for i, row := range judgeReq.TestcaseIDs {
		var testCaseResult model.TestcaseResult
		testCaseResult.ID = row
		resp.TestcaseResults[i] = testCaseResult
		resp.TestcaseResults[i].ExecutionMemory = int64(result.ExecutionMemory)
		resp.TestcaseResults[i].ExecutionTime   = int64(result.ExecutionTime)
	}

	resp.Status = model.StatusAC
	for i, row := range execResult {
		if !(row.Success) {
			resp.TestcaseResults[i].Status = model.StatusCE
			resp.Status = model.StatusCE
			resp.CompileMessage = &result.Stderr
		} else if row.Stderr != "" {
			resp.TestcaseResults[i].Status = model.StatusRE
			resp.TestcaseResults[i].Status = model.StatusRE
			resp.ErrorMessage = &result.Stderr
		} else if row.ExecutionTime.Milliseconds() > 2000 {
			resp.TestcaseResults[i].Status = model.StatusTLE
			resp.Status = model.StatusTLE
		} else if row.ExecutionMemory > 1024*1000 {
			resp.TestcaseResults[i].Status = model.StatusMLE
			resp.Status = model.StatusMLE
		} else if false { // TODO oleの条件
			resp.TestcaseResults[i].Status = model.StatusOLE
			resp.Status = model.StatusOLE
		} else {
			userAns := strings.Fields(row.Stdout)
			correctAns := strings.Fields(string(testCaseOut[i]))
			if reflect.DeepEqual(userAns, correctAns) {
				resp.TestcaseResults[i].Status = model.StatusAC
			} else {
				resp.TestcaseResults[i].Status = model.StatusWA
				resp.Status = model.StatusWA
			}
		}
	}
	return &resp, nil
}
