package server

import (
	"context"
	"fmt"
	"io"
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
	for _, testCaseID := range judgeReq.TestcaseIDs {
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

		file, err = os.Create(filepath.Join(testCasesDir, "out", testCaseID))
		if err != nil {
			return nil, err
		}

		_, err = io.Copy(file, r)
		if err != nil {
			return nil, err
		}
		file.Close()
	}

	// ソースコードをコンパイルする
	cmd := proglang.NewCommand(judgeReq.LanguageID, submitsDir)
	result, err := exec.RunCommand(cmd.CompileCommand, submitsDir)
	if err != nil {
		return nil, err
	}
	fmt.Println(result.ExecutionTime)

	// ソースコードを全てのテストケースに対して実行する
	var results []*exec.Result
	for _, testCaseID := range judgeReq.TestcaseIDs {
		execCmd := cmd.ExecuteCommand + "  <" + filepath.Join(testCasesDir, "in", testCaseID)
		result, err = exec.RunCommand(execCmd, submitsDir)
		results = append(results, result)
	}
	for i, row := range results {
		fmt.Println(strconv.Itoa(i) + ": " + row.Stdout)
	}

	// レスポンスを返す

	return nil, nil
}