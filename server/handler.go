package server

import (
	"context"
	"os"
	"path/filepath"

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
	err := os.MkdirAll(submitsDir, os.ModePerm)
	if err != nil {
		return nil, err
	}

	testCasesDir := filepath.Join("tmp", "test-cases", judgeReq.SubmitID)
	err = os.MkdirAll(testCasesDir, os.ModePerm)
	if err != nil {
		return nil, err
	}

	// ソースコードをGCPから取得
	err = saveGCSContentAsFile(ctx, bkt, filepath.Join("submits", judgeReq.SubmitID), filepath.Join(submitsDir, "Main.cpp"))
	if err != nil {
		return nil, err
	}

	// テストケースをGCSから取得
	correctAns := [][]byte{}
	for _, testCaseID := range judgeReq.TestcaseIDs {
		err = saveGCSContentAsFile(ctx, bkt, filepath.Join("testcases", judgeReq.SubmitID, "in", testCaseID), filepath.Join(testCasesDir, testCaseID))
		if err != nil {
			return nil, err
		}

		tmp, err := getGCSContentAsBytes(ctx, bkt, filepath.Join("testcases", judgeReq.SubmitID, "out", testCaseID))
		if err != nil {
			return nil, err
		}
		correctAns = append(correctAns, tmp)
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
		execCmd := cmd.ExecuteCommand + "  <" + filepath.Join(testCasesDir, testCaseID)
		result, err = exec.RunCommand(execCmd, submitsDir)
		execResult = append(execResult, result)
	}

	// 判定してレスポンスを返す。
	resp := *makeResp(judgeReq.TestcaseIDs, execResult, correctAns)

	return &resp, nil
}
