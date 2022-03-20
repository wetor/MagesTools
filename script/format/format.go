package format

const (
	LineBreak             = 0x00
	NameStart             = 0x01
	LineStart             = 0x02
	Present               = 0x03
	SetColor              = 0x04
	PresentUnknown05      = 0x05
	PresentResetAlignment = 0x08
	RubyBaseStart         = 0x09
	RubyTextStart         = 0x0A
	RubyTextEnd           = 0x0B
	SetFontSize           = 0x0C
	PrintInParallel       = 0x0E
	PrintInCenter         = 0x0F
	SetMarginTop          = 0x11
	SetMarginLeft         = 0x12
	GetHardcodedValue     = 0x13
	EvaluateExpression    = 0x15
	PresentUnknown18      = 0x18
	AutoForward           = 0x19
	AutoForward1A         = 0x1A
	RubyCenterPerChar     = 0x1E
	AltLineBreak          = 0x1F
	Terminator            = 0xFF
)

type Format interface {
	SetCharset(decode map[uint16]string, encode map[string]uint16)
	DecodeLine(data []byte) string
	EncodeLine(str string) []byte
}
