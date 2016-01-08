package QueryPq


import (
	"log"
	"github.com/jackc/pgx"
	"github.com/vycb/gotol/Parser"
	"github.com/vycb/gotol/DbClient"
	"strconv"
)

const INSERT_COUNT uint = 1

type Pq struct {
	Conn  *pgx.ConnPool
	ct    DbClient.Counter
	query string
}

func extractConfig() pgx.ConnConfig {
	var config pgx.ConnConfig
	//postgres://fumpaabc:4niio10DZ9IeiIgAK3V1IAyMn_9Fk6Ig@pellefant.db.elephantsql.com:5432/fumpaabc
	config.Host = "pellefant.db.elephantsql.com"
	config.Port = 5432
	config.User = "fumpaabc"
	config.Password = "4niio10DZ9IeiIgAK3V1IAyMn_9Fk6Ig"
	config.Database = "fumpaabc"
	return config
}
func afterConnect(conn *pgx.Conn) (err error) {
	_, _ = conn.Prepare("insert", `
    insert into tol(id, name, parent, othername, description)
     SELECT $1, $2, $3, $4, $5
     WHERE NOT EXISTS (SELECT 1 FROM tol WHERE id = $6)
  `)
	_, _ = conn.Prepare("create", `
    CREATE TABLE IF NOT EXISTS tol (
    id int PRIMARY KEY,
    name varchar(150) NULL,
    parent int,
    othername text NULL,
    description text NULL
    )
  `)
	return
}

func (p *Pq) Init() {
	p.ct = 0
	conn, err := pgx.NewConnPool(pgx.ConnPoolConfig{
		ConnConfig:  extractConfig(),
		AfterConnect: afterConnect,
	})
	if err != nil {
		log.Fatal("Unable to connection to database:", err)
	}
	p.Conn = conn

	//if _, err := p.Conn.Exec("create"); err == nil {
	//	log.Fatal(err)
	//}
	p.NewBatch()
}

func (p *Pq) SessionClose() {
}
func (p *Pq) NewBatch() {
	p.query = "insert into tol(id, name, parent, othername, description) values"
}

func (p *Pq) Save(n *Parser.Node) {
	d := n.ToDNode()
	if d.Id == 4 {
		var _ = d.Name
	}
	i := strconv.Itoa(d.Id)
	a := strconv.Itoa(d.Parent)

	p.query += "("+ i +",'"+ d.Name+"', "+ a +", '"+d.OtherName+"', '"+d.Description+"'+)"

	p.ct.CtNext()

	if p.ct.GetCt() >= INSERT_COUNT {

		if _, err := p.Conn.Exec("insert", d.Id, d.Name, d.Parent, d.OtherName, d.Description, d.Id); err != nil {
			log.Panic(err)
		}
		p.NewBatch()
	}

}