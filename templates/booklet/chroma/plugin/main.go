package plugin

import (
	booklet "../booklet"
	chroma "../chroma"
)

func init() {
	booklet.RegisterPlugin("chroma", chroma.NewPlugin)
}
