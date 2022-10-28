package server

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/szpp-dev-team/szpp-judge-judge/lib/exec"
	"github.com/szpp-dev-team/szpp-judge-judge/model"
	"github.com/szpp-dev-team/szpp-judge-judge/proglang"
)

func (srv *Server) HandleJudgeRequest(judgeReq *model.JudgeRequest) (*model.JudgeResponse, error) {
	// GCSの準備
	ctx := context.Background()
	bkt := srv.gcs.Bucket("szpp-judge")

	// tmp directory 作成
	tmpDirPath := filepath.Join("../tmp", judgeReq.SubmitID)
	if err := os.MkdirAll(tmpDirPath, os.ModePerm); err != nil {
		return nil, err
	}
	if err := os.MkdirAll(filepath.Join(tmpDirPath, "test-cases"), os.ModePerm); err != nil {
		return nil, err
	}

	// GCSからソースコード・テストケースを取得
	if err := saveGCSContentAsFile(ctx, bkt, filepath.Join("submits", judgeReq.SubmitID), filepath.Join(tmpDirPath, "Main.cpp")); err != nil {
		return nil, err
	}

	correctAns := [][]byte{}
	for _, testCaseID := range judgeReq.TestcaseIDs {
		if err := saveGCSContentAsFile(ctx, bkt, filepath.Join("testcases", judgeReq.TaskID, "in", testCaseID), filepath.Join(tmpDirPath, "test-cases", testCaseID)); err != nil {
			return nil, err
		}

		tmp, err := getGCSContentAsBytes(ctx, bkt, filepath.Join("testcases", judgeReq.TaskID, "out", testCaseID))
		if err != nil {
			return nil, err
		}
		correctAns = append(correctAns, tmp)
	}

	// ソースコードをコンパイルする
	cmd := proglang.NewCommand(judgeReq.LanguageID, tmpDirPath)
	fmt.Println(cmd.CompileCommand)
	result, err := exec.RunCommand(cmd.CompileCommand, tmpDirPath, exec.OptTimeLimit(60*time.Second))
	if err != nil {
		return nil, err
	}
	// コンパイル失敗してたらCEを返す
	if !(result.Success) {
		ans := makeCEresp(result.Stderr)
		return ans, nil
	}

	// ソースコードを全てのテストケースに対して実行する
	var execResult []*exec.Result
	for _, testCaseID := range judgeReq.TestcaseIDs {
		execCmd := cmd.ExecuteCommand + "  <" + filepath.Join("test-cases", testCaseID)
		result, err = exec.RunCommand(execCmd, tmpDirPath, exec.OptTimeLimit(3*time.Second))
		if err != nil {
			return nil, err
		}
		execResult = append(execResult, result)
	}

	// 判定してレスポンスを返す
	resp := *makeResp(judgeReq.TestcaseIDs, execResult, correctAns)

	return &resp, nil
}
