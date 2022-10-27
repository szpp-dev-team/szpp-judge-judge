package proglang

import (
	"fmt"
	"path"
)

type Command struct {
	ID             string
	CompileCommand string
	ExecuteCommand string
	Filename       string
}

func NewCommand(langID string, tmpDirPath string) *Command {
	if langID == "c(gcc)" {
		return &Command{
			ID:             "c(gcc)",
			CompileCommand: fmt.Sprintf("gcc -O2 -Wall -Wextra -o Main %s", path.Join(tmpDirPath, "Main.c")),
			ExecuteCommand: fmt.Sprintf("./%s", path.Join(tmpDirPath, "Main")),
		}
	}
	if langID == "cpp" {
		return &Command{
			ID:             "cpp",
			CompileCommand: fmt.Sprintf("g++ -o %s/Main %s >%s 2>%s", path.Join(tmpDirPath), path.Join(tmpDirPath, "Main.cpp"), path.Join(tmpDirPath, "stdout.txt"), path.Join(tmpDirPath, "stderr.txt")),
			ExecuteCommand: fmt.Sprintf("./%s >%s 2>%s", path.Join(tmpDirPath, "Main"), path.Join(tmpDirPath, "stdout.txt"), path.Join(tmpDirPath, "stderr.txt")),
		}
	}

	// unreachable
	return nil
}