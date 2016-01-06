package main

import (
	"flag"
)

func main() {
	flag.Parse()

	switch fmod {
	default:
	case "query":
		switch fwhat {
		case "mr":
			mr()
		case "name":
		default:
			query()
		}
	case "import":
		parse()
	}

}

var (
fmod, fwhat, fsearch, fimport string
)

func init() {
	flag.StringVar(&fimport, "f", "xml/tol.xml", "Input file path")
	flag.StringVar(&fmod, "m", "query", "exe mode import Db/query Db")
	flag.StringVar(&fwhat, "w", "mr", "what is function to perform query")
	flag.StringVar(&fsearch, "s", "Homopterus", "word to search Db")
}