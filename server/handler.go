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
	tmpDirPath := "../tmp"
	submitsDir := filepath.Join(tmpDirPath, "submits", judgeReq.SubmitID)
	err := os.MkdirAll(submitsDir, os.ModePerm)
	if err != nil {
		fmt.Println("fail to make submit dir")
		return nil, err
	}

	testCasesDir := filepath.Join(tmpDirPath, "test-cases", judgeReq.SubmitID)
	err = os.MkdirAll(testCasesDir, os.ModePerm)
	if err != nil {
		fmt.Println("fail make test-cases dir")
		return nil, err
	}

	// GCSからソースコード・テストケースを取得
	err = saveGCSContentAsFile(ctx, bkt, filepath.Join("submits", judgeReq.SubmitID), filepath.Join(submitsDir, "Main.cpp"))
	if err != nil {
		fmt.Println("fail to get submit")
		return nil, err
	}

	correctAns := [][]byte{}
	for _, testCaseID := range judgeReq.TestcaseIDs {
		err = saveGCSContentAsFile(ctx, bkt, filepath.Join("testcases", judgeReq.TaskID, "in", testCaseID), filepath.Join(testCasesDir, testCaseID))
		if err != nil {
			fmt.Println("fail to get testcase in")
			return nil, err
		}

		tmp, err := getGCSContentAsBytes(ctx, bkt, filepath.Join("testcases", judgeReq.TaskID, "out", testCaseID))
		if err != nil {
			fmt.Println("fail to get testcase out")
			return nil, err
		}
		correctAns = append(correctAns, tmp)
	}

	// ソースコードをコンパイルする
	cmd := proglang.NewCommand(judgeReq.LanguageID, submitsDir)
	result, err := exec.RunCommand(cmd.CompileCommand, submitsDir, exec.OptTimeLimit(60*time.Second))
	if err != nil {
		fmt.Println("fail to compile")
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
		execCmd := cmd.ExecuteCommand + "  <" + filepath.Join(testCasesDir, testCaseID)
		result, err = exec.RunCommand(execCmd, submitsDir, exec.OptTimeLimit(3*time.Second))
		if err != nil {
			fmt.Println("fail to execute")
			return nil, err
		}
		execResult = append(execResult, result)
	}

	// 判定してレスポンスを返す
	resp := *makeResp(judgeReq.TestcaseIDs, execResult, correctAns)

	return &resp, nil
}
