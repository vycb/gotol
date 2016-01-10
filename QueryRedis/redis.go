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

const INSERT_COUNT uint = 1000
const ANALISE_SCRIPT string =`
local res,keyp,keyid,count
local hashd,toroot,childes,key = {},{},{},{}

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

local function getParent(hash, pid)

	local p = findHash(pid, 'id')

	if isset(p) then
		table.insert(hash, hashToNode(p))
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
		table.insert(hashd, {i = getProp(hash, 'id'), n = getProp(hash,'name'), p = parent, o = getProp(hash, 'othername'), d = getProp(hash, 'description') })
	end
end


keyp = getProp(key, 'parent')

if isset(keyp) and tonumber(keyp) > 0 and table.getn(hashd) > 0 then

	getParent(toroot, keyp)
end


for j,v in ipairs(hashd) do
	count = 0

	for x = 1, #res, 1 do

		local hash = redis.call('HGETALL', res[x])
		local parent = getProp(hash, 'parent')

		if parent == v['i'] then
			count = count +1
		end
	end

	table.insert(childes, {id = v['i'], name = v['n'], parent = v['p'], othername = v['o'], description = v['d'], count = count} )
end

cjson.encode_sparse_array(true)
return cjson.encode({['key'] = hashToNode(key), ['childes'] = childes, ['parents'] = toroot})
`;
type Redis struct {
	client   *redis.Client
	pipeline *redis.Pipeline
	ct       DbClient.Counter
	idsSet   intsets.Sparse
	analiseScript *redis.Script
}

func (r *Redis) Init() {

	r.client = redis.NewClient(&redis.Options{
		Addr:"pub-redis-11548.us-east-1-3.2.ec2.garantiadata.com:11548",
	})
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

func (r *Redis  ) Save(n *Parser.Node) {

	dn := n.ToDNode()

	if r.idsSet.Has(dn.Id) {
		return
	}
	r.idsSet.Insert(dn.Id)

	key := strconv.Itoa(dn.Id) + ":" + strings.ToLower(dn.Name)

	r.pipeline.HMSet(key, "id", strconv.Itoa(dn.Id), "name", dn.Name, "parent", strconv.Itoa(dn.Parent), "othername", dn.OtherName, "description", dn.Description)

	r.ct.CtNext()

	if r.ct.GetCt() >= INSERT_COUNT {

		cmds, err := r.pipeline.Exec(); var _ = cmds
		if err != nil {
			log.Println("pipeline.Exec:", err)
		}

		r.NewBatch()
		r.ct.SetCt()
	}

}