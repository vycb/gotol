package main

import (
	"log"
	"encoding/xml"
	"strings"
	"strconv"
	"os"
)

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
)

func parse() {
	mc := new(Mongo)
	mc.Init()
	defer mc.Sess.Close()
	mc.NewBulk()


	xmlFile, err := os.Open(fimport)
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
				log.Println("NODES savae:", node)
				mc.Save(&node);
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
				mc.Save(&node);
			}
			if inEl == "NODES" {
				pnode = *pnode.Parent
			}
		}
	}
}
