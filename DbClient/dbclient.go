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

	Counter struct {
		Ct uint
	}
)

func (dc *Counter) CtNext() {
	dc.Ct++
}

func (dc *Counter) GetCt() uint {
	return dc.Ct
}

func (dc *Counter) SetCt() {
	dc.Ct = 0
}