package goxsdpkg

import (
	xsdt "github.com/metaleap/go-xsd/types"
)

//	The actual definition is
//	baseline | sub | super | <percentage> | <length> | inherit
//	not sure that union can do this
type BaselineShiftValueType string

//	Space-separated list of classes

//	<shape> | auto | inherit
type ClipValueType string

//	<uri> | none | inherit
type ClipPathValueType string

//	'clip-rule' or fill-rule property/attribute value
type ClipFillRuleType string

//	media type, as per [RFC2045]
//	media type, as per [RFC2045]
type ContentTypeType string

//	a <co-ordinate>
//	a coordinate, which is a number optionally followed immediately by a unit identifier. Perhaps it is possible to represent this as a union by declaring unit idenifiers as a type?
type CoordinateType string

//	a space separated list of CoordinateType. Punt to 'string' for now
type CoordinatesType string

//	a CSS2 Color
//	Color as defined in CSS2 and XSL 1.0 plus additional recognised color keyword names (the 'X11 colors')
type ColorType string

//	Value is an optional comma-separated list orf uri references followed by one token from an enumerated list.
//	[ [<uri> ,]* [ auto | crosshair | default | pointer | move | e-resize | ne-resize | nw-resize | n-resize | se-resize | sw-resize | s-resize | w-resize| text | wait | help ] ] | inherit
type CursorValueType string

//	accumulate | new [ <x> <y> <width> <height> ] | inherit
type EnableBackgroundValueType string

//	extension list specification
type ExtensionListType string

//	feature list specification
type FeatureListType string

//	<uri> | none | inherit
type FilterValueType string

//	[[ <family-name> | <generic-family> ],]* [<family-name> | <generic-family>] | inherit
//	'font-family' property/attribute value (i.e., list of fonts)
type FontFamilyValueType string

//	'font-size' property/attribute value
//	<absolute-size> | <relative-size> | <length> | <percentage> | inherit
type FontSizeValueType string

//	'font-size-adjust' property/attribute value
//	<number> | none | inherit
type FontSizeAdjustValueType string

//	'glyph-orientation-horizontal' property/attribute value (e.g., <angle>)
//	<angle> | inherit
type GlyphOrientationHorizontalValueType string

//	'glyph-orientation-vertical' property/attribute value (e.g., 'auto', <angle>)
//	auto | <angle> | inherit
type GlyphOrientationVerticalValueType string

//	'kerning' property/attribute value (e.g., auto | <length>)
//	auto | <length> | inherit
type KerningValue string

//	a language code, as per [RFC3066]
type LanguageCodeType string

//	a comma-separated list of language codes, as per [RFC3066]
type LanguageCodesType string

//	a <length>
type LengthType string

//	a list of <length>s
type LengthsType string

//	link to this target
type LinkTargetType string

//	'marker' property/attribute value (e.g., 'none', %URI;)
type MarkerValueType string

//	'mask' property/attribute value (e.g., 'none', %URI;)
//	<uri> | none | inherit
type MaskValueType string

//	comma-separated list of media descriptors.
type MediaDescType string

//	list of <number>s, but at least one and at most two
type NumberOptionalNumberType string

//	a <number> or a  <percentage>
type NumberOrPercentageType string

//	list of <number>s
type NumbersType string

//	opacity value (e.g., <number>)
//	<alphavalue> | inherit
type OpacityValueType string

//	a 'fill' or 'stroke' property/attribute value
type PaintType string

//	a path data specification
//	Yes, of course this was generated by a program!
type PathDataType string

//	a list of points
type PointsType string

//	'preserveAspectRatio' attribute specification
type PreserveAspectRatioSpecType string

//	script expression
type ScriptType string

//	'letter-spacing' or 'word-spacing' property/attribute value (e.g., normal | <length>)
type SpacingValueType string

//	'stroke-dasharray' property/attribute value (e.g., 'none', list of <number>s)
type StrokeDashArrayValueType string

//	'stroke-dashoffset' property/attribute value (e.g., 'none', >length>)
type StrokeDashOffsetValueType string

//	'stroke-miterlimit' property/attribute value (e.g., <number>)
type StrokeMiterLimitValueType string

//	'stroke-width' property/attribute value (e.g., <length>)
type StrokeWidthValueType string

//	style sheet data
type StyleSheetType string

//	An SVG color value (sRGB plus optional ICC)
type SVGColorType string

//	'text-decoration' property/attribute value (e.g., 'none', 'underline')
type TextDecorationValueType string

//	Yes, of course this was generated by a program!
//	list of transforms
type TransformListType string

//	'viewBox' attribute specification
type ViewBoxSpecType string
