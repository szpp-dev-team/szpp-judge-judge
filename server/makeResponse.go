package server

import (
	"reflect"
	"strings"

	"github.com/szpp-dev-team/szpp-judge-judge/lib/exec"
	"github.com/szpp-dev-team/szpp-judge-judge/model"
)

func makeResp(testCaseIDs []string, execResults []*exec.Result, correctAns [][]byte) *model.JudgeResponse {
	var resp model.JudgeResponse
	resp.TestcaseResults = make([]model.TestcaseResult, len(execResults))

	resp.Status = model.StatusAC

	for i, result := range execResults {
		var tcr model.TestcaseResult
		tcr.ID = testCaseIDs[i]
		tcr.ExecutionMemory = int64(result.ExecutionMemory)
		tcr.ExecutionTime = result.ExecutionTime.Milliseconds()

		if !result.Success { // RE
			tcr.Status = model.StatusRE
			resp.Status = model.StatusRE
			resp.ErrorMessage = &result.Stderr
		} else if result.ExecutionTime.Milliseconds() > 2000 { // TLE
			tcr.Status = model.StatusTLE
			resp.Status = model.StatusTLE
		} else if result.ExecutionMemory > 128*100 { // MLE
			tcr.Status = model.StatusMLE
			resp.Status = model.StatusMLE
		} else if false { // OLE
			tcr.Status = model.StatusOLE
			resp.Status = model.StatusOLE
		} else { // AC or WA
			userAns := strings.Fields(result.Stdout)
			correct := strings.Fields(string(correctAns[i]))
			if reflect.DeepEqual(userAns, correct) {
				tcr.Status = model.StatusAC
			} else {
				tcr.Status = model.StatusWA
				resp.Status = model.StatusWA
			}
		}

		resp.TestcaseResults[i] = tcr
	}

	return &resp
}

func makeCEresp(compileMessage string) *model.JudgeResponse {
	var ans model.JudgeResponse
	ans.Status = model.StatusCE
	ans.CompileMessage = &compileMessage
	return &ans
}
