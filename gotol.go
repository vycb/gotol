package main

import (
	"log"
	"os"
	"flag"
	"encoding/xml"
	"strings"
	"strconv"
)

var inputFile = flag.String("infile", "xml/tol.xml", "Input file path")

type(
	Node struct {
		Id                           int
		Name, OtherName, Description string
		Parent                       *Node
	}
	MNode struct {
		Id          int `bson:"_id,omitempty" json:"_id"`
		Parent      int `bson:"parent,omitempty" json:"parent"`
		Name        string `bson:"name" json:"name"`
		OtherName   string `bson:"othername" json:"othername"`
		Description string `bson:"description" json:"description"`
	}
/*	XMLID struct {
		ID string `xml:"ID,attr"`
	}
	XmlNode struct {
		ID          XMLID `xml:"NODE"`
	}
	XmlOtherName struct {
		Name        string
		OtherName   string
		Description string
	}*/
)

func main() {
	flag.Parse()

	dbcl := new(Mongo)
	dbcl.Init()
	defer dbcl.Sess.Close()

	xmlFile, err := os.Open(*inputFile)
	if err != nil {
		log.Println("Error opening file:", err)
		return
	}
	defer xmlFile.Close()

	decoder := xml.NewDecoder(xmlFile)
	var inEl, ct, pt string
	var node, pnode Node
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
				id, err := strconv.Atoi(ID.Value)
				var _ = err
				node = Node{Id:id, Name:"", Parent: &pnode, OtherName:"", Description:""}

			} else if inEl == "NODES" {
				log.Println("NODES savae:",node)
				dbcl.Save(&node);
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
				log.Println("NODE save:", node)
				dbcl.Save(&node);
			}
			if inEl == "NODES" {
				pnode = *pnode.Parent
			}
		}
	}

}