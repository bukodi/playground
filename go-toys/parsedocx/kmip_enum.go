package parsedocx

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
