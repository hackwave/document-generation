package plugin

import (
	"fmt"
	"path/filepath"

	ast "../ast"
	booklet "../booklet"
)

func init() {
	booklet.RegisterPlugin("base", NewPlugin)
}

func NewPlugin(section *booklet.Section) booklet.Plugin {
	return Plugin{
		section: section,
	}
}

type Plugin struct {
	section *booklet.Section
}

func (plugin Plugin) UsePlugin(name string) error {
	pluginFactory, found := booklet.LookupPlugin(name)
	if !found {
		return fmt.Errorf("unknown plugin '%s'", name)
	}

	plugin.section.UsePlugin(pluginFactory)

	return nil
}

func (plugin Plugin) Styled(name string) {
	plugin.section.Style = name
}

func (plugin Plugin) Title(title booklet.Content, tags ...string) {
	plugin.section.SetTitle(title, plugin.section.InvokeLocation, tags...)
}

func (plugin Plugin) Aux(content booklet.Content) booklet.Content {
	return booklet.Aux{
		Content: content,
	}
}

func (plugin Plugin) Section(node ast.Node) error {
	section, err := plugin.section.Processor.EvaluateNode(plugin.section, node, plugin.section.PluginFactories)
	if err != nil {
		return err
	}

	section.Location = plugin.section.InvokeLocation

	plugin.section.Children = append(plugin.section.Children, section)

	return nil
}

func (plugin Plugin) IncludeSection(path string) error {
	sectionPath := filepath.Join(filepath.Dir(plugin.section.FilePath()), path)

	section, err := plugin.section.Processor.EvaluateFile(plugin.section, sectionPath, []booklet.PluginFactory{NewPlugin})
	if err != nil {
		return err
	}

	plugin.section.Children = append(plugin.section.Children, section)

	return nil
}

func (plugin Plugin) SinglePage() {
	plugin.section.PreventSplitSections = true
}

func (plugin Plugin) SplitSections() {
	plugin.section.ResetDepth = true

	if !plugin.section.SplitSectionsPrevented() {
		plugin.section.SplitSections = true
	}
}

func (plugin Plugin) OmitChildrenFromTableOfContents() {
	plugin.section.OmitChildrenFromTableOfContents = true
}

func (plugin Plugin) TableOfContents() booklet.Content {
	return booklet.TableOfContents{
		Section: plugin.section,
	}
}

func (plugin Plugin) Code(content booklet.Content) booklet.Content {
	return booklet.Styled{
		Content: content,
		Style:   booklet.StyleVerbatim,
	}
}

func (plugin Plugin) Italic(content booklet.Content) booklet.Content {
	return booklet.Styled{
		Content: content,
		Style:   booklet.StyleItalic,
	}
}

func (plugin Plugin) Bold(content booklet.Content) booklet.Content {
	return booklet.Styled{
		Content: content,
		Style:   booklet.StyleBold,
	}
}

func (plugin Plugin) Larger(content booklet.Content) booklet.Content {
	return booklet.Styled{
		Content: content,
		Style:   booklet.StyleLarger,
	}
}

func (plugin Plugin) Smaller(content booklet.Content) booklet.Content {
	return booklet.Styled{
		Content: content,
		Style:   booklet.StyleSmaller,
	}
}

func (plugin Plugin) Strike(content booklet.Content) booklet.Content {
	return booklet.Styled{
		Content: content,
		Style:   booklet.StyleStrike,
	}
}

func (plugin Plugin) Superscript(content booklet.Content) booklet.Content {
	return booklet.Styled{
		Content: content,
		Style:   booklet.StyleSuperscript,
	}
}

func (plugin Plugin) Subscript(content booklet.Content) booklet.Content {
	return booklet.Styled{
		Content: content,
		Style:   booklet.StyleSubscript,
	}
}

func (plugin Plugin) Inset(content booklet.Content) booklet.Content {
	return booklet.Styled{
		Content: content,
		Style:   booklet.StyleInset,
	}
}

func (plugin Plugin) Aside(content booklet.Content) booklet.Content {
	return booklet.Styled{
		Content: content,
		Style:   booklet.StyleAside,
	}
}

func (plugin Plugin) Link(content booklet.Content, target string) booklet.Content {
	return booklet.Link{
		Content: content,
		Target:  target,
	}
}

func (plugin Plugin) Reference(tag string, content ...booklet.Content) booklet.Content {
	ref := &booklet.Reference{
		TagName: tag,

		Location: plugin.section.InvokeLocation,
	}

	if len(content) > 0 {
		ref.Content = content[0]
	}

	return ref
}

func (plugin Plugin) Target(tag string, titleAndContent ...booklet.Content) booklet.Content {
	ref := &booklet.Target{
		TagName:  tag,
		Location: plugin.section.InvokeLocation,
	}

	switch len(titleAndContent) {
	case 2:
		ref.Title = titleAndContent[0]
		ref.Content = titleAndContent[1]
	case 1:
		ref.Title = titleAndContent[0]
	default:
		ref.Title = booklet.String(tag)
	}

	return ref
}

func (plugin Plugin) List(items ...booklet.Content) booklet.Content {
	return booklet.List{
		Items: items,
	}
}

func (plugin Plugin) OrderedList(items ...booklet.Content) booklet.Content {
	return booklet.List{
		Items:   items,
		Ordered: true,
	}
}

func (plugin Plugin) Image(path string, description ...string) booklet.Content {
	img := booklet.Image{
		Path: path,
	}

	if len(description) > 0 {
		img.Description = description[0]
	}

	return img
}

func (plugin Plugin) SetPartial(name string, content booklet.Content) {
	plugin.section.SetPartial(name, content)
}

func (plugin Plugin) Table(rows ...booklet.Content) (booklet.Content, error) {
	table := booklet.Table{}

	for _, row := range rows {
		list, ok := row.(booklet.List)
		if !ok {
			return nil, fmt.Errorf("table row is not a list: %s", row)
		}

		table.Rows = append(table.Rows, list.Items)
	}

	return table, nil
}

func (plugin Plugin) TableRow(cols ...booklet.Content) booklet.Content {
	return plugin.List(cols...)
}

func (plugin Plugin) Definitions(items ...booklet.Content) (booklet.Content, error) {
	defs := booklet.Definitions{}
	for _, item := range items {
		list, ok := item.(booklet.List)
		if !ok {
			return nil, fmt.Errorf("definition item is not a list: %s", item)
		}

		if len(list.Items) != 2 {
			return nil, fmt.Errorf("definition item must have two entries: %s", item)
		}

		defs = append(defs, booklet.Definition{
			Subject:    list.Items[0],
			Definition: list.Items[1],
		})
	}

	return defs, nil
}

func (plugin Plugin) Definition(subject booklet.Content, definition booklet.Content) booklet.Content {
	return plugin.List(subject, definition)
}
