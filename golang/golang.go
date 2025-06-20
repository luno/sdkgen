package golang

import (
	"bytes"
	"embed"
	"fmt"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/luno/sdkgen/clientgen"
	"golang.org/x/tools/imports"
)

//go:embed templates/*
var templates embed.FS

// Generate generates a Go API client.
func Generate(api clientgen.API) ([]clientgen.File, error) {
	var fl []clientgen.File

	// Generate api.go.
	f, err := generateGofile("api.go", "api.go.tpl", api)
	if err != nil {
		return nil, err
	}
	fl = append(fl, f)

	// Generate types.go.
	f, err = generateGofile("types.go", "types.go.tpl", api)
	if err != nil {
		return nil, err
	}
	fl = append(fl, f)

	return fl, nil
}

func generateGofile(fileName, tplName string, api clientgen.API) (
	clientgen.File, error) {

	funcMap := template.FuncMap{
		"opname":    opname,
		"propname":  propname,
		"typename":  typename,
		"enumvalue": enumvalue,
	}

	tpl, err := template.New(tplName).Funcs(funcMap).ParseFS(templates,
		filepath.Join("templates", tplName))
	if err != nil {
		return clientgen.File{}, err
	}

	f := clientgen.NewFile(fileName)

	c := struct {
		API clientgen.API
	}{
		API: api,
	}

	raw := bytes.NewBuffer(nil)
	if err := tpl.Execute(raw, c); err != nil {
		return clientgen.File{}, err
	}

	formatted, err := imports.Process("", raw.Bytes(), nil)
	if err != nil {
		return clientgen.File{}, err
	}

	_, err = f.Write(formatted)
	if err != nil {
		return clientgen.File{}, err
	}

	return f, nil
}

func opname(s string) string {
	return strings.Title(s)
}

func propname(s string) string {
	b := bytes.NewBuffer(nil)
	for i := 0; i < len(s); i++ {
		r := s[i]
		if r == '_' {
			if i < len(s)-1 {
				fmt.Fprintf(b, strings.ToUpper(string(s[i+1])))
				i++
			}
			continue
		}
		if i == 0 {
			fmt.Fprint(b, strings.ToUpper(string(r)))
			continue
		}
		fmt.Fprint(b, string(r))
	}
	return string(b.Bytes())
}

func typename(t clientgen.Type) string {
	switch t.Kind {
	case clientgen.KindFloat:
		return "float64"
	case clientgen.KindInteger:
		return "int64"
	case clientgen.KindDecimal:
		return "decimal.Decimal"
	case clientgen.KindString:
		return "string"
	case clientgen.KindBoolean:
		return "bool"
	case clientgen.KindEnum, clientgen.KindStruct:
		return t.Name
	case clientgen.KindArray:
		return "[]" + typename(t.ArrayProps.Type)
	case clientgen.KindTimestamp:
		return "Time"
	case clientgen.KindObject:
		return "map[string]" + typename(t.ObjectProps.ValueType)
	}
	return ""
}

func enumvalue(t, v string) string {
	return t + strings.Title(strings.ToLower(clientgen.SanitiseVariableName(v)))
}
