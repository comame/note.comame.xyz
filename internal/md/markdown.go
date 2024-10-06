package md

import (
	"regexp"
	"strings"
	"unicode"
)

func ToHTML(md string) string {
	return blockElementsToHTML(parseBlock(md))
}

func parseBlock(s string) []blockElement {
	var ret []blockElement

	curr := inlineElement{
		kind: inlineElementKindRoot,
	}

	var isCodeBlock bool
	var codeBlockName string
	var codeBlockLines []string

	var isDetails bool
	var isDetailsSummaryParsed bool
	var detailsSummary string
	var detailsContentLines []string

	for _, l := range strings.Split(s, "\n") {
		// capture curr, ret
		flush := func() {
			if len(curr.children) > 0 {
				ret = append(ret, blockElement{
					kind:     blockElementKindParagraph,
					children: curr,
				})
				curr = inlineElement{
					kind: inlineElementKindRoot,
				}
			}
		}

		if isCodeBlock {
			if l == "```" {
				isCodeBlock = false

				ret = append(ret, blockElement{
					kind:     blockElementKindCodeBlock,
					codeText: strings.Join(codeBlockLines, "\n"),
					codeName: codeBlockName,
				})

				codeBlockLines = nil
				continue
			}

			codeBlockLines = append(codeBlockLines, l)
			continue
		}

		if isDetails {
			summaryPattern := regexp.MustCompile(`^<summary>(.+)<\/summary>$`)
			if !isDetailsSummaryParsed {
				if m := summaryPattern.FindStringSubmatch(l); len(m) > 0 {
					detailsSummary = m[1]
					isDetailsSummaryParsed = true
					continue
				}
			}
			if l == "</details>" {
				ret = append(ret, blockElement{
					kind:               blockElementDetails,
					detailsSummary:     detailsSummary,
					detailsContentHTML: ToHTML(strings.Join(detailsContentLines, "\n")),
				})

				isDetails = false
				isDetailsSummaryParsed = false
				detailsSummary = ""
				detailsContentLines = nil
				continue
			}

			detailsContentLines = append(detailsContentLines, l)
			continue
		}

		// コードブロック中は Markdown として解釈してはならないので、ここより上で処理する必要がある
		l = strings.TrimRightFunc(l, unicode.IsSpace)

		codeStartPattern := regexp.MustCompile("^```(.*)$")
		if m := codeStartPattern.FindStringSubmatch(l); len(m) > 0 {
			flush()

			codeBlockName = m[1]

			isCodeBlock = true
			codeBlockLines = nil
			continue
		}

		if l == "<details>" {
			flush()

			isDetails = true
			continue
		}

		// 簡単のため、リストのインデントは常にスペース2つとする
		// checkboxList は有効な list なので、list より前に検証する必要がある
		checkboxListPattern := regexp.MustCompile(`^((?:  )*)- \[([ x])\](.+)$`)
		if m := checkboxListPattern.FindStringSubmatch(l); len(m) > 0 {
			flush()

			d := m[1]
			ch := m[2]
			c := m[3]

			checked := false
			if ch == "x" {
				checked = true
			}

			ret = append(ret, blockElement{
				kind:              blockElementKindList,
				children:          parseInlineTree(c),
				listLevel:         len(d)/2 + 1,
				checkboxList:      true,
				checkboxIsChecked: checked,
			})
			continue
		}

		// 簡単のため、リストのインデントは常にスペース2つとする
		listPattern := regexp.MustCompile((`^((?:  )*)- (.+)$`))
		if m := listPattern.FindStringSubmatch(l); len(m) > 0 {
			flush()

			d := m[1]
			c := m[2]

			ret = append(ret, blockElement{
				kind:      blockElementKindList,
				children:  parseInlineTree(c),
				listLevel: len(d)/2 + 1,
			})
			continue
		}

		headPattern := regexp.MustCompile(`^(##?#?) +(.+)$`)
		if m := headPattern.FindStringSubmatch(l); len(m) > 0 {
			flush()

			hash := m[1]
			title := m[2]

			k := blockElementKindHeading1
			switch hash {
			case "##":
				k = blockElementKindHeading2
			case "###":
				k = blockElementKindHeading3
			}

			ret = append(ret, blockElement{
				kind: k,
				children: inlineElement{
					// TODO: 実はここも inlineElement として解釈したほうがよい説？
					kind: inlineElementKindText,
					s:    title,
				},
			})
			continue
		}

		imagePattern := regexp.MustCompile(`^!\[(.+)\]\((https:\/\/[\w/.\-_]+)\)$`)
		if m := imagePattern.FindStringSubmatch(l); len(m) > 0 {
			flush()

			caption := m[1]
			src := m[2]

			ret = append(ret, blockElement{
				kind:         blockElementKindImage,
				imageSrc:     src,
				imageCaption: caption,
			})
			continue
		}

		if l == "" {
			flush()
			ret = append(ret, blockElement{
				kind: blockElementKindEmpty,
			})
			continue
		}

		curr.children = append(curr.children, parseInlineTree(l))
	}

	if len(curr.children) > 0 {
		ret = append(ret, blockElement{
			kind:     blockElementKindParagraph,
			children: curr,
		})
		curr = inlineElement{
			kind: inlineElementKindRoot,
		}
	}

	if isCodeBlock && len(codeBlockLines) > 0 {
		ret = append(ret, blockElement{
			kind:     blockElementKindCodeBlock,
			codeText: strings.Join(codeBlockLines, "\n"),
		})
	}

	if isDetails && len(detailsContentLines) > 0 {
		ret = append(ret, blockElement{
			kind:               blockElementDetails,
			detailsSummary:     detailsSummary,
			detailsContentHTML: ToHTML(strings.Join(detailsContentLines, "\n")),
		})
	}

	return ret
}

func parseInlineTree(s string) inlineElement {
	tokens := tokenize(s)
	return parseTokens(inlineElement{kind: inlineElementKindRoot}, tokens)
}
