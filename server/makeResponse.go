package server

import (
	"reflect"
	"strings"

	"github.com/szpp-dev-team/szpp-judge-judge/lib/exec"
	"github.com/szpp-dev-team/szpp-judge-judge/model"
)

func makeResp(testCaseIDs []string, execResults []*exec.Result, correctAns [][]byte) *model.JudgeResponse {
	var ans model.JudgeResponse
	ans.TestcaseResults = make([]model.TestcaseResult, len(execResults))

	ans.Status = model.StatusAC

	for i, result := range execResults {
		var tcr model.TestcaseResult
		tcr.ID = testCaseIDs[i]
		tcr.ExecutionMemory = int64(result.ExecutionMemory)
		tcr.ExecutionTime = result.ExecutionTime.Milliseconds()

		if result.Stderr != "" { // RE
			tcr.Status = model.StatusRE
			ans.Status = model.StatusRE
			ans.ErrorMessage = &result.Stderr
		} else if result.ExecutionTime.Milliseconds() > 2000 { // TLE
			tcr.Status = model.StatusTLE
			ans.Status = model.StatusTLE
		} else if result.ExecutionMemory > 128*100 { // MLE
			tcr.Status = model.StatusMLE
			ans.Status = model.StatusMLE
		} else if false { // OLE
			tcr.Status = model.StatusOLE
			ans.Status = model.StatusOLE
		} else { // AC or WA
			userAns := strings.Fields(result.Stdout)
			correct := strings.Fields(string(correctAns[i]))
			if reflect.DeepEqual(userAns, correct) {
				tcr.Status = model.StatusAC
			} else {
				tcr.Status = model.StatusWA
				ans.Status = model.StatusWA
			}
		}

		ans.TestcaseResults[i] = tcr
	}

	return &ans
}

func makeCEresp(compileMessage string) *model.JudgeResponse {
	var ans model.JudgeResponse
	ans.Status = model.StatusCE
	ans.CompileMessage = &compileMessage
	return &ans
}
