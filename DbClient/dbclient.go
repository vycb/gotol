package DbClient

import (
	"github.com/vycb/gotol/Parser"
)

type (

	Db interface {
		Init()
		NewBatch()
		SessionClose()
		Save(n *Parser.Node)
	}

	Counter uint
)

func (dc *Counter) CtNext() {
	*dc++
}

func (dc *Counter) GetCt() uint {
	return uint(*dc)
}

func (dc *Counter) SetCt() {
	*dc = 0
}