package booklet

import (
	"fmt"
	"regexp"
	"strings"

	ast "../ast"
	levenshtein "../levenshtein"
)

type Section struct {
	Path string

	Title Content
	Body  Content

	PrimaryTag Tag
	Tags       []Tag

	Parent   *Section
	Children []*Section

	Style    string
	Partials Partials

	SplitSections        bool
	PreventSplitSections bool

	ResetDepth bool

	OmitChildrenFromTableOfContents bool

	Processor       SectionProcessor
	PluginFactories []PluginFactory
	Plugins         []Plugin

	Location       ast.Location
	InvokeLocation ast.Location
}

type Partials map[string]Content

type SectionProcessor interface {
	EvaluateFile(*Section, string, []PluginFactory) (*Section, error)
	EvaluateNode(*Section, ast.Node, []PluginFactory) (*Section, error)
}

type Tag struct {
	Name  string
	Title Content

	Section  *Section
	Location ast.Location
	Anchor   string

	Content Content
}

func (con *Section) String() string {
	return fmt.Sprintf("{section (%s): %s}", con.Path, con.Title)
}

func (con *Section) FilePath() string {
	if con.Path != "" {
		return con.Path
	}

	if con.Parent != nil {
		return con.Parent.FilePath()
	}

	return ""
}

func (con *Section) IsFlow() bool {
	return false
}

func (con *Section) Visit(visitor Visitor) error {
	return visitor.VisitSection(con)
}

func (con *Section) SetTitle(title Content, loc ast.Location, tags ...string) {
	if len(tags) == 0 {
		tags = []string{con.defaultTag(title)}
	}

	con.Tags = []Tag{}
	for _, name := range tags {
		con.SetTag(name, title, loc)
	}

	con.Title = title
	con.PrimaryTag = con.Tags[0]
}

func (con *Section) SetTag(name string, title Content, loc ast.Location, optionalAnchor ...string) {
	anchor := ""
	if len(optionalAnchor) > 0 {
		anchor = optionalAnchor[0]
	}

	con.Tags = append(con.Tags, Tag{
		Section:  con,
		Location: loc,

		Name:   name,
		Title:  title,
		Anchor: anchor,
	})
}

func (con *Section) SetTagAnchored(name string, title Content, loc ast.Location, content Content, anchor string) {
	con.Tags = append(con.Tags, Tag{
		Section:  con,
		Location: loc,

		Name:   name,
		Title:  title,
		Anchor: anchor,

		Content: content,
	})
}

func (con *Section) Number() string {
	if con.Parent == nil {
		return ""
	}

	parentNumber := con.Parent.Number()
	selfIndex := 1
	for _, child := range con.Parent.Children {
		if child == con {
			break
		}

		selfIndex++
	}

	if parentNumber == "" {
		return fmt.Sprintf("%d", selfIndex)
	}

	return fmt.Sprintf("%s.%d", parentNumber, selfIndex)
}

func (con *Section) HasAnchors() bool {
	for _, tag := range con.Tags {
		if tag.Anchor != "" {
			return true
		}
	}

	if con.SplitSections {
		return false
	}

	for _, child := range con.Children {
		if child.HasAnchors() {
			return true
		}
	}

	return false
}

func (con *Section) AnchorTags() []Tag {
	tags := []Tag{}

	for _, tag := range con.Tags {
		if tag.Anchor == "" {
			continue
		}

		tags = append(tags, tag)
	}

	return tags
}

func (con *Section) Top() *Section {
	if con.Parent != nil {
		return con.Parent.Top()
	}

	return con
}

func (con *Section) Contains(sub *Section) bool {
	if con == sub {
		return true
	}

	for _, child := range con.Children {
		if child.Contains(sub) {
			return true
		}
	}

	return false
}

func (con *Section) IsOrHasChild(sub *Section) bool {
	if con == sub {
		return true
	}

	for _, child := range con.Children {
		if child == sub {
			return true
		}
	}

	return false
}

func (con *Section) Prev() *Section {
	if con.Parent == nil {
		return nil
	}

	var lastChild *Section
	for _, child := range con.Parent.Children {
		if lastChild != nil && child == con {
			return lastChild
		}

		lastChild = child
	}

	return con.Parent
}

func (con *Section) Next() *Section {
	if con.SplitSections {
		if len(con.Children) > 0 {
			return con.Children[0]
		}
	}

	return con.NextSibling()
}

func (con *Section) NextSibling() *Section {
	if con.Parent == nil {
		return nil
	}

	var sawSelf bool
	for _, child := range con.Parent.Children {
		if sawSelf {
			return child
		}

		if child == con {
			sawSelf = true
		}
	}

	return con.Parent.NextSibling()
}

func (con *Section) FindTag(tagName string) []Tag {
	return con.filterTags(true, nil, func(other string) bool {
		return other == tagName
	})
}

func (con *Section) SimilarTags(tagName string) []Tag {
	return con.filterTags(true, nil, func(other string) bool {
		return levenshtein.Match(tagName, other, nil) > 0.5
	})
}

func (con *Section) SetPartial(name string, value Content) {
	if con.Partials == nil {
		con.Partials = Partials{}
	}

	con.Partials[name] = value
}

func (con *Section) Partial(name string) Content {
	return con.Partials[name]
}

func (con *Section) UsePlugin(pf PluginFactory) {
	con.PluginFactories = append(con.PluginFactories, pf)
	con.Plugins = append(con.Plugins, pf(con))
}

func (con *Section) Depth() int {
	if con.Parent == nil {
		return 0
	}

	return con.Parent.Depth() + 1
}

func (con *Section) PageDepth() int {
	if con.Parent == nil || con.Parent.ResetDepth {
		return 0
	}

	return con.Parent.PageDepth() + 1
}

func (con *Section) SplitSectionsPrevented() bool {
	if con.PreventSplitSections {
		return true
	}

	if con.Parent != nil && con.Parent.SplitSectionsPrevented() {
		return true
	}

	return false
}

func (con *Section) filterTags(up bool, exclude *Section, match func(string) bool) []Tag {
	tags := []Tag{}

	if match(con.Title.String()) {
		tags = append(tags, con.PrimaryTag)
	}

	for _, t := range con.Tags {
		if match(t.Name) {
			tags = append(tags, t)
		}
	}

	for _, sub := range con.Children {
		if sub != exclude {
			tags = append(tags, sub.filterTags(false, nil, match)...)
		}
	}

	if up && con.Parent != nil {
		tags = append(tags, con.Parent.filterTags(true, con, match)...)
	}

	return tags
}

var whitespaceRegexp = regexp.MustCompile(`\s+`)
var specialCharsRegexp = regexp.MustCompile(`[^[:alnum:]_\-]`)

func (con *Section) defaultTag(title Content) string {
	return strings.ToLower(
		specialCharsRegexp.ReplaceAllString(
			whitespaceRegexp.ReplaceAllString(
				strings.Replace(
					StripAux(title).String(),
					" & ",
					" and ",
					-1,
				),
				"-",
			),
			"",
		),
	)
}
