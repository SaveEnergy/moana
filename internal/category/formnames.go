package category

// FieldName, FieldIcon, FieldID, FieldColor are HTML form control names for category HTTP handlers and tests.
const (
	FieldName        = "name"
	FieldIcon        = "icon"
	FieldID          = "id"
	FieldColor       = "color"
	FieldColorCustom = "color_custom"
)

// ColorPickCustom is the FieldColor value that selects FieldColorCustom (preset swatches vs custom hex).
const ColorPickCustom = "custom"
