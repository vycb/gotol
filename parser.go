package main

import (
	"fmt"
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
	DNode struct {
		Id          int `bson:"_id,omitempty" json:"_id"`
		Parent      int `bson:"parent,omitempty" json:"parent"`
		Name        string `bson:"name" json:"name"`
		OtherName   string `bson:"othername" json:"othername"`
		Description string `bson:"description" json:"description"`
	}
)

func (n *Node)ToDNode() *DNode {
	return &DNode{Id:n.Id, Name:n.Name, Parent: n.Parent.Id, OtherName:n.OtherName, Description:n.Description}
}
func parse() {
	var dc Db

	if fdbcl == "mongo" {
		dc = &Mongo{}
	}else{
		dc = &Cassandra{}
	}
	dc.Init()
	dc.NewBatch()
	defer dc.SessionClose()

	xmlFile, err := os.Open(fimport)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer xmlFile.Close()

	decoder := xml.NewDecoder(xmlFile)
	var inEl, ct, pt string
	var node *Node
	pnode := new(Node)
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

				node = &Node{Id:id, Name:"", Parent: pnode, OtherName:"", Description:""}

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
