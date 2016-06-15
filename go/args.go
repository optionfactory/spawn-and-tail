package main

import (
	"flag"
	"fmt"
	"os"
	"regexp"
	"strings"
)

var loggerNameRe = regexp.MustCompile(`(.*/)?(.*)(\..*?)?`)

type FilesToTail []FileToTail

func (self *FilesToTail) String() string {
	return fmt.Sprintf("%v", *self)
}

func (self *FilesToTail) Set(value string) error {

	split := strings.SplitN(value, ":", 2)

	path := split[0]
	pf := ""

	if len(split) == 2 {
		path = split[1]
		pf = split[0]
	} else {
		pf = loggerNameRe.FindStringSubmatch(split[0])[2]
	}

	*self = append(*self, FileToTail{
		PrefixFormat: pf,
		Path:         path,
	})
	return nil
}

func ParseConf() Conf {
	var files FilesToTail
	flag.Var(&files, "file", "File. can be repeated multiple times")
	commonPrefix := flag.String("prefix", "", "Common prefix. Defaults to command name")
	outPrefixFormat := flag.String("prefix-stdout", "[out] ", "Stdout prefix format string")
	errPrefixFormat := flag.String("prefix-stderr", "[err] ", "Stderr prefix format string")
	outSuppression := flag.Bool("suppress-stdout", false, "Standard output suppression")
	errSuppression := flag.Bool("suppress-stderr", false, "Standard error suppression")
	flag.Parse()
	args := flag.Args()
	if len(args) == 0 {
		flag.Usage()
		os.Exit(1)
	}
	if *commonPrefix == "" {
		defaultPrefix := fmt.Sprintf("[%s]", args[0])
		commonPrefix = &defaultPrefix
	}

	return Conf{
		FilesToTail:        files,
		StdOutPrefixFormat: *outPrefixFormat,
		StdOutSuppression:  *outSuppression,
		StdErrPrefixFormat: *errPrefixFormat,
		StdErrSuppression:  *errSuppression,
		CommonPrefixFormat: *commonPrefix,
		Command:            args[0],
		CommandArgs:        args[1:],
	}
}
