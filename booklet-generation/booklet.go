package booklet

import (
	"fmt"
)

type Logo struct {
	SVG string
}

type Organization struct {
	Name  string
	Email string
	Phone string
	Logo  *Logo
}

type PageType int

type FrontMatter PageType

const (
	HalfTitle FrontMatter = iota
	Frontispiece
	TitlePage
	EditionNotice
	Dedication
	Epigraph
	TableOfContents
	Foreward
	Preface
	Acknowledgements
	Introduction
	Prologue
	PrintersMark
)

type BodyMatter PageType

const (
	Body BodyMatter = iota
	ChapterPage
	PartPage
	SectionPage
	TippedInPage
)

type BackMatter PageType

const (
	Afterword BackMatter = iota
	Conclusion
	Epilogue
	Postscript
	Appendix
	Glossary
	Bibliography
	Index
	Errata
	Colophon
	Postface
)

// Aliasing for more epxressive API
const (
	Outro        = Epilogue
	Addendum     = Appendix
	BastardTitle = HalfTitle
)

// TODO: Shold it be lines?
type Page struct {
	Number  int
	Columns int
	Content string
}

type Definition struct {
	Word        string
	Description string // LOL WHAT IS THIS CALLED ITS IMPOSSIBLE TO SEARCH FOR
}

type Topic struct {
	Topic string
	Pages []int
}

type Author struct {
	Name      string
	Email     string
	PublicKey string
}

type Booklet struct {
	Title    string
	DOI      string
	Authors  []*Author
	Subtitle string
	Pages    []*Page
	Glossary []*Definition
	Index    []*Topic
}

func Booklet() {
	fmt.Println("vim-go")
}
