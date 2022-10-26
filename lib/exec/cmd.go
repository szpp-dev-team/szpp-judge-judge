package exec

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"os/exec"
	pkgexec "os/exec"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
)

const (
	TmpRootDirPath       = "./tmp/exec"
	GnuTimeStdoutBufSize = 128_000    // 128KB
	StdoutSizeLimit      = 10_000_000 // 10MB
	StderrSizeLimit      = 10_000_000 // 10MB
)

type Result struct {
	Success         bool
	ExecutionTime   time.Duration
	ExecutionMemory int
	Stdout          string
	Stderr          string
}

func RunCommand(command string, tmpDirPath string, optFuncs ...OptionFunc) (*Result, error) {
	// コマンドのバリデーション
	tokens := strings.Fields(command)
	if len(tokens) == 0 {
		return nil, errors.New("the length of command must not be 0")
	}

	// オプション
	opt := DefaultOption()
	for _, optFunc := range optFuncs {
		optFunc(opt)
	}

	// Command 構造体の build
	cmd := pkgexec.Command(GnuTimeCommandPath, "-v", "sh", "-c", command)
	gtimeStderrBuf := bytes.NewBuffer(make([]byte, GnuTimeStdoutBufSize))
	cmd.Stderr = gtimeStderrBuf

	// コマンド実行
	if err := cmd.Start(); err != nil {
		fmt.Println("check cmd start")
		return nil, err
	}
	tc := time.NewTicker(opt.TimeLimit) // TimeLimit の時間が経ったら chan を send する
	beginTime := time.Now()             // 計測開始
	pid := cmd.Process.Pid

	// 並列処理でコマンド終了を監視する
	cmdExitChan := make(chan error, 1)
	go func() {
		cmdExitChan <- cmd.Wait()
	}()

	// TimeLimit とコマンド終了の速い方を選択する
	select {
	case <-tc.C:
		log.Println("time limit exceeed")
		if err := killChildProcesses(pid); err != nil {
			return nil, err
		}

	case err := <-cmdExitChan:
		log.Println("exited")
		if err != nil {
			exitError := &pkgexec.ExitError{}
			if !errors.As(err, &exitError) {
				return nil, err
			}
		}
	}

	// 出力等の read
	exectionTime := time.Since(beginTime)
	gtimeRes, err := ParseGnuTimeOutput(gtimeStderrBuf)
	if err != nil {
		return nil, err
	}

	stdoutBytes, err := readFileFull(path.Join(tmpDirPath, "stdout.txt"), StdoutSizeLimit)
	if err != nil {
		// コンパイル時は stdout.txt は生成されないため、ErrNotExist は無視する
		if errors.Is(err, os.ErrNotExist) {
			stdoutBytes = make([]byte, 0)
		} else {
			return nil, err
		}
	}
	stderrBytes, err := readFileFull(path.Join(tmpDirPath, "stderr.txt"), StderrSizeLimit)
	if err != nil {
		// コンパイル時は stderr.txt は生成されないため、ErrNotExist は無視する
		if errors.Is(err, os.ErrNotExist) {
			stderrBytes = make([]byte, 0)
		} else {
			return nil, err
		}
	}

	var success bool
	if cmd.ProcessState != nil {
		success = cmd.ProcessState.Success()
	}

	return &Result{
		Success:         success,
		ExecutionTime:   exectionTime,
		ExecutionMemory: gtimeRes.MaximumResidentSetSize,
		Stdout:          string(stdoutBytes), // TODO: メモリコピーが走るので、unsafe なりを使いたい
		Stderr:          string(stderrBytes),
	}, nil
}

type Option struct {
	TimeLimit   time.Duration
	MemoryLimit int64 // KB
}

type OptionFunc func(*Option)

func DefaultOption() *Option {
	return &Option{
		TimeLimit: 20 * time.Second,
	}
}

func OptTimeLimit(limit time.Duration) OptionFunc {
	return func(o *Option) {
		o.TimeLimit = limit
	}
}

var ErrTooBigFile = errors.New("the specified file is too big to open")

func readFileFull(filename string, limit int) ([]byte, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	fileInfo, err := f.Stat()
	if err != nil {
		return nil, err
	}
	fileSizeInt64 := fileInfo.Size()
	if fileSizeInt64 > math.MaxInt {
		return nil, ErrTooBigFile
	}
	fileSizeInt := int(fileSizeInt64)
	if fileSizeInt < limit {
		limit = fileSizeInt
	}

	buf := make([]byte, limit)
	if _, err := io.ReadFull(f, buf); err != nil {
		return nil, err
	}

	return buf, nil
}

func killChildProcesses(parentPid int) error {
	fmt.Println("KillChildProcess start. pid :" + strconv.Itoa(parentPid))
	stdoutBuf := &bytes.Buffer{}
	pgrepCmd := pkgexec.Command("pgrep", "-P", strconv.Itoa(parentPid))
	pgrepCmd.Stdout = stdoutBuf
	if err := pgrepCmd.Run(); err != nil {
		exitCode := err.(*exec.ExitError).ProcessState.ExitCode()
		if exitCode == 1 {
			// parentPid のプロセスに子プロセスがなかった
			return nil
		} else {
			return err
		}
	}

	sc := bufio.NewScanner(stdoutBuf)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		pid, err := strconv.ParseInt(line, 10, 64)
		if err != nil {
			return err
		}

		// 子プロセスを先に kill する
		err = killChildProcesses(int(pid))
		if err != nil {
			return err
		}

		// kill したいプロセスの存在確認
		tmp, err := checkProcessIsExit(int(pid))
		if tmp {
			err = pkgexec.Command("kill", "-9", strconv.Itoa(int(pid))).Run()
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func checkProcessIsExit(pid int) (bool, error) {
	cmd := pkgexec.Command("ps", "-p", strconv.Itoa(pid))
	err := cmd.Run()
	exitCode := cmd.ProcessState.ExitCode()
	if exitCode == 0 {
		return true, nil
	} else if exitCode == 1 {
		return false, nil
	} else {
		return false, err
	}
}
