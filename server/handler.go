package server

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"

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
	fmt.Println(result.ExecutionTime)

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
	for i, testCaseID := range judgeReq.TestcaseIDs {
		var testCaseResult model.TestcaseResult
		testCaseResult.ID = testCaseID
		resp.TestcaseResults[i] = testCaseResult
	}
	for _, testCaseResult := range resp.TestcaseResults {
		testCaseResult.ExecutionMemory = int64(result.ExecutionMemory)
		testCaseResult.ExecutionTime = int64(result.ExecutionTime)
	}
	resp.Status = model.StatusAC
	for i, row := range resp.TestcaseResults {
		tmp := execResult[i]
		if !(tmp.Success) {
			row.Status = model.StatusCE
			resp.Status = model.StatusCE
			resp.CompileMessage = &result.Stderr
		} else if tmp.Stderr != "" {
			row.Status = model.StatusRE
			resp.Status = model.StatusRE
			resp.ErrorMessage = &result.Stderr
		} else if tmp.ExecutionTime.Milliseconds() > 2000 {
			row.Status = model.StatusTLE
			resp.Status = model.StatusTLE
		} else if tmp.ExecutionMemory > 1024*1000 {
			row.Status = model.StatusMLE
			resp.Status = model.StatusMLE
		} else if false {
			row.Status = model.StatusOLE
			resp.Status = model.StatusOLE
		} else {
			fmt.Println(strconv.Itoa(i))
			fmt.Println([]byte(tmp.Stdout))
			fmt.Println(testCaseOut[i])
			if bytes.Equal([]byte(tmp.Stdout+"\n"), testCaseOut[i]) {
				row.Status = model.StatusAC
			} else {
				row.Status = model.StatusWA
				resp.Status = model.StatusWA
			}
		}
	}
	return &resp, nil
}
