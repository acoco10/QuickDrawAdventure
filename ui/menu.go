package ui

import "github.com/ebitenui/ebitenui/widget"

type dialogueType uint8

const (
	Good dialogueType = iota
	Hostile
)

type Menu struct {
	MenuContainer    *widget.Container
	Buttons          []*widget.Button
	ButtonVisibility bool
}

type DialogeMenu struct {
	DialogueOptions map[dialogueType]*widget.Button
	ButtonText      string
	ButtonTextIndex map[dialogueType]int
	ButtonTextList  map[dialogueType][]string
}

func (m *Menu) DisableButtonVisibility() {
	for _, b := range m.Buttons {
		b.GetWidget().Visibility = widget.Visibility_Hide
	}
}

func (m *Menu) DisableButtons() {
	for _, b := range m.Buttons {
		b.GetWidget().Disabled = true
	}
}

func (m *Menu) EnableButtonVisibility() {
	for _, b := range m.Buttons {
		b.GetWidget().Visibility = widget.Visibility_Show
	}
}
