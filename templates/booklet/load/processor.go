package load

import (
	"os"
	"sync"
	"time"

	ast "../ast"
	booklet "../booklet"
	stages "../stages"

	"github.com/sirupsen/logrus"
)

type Processor struct {
	AllowBrokenReferences bool

	parsed  map[string]parsedNode
	parsedL sync.Mutex
}

type parsedNode struct {
	Node    ast.Node
	ModTime time.Time
}

func (processor *Processor) LoadFile(path string, pluginFactories []booklet.PluginFactory) (*booklet.Section, error) {
	return processor.LoadFileIn(nil, path, pluginFactories)
}

func (processor *Processor) LoadFileIn(parent *booklet.Section, path string, pluginFactories []booklet.PluginFactory) (*booklet.Section, error) {
	section, err := processor.EvaluateFile(parent, path, pluginFactories)
	if err != nil {
		return nil, err
	}

	return processor.runStages(section)
}

func (processor *Processor) EvaluateFile(parent *booklet.Section, path string, pluginFactories []booklet.PluginFactory) (*booklet.Section, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	modTime := info.ModTime()

	processor.parsedL.Lock()
	if processor.parsed == nil {
		processor.parsed = map[string]parsedNode{}
	}
	parsed, found := processor.parsed[path]
	processor.parsedL.Unlock()

	log := logrus.WithFields(logrus.Fields{
		"path": path,
	})

	var node ast.Node
	if found && !modTime.After(parsed.ModTime) {
		log.Debug("already parsed section")
		node = parsed.Node
	} else {
		log.Debug("parsing section")

		file, err := os.Open(path)
		if err != nil {
			return nil, err
		}

		result, err := ast.ParseReader(path, file)
		if err != nil {
			err, loc, ok := ast.UnpackError(err)
			if !ok {
				return nil, err
			}

			return nil, booklet.ParseError{
				Err: err,
				ErrorLocation: booklet.ErrorLocation{
					FilePath:     path,
					NodeLocation: loc,
					Length:       1,
				},
			}
		}

		err = file.Close()
		if err != nil {
			return nil, err
		}

		node = result.(ast.Node)
	}

	section := &booklet.Section{
		Parent: parent,

		Path: path,

		Title: booklet.Empty,
		Body:  booklet.Empty,

		Processor: processor,
	}

	err = processor.evaluateSection(section, node, pluginFactories)
	if err != nil {
		return nil, err
	}

	processor.parsedL.Lock()
	processor.parsed[path] = parsedNode{
		Node:    node,
		ModTime: modTime,
	}
	processor.parsedL.Unlock()

	return section, nil
}

func (processor *Processor) EvaluateNode(parent *booklet.Section, node ast.Node, pluginFactories []booklet.PluginFactory) (*booklet.Section, error) {
	section := &booklet.Section{
		Parent: parent,

		Title: booklet.Empty,
		Body:  booklet.Empty,

		Processor: processor,
	}

	err := processor.evaluateSection(section, node, pluginFactories)
	if err != nil {
		return nil, err
	}

	return section, nil
}

func (processor *Processor) evaluateSection(section *booklet.Section, node ast.Node, pluginFactories []booklet.PluginFactory) error {
	for _, pf := range pluginFactories {
		section.UsePlugin(pf)
	}

	evaluator := &stages.Evaluate{
		Section: section,
	}

	err := node.Visit(evaluator)
	if err != nil {
		return err
	}

	if evaluator.Result != nil {
		section.Body = evaluator.Result
	}

	return nil
}

func (processor *Processor) runStages(section *booklet.Section) (*booklet.Section, error) {
	collector := &stages.Collect{
		Section: section,
	}

	err := section.Visit(collector)
	if err != nil {
		return nil, err
	}

	resolver := &stages.Resolve{
		AllowBrokenReferences: processor.AllowBrokenReferences,

		Section: section,
	}

	err = section.Visit(resolver)
	if err != nil {
		return nil, err
	}

	return section, nil
}
