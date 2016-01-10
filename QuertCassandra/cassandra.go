package QueryCassandra

import(
	"github.com/vycb/gotol/DbClient"
	"github.com/vycb/gotol/Parser"
	"log"
	"github.com/gocql/gocql"
	"golang.org/x/tools/container/intsets"
)
const INSERT_COUNT uint = 600
const INSERTST string = `INSERT INTO tol (id, name, parent, othername, description) VALUES(?, ?, ?, ?, ?)`

type Cassandra struct {
	Session *gocql.Session
	ct      DbClient.Counter
	batch   *gocql.Batch
	idsSet  intsets.Sparse
}

func (c *Cassandra)Init() {
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

func (c *Cassandra) Save(n *Parser.Node) {
	d := n.ToDNode()

	if c.idsSet.Has(d.Id) {
		return
	}
	c.idsSet.Insert(d.Id)

	c.batch.Query(INSERTST, d.Id, d.Name, d.Parent, d.OtherName, d.Description)
	c.ct.CtNext()

	if c.ct.GetCt() >= INSERT_COUNT {
		err := c.Session.ExecuteBatch(c.batch)
		if err != nil {
			log.Panic(err)
		}
		c.NewBatch()
	}

}
