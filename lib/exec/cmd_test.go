package exec

import (
	"testing"
	"time"

	"github.com/szpp-dev-team/szpp-judge-judge/lib/test"
)

func TestRunCommandSuccess(t *testing.T) {
	res, err := RunCommand("ls > ./tmp/exec/1/stdout.txt 2> ./tmp/exec/1/stderr.txt", "./tmp/exec/1")
	if err != nil {
		t.Fatal(err)
	}
	test.AssertNeq(t, res.ExecutionMemory, 0)
	test.Assert(t, res.Success)
	test.Assert(t, len(res.Stdout) != 0)
}

func TestRunCommandFailure(t *testing.T) {
	timeLimit := 5 * time.Second
	res, err := RunCommand("while true; do :; done", "./tmp/exec/1", OptTimeLimit(timeLimit))
	if err != nil {
		t.Fatal(err)
	}
	test.AssertGt(t, res.ExecutionTime, timeLimit)
	test.AssertNeq(t, res.ExecutionMemory, 0)
	test.Assert(t, !res.Success)
}
