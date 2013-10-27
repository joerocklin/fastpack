package main

import (
	"flag"
	"log"
	"os"
)

var flag_help bool
var flag_list bool
var flag_add bool

var archive_filename string
var files []string
var logger = log.New(os.Stderr, "fastpack", log.Flags())
var tempdir string

func init() {
	flag.BoolVar(&flag_help, "h", false, "help")
	flag.BoolVar(&flag_list, "l", false, "List the archive contents")
	flag.BoolVar(&flag_add, "a", false, "Add one or more files to the archive")
	flag.Parse()

	args := flag.Args()

	if len(args) > 0 {
		archive_filename = args[0]
	}

	if len(args) > 1 {
		files = args[1:]
	}

	if len(files) == 0 && flag_add {
		flag_help = true
	}
}

func main() {
	if flag_help {
		flag.Usage()
		return
	}

	if flag_add {
		log.Println("Adding files")
		cmd_add_files(archive_filename, files)
	} else if flag_list {
		log.Println("Listing archive")
		cmd_list_files(archive_filename)
	}

}
