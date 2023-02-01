package ui

type Component interface {
	Draw()
	Redraw()
	Update()
}
