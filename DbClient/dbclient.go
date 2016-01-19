package DbClient

import (
	."github.com/vycb/gotol/Node"
)

type (

	Db interface {
		Init()
		NewBatch()
		SessionClose()
		Save(n *Node)
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