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

	for i, r := range execResults {
		var tcr model.TestcaseResult
		tcr.ID = testCaseIDs[i]
		tcr.ExecutionMemory = int64(r.ExecutionMemory)
		tcr.ExecutionTime = r.ExecutionTime.Milliseconds()

		if !(r.Success) { // CE
			tcr.Status = model.StatusCE
			ans.Status = model.StatusCE
			ans.CompileMessage = &r.Stderr
		} else if r.Stderr != "" { // RE
			tcr.Status = model.StatusRE
			ans.Status = model.StatusCE
			ans.ErrorMessage = &r.Stderr
		} else if r.ExecutionTime.Milliseconds() > 2000 { // TLE
			tcr.Status = model.StatusTLE
			ans.Status = model.StatusTLE
		} else if r.ExecutionMemory > 1024*100 { // MLE
			tcr.Status = model.StatusMLE
			ans.Status = model.StatusMLE
		} else if false { // OLE
			tcr.Status = model.StatusOLE
			ans.Status = model.StatusOLE
		} else { // AC or WA
			userAns := strings.Fields(r.Stdout)
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
