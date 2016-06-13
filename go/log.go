package main

import (
	"fmt"
	"sync"
)

type Log struct {
	mu                 *sync.Mutex
	commonPrefixFormat string
	prefixFormat       string
}

func Root(prefix string) *Log {
	return &Log{
		mu:                 &sync.Mutex{},
		commonPrefixFormat: prefix,
	}
}

func (self *Log) Child(prefix string) *Log {
	return &Log{
		self.mu, self.commonPrefixFormat, prefix,
	}
}

func (self *Log) Log(text interface{}) {
	self.mu.Lock()
	defer self.mu.Unlock()
	cp := fmt.Sprintf(self.commonPrefixFormat)
	p := fmt.Sprintf(self.prefixFormat)

	fmt.Printf("%s%s%v\n", cp, p, text)
}
