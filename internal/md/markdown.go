package md

import (
	"regexp"
	"strings"
	"unicode"
)

var codeStartPattern = regexp.MustCompile("^```(.*)$")
var checkboxListPattern = regexp.MustCompile(`^((?:  )*)- \[([ x])\] (.+)$`)
var listPattern = regexp.MustCompile((`^((?:  )*)- (.+)$`))
var headPattern = regexp.MustCompile(`^(##?#?) +(.+)$`)
var imagePattern = regexp.MustCompile(`^!\[(.+)\]\((https:\/\/[\w/.\-_]+)\)$`)
var summaryPattern = regexp.MustCompile(`^<summary>(.+)<\/summary>$`)
var customDetailsPattern = regexp.MustCompile("^:::details (.+)$")
var descriptionTermPattern = regexp.MustCompile(`^; +(.+)$`)
var descriptionDetailsPattern = regexp.MustCompile(`^: +(.+)$`)

func ToHTML(md string) string {
	return blockElementsToHTML(parseBlock(md))
}

func parseBlock(s string) []blockElement {
	lines := strings.Split(s, "\n")
	b, _ := parseBlockInternal(lines, 0, "")
	return b
}

func parseBlockInternal(lines []string, startIndex int, terminate string) ([]blockElement, bool) {
	var ret []blockElement

	var paragraphBuffer string

	var isCodeBlock bool
	var codeBlockName string
	var codeBlockLines []string

	var isDetails bool
	var isCustomDetails bool
	var isDetailsSummaryParsed bool
	var detailsSummary string
	var detailsContentLines []string

	var isDescription bool
	var descriptionTerm string
	var descriptionTermOriginalLine string

	for _, l := range lines[startIndex:] {
		// バッファに溜まってる文字を通常の段落として書き出す
		flush := func() {
			if paragraphBuffer != "" {
				ret = append(ret, blockElement{
					kind:     blockElementKindParagraph,
					children: parseInlineTree(paragraphBuffer),
				})
				paragraphBuffer = ""
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

		if isDetails && isCustomDetails {
			if l == ":::" {
				ret = append(ret, blockElement{
					kind:               blockElementDetails,
					detailsSummary:     detailsSummary,
					detailsContentHTML: ToHTML(strings.Join(detailsContentLines, "\n")),
				})

				isDetails = false
				isCustomDetails = false
				isDetailsSummaryParsed = false
				detailsSummary = ""
				detailsContentLines = nil
				continue
			}

			detailsContentLines = append(detailsContentLines, l)
			continue
		}

		if isDetails {
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

		if isDescription {
			m := descriptionDetailsPattern.FindStringSubmatch(l)
			if len(m) == 0 {
				// 「; 定義」の次行に「: 説明」が来なければ、ただの段落として扱う
				if paragraphBuffer != "" {
					paragraphBuffer += "\n"
				}
				paragraphBuffer += descriptionTermOriginalLine
				paragraphBuffer += "\n"
				paragraphBuffer += l
				isDescription = false
				continue
			}

			details := m[1]
			ret = append(ret, blockElement{
				kind:                   blockElementDescriptionListItem,
				descriptionTerm:        descriptionTerm,
				descriptionDescription: details,
			})

			isDescription = false
			continue
		}

		// コードブロック中は Markdown として解釈してはならないので、ここより上で処理する必要がある
		l = strings.TrimRightFunc(l, unicode.IsSpace)

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

		if m := customDetailsPattern.FindStringSubmatch(l); len(m) > 0 {
			flush()

			detailsSummary = m[1]

			isDetails = true
			isCustomDetails = true
			continue
		}

		// 簡単のため、リストのインデントは常にスペース2つとする
		// checkboxList は有効な list なので、list より前に検証する必要がある
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
					kind: inlineElementKindText,
					s:    title,
				},
			})
			continue
		}

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

		if m := descriptionTermPattern.FindStringSubmatch(l); len(m) > 0 {
			flush()

			isDescription = true
			descriptionTerm = m[1]
			descriptionTermOriginalLine = l
			continue
		}

		if l == "" {
			flush()
			ret = append(ret, blockElement{
				kind: blockElementKindEmpty,
			})
			continue
		}

		if paragraphBuffer != "" {
			paragraphBuffer += "\n"
		}
		paragraphBuffer += l
	}

	if paragraphBuffer != "" {
		ret = append(ret, blockElement{
			kind:     blockElementKindParagraph,
			children: parseInlineTree(paragraphBuffer),
		})
	}

	if isCodeBlock && len(codeBlockLines) > 0 {
		ret = append(ret, blockElement{
			kind:     blockElementKindCodeBlock,
			codeName: codeBlockName,
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

	// ループが回り切ったということは、terminateが見つからなかったということ
	return ret, false
}

func parseInlineTree(s string) inlineElement {
	tokens := tokenize(s)
	return parseTokens(inlineElement{kind: inlineElementKindRoot}, tokens)
}
