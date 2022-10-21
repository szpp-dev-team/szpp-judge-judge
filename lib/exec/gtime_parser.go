package exec

import (
	"bufio"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

/*
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
*/

type GnuTimeResult struct {
	CommandBeginTimed          string
	UserTime                   time.Duration
	SystemTime                 time.Duration
	PercentOfCPU               int
	Elapsed                    time.Duration
	AverageSharedTextSize      int
	AverageUnsharedDataSize    int
	AverageStackSize           int
	AverageTotalSize           int
	MaximumResidentSetSize     int
	AverageResidentSetSize     int
	MajorPageFaults            int
	MinorPageFaults            int
	VoluntaryContextSwitches   int
	InvoluntaryContextSwitches int
	Swaps                      int
	FileSystemInputs           int
	FileSystemOutputs          int
	SocketMessagesSent         int
	SocketMessagesReceived     int
	SignalsDelivered           int
	PageSize                   int
	ExitStatus                 int
}

func ParseGnuTimeOutput(r io.Reader) (*GnuTimeResult, error) {
	res := &GnuTimeResult{}

	sc := bufio.NewReader(r)
	for {
		line, _, err := sc.ReadLine()
		if err != nil {
			if err != io.EOF {
				return nil, err
			}
			break
		}
		log.Println(string(line))
		key, value, _ := strings.Cut(strings.TrimSpace(string(line)), ":")
		value = strings.TrimSpace(value)
		switch key {
		case "Command being timed":
			res.CommandBeginTimed = value
		case "User time (seconds)":
			res.UserTime = time.Duration(unwrap(strconv.ParseFloat(value, 64)) * float64(time.Second))
		case "System time (seconds)":
			res.SystemTime = time.Duration(unwrap(strconv.ParseFloat(value, 64)) * float64(time.Second))
		case "Percent of CPU this job got":
			res.PercentOfCPU = unwrap(strconv.Atoi(strings.TrimSuffix(value, "%")))
		case "Elapsed (wall clock) time (h:mm:ss or m:ss)":
			res.Elapsed = unwrap(time.ParseDuration(value))
		case "Average shared text size (kbytes)":
			res.AverageSharedTextSize = unwrap(strconv.Atoi(value))
		case "Average unshared data size (kbytes)":
			res.AverageUnsharedDataSize = unwrap(strconv.Atoi(value))
		case "Average stack size (kbytes)":
			res.AverageStackSize = unwrap(strconv.Atoi(value))
		case "Average total size (kbytes)":
			res.AverageTotalSize = unwrap(strconv.Atoi(value))
		case "Maximum resident set size (kbytes)":
			res.MaximumResidentSetSize = unwrap(strconv.Atoi(value))
		case "Average resident set size (kbytes)":
			res.AverageResidentSetSize = unwrap(strconv.Atoi(value))
		case "Major (requiring I/O) page faults":
			res.MajorPageFaults = unwrap(strconv.Atoi(value))
		case "Minor (reclaiming a frame) page faults":
			res.MinorPageFaults = unwrap(strconv.Atoi(value))
		case "Voluntary context switches":
			res.VoluntaryContextSwitches = unwrap(strconv.Atoi(value))
		case "Involuntary context switches":
			res.InvoluntaryContextSwitches = unwrap(strconv.Atoi(value))
		case "Swaps":
			res.Swaps = unwrap(strconv.Atoi(value))
		case "File system inputs":
			res.FileSystemInputs = unwrap(strconv.Atoi(value))
		case "File system outputs":
			res.FileSystemOutputs = unwrap(strconv.Atoi(value))
		case "Socket messages sent":
			res.SocketMessagesSent = unwrap(strconv.Atoi(value))
		case "Socket messages received":
			res.SocketMessagesReceived = unwrap(strconv.Atoi(value))
		case "Signals delivered":
			res.SignalsDelivered = unwrap(strconv.Atoi(value))
		case "Page size (bytes)":
			res.PageSize = unwrap(strconv.Atoi(value))
		case "Exit status":
			res.ExitStatus = unwrap(strconv.Atoi(value))
		}
	}

	return res, nil
}

func ParseGnuTimeOutputFromFile(path string) (*GnuTimeResult, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return ParseGnuTimeOutput(f)
}

func unwrap[T any](returnValue T, err error) T {
	return returnValue
}
