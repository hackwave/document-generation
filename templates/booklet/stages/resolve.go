package stages

import (
	"fmt"

	booklet "../booklet"
)

type Resolve struct {
	AllowBrokenReferences bool

	Section *booklet.Section
}

func (resolve *Resolve) VisitString(booklet.String) error {
	return nil
}

func (resolve *Resolve) VisitSequence(con booklet.Sequence) error {
	for _, c := range con {
		err := c.Visit(resolve)
		if err != nil {
			return err
		}
	}

	return nil
}

func (resolve *Resolve) VisitParagraph(con booklet.Paragraph) error {
	for _, c := range con {
		err := c.Visit(resolve)
		if err != nil {
			return err
		}
	}

	return nil
}

func (resolve *Resolve) VisitPreformatted(con booklet.Preformatted) error {
	for _, c := range con {
		err := c.Visit(resolve)
		if err != nil {
			return err
		}
	}

	return nil
}

func (resolve *Resolve) VisitReference(con *booklet.Reference) error {
	tags := resolve.Section.FindTag(con.TagName)

	var err error
	switch len(tags) {
	case 0:
		err = booklet.UnknownTagError{
			TagName:     con.TagName,
			SimilarTags: resolve.Section.SimilarTags(con.TagName),
			ErrorLocation: booklet.ErrorLocation{
				FilePath:     resolve.Section.FilePath(),
				NodeLocation: con.Location,
				Length:       len("\\reference"),
			},
		}
	case 1:
		con.Tag = &tags[0]
	default:
		locs := []booklet.ErrorLocation{}
		for _, t := range tags {
			locs = append(locs, booklet.ErrorLocation{
				FilePath:     t.Section.FilePath(),
				NodeLocation: t.Location,
			})
		}

		err = booklet.AmbiguousReferenceError{
			TagName:          con.TagName,
			DefinedLocations: locs,
			ErrorLocation: booklet.ErrorLocation{
				FilePath:     resolve.Section.FilePath(),
				NodeLocation: con.Location,
				Length:       len("\\reference"),
			},
		}
	}

	if err == nil {
		return nil
	}

	if resolve.AllowBrokenReferences {
		fmt.Println("broken reference")

		con.Tag = &booklet.Tag{
			Name:     con.TagName,
			Anchor:   "broken",
			Title:    booklet.String(fmt.Sprintf("{broken reference: %s}", con.TagName)),
			Section:  resolve.Section,
			Location: con.Location,
		}

		return nil
	}

	return err
}

func (resolve *Resolve) VisitSection(con *booklet.Section) error {
	err := con.Title.Visit(resolve)
	if err != nil {
		return err
	}

	err = con.Body.Visit(resolve)
	if err != nil {
		return err
	}

	for _, p := range con.Partials {
		err = p.Visit(resolve)
		if err != nil {
			return err
		}
	}

	// TODO: this probably does redundant resolving, since i think the section
	// was loaded via a processor in the first place
	for _, child := range con.Children {
		subResolver := &Resolve{
			AllowBrokenReferences: resolve.AllowBrokenReferences,
			Section:               child,
		}

		err := child.Visit(subResolver)
		if err != nil {
			return err
		}
	}

	return nil
}

func (resolve *Resolve) VisitTableOfContents(booklet.TableOfContents) error {
	return nil
}

func (resolve *Resolve) VisitStyled(con booklet.Styled) error {
	err := con.Content.Visit(resolve)
	if err != nil {
		return err
	}

	for _, v := range con.Partials {
		if v == nil {
			continue
		}

		err := v.Visit(resolve)
		if err != nil {
			return err
		}
	}

	return nil
}

func (resolve *Resolve) VisitTarget(con booklet.Target) error {
	err := con.Title.Visit(resolve)
	if err != nil {
		return err
	}

	if con.Content != nil {
		err := con.Content.Visit(resolve)
		if err != nil {
			return err
		}
	}

	return nil
}

func (resolve *Resolve) VisitImage(con booklet.Image) error {
	return nil
}

func (resolve *Resolve) VisitList(con booklet.List) error {
	for _, c := range con.Items {
		err := c.Visit(resolve)
		if err != nil {
			return err
		}
	}

	return nil
}

func (resolve *Resolve) VisitLink(con booklet.Link) error {
	return con.Content.Visit(resolve)
}

func (resolve *Resolve) VisitTable(con booklet.Table) error {
	for _, row := range con.Rows {
		for _, c := range row {
			err := c.Visit(resolve)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (resolve *Resolve) VisitDefinitions(con booklet.Definitions) error {
	for _, def := range con {
		err := def.Subject.Visit(resolve)
		if err != nil {
			return err
		}

		err = def.Definition.Visit(resolve)
		if err != nil {
			return err
		}
	}

	return nil
}
