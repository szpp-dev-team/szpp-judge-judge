package server

import (
	"context"
	"io"
	"os"
	"path/filepath"

	"github.com/szpp-dev-team/szpp-judge-judge/model"
)

func (srv *Server) HandleJudgeRequest(judgeReq *model.JudgeRequest) (*model.JudgeResponse, error) {
	// GCSの準備
	ctx := context.Background()
	bkt := srv.gcs.Bucket("szpp-judge")

	// tmp directory 作成
	tmpDir := filepath.Join("tmp", "submits", judgeReq.SubmitID)
	os.Chmod(tmpDir, os.ModePerm)

	err := os.MkdirAll(tmpDir, os.ModePerm)
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

	file, err := os.Create(filepath.Join(tmpDir, "Main.cpp"))
	if err != nil {
		return nil, err
	}
	defer file.Close()

	_, err = io.Copy(file, r)
	if err != nil {
		return nil, err
	}

	// テストケースをGCSから取得

	// ソースコードをコンパイルする

	// ソースコードを全てのテストケースに対して実行する

	// レスポンスを返す

	return nil, nil
}