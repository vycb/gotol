package main

import (
	"github.com/vycb/gotol/Parser"
	"github.com/vycb/gotol/DbClient"
	"github.com/vycb/gotol/QuertCassandra"
	"github.com/vycb/gotol/QueryMongo"
	"flag"
	"fmt"
	"encoding/xml"
	"strings"
	"strconv"
	"os"
)

func main() {
	flag.Parse()

	switch fmod {
	default:
	case "query":

		switch fwhat {
		case "toroot":
			QueryCassandra.Toroot(&fsearch)

		case "mr":

			switch fdbcl {
			default:
			case "mongo":
				QueryMongo.MapReduce(&fsearch, &fwhat)
			case "cassandra":
				QueryCassandra.Query(&fsearch)
			}

		case "name":

			switch fdbcl {
			default:
			case "mongo":
				QueryMongo.Query(&fsearch, &fwhat)
			case "cassandra":
				QueryCassandra.Query(&fsearch)
			}
		}

	case "import":
		Parse(&fimport, &fdbcl)
	}

}

func Parse(fimport *string, fdbcl *string) {
	var dc DbClient.Db

	if *fdbcl == "mongo" {
		dc = &QueryMongo.Mongo{}
	}else{
		dc = &QueryCassandra.Cassandra{}
	}
	dc.Init()
	dc.NewBatch()
	defer dc.SessionClose()

	xmlFile, err := os.Open(*fimport)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer xmlFile.Close()

	decoder := xml.NewDecoder(xmlFile)
	var inEl, ct, pt string
	var node *Parser.Node
	pnode := new(Parser.Node)
	Loop:
	for {
		t, _ := decoder.Token()
		if t == nil {
			break
		}
		switch se := t.(type) {
		case xml.StartElement:
			inEl = se.Name.Local
			pt = ct
			ct = inEl

			if inEl == "NODE" {
				if pt == "NODES" {
					pnode = node
				}
				ID := se.Attr[1]
				id, _ := strconv.Atoi(ID.Value)

				node = &Parser.Node{Id:id, Name:"", Parent: pnode, OtherName:"", Description:""}

			} else if inEl == "NODES" {

				fmt.Println(node.Id, node.Name, node.Parent.Id, node.OtherName, node.Description)
				dc.Save(node);
			}

		case xml.CharData:
			chd := strings.TrimSpace(string(se.Copy()))
			if (chd == "") {
				continue Loop
			}
			if ct == "NAME" && pt == "NODE" {
				node.Name = chd
			}else if ct == "DESCRIPTION" {
				node.Description = chd
			}else if ct == "NAME" && pt == "OTHERNAME" {
				if node.OtherName != "" {
					node.OtherName += ", " + chd
				}else {
					node.OtherName += chd
				}
			}
		case xml.EndElement:
			inEl = se.Name.Local
			if inEl == "NODE" {

				fmt.Println(node.Id, node.Name, node.Parent.Id, node.OtherName, node.Description)
				dc.Save(node);
			}
			if inEl == "NODES" {
				pnode = pnode.Parent
			}
		}
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