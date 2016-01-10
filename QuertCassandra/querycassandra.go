package QueryCassandra

import (
	"strings"
	"fmt"
	"github.com/gocql/gocql"
)

func count(p int) int {
	dc := &Cassandra{}
	dc.Init()
	defer dc.SessionClose()

	var cnt int

	if err := dc.Session.Query("SELECT count(*) AS count FROM tol WHERE parent = ?", &p).Consistency(gocql.One).Scan(&cnt); err != nil {
		panic(err)
	}
	return cnt
}

func Query(fsearch string) {
	dc := &Cassandra{}
	dc.Init()
	defer dc.SessionClose()

	childs := func(p int) {
		iter := dc.Session.Query("SELECT id, name, parent, othername, description  FROM tol WHERE parent = ?", &p).Iter()
		var id, parent int
		var name, othername, description string

		for iter.Scan(&id, &name, &parent, &othername, &description) {
			fmt.Println(id, name, parent, othername, description, "c:", count(id))
		}
		fmt.Println("--------------------------------\n")
	}

	iter := dc.Session.Query("SELECT id, name, parent, othername, description  FROM tol").Iter()
	var id, parent, ct int
	var name, othername, description string

	for iter.Scan(&id, &name, &parent, &othername, &description) {

		if strings.Contains(strings.ToLower(name), strings.ToLower(fsearch)) {
			ct = count(id)

			fmt.Println(id, name, parent, othername, description, "c:", ct)

			if ct > 0 {

				fmt.Println(">")

				childs(id)
			}

		}
	}
}

func Toroot(fsearch string) {
	dc := &Cassandra{}
	dc.Init()
	defer dc.SessionClose()

	var id, parent int
	var name, othername, description string

	if er := dc.Session.Query("SELECT id, name, parent, othername, description  FROM tol WHERE name = ?", &fsearch).Consistency(gocql.One).Scan(&id, &name, &parent, &othername, &description); er != nil {
		panic(er)
	}
	fmt.Println(id, name, parent, othername, description, "c:", count(id))

	getparent := func(sid int) int {
		var id, parent int
		var name, othername, description string

		if er := dc.Session.Query("SELECT id, name, parent, othername, description  FROM tol WHERE id = ?", &sid).Consistency(gocql.One).Scan(&id, &name, &parent, &othername, &description); er!=nil{
			panic(er)
		}
		fmt.Println(id, name, parent, othername, description, "c:", count(id))

		return parent
	}

	for p:= getparent(parent); p >0; {
		p = getparent(p)
	}
}