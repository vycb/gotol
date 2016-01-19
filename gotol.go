package main

import (
	"flag"
	."github.com/vycb/gotol/Parser"
	"github.com/vycb/gotol/QuertCassandra"
	"github.com/vycb/gotol/QueryMongo"
	"github.com/vycb/gotol/QueryRedis"
)

func main() {

	switch fmod {
	default:
	case "query":

		switch fwhat {
		case "toroot":
			QueryCassandra.Toroot(fsearch)

		case "mr":

			switch fdbcl {
			default:
			case "mongo":
			case "cassandra":
				QueryCassandra.Query(fsearch)
			}

		case "name":

			switch fdbcl {
			default:
			case "mongo":
				QueryMongo.Query(fsearch, fwhat)
			case "cassandra":
				QueryCassandra.Query(fsearch)

			case "redis":
				qc := &QueryRedis.QueryClient{}

				switch fqmeth {
				case "lua":
					qc.Query(fsearch)
				case "scan":
					qc.Scan(fsearch)
				}
			}
		}

	case "import":
		Parse(fsearch, fdbcl)
	}

}



var (
	fmod, fdbcl, fwhat, fqmeth, fsearch, fimport string
)

func init() {
	flag.StringVar(&fsearch, "s", "*", "word to search Db")
	flag.StringVar(&fwhat, "w", "name", "mr/name: what is a function to query")
	flag.StringVar(&fqmeth, "qm", "scan", "scan/lua: query method")
	flag.StringVar(&fdbcl, "dc", "redis", "redis/pq/mongo/cassandra: Db Client type")
	flag.StringVar(&fmod, "m", "query", "query/import: exe mode import Db/query Db")
	flag.StringVar(&fimport, "f", "xml/tol.xml", "Input file path")
	flag.Parse()
}
