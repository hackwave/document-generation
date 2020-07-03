package render

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	booklet "../booklet"
	text "../render/text"
)

var initTextTmpl *template.Template

func init() {
	initTextTmpl = template.New("engine").Funcs(template.FuncMap{
		"url": func(ext string, tag booklet.Tag) (string, error) {
			return "", errors.New("url stubbed out")
		},

		"htmlURL": func(tag booklet.Tag) string {
			return sectionURL("html", tag.Section, tag.Anchor)
		},

		"stripAux": booklet.StripAux,

		"render": func(booklet.Content) (string, error) {
			return "", errors.New("render stubbed out")
		},

		"walkContext": func(current *booklet.Section, section *booklet.Section) WalkContext {
			return WalkContext{
				Current: current,
				Section: section,
			}
		},

		"headerDepth": func(con *booklet.Section) int {
			depth := con.PageDepth() + 1
			if depth > 6 {
				depth = 6
			}

			return depth
		},

		"joinLines": func(prefix string, str string) string {
			return strings.Join(strings.Split(str, "\n"), "\n"+prefix)
		},
	})

	for _, asset := range text.AssetNames() {
		info, err := text.AssetInfo(asset)
		if err != nil {
			panic(err)
		}

		content := strings.TrimRight(string(text.MustAsset(asset)), "\n")

		_, err = initTextTmpl.New(filepath.Base(info.Name())).Parse(content)
		if err != nil {
			panic(err)
		}
	}
}

type TextRenderingEngine struct {
	fileExtension string

	tmpl         *template.Template
	tmplModTimes map[string]time.Time

	template *template.Template
	data     interface{}
}

func NewTextRenderingEngine(fileExtension string) *TextRenderingEngine {
	engine := &TextRenderingEngine{
		fileExtension: fileExtension,

		tmplModTimes: map[string]time.Time{},
	}

	engine.resetTmpl()

	return engine
}

func (engine *TextRenderingEngine) resetTmpl() {
	engine.tmpl = template.Must(initTextTmpl.Clone())
	engine.tmpl.Funcs(template.FuncMap{
		"render": engine.subRender,
		"url": func(tag booklet.Tag) string {
			return sectionURL(engine.FileExtension(), tag.Section, tag.Anchor)
		},
	})
}

func (engine *TextRenderingEngine) LoadTemplates(templatesDir string) error {
	templates, err := filepath.Glob(filepath.Join(templatesDir, "*.tmpl"))
	if err != nil {
		return err
	}

	var shouldReload bool
	for _, path := range templates {
		info, err := os.Stat(path)
		if err != nil {
			return err
		}

		modTime := info.ModTime()

		lastModTime, found := engine.tmplModTimes[path]
		if !found || modTime.After(lastModTime) {
			shouldReload = true
		}

		engine.tmplModTimes[path] = modTime
	}

	if !shouldReload {
		return nil
	}

	engine.resetTmpl()

	for _, path := range templates {
		content, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}

		trimmed := strings.TrimRight(string(content), "\n")

		_, err = engine.tmpl.New(filepath.Base(path)).Parse(trimmed)
		if err != nil {
			return err
		}
	}

	return nil
}

func (engine *TextRenderingEngine) FileExtension() string {
	return engine.fileExtension
}

func (engine *TextRenderingEngine) URL(tag booklet.Tag) string {
	return sectionURL(engine.FileExtension(), tag.Section, tag.Anchor)
}

func (engine *TextRenderingEngine) RenderSection(out io.Writer, con *booklet.Section) error {
	tmpl := "page"
	if con.Style != "" {
		tmpl = con.Style + "-page"
	}

	engine.data = con

	err := engine.setTmpl(tmpl)
	if err != nil {
		return err
	}

	return engine.render(out)
}

func (engine *TextRenderingEngine) VisitString(con booklet.String) error {
	engine.data = con
	return engine.setTmpl("string")
}

func (engine *TextRenderingEngine) VisitReference(con *booklet.Reference) error {
	engine.data = con
	return engine.setTmpl("reference")
}

func (engine *TextRenderingEngine) VisitSection(con *booklet.Section) error {
	tmpl := "section"
	if con.Style != "" {
		tmpl = con.Style
	}

	engine.data = con
	return engine.setTmpl(tmpl)
}

func (engine *TextRenderingEngine) VisitSequence(con booklet.Sequence) error {
	engine.data = con
	return engine.setTmpl("sequence")
}

func (engine *TextRenderingEngine) VisitParagraph(con booklet.Paragraph) error {
	engine.data = con
	return engine.setTmpl("paragraph")
}

func (engine *TextRenderingEngine) VisitPreformatted(con booklet.Preformatted) error {
	engine.data = con
	return engine.setTmpl("preformatted")
}

func (engine *TextRenderingEngine) VisitTableOfContents(con booklet.TableOfContents) error {
	engine.data = con.Section
	return engine.setTmpl("toc")
}

func (engine *TextRenderingEngine) VisitStyled(con booklet.Styled) error {
	engine.data = con
	return engine.setTmpl(string(con.Style))
}

func (engine *TextRenderingEngine) VisitTarget(con booklet.Target) error {
	engine.data = con
	return engine.setTmpl("target")
}

func (engine *TextRenderingEngine) VisitImage(con booklet.Image) error {
	engine.data = con
	return engine.setTmpl("image")
}

func (engine *TextRenderingEngine) VisitList(con booklet.List) error {
	engine.data = con
	return engine.setTmpl("list")
}

func (engine *TextRenderingEngine) VisitLink(con booklet.Link) error {
	engine.data = con
	return engine.setTmpl("link")
}

func (engine *TextRenderingEngine) VisitTable(con booklet.Table) error {
	engine.data = con
	return engine.setTmpl("table")
}

func (engine *TextRenderingEngine) VisitDefinitions(con booklet.Definitions) error {
	engine.data = con
	return engine.setTmpl("definitions")
}

func (engine *TextRenderingEngine) setTmpl(name string) error {
	tmpl := engine.tmpl.Lookup(name + ".tmpl")

	if tmpl == nil {
		return fmt.Errorf("template '%s.tmpl' not found", name)
	}

	engine.template = tmpl

	return nil
}

func (engine *TextRenderingEngine) render(out io.Writer) error {
	if engine.template == nil {
		return fmt.Errorf("unknown template for '%s' (%T)", engine.data, engine.data)
	}

	return engine.template.Execute(out, engine.data)
}

func (engine *TextRenderingEngine) subRender(content booklet.Content) (string, error) {
	buf := new(bytes.Buffer)

	subEngine := NewTextRenderingEngine(engine.fileExtension)
	subEngine.tmpl = engine.tmpl

	err := content.Visit(subEngine)
	if err != nil {
		return "", err
	}

	err = subEngine.render(buf)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}
