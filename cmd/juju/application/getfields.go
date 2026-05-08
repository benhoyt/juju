package application

import (
	"reflect"
	"slices"
	"strings"
)

func GetFields() map[string][]FieldInfo {
	m := make(map[string][]FieldInfo)
	getFields(reflect.TypeOf(UnitInfo{}), m, "UnitInfo", 1)
	return m
}

func getFields(t reflect.Type, m map[string][]FieldInfo, typeName string, level int) {
	if _, ok := m[typeName]; ok {
		return
	}

	// Bit of a hack, but handle types like map[string][]Foo
	switch t.Kind() {
	case reflect.Map:
		t = t.Elem()
		typeName = strings.TrimPrefix(typeName, "map[string]")
	case reflect.Slice:
		t = t.Elem()
		typeName = strings.TrimPrefix(typeName, "[]")
	case reflect.Pointer:
		t = t.Elem()
		typeName = strings.TrimPrefix(typeName, "*")
	}
	if t.Kind() != reflect.Struct {
		return
	}

	m[typeName] = nil
	var result []FieldInfo
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		jsonTag := field.Tag.Get("json")
		if jsonTag == "" {
			jsonTag = field.Name
		}
		tagFields := strings.Split(jsonTag, ",")
		jsonField := tagFields[0]
		if jsonField == "-" {
			jsonField = ""
		}
		fieldType := field.Type.String()
		niceName := getNiceName(fieldType)
		result = append(result, FieldInfo{
			Name:      field.Name,
			Type:      niceName,
			JSONField: jsonField,
			Pointer:   fieldType[0] == '*',
			OmitEmpty: slices.Contains(tagFields[1:], "omitempty"),
		})
		if jsonField == "" {
			continue
		}
		switch field.Type.Kind() {
		case reflect.Struct:
			getFields(field.Type, m, niceName, level+1)
		case reflect.Map:
			elemType := field.Type.Elem()
			niceElemName := getNiceName(elemType.String())
			getFields(elemType, m, niceElemName, level+1)
		case reflect.Slice:
			elemType := field.Type.Elem()
			niceElemName := getNiceName(elemType.String())
			getFields(elemType, m, niceElemName, level+1)
		case reflect.Pointer:
			elemType := field.Type.Elem()
			niceElemName := getNiceName(elemType.String())
			getFields(elemType, m, niceElemName, level+1)
		}
	}
	m[typeName] = result
}

func getNiceName(s string) string {
	return strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(s, "common.", ""), "*", ""), "application.", "")
}

type FieldInfo struct {
	Name      string `json:"name"`
	Type      string `json:"type"`
	JSONField string `json:"json_field"`
	Pointer   bool   `json:"pointer"`
	OmitEmpty bool   `json:"omit_empty"`
}
