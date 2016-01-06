package main

const (
	INSERT_COUNT uint = 1000
	CINSERT_COUNT uint = 600
//MaxOutstanding uint = 1
)

type (
	Db interface {
		Init()
		NewBatch()
		SessionClose()
		Save(n *Node)
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