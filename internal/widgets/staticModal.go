package widgets

import (
	"github.com/rivo/tview"
)

type StaticModal struct {
	Modal   tview.Primitive
	Content tview.Primitive
}

func NewStaticModal(content tview.Primitive) *StaticModal {
	m := StaticModal{
		Content: content,
	}

	m.Modal = NewGridModal(m.Content, 40, 10)

	return &m
}

func (m *StaticModal) Primitive() tview.Primitive {
	return m.Modal
}
