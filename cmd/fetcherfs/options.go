package main

import (
	"os"

	"yunion.io/x/log"
	"yunion.io/x/structarg"
)

const BLOCK_SIZE = 8

type Options struct {
	Url        string `help:"destination to fetch content" required:"true"`
	Tmpdir     string `help:"temporary dir save fs content" required:"true"`
	Token      string `help:"authentication information to access given url" required:"true"`
	Blocksize  int    `help:"block size of content file system(MB)"`
	MountPoint string `help:"mount path of fuse fs" required:"true"`
	Debug      bool   `help:"enable debug go fuse"`
}

var opt = &Options{}

func init() {
	// structarg.NewArgumentParser(&BaseOptions{}
	parser, err := structarg.NewArgumentParser(opt, "", "", "")
	if err != nil {
		log.Fatalf("Error define argument parser: %v", err)
	}
	err = parser.ParseArgs2(os.Args[1:], true, true)
	if err != nil {
		log.Fatalf("Failed parse args %s", err)
	}
	log.Errorf("%v", opt)
	if opt.Blocksize <= 0 {
		opt.Blocksize = BLOCK_SIZE
	}
}