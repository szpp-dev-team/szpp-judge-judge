package server

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
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
	tmpDirPath := filepath.Join("../tmp", strconv.Itoa(judgeReq.SubmitID))
	if err := os.MkdirAll(filepath.Join(tmpDirPath, "test-cases"), os.ModePerm); err != nil {
		fmt.Println("Error: make directory")
		return nil, err
	}

	// GCSからソースコードを取得
	if err := saveGCSContentAsFile(ctx, bkt, filepath.Join("submits", strconv.Itoa(judgeReq.SubmitID)), filepath.Join(tmpDirPath, "Main.cpp")); err != nil {
		fmt.Println("Error: get src code from gcs")
		return nil, err
	}

	// GCSからテストケースを取得
	correctAns := [][]byte{}
	for _, testCase := range judgeReq.Testcases {
		testCaseName := testCase.Name
		if err := saveGCSContentAsFile(ctx, bkt, filepath.Join("testcases", strconv.Itoa(judgeReq.TaskID), "in", testCaseName), filepath.Join(tmpDirPath, "test-cases", testCaseName)); err != nil {
			fmt.Println("Error: get testcase (in) from gcs")
			return nil, err
		}

		tmp, err := getGCSContentAsBytes(ctx, bkt, filepath.Join("testcases", strconv.Itoa(judgeReq.TaskID), "out", testCaseName))
		if err != nil {
			fmt.Println("Error: get testcase (out) from gcs")
			return nil, err
		}
		correctAns = append(correctAns, tmp)
	}

	// ソースコードをコンパイルする
	cmd := proglang.NewCommand(judgeReq.LanguageID, tmpDirPath)
	fmt.Println(cmd.CompileCommand)
	result, err := exec.RunCommand(cmd.CompileCommand, tmpDirPath, exec.OptTimeLimit(60*time.Second))
	if err != nil {
		fmt.Println("Error: fail to compile")
		return nil, err
	}
	// コンパイル失敗してたらCEを返す
	if !result.Success {
		ans := makeCEresp(result.Stderr)
		return ans, nil
	}

	// ソースコードを全てのテストケースに対して実行する
	var execResult []*exec.Result
	for _, testCase := range judgeReq.Testcases {
		testCaseName := testCase.Name
		execCmd := cmd.ExecuteCommand + "  <" + filepath.Join("test-cases", testCaseName)
		result, err = exec.RunCommand(execCmd, tmpDirPath, exec.OptTimeLimit(3*time.Second))
		if err != nil {
			fmt.Println("Error: get testcase in from gcs")
			return nil, err
		}
		execResult = append(execResult, result)
	}

	// 判定してレスポンスを返す
	var testCaseIDs []int
	for _, row := range judgeReq.Testcases {
		testCaseIDs = append(testCaseIDs, row.ID)
	}
	resp, err := makeResp(testCaseIDs, execResult, correctAns)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
