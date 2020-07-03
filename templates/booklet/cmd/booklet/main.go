package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"

	booklet "../../booklet"
	load "../../load"
	plugin "../../plugin"
	render "../../render"

	flags "github.com/jessevdk/go-flags"
)

func main() {
	cmd := &Command{}
	cmd.Version = func() {
		fmt.Println("v0.1.0")
		os.Exit(0)
	}

	parser := flags.NewParser(cmd, flags.Default)
	parser.NamespaceDelimiter = "-"

	args, err := parser.Parse()
	if err != nil {
		if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
			fmt.Println(err)
			os.Exit(0)
		} else {
			os.Exit(1)
		}
	}

	err = cmd.Execute(args)
	if err != nil {
		if prettyErr, ok := err.(booklet.PrettyError); ok {
			prettyErr.PrettyPrint(os.Stderr)
		} else {
			fmt.Fprintln(os.Stderr, err)
		}

		os.Exit(1)
	}
}

type Command struct {
	Version func() `short:"v" long:"version" description:"Print the version of Boooklit and exit."`

	In  string `long:"in"  short:"i" required:"true" description:"Input .lit file to load."`
	Out string `long:"out" short:"o" description:"Directory into which sections will be rendered."`

	SectionTag  string `long:"section-tag"  description:"Section tag to render."`
	SectionPath string `long:"section-path" description:"Section path to load and render with --in as its parent."`

	SaveSearchIndex bool `long:"save-search-index" description:"Save a search index JSON file in the destination."`

	ServerPort int `long:"serve" short:"s" description:"Start an HTTP server on the given port."`

	Plugins []string `long:"plugin" short:"p" description:"Package to import, providing a plugin."`

	Debug bool `long:"debug" short:"d" description:"Log at debug level."`

	AllowBrokenReferences bool `long:"allow-broken-references" description:"Replace broken references with a bogus tag."`

	HTMLEngine struct {
		Templates string `long:"templates" description:"Directory containing .tmpl files to load."`
	} `group:"HTML Rendering Engine" namespace:"html"`

	TextEngine struct {
		FileExtension string `long:"file-extension" description:"File extension to use for generated files."`
		Templates     string `long:"templates"      description:"Directory containing .tmpl files to load."`
	} `group:"Text Rendering Engine" namespace:"text"`
}

func (cmd *Command) Execute(args []string) error {
	isReexec := os.Getenv("BOOKLIT_REEXEC") != ""
	if !isReexec && len(cmd.Plugins) > 0 {
		fmt.Println("plugins configured; reexecing")
		return cmd.reexec()
	}

	if cmd.ServerPort != 0 {
		return cmd.Serve()
	} else {
		return cmd.Build()
	}
}

func (cmd *Command) Serve() error {
	http.Handle("/", &Server{
		In: cmd.In,
		Processor: &load.Processor{
			AllowBrokenReferences: cmd.AllowBrokenReferences,
		},

		Templates:  cmd.HTMLEngine.Templates,
		Engine:     render.NewHTMLRenderingEngine(),
		FileServer: http.FileServer(http.Dir(cmd.Out)),
	})

	fmt.Println("port ", cmd.ServerPort, " listening")

	return http.ListenAndServe(fmt.Sprintf(":%d", cmd.ServerPort), nil)
}

var basePluginFactories = []booklet.PluginFactory{
	plugin.NewPlugin,
}

func (cmd *Command) Build() error {
	processor := &load.Processor{
		AllowBrokenReferences: cmd.AllowBrokenReferences,
	}

	var engine render.RenderingEngine
	if cmd.TextEngine.FileExtension != "" {
		textEngine := render.NewTextRenderingEngine(cmd.TextEngine.FileExtension)

		if cmd.TextEngine.Templates != "" {
			err := textEngine.LoadTemplates(cmd.TextEngine.Templates)
			if err != nil {
				return err
			}
		}

		engine = textEngine
	} else {
		htmlEngine := render.NewHTMLRenderingEngine()

		if cmd.HTMLEngine.Templates != "" {
			err := htmlEngine.LoadTemplates(cmd.HTMLEngine.Templates)
			if err != nil {
				return err
			}
		}

		engine = htmlEngine
	}

	section, err := processor.LoadFile(cmd.In, basePluginFactories)
	if err != nil {
		return err
	}

	sectionToRender := section
	if cmd.SectionTag != "" {
		tags := section.FindTag(cmd.SectionTag)
		if len(tags) == 0 {
			return fmt.Errorf("unknown tag: %s", cmd.SectionTag)
		}

		sectionToRender = tags[0].Section
	} else if cmd.SectionPath != "" {
		sectionToRender, err = processor.LoadFileIn(section, cmd.SectionPath, basePluginFactories)
		if err != nil {
			return err
		}
	}

	if cmd.Out == "" {
		return engine.RenderSection(os.Stdout, sectionToRender)
	}

	err = os.MkdirAll(cmd.Out, 0755)
	if err != nil {
		return err
	}

	writer := render.Writer{
		Engine:      engine,
		Destination: cmd.Out,
	}

	err = writer.WriteSection(sectionToRender)
	if err != nil {
		return err
	}

	if cmd.SaveSearchIndex {
		err = writer.WriteSearchIndex(section, "search_index.json")
		if err != nil {
			return err
		}
	}

	return nil
}

