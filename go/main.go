package main

import (
	"bufio"
	"fmt"
	"github.com/hpcloud/tail"
	"io"
	"os"
	"os/exec"
	"runtime"
)

var version string

type FileToTail struct {
	PrefixFormat string
	Path         string
}

type Conf struct {
	FilesToTail        []FileToTail
	StdOutPrefixFormat string
	StdOutSuppression  bool
	StdErrPrefixFormat string
	StdErrSuppression  bool
	CommonPrefixFormat string
	Command            string
	CommandArgs        []string
}

func logFromPipe(parentLogger *Log, pipe io.ReadCloser, prefix string) {
	scanner := bufio.NewScanner(pipe)
	logger := parentLogger.Child(prefix)
	for scanner.Scan() {
		logger.Log(scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		logger.Log(fmt.Sprintf("error scanning pipe: %v", err))
	}
}

func logFromFile(parentLogger *Log, f FileToTail) {
	logger := parentLogger.Child(f.PrefixFormat)
	t, err := tail.TailFile(f.Path, tail.Config{Follow: true, ReOpen: true, Logger: tail.DiscardingLogger})
	if err != nil {
		logger.Log(fmt.Sprintf("error tailing file: %v", err))
	}
	for line := range t.Lines {
		logger.Log(line.Text)
	}
}

func main() {

	conf := ParseConf()
	logger := Root(conf.CommonPrefixFormat)
	logger.Log(fmt.Sprintf("starting %s version %v. max procs:%v", os.Args[0], version, runtime.GOMAXPROCS(0)))

	cmd := exec.Command(conf.Command, conf.CommandArgs...)
	if !conf.StdOutSuppression {
		cmd.Stdin = os.Stdin
		stdout, err := cmd.StdoutPipe()
		if err != nil {
			panic(err)
		}
		go logFromPipe(logger, stdout, conf.StdOutPrefixFormat)
	}
	if !conf.StdErrSuppression {
		stderr, err := cmd.StderrPipe()
		if err != nil {
			panic(err)
		}
		go logFromPipe(logger, stderr, conf.StdErrPrefixFormat)
	}

	for _, file := range conf.FilesToTail {
		go logFromFile(logger, file)
	}

	if err := cmd.Start(); err != nil {
		panic(err)
	}

	if err := cmd.Wait(); err != nil {
		panic(err)
	}

}
