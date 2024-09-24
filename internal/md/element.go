package md

type blockElement struct {
	kind     blockElementKind
	children inlineElement
	// 1 or greater than 1
	listLevel         int
	imageSrc          string
	imageCaption      string
	codeName          string
	codeText          string
	checkboxList      bool
	checkboxIsChecked bool
}

type blockElementKind int

const (
	blockElementKindParagraph blockElementKind = iota
	blockElementKindList
	blockElementKindImage
	blockElementKindCodeBlock
	blockElementKindHeading1
	blockElementKindHeading2
	blockElementKindHeading3
)

type inlineElementKind int

const (
	inlineElementKindRoot inlineElementKind = iota
	inlineElementKindText
	inlineElementKindBold
	inlineElementKindCode
	inlineElementKindLink
)

type inlineElement struct {
	kind     inlineElementKind
	s        string
	children []inlineElement

	linkHref string
}

type token struct {
	// Markdown のキーワードか否か
	r bool
	// トークンに含まれる文字列
	s string
}
