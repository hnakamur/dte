package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"flag"
	"io"
	"os"

	"github.com/BurntSushi/toml"
	"gopkg.in/flosch/pongo2.v3"
)

const Version = "0.2"

var InvalidContextError = errors.New("Invalid pongo2 context")

var (
	jsonFilename   string
	tomlFilename   string
	outputFilename string
	versionFlag    bool
)

func init() {
	flag.StringVar(&jsonFilename, "j", "", "json filename (use - for stdin)")
	flag.StringVar(&outputFilename, "o", "-", "output filename (default stdout)")
	flag.StringVar(&tomlFilename, "t", "", "toml filename (use - for stdin)")
	flag.BoolVar(&versionFlag, "v", false, "show version and exit")
}

func main() {
	flag.Parse()

	if versionFlag {
		println(Version)
		return
	}

	if flag.NArg() == 0 {
		showErrAndExit(errors.New("Please pass template filename as arguments"))
	}
	if jsonFilename == "" && tomlFilename == "" {
		showErrAndExit(errors.New("Please specify either json or toml."))
	}
	if jsonFilename != "" && tomlFilename != "" {
		showErrAndExit(errors.New("Please specify either json or toml, not both."))
	}

	var err error
	tpl, err := buildTemplateFromFile(flag.Arg(0))
	if err != nil {
		showErrAndExit(err)
	}

	var ctx pongo2.Context
	if jsonFilename != "" {
		if jsonFilename == "-" {
			ctx, err = buildContextFromJsonReader(os.Stdin)
		} else {
			ctx, err = buildContextFromJsonFile(jsonFilename)
		}
	} else {
		if tomlFilename == "-" {
			ctx, err = buildContextFromTomlReader(os.Stdin)
		} else {
			ctx, err = buildContextFromTomlFile(tomlFilename)
		}
	}

	if err != nil {
		showErrAndExit(err)
	}

	if outputFilename == "-" {
		err = executeTemplateWriter(tpl, ctx, os.Stdout)
	} else {
		err = executeTemplateFile(tpl, ctx, outputFilename)
	}
	if err != nil {
		showErrAndExit(err)
	}
}

func showErrAndExit(err error) {
	println(err.Error())
	os.Exit(1)
}

func buildContextFromJsonFile(filename string) (pongo2.Context, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return buildContextFromJsonReader(f)
}

func buildContextFromJsonReader(r io.Reader) (pongo2.Context, error) {
	var v interface{}
	err := json.NewDecoder(bufio.NewReader(r)).Decode(&v)
	if err != nil {
		return nil, err
	}
	return convertToContext(v)
}

func buildContextFromTomlFile(filename string) (pongo2.Context, error) {
	var v interface{}
	_, err := toml.DecodeFile(filename, &v)
	if err != nil {
		return nil, err
	}
	return convertToContext(v)
}

func buildContextFromTomlReader(r io.Reader) (pongo2.Context, error) {
	var v interface{}
	_, err := toml.DecodeReader(bufio.NewReader(r), &v)
	if err != nil {
		return nil, err
	}
	return convertToContext(v)
}

func convertToContext(v interface{}) (pongo2.Context, error) {
	c, ok := v.(map[string]interface{})
	if !ok {
		return nil, InvalidContextError
	}
	return c, nil
}

func buildTemplateFromFile(filename string) (*pongo2.Template, error) {
	return pongo2.FromFile(filename)
}

func executeTemplateFile(tpl *pongo2.Template, ctx pongo2.Context, filename string) error {
	w, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer w.Close()
	return executeTemplateWriter(tpl, ctx, w)
}

func executeTemplateWriter(tpl *pongo2.Template, ctx pongo2.Context, writer io.Writer) error {
	w := bufio.NewWriter(writer)
	err := tpl.ExecuteWriter(ctx, w)
	if err != nil {
		return err
	}
	return w.Flush()
}
