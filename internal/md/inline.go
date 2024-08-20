package md

import (
	"strings"
)

// トークンに分割する
func tokenize(str string) []token {
	var ret []token
	buf := ""

	s := []rune(str)

	for i := 0; i < len(s); i++ {
		if i >= len(s) {
			break
		}

		// capture ret, buf
		flush := func() {
			if len(buf) == 0 {
				return
			}
			ret = append(ret, token{s: buf})
			buf = ""
		}
		takeTwo := func() string {
			if i >= len(s)-1 {
				return string(s[i])
			}
			return string(s[i]) + string(s[i+1])
		}

		c := string(s[i])

		if c == "\\" {
			if takeTwo() == "\\\\" {
				buf += "\\"
				i++
				continue
			}
			if len(takeTwo()) == 2 {
				buf += string([]rune(takeTwo())[1])
				i++
				continue
			}
			continue
		}

		switch c {
		case "<":
			flush()
			ret = append(ret, token{r: true, s: "<"})
			continue
		case ">":
			flush()
			ret = append(ret, token{r: true, s: ">"})
			continue
		case "[":
			flush()
			ret = append(ret, token{r: true, s: "["})
			continue
		case "]":
			flush()
			ret = append(ret, token{r: true, s: "]"})
			continue
		case "(":
			flush()
			ret = append(ret, token{r: true, s: "("})
			continue
		case ")":
			flush()
			ret = append(ret, token{r: true, s: ")"})
			continue
		case "`":
			flush()
			ret = append(ret, token{r: true, s: "`"})
			continue
		}

		switch takeTwo() {
		case "**":
			flush()
			ret = append(ret, token{r: true, s: "**"})
			i++
			continue
		}

		buf += c
	}

	if len(buf) != 0 {
		ret = append(ret, token{s: buf})
	}

	return ret
}

func findNextReservedToken(startIndex int, token string, tokens []token) int {
	if startIndex+1 > len(tokens)-1 {
		return -1
	}
	for i := startIndex + 1; i < len(tokens); i++ {
		if tokens[i].r && tokens[i].s == token {
			return i
		}
	}
	return -1
}

func parseTokens(tree inlineElement, tokens []token) inlineElement {
	for i := 0; i < len(tokens); i++ {
		t := tokens[i]

		if t.r && t.s == "**" {
			c := findNextReservedToken(i, "**", tokens)
			// 閉じタグが無ければ、通常の文字列として扱う
			if c < 0 {
				tree.children = append(tree.children, inlineElement{
					kind: inlineElementKindText,
					s:    t.s,
				})
				continue
			}
			tree.children = append(
				tree.children,
				parseTokens(inlineElement{kind: inlineElementKindBold}, tokens[i+1:c]),
			)
			i += c - i
			continue
		}

		if t.r && t.s == "`" {
			c := findNextReservedToken(i, "`", tokens)
			// 閉じタグがなければ、通常の文字列として扱う
			if c < 0 {
				tree.children = append(tree.children, inlineElement{
					kind: inlineElementKindText,
					s:    t.s,
				})
				continue
			}
			tree.children = append(
				tree.children,
				parseTokens(inlineElement{kind: inlineElementKindCode}, tokens[i+1:c]),
			)
			i += c - i
			continue
		}

		if t.r && t.s == "<" {
			cl := findNextReservedToken(i, ">", tokens)
			// 閉じタグがなければ、通常の文字列として扱う
			if cl < 0 {
				tree.children = append(tree.children, inlineElement{
					kind: inlineElementKindText,
					s:    t.s,
				})
				continue
			}

			// <...> の中身は URL の可能性があるので、普通の文字列として取得する
			href := ""
			for _, t := range tokens[i+1 : cl] {
				href += t.s
			}

			// <...> の中身が URL ではなかったら、リンクでは無かったとして扱う
			if !(strings.HasPrefix(href, "https://") || strings.HasPrefix(href, "http://")) {
				tree.children = append(tree.children, inlineElement{
					kind: inlineElementKindText,
					s:    t.s,
				})
				continue
			}

			tree.children = append(
				tree.children,
				// URL リンクの URL 部分はインライン要素としてパースしない、なぜなら URL なので
				inlineElement{kind: inlineElementKindLink, linkHref: href, children: []inlineElement{
					{kind: inlineElementKindText, s: href},
				}},
			)
			i += cl - i
			continue
		}

		if t.r && t.s == "[" {
			i1 := findNextReservedToken(i, "]", tokens)
			// キーワードが順番に並んでいなければ、通常の文字列として扱う
			if i1 < 0 || i1-i == 1 {
				tree.children = append(tree.children, inlineElement{
					kind: inlineElementKindText,
					s:    t.s,
				})
				continue
			}
			i2 := findNextReservedToken(i1, "(", tokens)
			if i2 < 0 || i2-i1 != 1 {
				tree.children = append(tree.children, inlineElement{
					kind: inlineElementKindText,
					s:    t.s,
				})
				continue
			}
			i3 := findNextReservedToken(i2, ")", tokens)
			if i3 < 0 || i3-i2 == 1 {
				tree.children = append(tree.children, inlineElement{
					kind: inlineElementKindText,
					s:    t.s,
				})
				continue
			}

			// (...) の中身は URL の可能性があるので、普通の文字列として取得する
			href := ""
			for _, t := range tokens[i2+1 : i3] {
				href += t.s
			}

			// (...) の中身が URL ではなかったら、リンクで無かったとして扱う
			if !(strings.HasPrefix(href, "https://") || strings.HasPrefix(href, "http://")) {
				tree.children = append(tree.children, inlineElement{
					kind: inlineElementKindText,
					s:    t.s,
				})
				continue
			}

			name := parseTokens(inlineElement{kind: inlineElementKindRoot}, tokens[i+1:i1]).children

			tree.children = append(
				tree.children,
				inlineElement{
					kind:     inlineElementKindLink,
					linkHref: href,
					children: name,
				},
			)

			i += i3 - i
			continue
		}

		tree.children = append(tree.children, inlineElement{
			kind: inlineElementKindText,
			s:    t.s,
		})
	}
	return tree
}
