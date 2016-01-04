package main

import (
	"gopkg.in/mgo.v2"
//"gopkg.in/mgo.v2/bson"
	"log"
)

const INSERT_COUNT uint = 1000
const MaxOutstanding uint = 1

var sem = make(chan uint, MaxOutstanding)

type(
	Mongo struct {
		Db    *mgo.Database
		Sess  *mgo.Session
		Tol   *mgo.Collection
		Bulk  *mgo.Bulk
		Ct    uint
	}
)

func (m *Mongo) CtNext() {
	m.Ct++
}

func (m *Mongo) GetCt() uint {
	return m.Ct
}

func (m *Mongo) SetCt() {
	m.Ct = 0
}

func (n *Node)ToMNode() *MNode {
	return &MNode{Id:n.Id, Name:n.Name, Parent:n.Parent.Id, OtherName:n.OtherName, Description:n.Description}
}

func (m *Mongo) Init() {

	sess, err := mgo.Dial("mongodb://vycb:123@ds029541.mongolab.com:29541/blog")
	if err != nil {
		panic(err)
	}

	m.Sess = sess

	m.Sess.SetMode(mgo.Monotonic, true)

	m.Db = m.Sess.DB("blog")
	m.Tol = m.Db.C("tol")
	m.Bulk = m.Tol.Bulk()
	//m.Nodes = []*Node
}

func (m *Mongo) NewBulk() {
	m.Bulk = m.Tol.Bulk()
}

func (m *Mongo  ) Save(node *Node) {

	mn := node.ToMNode()

	m.Bulk.Insert(&mn)
	m.CtNext()

	if m.GetCt() >= INSERT_COUNT {
		sem <- 1
		go func() {
			r, err := m.Bulk.Run()
			if err != nil {
				log.Println(err, r)
			}
			m.SetCt()
			m.NewBulk()
			<-sem
		}()
	}

}
