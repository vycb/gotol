package QueryRedis

import (
	"gopkg.in/redis.v3"
	"github.com/vycb/gotol/DbClient"
	"golang.org/x/tools/container/intsets"
	"github.com/vycb/gotol/Parser"
	"strconv"
	"strings"
	"log"
)

const INSERT_COUNT uint = 500
const MaxOutstanding uint = 1
const ANALISE_SCRIPT string = `
local res,keyp,keyid,count
local toroot,childes,key = {},{},{}

local function isset(val)
	return  val ~= nil and val ~= ''
end

local function getProp(hash, prop)
   for i = 1, #hash, 2 do
			if hash[i] == prop then
				return hash[i+1]
			end
		end
		return ''
end

local function hashToNode(hash)
	return {id = getProp(hash,'id'), name = getProp(hash,'name'), parent = getProp(hash, 'parent'), othername = getProp(hash, 'othername'), description = getProp(hash, 'description')}
end

local function findHash(find, prop)

	for x = 1, #res, 1 do
		local hash = redis.call('HGETALL', res[x])
		local pval = getProp(hash, prop)

		if pval == find then
			return hash
		end
	end
	return ''
end

local function countItems(find, prop)
	count = 0
	for x = 1, #res, 1 do

		local hash = redis.call('HGETALL', res[x])
		local item = getProp(hash, prop)

		if item == find then
			count = count +1
		end
	end
	return count
end

local function getParent(hash, pid)
	local p = findHash(pid, 'id')

	if isset(p) then
		local ip = getProp(p, 'id')
		local count = countItems(ip,'parent')
		local node = hashToNode(p)
		node['count'] = count
		table.insert(hash, node)
	end

	local pi = getProp(p, 'parent')

	if isset(pi) and tonumber(pi) > 0 then
		getParent(hash, pi)
	end
end


key = redis.call('HGETALL', ARGV[1])
keyid = getProp(key, 'id')

res = redis.call('KEYS','*')

for idx = 1, #res, 1 do
	local hash = redis.call('HGETALL', res[idx])

	local parent = getProp(hash, 'parent')

	if parent == keyid then
		local hi = getProp(hash, 'id')
		local count = countItems(hi,'parent')

		table.insert(childes, {id = hi, name = getProp(hash,'name'), parent = parent, othername = getProp(hash, 'othername'), description = getProp(hash, 'description'), count = count })
	end
end

keyp = getProp(key, 'parent')
if isset(keyp) and tonumber(keyp) > 0 and table.getn(childes) > 0 then

	getParent(toroot, keyp)
end

cjson.encode_sparse_array(true)
return cjson.encode({key = hashToNode(key), childes = childes, parents = toroot})
`;

type Redis struct {
	client        *redis.Client
	pipeline      *redis.Pipeline
	//pipeline      []*redis.Pipeline
	ct            DbClient.Counter
	idsSet        intsets.Sparse
	analiseScript *redis.Script
	sem           chan int
}

func (r *Redis) Init() {

	r.client = redis.NewClient(&redis.Options{
		Addr:"pub-redis-11548.us-east-1-3.2.ec2.garantiadata.com:11548",
	})

	r.sem = make(chan int, MaxOutstanding)
}

func (r *Redis) initScript() {
	r.analiseScript = redis.NewScript(ANALISE_SCRIPT)
}

func (r *Redis)SessionClose() {
	defer r.client.Close()
}

func (r *Redis) NewBatch() {
	r.pipeline = r.client.Pipeline()
}
//func (r *Redis) NewBatch() {
//	r.pipeline = append(r.pipeline, r.client.Pipeline())
//}

func (r *Redis) getBatch(cb int) *redis.Pipeline {
	return r.pipeline
}
//func (r *Redis) getBatch(cb int) *redis.Pipeline {
//	return r.pipeline[cb]
//}

func (r *Redis) Crp() int{
	return 1
}
//func (r *Redis) Crp() int{
//	return len(r.pipeline) -1
//}

func (r *Redis  ) Save(n *Parser.Node) {

	dn := n.ToDNode()

	if r.idsSet.Has(dn.Id) {
		return
	}
	r.idsSet.Insert(dn.Id)

	key := strconv.Itoa(dn.Id) + ":" + strings.ToLower(dn.Name)

	r.getBatch(r.Crp()).HMSet(key, "id", strconv.Itoa(dn.Id), "name", dn.Name, "parent", strconv.Itoa(dn.Parent), "othername", dn.OtherName, "description", dn.Description)

	r.ct.CtNext()

	if r.ct.GetCt() >= INSERT_COUNT {

		//r.sem <-1

		//go func(cb int) {

			cmds, err := r.pipeline.Exec(); var _ = cmds
//			cmds, err := r.getBatch(r.Crp()).Exec(); var _ = cmds
			if err != nil {
				log.Println("pipeline.Exec:", err)
			}
			//<-r.sem

		//}(r.Crp())

		r.NewBatch()
		r.ct.SetCt()
	}

}
