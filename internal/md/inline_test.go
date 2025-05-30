package md

import (
	"testing"

	"github.com/comame/note.comame.xyz/internal/test"
)

func TestTokenize(t *testing.T) {
	var got []token
	var expect []token

	got = tokenize("abc**def*gh<fooo>abc")
	expect = []token{
		{s: "abc"},
		{r: true, s: "**"},
		{s: "def*gh"},
		{r: true, s: "<"},
		{s: "fooo"},
		{r: true, s: ">"},
		{s: "abc"}}
	test.AssertEquals(t, got, expect)

	got = tokenize("*****")
	expect = []token{
		{r: true, s: "**"},
		{r: true, s: "**"},
		{s: "*"},
	}
	test.AssertEquals(t, got, expect)

	got = tokenize("")
	expect = nil
	test.AssertEquals(t, got, expect)

	got = tokenize("\\")
	expect = []token(nil)
	test.AssertEquals(t, got, expect)

	got = tokenize("\\\\")
	expect = []token{{s: "\\"}}
	test.AssertEquals(t, got, expect)

	got = tokenize("\\***")
	expect = []token{{s: "*"}, {r: true, s: "**"}}
	test.AssertEquals(t, got, expect)

	got = tokenize("日本語だよ😄")
	expect = []token{{s: "日本語だよ😄"}}
	test.AssertEquals(t, got, expect)

	got = tokenize("line\nline")
	expect = []token{
		{s: "line"},
		{r: true, s: "\n"},
		{s: "line"},
	}
	test.AssertEquals(t, got, expect)
}

