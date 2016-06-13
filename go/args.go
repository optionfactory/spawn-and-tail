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
	flag.Var(&files, "f", "Files")
	outPrefixFormat := flag.String("out", "[out] ", "Stdout prefix format string")
	errPrefixFormat := flag.String("err", "[err] ", "Stderr prefix format string")
	commonPrefix := flag.String("p", "", "Common prefix")
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
		StdErrPrefixFormat: *errPrefixFormat,
		CommonPrefixFormat: *commonPrefix,
		Command:            args[0],
		CommandArgs:        args[1:],
	}
}
