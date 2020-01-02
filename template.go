package goss

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"
	"text/template"
)

func mkSlice(args ...interface{}) []interface{} {
	return args
}

func readFile(f string) (string, error) {
	b, err := ioutil.ReadFile(f)
	if err != nil {
		return "", err

	}
	return strings.TrimSpace(string(b)), nil
}

func getEnv(key string, def ...string) string {
	val := os.Getenv(key)
	if val == "" && len(def) > 0 {
		return def[0]
	}

	return os.Getenv(key)
}

func regexMatch(re, s string) (bool, error) {
	compiled, err := regexp.Compile(re)
	if err != nil {
		return false, err
	}

	return compiled.MatchString(s), nil
}

var funcMap = template.FuncMap{
	"mkSlice":    mkSlice,
	"readFile":   readFile,
	"getEnv":     getEnv,
	"regexMatch": regexMatch,
	"toUpper":    strings.ToUpper,
	"toLower":    strings.ToLower,
}

func NewTemplateFilter(varsFile string, varsInline string) func([]byte) []byte {
	vars, err := varsFromFile(varsFile)
	if err != nil {
		fmt.Printf("Error: loading vars file '%s'\n%v\n", varsFile, err)
		os.Exit(1)
	}

	varsExtra, err := varsFromString(varsInline)
	if err != nil {
		fmt.Printf("Error: loading inline vars\n%v\n", err)
		os.Exit(1)
	}

	for k, v := range varsExtra {
		vars[k] = v
	}

	tVars := &TmplVars{Vars: vars}

	f := func(data []byte) []byte {
		funcMap := funcMap
		t := template.New("test").Funcs(funcMap)
		tmpl, err := t.Parse(string(data))
		if err != nil {
			log.Fatal(err)
		}
		tmpl.Option("missingkey=error")
		var doc bytes.Buffer
		err = tmpl.Execute(&doc, tVars)
		if err != nil {
			log.Fatal(err)
		}
		return doc.Bytes()
	}
	return f
}
