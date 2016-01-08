package QueryMongo

import (
	"github.com/vycb/gotol/Parser"
	"github.com/vycb/gotol/DbClient"
	"gopkg.in/mgo.v2"
	"log"
	"golang.org/x/tools/container/intsets"
)
const INSERT_COUNT uint = 100

type Mongo struct {
	Db   *mgo.Database
	Sess *mgo.Session
	Tol  *mgo.Collection
	Bulk *mgo.Bulk
	ct   DbClient.Counter
	idsSet intsets.Sparse
}

func (m *Mongo) Init() {
	m.ct = 0

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
}

func (m *Mongo)SessionClose() {
	defer m.Sess.Close()
}

func (m *Mongo) NewBatch() {
	m.Bulk = m.Tol.Bulk()
	m.Bulk.Unordered()
}

func (m *Mongo  ) Save(n *Parser.Node) {

	dn := n.ToDNode()

	if m.idsSet.Has(dn.Id) {
		return
	}
	m.idsSet.Insert(dn.Id)

	m.Bulk.Insert(&dn)
	m.ct.CtNext()

	if m.ct.GetCt() >= INSERT_COUNT {
		_, err := m.Bulk.Run()
		if err != nil {
			log.Println("Bulk.Run:", err)
		}
		m.NewBatch()
		m.ct.SetCt()
	}

}
