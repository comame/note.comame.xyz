package md

import (
	"fmt"
	"html"
)

func blockElementsToHTML(elements []blockElement) string {
	ret := ""

	previousListLevel := 0
	for i := range len(elements) {
		if elements[i].kind != blockElementKindList {
			for previousListLevel > 0 {
				previousListLevel--
				ret += "</ul>"
			}
		}

		c := inlineElementToHTML(elements[i].children)

		switch elements[i].kind {
		case blockElementKindParagraph:
			ret += "<p>" + c + "</p>"
		case blockElementKindList:
			if previousListLevel == elements[i].listLevel {
				ret += "<li>" + c + "</li>"
			} else if previousListLevel < elements[i].listLevel {
				previousListLevel++
				ret += "<ul><li>" + c + "</li>"
			} else {
				for previousListLevel > elements[i].listLevel {
					previousListLevel--
					ret += "</ul>"
				}
				ret += "<li>" + c + "</li>"
			}
		case blockElementKindImage:
			ret += fmt.Sprintf(
				"<figure><img src=\"%s\" alt=\"%s\"><figcaption>%s</figcaption></figure>",
				html.EscapeString(elements[i].imageSrc),
				html.EscapeString(elements[i].imageCaption),
				html.EscapeString(elements[i].imageCaption),
			)
		case blockElementKindHeading1:
			ret += "<h1>" + c + "</h1>"
		case blockElementKindHeading2:
			ret += "<h2>" + c + "</h2>"
		case blockElementKindHeading3:
			ret += "<h3>" + c + "</h3>"
		case blockElementKindCodeBlock:
			ret += "<pre><code>" + html.EscapeString(elements[i].codeText) + "</code></pre>"
		default:
			panic("invalid blockElementKind")
		}
	}

	for previousListLevel > 0 {
		previousListLevel--
		ret += "</ul>"
	}

	return ret
}

func inlineElementToHTML(tree inlineElement) string {
	c := ""
	for _, v := range tree.children {
		c += inlineElementToHTML(v)
	}

	switch tree.kind {
	case inlineElementKindRoot:
		return c
	case inlineElementKindText:
		return html.EscapeString(tree.s)
	case inlineElementKindBold:
		return "<b>" + c + "</b>"
	case inlineElementKindCode:
		return "<code>" + c + "</code>"
	case inlineElementKindLink:
		return fmt.Sprintf("<a href=\"%s\">%s</a>", html.EscapeString(tree.linkHref), c)
	}

	panic("unknown inlineElementKind")
}
