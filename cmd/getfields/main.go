package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"maps"
	"os"
	"slices"
	"strings"

	"github.com/juju/juju/cmd/juju/status"
)

func main() {
	structs := status.GetFields()

	classes := make(map[string]string)
	successors := make(map[string][]string)

	structNames := slices.Sorted(maps.Keys(structs))
	for _, name := range structNames {
		var buf bytes.Buffer

		fmt.Fprintln(&buf, "\n\n@dataclasses.dataclass")
		className := strings.Title(name)
		className = strings.ReplaceAll(className, "Application", "App")
		fmt.Fprintf(&buf, "class %s:\n", className)
		successors[className] = nil
		var required []string
		var optional []string
		for _, field := range structs[name] {
			pythonField := getPythonField(field.JSONField)
			pythonType := ""
			if _, ok := structs[field.Type]; ok {
				pythonType, _ = getPythonType(structs, field.Type)
				successors[className] = append(successors[className], pythonType)
			} else {
				var base string
				pythonType, base = getPythonType(structs, field.Type)
				if base != "" {
					successors[className] = append(successors[className], base)
				}
			}
			if field.OmitEmpty {
				pythonType += " | None = None"
			}
			line := fmt.Sprintf("    %s: %s\n", pythonField, pythonType)
			if field.OmitEmpty {
				optional = append(optional, line)
			} else {
				required = append(required, line)
			}
		}
		for _, line := range required {
			fmt.Fprint(&buf, line)
		}
		if len(required) > 0 && len(optional) > 0 {
			fmt.Fprintln(&buf)
		}
		for _, line := range optional {
			fmt.Fprint(&buf, line)
		}

		fmt.Fprintf(&buf, "\n    @classmethod\n")
		fmt.Fprintf(&buf, "    def from_dict(cls, d: dict[str, Any]) -> %s:\n", className)
		fmt.Fprintln(&buf, "        return cls(")
		for _, field := range structs[name] {
			pythonField := getPythonField(field.JSONField)
			pythonType, _ := getPythonType(structs, field.Type)
			dictGetter := getDictGetter(pythonType, field.JSONField, field.OmitEmpty)
			fmt.Fprintf(&buf, "            %s=%s,\n", pythonField, dictGetter)
		}
		fmt.Fprintln(&buf, "        )")

		classes[className] = buf.String()
		buf.Reset()
	}

	var order []string
	for _, names := range tarjanSort(successors) {
		if len(names) > 1 {
			panic(fmt.Sprintf("dependency loop: %s", strings.Join(names, ", ")))
		}
		order = append(order, names[0])
	}
	fmt.Println("from __future__ import annotations")
	fmt.Println("import dataclasses")
	fmt.Println("from typing import Any")
	for _, name := range order {
		fmt.Print(classes[name])
	}

	f, err := os.Create("structs.json")
	if err != nil {
		log.Fatal(err)
	}
	b, err := json.MarshalIndent(structs, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	_, err = f.Write(b)
	if err != nil {
		log.Fatal(err)
	}
	err = f.Close()
	if err != nil {
		log.Fatal(err)
	}
}

func getPythonField(s string) string {
	s = strings.ReplaceAll(s, "-", "_")
	s = strings.ReplaceAll(s, "application", "app")
	return s
}

func getPythonType(structs map[string][]status.FieldInfo, goType string) (string, string) {
	switch {
	case strings.HasPrefix(goType, "[]"):
		inner, base := getPythonType(structs, goType[2:])
		return "list[" + inner + "]", base
	case strings.HasPrefix(goType, "map[string]"):
		inner, base := getPythonType(structs, goType[11:])
		return "dict[str, " + inner + "]", base
	case goType == "string":
		return "str", ""
	case goType == "bool":
		return "bool", ""
	case goType == "int" || goType == "uint64":
		return "int", ""
	default:
		if _, ok := structs[goType]; !ok {
			fmt.Fprintf(os.Stderr, "# unhandled Go type: %s\n", goType)
		}
		pythonType := strings.Title(goType)
		pythonType = strings.ReplaceAll(pythonType, "Application", "App")
		return pythonType, pythonType
	}
}

func getDictGetter(pythonType string, jsonField string, omitEmpty bool) string {
	s := fmt.Sprintf("d['%s']", jsonField)
	orig := s
	s = doType(pythonType, s)
	if omitEmpty {
		if s == orig {
			// shortcut for simple value lookup
			return fmt.Sprintf("d.get('%s')", jsonField)
		}
		s += fmt.Sprintf(" if '%s' in d else None", jsonField)
	}
	return s
}

func doType(pythonType string, value string) string {
	switch {
	case strings.HasPrefix(pythonType, "list["):
		t := pythonType[5 : len(pythonType)-1]
		inner := doType(t, "x")
		if inner == "x" {
			return value
		}
		return fmt.Sprintf("[%s for x in %s]", inner, value)
	case strings.HasPrefix(pythonType, "dict[str, "):
		t := pythonType[10 : len(pythonType)-1]
		inner := doType(t, "v")
		if inner == "v" {
			return value
		}
		return fmt.Sprintf("{k: %s for k, v in %s.items()}", inner, value)
	case pythonType == "str" || pythonType == "int" || pythonType == "bool":
		return value
	default:
		return fmt.Sprintf("%s.from_dict(%s)", pythonType, value)
	}
}
