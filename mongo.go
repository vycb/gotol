package main

import (
	"gopkg.in/mgo.v2"
//"gopkg.in/mgo.v2/bson"
	"log"
)

//var sem = make(chan uint, MaxOutstanding)

type Mongo struct {
	Db   *mgo.Database
	Sess *mgo.Session
	Tol  *mgo.Collection
	Bulk *mgo.Bulk
	Ct   *Counter
}

func (m *Mongo) Init() {
	m.Ct = &Counter{}
	sess, err := mgo.Dial("mongodb://vycb:123@ds029541.mongolab.com:29541/blog")
	if err != nil {
		panic(err)
	}

	m.Sess = sess

	m.Sess.SetMode(mgo.Monotonic, true)

	m.Db = m.Sess.DB("blog")
	m.Tol = m.Db.C("tol")

	index := mgo.Index{
		Key:        []string{"name", "parent"},
		Background: true,
		Sparse:     true,
	}

	err = m.Tol.EnsureIndex(index)
	if err != nil {
		panic(err)
	}
	//m.Nodes = []*Node
}

func (m *Mongo)SessionClose() {
	defer m.Sess.Close()
}

func (m *Mongo) NewBatch() {
	m.Bulk = m.Tol.Bulk()
	m.Bulk.Unordered()
}

func (m *Mongo  ) Save(n *Node) {

	dn := n.ToDNode()

	m.Bulk.Insert(&dn)
	m.Ct.CtNext()

	if m.Ct.GetCt() >= INSERT_COUNT {
		//sem <- 1
		//go func() {
		_, err := m.Bulk.Run()
		if err != nil {
			log.Println("Bulk.Run:", err)
		}
		m.NewBatch()
		m.Ct.SetCt()
		//<-sem
		//}()
	}

}
