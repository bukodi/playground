package parsedocx

import (
	"strings"
)

type RNGModeEnum uint32

const (
	RngModeUnspecified RNGModeEnum = 0x00000001
)

var RNGModeEnumNames = map[RNGModeEnum]string{
	RngModeUnspecified: "Unspecified",
}

type EnumLiteral struct {
	Comment string
	Name    string
	Value   uint32
}

type EnumType struct {
	Name     string
	Literals []EnumLiteral
}

func ParseEnumSection(section *Section) (*EnumType, error) {
	enum := &EnumType{
		Name: section.Heading,
	}
	for _, child := range section.Children {
		enum.Literals = append(enum.Literals, EnumLiteral{
			Name:  child.Heading,
			Value: uint32(child.Index),
		})
	}
	return enum, nil
}

func toCamelCase(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return s
	}

	words := strings.Fields(s)
	for i := 0; i < len(words); i++ {
		words[i] = strings.Title(words[i])
	}

	result := strings.Join(words, "")
	//result = string(unicode.ToLower(rune(result[0]))) + result[1:]
	return result
}
