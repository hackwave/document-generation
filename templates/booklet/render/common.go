package render

import (
	booklet "../booklet"
)

type WalkContext struct {
	Current *booklet.Section
	Section *booklet.Section
}

func sectionPageOwner(section *booklet.Section) *booklet.Section {
	if section.Parent == nil {
		return section
	}

	if section.Parent.SplitSections {
		return section
	}

	return sectionPageOwner(section.Parent)
}

func sectionURL(ext string, section *booklet.Section, anchor string) string {
	owner := sectionPageOwner(section)

	if owner != section {
		if anchor == "" {
			anchor = section.PrimaryTag.Name
		}

		return sectionURL(ext, owner, anchor)
	}

	filename := section.PrimaryTag.Name + "." + ext

	if anchor != "" {
		filename += "#" + anchor
	}

	return filename
}