func (cmd *Command) reexec() error {
	tmpdir, err := ioutil.TempDir("", "booklet-reexec")
	if err != nil {
		return err
	}

	defer func() {
		_ = os.RemoveAll(tmpdir)
	}()

	src := filepath.Join(tmpdir, "main.go")
	bin := filepath.Join(tmpdir, "main")

	goSrc := "package main\n"
	goSrc += "import \"github.com/vito/booklet/bookletcmd\"\n"
	for _, p := range cmd.Plugins {
		goSrc += "import _ \"" + p + "\"\n"
	}
	goSrc += "func main() {\n"
	goSrc += "	bookletcmd.Main()\n"
	goSrc += "}\n"

	err = ioutil.WriteFile(src, []byte(goSrc), 0644)
	if err != nil {
		return err
	}

	build := exec.Command("go", "install", src)
	build.Env = append(os.Environ(), "GOBIN="+tmpdir)
	build.Stdout = os.Stdout
	build.Stderr = os.Stderr

	fmt.Println("building reexec binary")

	err = build.Run()
	if err != nil {
		return fmt.Errorf("build failed: %w", err)
	}

	run := exec.Command(bin, os.Args[1:]...)
	run.Env = append(os.Environ(), "BOOKLET_REEXEC=1")
	run.Stdout = os.Stdout
	run.Stderr = os.Stderr

	fmt.Println("reexecing")

	err = run.Run()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			os.Exit(exitErr.ExitCode())
			return nil
		}

		return fmt.Errorf("reexec failed: %w", err)
	}

	return nil
}

type Server struct {
	In        string
	Processor *load.Processor

	Templates string
	Engine    *render.HTMLRenderingEngine

	FileServer http.Handler

	buildLock sync.Mutex
}

func (server *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println("request @ ", r.URL.Path)

	fmt.Println("serving")

	section, found, err := server.loadRequestedSection(r.URL.Path)
	if err != nil {
		fmt.Errorf("failed to load section: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		booklet.ErrorPage(err, w)
		return
	}

	if !found {
		server.FileServer.ServeHTTP(w, r)
		return
	}

	server.buildLock.Lock()
	defer server.buildLock.Unlock()

	if len(server.Templates) != 0 {
		err := server.Engine.LoadTemplates(server.Templates)
		if err != nil {
			fmt.Errorf("failed to load templates: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			booklet.ErrorPage(err, w)
			return
		}
	}

	fmt.Println("on section", section.Path)

	fmt.Println("rendering")

	err = server.Engine.RenderSection(w, section)
	if err != nil {
		fmt.Errorf("failed to render: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		booklet.ErrorPage(err, w)
		return
	}

	return
}

func (server *Server) loadRequestedSection(path string) (*booklet.Section, bool, error) {
	ext := server.Engine.FileExtension()

	if path == "/" {
		path = "/index." + ext
	}

	if !strings.HasSuffix(path, "."+ext) {
		return nil, false, nil
	}

	tagName := strings.TrimSuffix(strings.TrimPrefix(path, "/"), "."+ext)

	fmt.Println("loading section", server.In)

	rootSection, err := server.Processor.LoadFile(server.In, basePluginFactories)
	if err != nil {
		return nil, false, err
	}

	tags := rootSection.FindTag(tagName)
	if len(tags) == 0 {
		return nil, false, nil
	}

	return tags[0].Section, true, nil
}
