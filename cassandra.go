package main

import (
	"log"
	"github.com/gocql/gocql"
)

const INSERTST string = `INSERT INTO tol (id, name, parent, othername, description) VALUES(?, ?, ?, ?, ?)`

type Cassandra struct {
	Session *gocql.Session
	Ct      *Counter
	batch   *gocql.Batch
}

func (c *Cassandra)Init() {
	c.Ct = &Counter{}
	cluster := gocql.NewCluster("localhost")
	cluster.Keyspace = "tol_keyspace"
	cluster.Consistency = gocql.Quorum
	session, _ := cluster.CreateSession()
	c.Session = session
}

func (c *Cassandra) SessionClose() {
	defer c.Session.Close()
}
func (c *Cassandra) NewBatch() {
	c.batch = gocql.NewBatch(gocql.LoggedBatch)
}

func (c *Cassandra) Save(n *Node) {
	d := n.ToDNode()
	if d.Id == 4 {
		var _ = d.Name
	}
	c.batch.Query(INSERTST, d.Id, d.Name, d.Parent, d.OtherName, d.Description)
	c.Ct.CtNext()

	if c.Ct.GetCt() >= CINSERT_COUNT {
		err := c.Session.ExecuteBatch(c.batch)
		if err != nil {
			log.Panic(err)
		}
		c.NewBatch()
	}

}