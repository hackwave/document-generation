package chroma

import (
	"bytes"

	booklet "../booklet"

	"github.com/alecthomas/chroma/formatters/html"
	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"
)

func NewPlugin(section *booklet.Section) booklet.Plugin {
	return Plugin{
		section: section,
	}
}

type Plugin struct {
	section *booklet.Section
}

func (plugin Plugin) Syntax(language string, code booklet.Content, styleName ...string) (booklet.Content, error) {
	lexer := lexers.Get(language)
	if lexer == nil {
		lexer = lexers.Fallback
	}

	iterator, err := lexer.Tokenise(nil, code.String())
	if err != nil {
		return nil, err
	}

	formatter := html.New(html.PreventSurroundingPre(code.IsFlow()))

	chromaStyle := styles.Fallback
	if len(styleName) > 0 {
		chromaStyle = styles.Get(styleName[0])
	}

	buf := new(bytes.Buffer)
	err = formatter.Format(buf, chromaStyle, iterator)
	if err != nil {
		return nil, err
	}

	var style booklet.Style
	if code.IsFlow() {
		style = "inline-code"
	} else {
		style = "code-block"
	}

	return booklet.Styled{
		Style:   style,
		Block:   !code.IsFlow(),
		Content: code,
		Partials: booklet.Partials{
			"Language": booklet.String(language),
			"HTML":     booklet.String(buf.String()),
		},
	}, nil
}
