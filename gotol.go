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
		case "toroot":
			toroot()

		case "mr":

			switch fdbcl {
			default:
			case "mongo":
				mr()
			case "cassandra":
				querycassandra()
			}

		case "name":

			switch fdbcl {
			default:
			case "mongo":
				querymongo()
			case "cassandra":
				querycassandra()
			}
		}

	case "import":
		parse()
	}

}

var (
	fmod, fdbcl, fwhat, fsearch, fimport string
)

func init() {
	flag.StringVar(&fimport, "f", "xml/tol.xml", "Input file path")
	flag.StringVar(&fmod, "m", "query", "exe mode import Db/query Db")
	flag.StringVar(&fdbcl, "dc", "cassandra", "Db Client type")
	flag.StringVar(&fwhat, "w", "toroot", "what is a function to query")
	flag.StringVar(&fsearch, "s", "Eudorylas", "word to search Db")
}