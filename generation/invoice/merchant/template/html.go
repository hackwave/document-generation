package template

import (
	"fmt"
	"strings"
)

type Element struct {
	Path        string
	ID          string
	Classes     []string
	InlineStyle string
	Width       int
	Height      int
}

func (self Element) Size(width, height int) Element {
	self.Width = width
	self.Height = height
	return self
}

func Image(path string) Element {
	return Element{
		Path:   path,
		Width:  0,
		Height: 0,
	}
}

func (self Element) Style(style string) Element {
	self.InlineStyle = style
	return self
}

func (self Element) HTML() string {
	var attributes string

	if 0 < self.Width {
		attributes += fmt.Sprintf(` width="%v"`, self.Width)
	}
	if 0 < self.Height {
		attributes += fmt.Sprintf(` height="%v"`, self.Height)
	}
	if 0 < len(self.ID) {
		attributes += fmt.Sprintf(` id="%v"`, self.ID)
	}
	if 0 < len(self.Classes) {
		attributes += fmt.Sprintf(` class="%v"`, strings.Join(self.Classes, " "))
	}
	if 0 < len(self.InlineStyle) {
		attributes += fmt.Sprintf(` style="%v"`, self.InlineStyle)
	}

	return `<image src="` + self.Path + `" ` + attributes + ` />`
}
