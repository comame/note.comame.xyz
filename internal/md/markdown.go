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
var customDetailsPattern = regexp.MustCompile("^(:{3,})+details (.+)$")
var descriptionTermPattern = regexp.MustCompile(`^; +(.+)$`)
var descriptionDetailsPattern = regexp.MustCompile(`^: +(.+)$`)

func ToHTML(md string) string {
	return blockElementsToHTML(parseBlock(md))
}

func parseBlock(s string) []blockElement {
	lines := strings.Split(s, "\n")
	return parseBlockInternal(lines)
}

func parseBlockInternal(lines []string) (ret []blockElement) {
	var paragraphBuffer string

	var isCodeBlock bool
	var codeBlockName string
	var codeBlockLines []string

	var isDescription bool
	var descriptionTerm string
	var descriptionTermOriginalLine string

	for i := 0; i < len(lines); i++ {
		l := lines[i]

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

		if m := customDetailsPattern.FindStringSubmatch(l); len(m) > 0 {
			colons := m[1]

			terminate := findLine(lines, i+1, colons) // 開始記号と同じ数のコロンを終端記号とする
			// details が閉じられていない場合は、段落として扱う
			if terminate < 0 {
				if paragraphBuffer != "" {
					paragraphBuffer += "\n"
				}
				paragraphBuffer += l
				continue
			}

			flush()

			detailsSummary := m[2]

			blocksInDetails := parseBlockInternal(lines[i+1 : terminate])
			ret = append(ret, blockElement{
				kind:           blockElementDetails,
				detailsSummary: detailsSummary,
				detailsContent: blocksInDetails,
			})

			i = terminate
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

	return ret
}

func findLine(lines []string, startIndex int, line string) int {
	if startIndex >= len(lines) {
		return -1
	}

	for i := startIndex; i < len(lines); i++ {
		if lines[i] == line {
			return i
		}
	}

	return -1
}

func parseInlineTree(s string) inlineElement {
	tokens := tokenize(s)
	return parseTokens(inlineElement{kind: inlineElementKindRoot}, tokens)
}
