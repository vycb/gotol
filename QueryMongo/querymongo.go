package QueryMongo

import (
	"github.com/vycb/gotol/Parser"
	"log"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"strconv"
)

func MapReduce(fsearch *string, fwhat *string) {
	mc := new(Mongo)
	mc.Init()
	defer mc.Sess.Close()

	job := &mgo.MapReduce{
		Map:      "function() { emit(this.parent, 1) }",
		Reduce:   "function(key, values) { return Array.sum(values) }",
		Finalize: "function(key, count) { return {count: count} }",
	}
	var result []struct {
		Id    int `bson:"_id" json:"_id"`
		Value struct{ Count int }
	}

	_, err := mc.Tol.Find(bson.M{"name": bson.RegEx{*fsearch, ""}}).MapReduce(job, result)

	if err != nil {
		log.Println(err)
	}

	*fsearch = strconv.Itoa(result[0].Id)
	*fwhat = "id"
	Query(fsearch, fwhat)

	for _, n := range result {
		log.Println(n.Id, n.Value.Count)
	}
}

func Query(fsearch *string, fwhat *string) {
	mc := new(Mongo)
	mc.Init()
	defer mc.Sess.Close()

	var (
		nodes []Parser.DNode
		fq bson.M
		ires int
		e error
		fld string
	)
	switch *fwhat {
	case "parent":
	case "id":
		ires, e = strconv.Atoi(*fsearch); var _ = e
		if *fwhat == "parent" {
			fld = "parent"
		}else{
			fld = "_id"
		}
		fq = bson.M{fld: bson.M{"$eq":ires}}
	default:
		fq = bson.M{"name":  bson.RegEx{*fsearch, ""}}
	}

	e = mc.Tol.Find(fq).Sort("name").All(&nodes)

	if e != nil {
		log.Println(e)
	}

	for _, n := range nodes {
		log.Println(n.Id, n.Name, n.Parent, n.OtherName, n.Description)
	}
}

