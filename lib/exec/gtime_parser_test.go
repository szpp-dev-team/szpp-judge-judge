package exec

import (
	"strings"
	"testing"

	"github.com/szpp-dev-team/szpp-judge-judge/lib/test"
)

func TestParseGnuTimeOutput(t *testing.T) {
	r := strings.NewReader(testdata)
	res, err := ParseGnuTimeOutput(r)
	if err != nil {
		t.Fatal(err)
	}
	test.AssertNeq(t, res.MaximumResidentSetSize, 0)
}

var testdata = `
Command being timed: "bash ./test/infinity.sh"
User time (seconds): 44.94
System time (seconds): 0.18
Percent of CPU this job got: 99%
Elapsed (wall clock) time (h:mm:ss or m:ss): 0:45.18
Average shared text size (kbytes): 0
Average unshared data size (kbytes): 0
Average stack size (kbytes): 0
Average total size (kbytes): 0
Maximum resident set size (kbytes): 2208
Average resident set size (kbytes): 0
Major (requiring I/O) page faults: 0
Minor (reclaiming a frame) page faults: 210
Voluntary context switches: 3
Involuntary context switches: 1856
Swaps: 0
File system inputs: 0
File system outputs: 0
Socket messages sent: 0
Socket messages received: 0
Signals delivered: 0
Page size (bytes): 16384
Exit status: 0
`