func TestParseTokens(t *testing.T) {
	// ただの文字
	test.AssertEquals(
		t,
		parseTokens(
			inlineElement{
				kind: inlineElementKindRoot,
			},
			[]token{
				{s: "normal"},
			},
		),
		inlineElement{
			kind: inlineElementKindRoot,
			children: []inlineElement{
				{
					kind: inlineElementKindText,
					s:    "normal",
				},
			},
		},
	)

	// マルチバイト文字
	test.AssertEquals(
		t,
		parseTokens(
			inlineElement{
				kind: inlineElementKindRoot,
			},
			[]token{
				{s: "日本語だよ😄"},
			},
		),
		inlineElement{
			kind: inlineElementKindRoot,
			children: []inlineElement{
				{
					kind: inlineElementKindText,
					s:    "日本語だよ😄",
				},
			},
		},
	)

	// 太字と思わせてただの文字列
	test.AssertEquals(
		t,
		parseTokens(
			inlineElement{
				kind: inlineElementKindRoot,
			},
			[]token{
				{r: true, s: "**"},
				{s: "normal"},
			},
		),
		inlineElement{
			kind: inlineElementKindRoot,
			children: []inlineElement{
				{
					kind: inlineElementKindText,
					s:    "**",
				},
				{
					kind: inlineElementKindText,
					s:    "normal",
				},
			},
		})

	// 太字
	test.AssertEquals(
		t,
		parseTokens(
			inlineElement{
				kind: inlineElementKindRoot,
			},
			[]token{
				{s: "normal1"},
				{r: true, s: "**"},
				{s: "bold1"},
				{r: true, s: "**"},
				{s: "normal2"},
				{r: true, s: "**"},
				{s: "bold2"},
				{r: true, s: "**"},
				{s: "normal3"},
			},
		),
		inlineElement{
			kind: inlineElementKindRoot,
			children: []inlineElement{
				{
					kind: inlineElementKindText,
					s:    "normal1",
				},
				{
					kind: inlineElementKindBold,
					children: []inlineElement{
						{
							kind: inlineElementKindText,
							s:    "bold1",
						},
					},
				},
				{
					kind: inlineElementKindText,
					s:    "normal2",
				},
				{
					kind: inlineElementKindBold,
					children: []inlineElement{
						{
							kind: inlineElementKindText,
							s:    "bold2",
						},
					},
				},
				{
					kind: inlineElementKindText,
					s:    "normal3",
				},
			},
		},
	)

	// インラインコード
	test.AssertEquals(
		t,
		parseTokens(
			inlineElement{
				kind: inlineElementKindRoot,
			},
			[]token{
				{s: "normal1"},
				{r: true, s: "`"},
				{s: "code1"},
				{r: true, s: "`"},
				{s: "normal2"},
				{r: true, s: "`"},
				{s: "code2"},
				{r: true, s: "`"},
				{s: "normal3"},
			},
		),
		inlineElement{
			kind: inlineElementKindRoot,
			children: []inlineElement{
				{
					kind: inlineElementKindText,
					s:    "normal1",
				},
				{
					kind: inlineElementKindCode,
					children: []inlineElement{
						{
							kind: inlineElementKindText,
							s:    "code1",
						},
					},
				},
				{
					kind: inlineElementKindText,
					s:    "normal2",
				},
				{
					kind: inlineElementKindCode,
					children: []inlineElement{
						{
							kind: inlineElementKindText,
							s:    "code2",
						},
					},
				},
				{
					kind: inlineElementKindText,
					s:    "normal3",
				},
			},
		},
	)

	// URLのみのリンク
	test.AssertEquals(
		t,
		parseTokens(
			inlineElement{
				kind: inlineElementKindRoot,
			},
			[]token{
				{s: "a"},
				{r: true, s: "<"},
				{s: "https://example.com"},
				{r: true, s: ">"},
				{s: "a"},
			},
		),
		inlineElement{
			kind: inlineElementKindRoot,
			children: []inlineElement{
				{
					kind: inlineElementKindText,
					s:    "a",
				},
				{
					kind:     inlineElementKindLink,
					linkHref: "https://example.com",
					children: []inlineElement{
						{kind: inlineElementKindText, s: "https://example.com"},
					},
				},
				{
					kind: inlineElementKindText,
					s:    "a",
				},
			},
		},
	)

	// URLのみのリンクと思わせてただの文字列
	test.AssertEquals(
		t,
		parseTokens(
			inlineElement{
				kind: inlineElementKindRoot,
			},
			[]token{
				{r: true, s: "<"},
				{r: true, s: "**"},
				{s: "not-link"},
				{r: true, s: "**"},
				{r: true, s: ">"},
			},
		),
		inlineElement{
			kind: inlineElementKindRoot,
			children: []inlineElement{
				{
					kind: inlineElementKindText,
					s:    "<",
				},
				{
					kind: inlineElementKindBold,
					children: []inlineElement{
						{
							kind: inlineElementKindText,
							s:    "not-link",
						},
					},
				},
				{
					kind: inlineElementKindText,
					s:    ">",
				},
			},
		},
	)

	// リンク
	test.AssertEquals(
		t,
		parseTokens(
			inlineElement{
				kind: inlineElementKindRoot,
			},
			[]token{
				{s: "a"},
				{r: true, s: "["},
				{s: "name"},
				{r: true, s: "]"},
				{r: true, s: "("},
				{s: "https://example.com/"},
				{r: true, s: ")"},
				{s: "a"},
			},
		),
		inlineElement{
			kind: inlineElementKindRoot,
			children: []inlineElement{
				{
					kind: inlineElementKindText,
					s:    "a",
				},
				{
					kind:     inlineElementKindLink,
					linkHref: "https://example.com/",
					children: []inlineElement{
						{kind: inlineElementKindText, s: "name"},
					},
				},
				{
					kind: inlineElementKindText,
					s:    "a",
				},
			},
		},
	)

	// リンク内太字
	test.AssertEquals(
		t,
		parseTokens(
			inlineElement{kind: inlineElementKindRoot},
			[]token{
				{r: true, s: "["},
				{r: true, s: "**"},
				{s: "name"},
				{r: true, s: "**"},
				{r: true, s: "]"},
				{r: true, s: "("},
				{s: "https://example.com/"},
				{r: true, s: ")"},
			},
		),
		inlineElement{
			kind: inlineElementKindRoot,
			children: []inlineElement{
				{
					kind:     inlineElementKindLink,
					linkHref: "https://example.com/",
					children: []inlineElement{
						{
							kind: inlineElementKindBold,
							children: []inlineElement{
								{kind: inlineElementKindText, s: "name"},
							},
						},
					},
				},
			},
		},
	)

	// 改行
	test.AssertEquals(
		t,
		parseTokens(
			inlineElement{kind: inlineElementKindRoot},
			[]token{
				{s: "line"},
				{r: true, s: "\n"},
				{s: "line"},
			},
		),
		inlineElement{
			kind: inlineElementKindRoot,
			children: []inlineElement{
				{
					kind: inlineElementKindText,
					s:    "line",
				},
				{
					kind: inlineElementKindBreak,
				},
				{
					kind: inlineElementKindText,
					s:    "line",
				},
			},
		},
	)
}
