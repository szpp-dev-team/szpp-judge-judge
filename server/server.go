package server

import (
	"encoding/json"
	"log"
	"net/http"

	"cloud.google.com/go/storage"
	"github.com/szpp-dev-team/szpp-judge-judge/model"
)

type Server struct {
	gcs *storage.Client
}

func New(gcs *storage.Client) *Server {
	return &Server{gcs}
}

func (srv *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var judgeReq model.JudgeRequest
	if err := json.NewDecoder(r.Body).Decode(&judgeReq); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	b, _ := json.Marshal(judgeReq)
	log.Println(string(b))
	judgeResp, err := srv.HandleJudgeRequest(&judgeReq)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	b, _ = json.Marshal(judgeResp)
	log.Println(string(b))
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(judgeResp); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
