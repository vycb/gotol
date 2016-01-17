package QueryRedis

import (
	"github.com/vycb/gotol/Parser"
	"fmt"
	"log"
	"encoding/json"
)

type QueryClient struct {
	rc *Redis
}

func (q *QueryClient)analise(key string) (int, *Parser.AData) {

	dat, err := q.rc.analiseScript.Run(q.rc.client, []string{"0"}, []string{key}).Result()
	if err != nil {
		log.Panic("analiseScript:", err)
	}
	js := dat.(string)
	//fmt.Println(js)

	var adata Parser.AData
	if len(js) > 2 {

		err = json.Unmarshal([]byte(js), &adata)
		if err != nil {
			//log.Println("Unmarshal:", err)
		}
	}

	return len(adata.Childes), &adata
}

func (q *QueryClient)Query(fsearch string) {
	q.rc = new(Redis)
	q.rc.Init()
	q.rc.initScript()
	defer q.rc.client.Close()

	size := q.rc.client.DbSize()
	fmt.Println(size)

	var cursor int64

	for {
		var keys []string
		var err error
		cursor, keys, err = q.rc.client.Scan(cursor, fsearch, 2000).Result()
		if err != nil {
			log.Println(err)
		}
		if cursor == 0 {
			break
		}

		for _, k := range keys {

			cnt, adata := q.analise(k)

			key := adata.Key

			fmt.Println(key.Id, key.Name, key.Parent, key.Othername, key.Description, cnt)

			if cnt > 0 {
				fmt.Println(">")
			}

			for _, c := range adata.Childes {

				fmt.Println(c.Id, c.Name, c.Parent, c.Othername, c.Description, c.Count)
			}

			if cnt > 0 && len(adata.Parents) >0 {
				fmt.Println("->")

				for _, c := range adata.Parents {

					fmt.Println(c.Id, c.Name, c.Parent, c.Othername, c.Description, c.Count)
				}
			}

			fmt.Println("----------------------\n")
		}
	}
}

