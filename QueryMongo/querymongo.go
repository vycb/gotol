package QueryMongo

import (
	."github.com/vycb/gotol/Node"
	"fmt"
//"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"strconv"
	"log"
)

func parent(id int) int {
	if id <=0 {
		return 0
	}
	mc := new(Mongo)
	mc.Init()
	defer mc.Sess.Close()

	var node DNode
	err := mc.Tol.Find(bson.M{"_id": id}).One(&node)

	if err != nil {
		log.Println(err)
	}

	fmt.Println(node.Id, node.Name, node.Parent, node.OtherName, node.Description)

	return node.Parent
}

func childes(id int) (count int, nodes []DNode) {
	mc := new(Mongo)
	mc.Init()
	defer mc.Sess.Close()

	/*job := &mgo.MapReduce{
		Map:      "function() { emit(this.parent, 1) }",
		Reduce:   "function(key, values) { return Array.sum(values) }",
		Finalize: "function(key, count) { return {count: count} }",
	}
	var result []struct {
		Id    int `bson:"_id" json:"_id"`
		Value struct{ Count int }
	}
	//_, err := mc.Tol.Find(bson.M{"parent": id}).MapReduce(job, &result)

	*/

	err := mc.Tol.Find(bson.M{"parent": id}).All(&nodes)

	if err != nil {
		log.Println(err)
	}

	return len(nodes), nodes
}

var nodes []DNode

func Query(fsearch string, fwhat string) {
	mc := new(Mongo)
	mc.Init()
	defer mc.Sess.Close()

	var (
	fq bson.M
	ires int
	e error
	fld string
	)
	switch fwhat {
	case "parent":
	case "id":
		ires, e = strconv.Atoi(fsearch); var _ = e
		if fwhat == "parent" {
			fld = "parent"
		}else {
			fld = "_id"
		}
		fq = bson.M{fld: bson.M{"$eq":ires}}
	default:
		fq = bson.M{"name":  bson.RegEx{Pattern:fsearch, Options:"i"}}
	}

	e = mc.Tol.Find(fq).Sort("name").All(&nodes)

	if e != nil {
		log.Println(e)
	}

	for _, n := range nodes {

		count, childes := childes(n.Id)

		fmt.Println(n.Id, n.Name, n.Parent, n.OtherName, n.Description, count)

		if count > 0 {
			fmt.Println(">")
		}

		for _, c := range childes {

			fmt.Println(c.Id, c.Name, c.Parent, c.OtherName, c.Description)
		}

		if count > 0 {
			fmt.Println("->")

			for p := parent(n.Id); p >0; p = parent(p)  {

			}

			fmt.Println("----------------------")
		}
		fmt.Println("")
	}
}

