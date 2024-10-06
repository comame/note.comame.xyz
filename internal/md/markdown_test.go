package md

import (
	"testing"

	"github.com/comame/note.comame.xyz/internal/test"
)

func TestParseBlock(t *testing.T) {
	var expect []blockElement
	var got []blockElement

	// FIXME: inlineElementKindRoot が 2 重になってる
	doubleRootInline := inlineElement{
		kind: inlineElementKindRoot,
		children: []inlineElement{
			{
				kind: inlineElementKindRoot,
				children: []inlineElement{
					{
						kind: inlineElementKindText,
						s:    "inline",
					},
				},
			},
		},
	}
	inline := inlineElement{
		kind: inlineElementKindRoot,
		children: []inlineElement{
			{
				kind: inlineElementKindText,
				s:    "inline",
			},
		},
	}

	// ただの文章
	got = parseBlock(`inline

inline`)
	expect = []blockElement{
		{
			kind:     blockElementKindParagraph,
			children: doubleRootInline,
		},
		{
			kind: blockElementKindEmpty,
		},
		{
			kind:     blockElementKindParagraph,
			children: doubleRootInline,
		},
	}
	test.AssertEquals(t, got, expect)

	// リスト
	got = parseBlock(`- inline
  - inline
    - inline
- inline

- inline`)
	expect = []blockElement{
		{
			kind:      blockElementKindList,
			children:  inline,
			listLevel: 1,
		},
		{
			kind:      blockElementKindList,
			children:  inline,
			listLevel: 2,
		},
		{
			kind:      blockElementKindList,
			children:  inline,
			listLevel: 3,
		},
		{
			kind:      blockElementKindList,
			children:  inline,
			listLevel: 1,
		},
		{
			kind: blockElementKindEmpty,
		},
		{
			kind:      blockElementKindList,
			children:  inline,
			listLevel: 1,
		},
	}
	test.AssertEquals(t, got, expect)

	// チェックボックス
	// FIXME: スペースがないのに受理している
	got = parseBlock(`- [ ]inline
- [x]inline`)
	expect = []blockElement{
		{
			kind:              blockElementKindList,
			checkboxList:      true,
			checkboxIsChecked: false,
			listLevel:         1,
			children:          inline,
		},
		{
			kind:              blockElementKindList,
			checkboxList:      true,
			checkboxIsChecked: true,
			listLevel:         1,
			children:          inline,
		},
	}
	test.AssertEquals(t, got, expect)

	// タイトル
	got = parseBlock(`# heading 1
## heading 2
### heading 3`)
	expect = []blockElement{
		{
			kind: blockElementKindHeading1,
			children: inlineElement{
				kind: inlineElementKindText,
				s:    "heading 1",
			},
		},
		{
			kind: blockElementKindHeading2,
			children: inlineElement{
				kind: inlineElementKindText,
				s:    "heading 2",
			},
		},
		{
			kind: blockElementKindHeading3,
			children: inlineElement{
				kind: inlineElementKindText,
				s:    "heading 3",
			},
		},
	}
	test.AssertEquals(t, got, expect)

	// 画像
	got = parseBlock(`![caption](https://example.com)`)
	expect = []blockElement{
		{
			kind:         blockElementKindImage,
			imageSrc:     "https://example.com",
			imageCaption: "caption",
		},
	}
	test.AssertEquals(t, got, expect)

	// コードブロック
	got = parseBlock("```file\nsource code\n```")
	expect = []blockElement{
		{
			kind:     blockElementKindCodeBlock,
			codeName: "file",
			codeText: "source code",
		},
	}
	test.AssertEquals(t, got, expect)
	got = parseBlock("```file\nsource code\n")
	expect = []blockElement{
		{
			kind: blockElementKindCodeBlock,
			// FIXME: 途中でコードブロックが切れたときにタグが入らない
			// codeName: "file",
			codeText: "source code\n",
		},
	}
	test.AssertEquals(t, got, expect)

	// トグル (HTML)
	got = parseBlock(`<details>
<summary>Summary</summary>
Hello, world!
- list
- list
</details>`)
	expect = []blockElement{
		{
			kind:               blockElementDetails,
			detailsSummary:     "Summary",
			detailsContentHTML: "<p>Hello, world!</p><ul><li>list</li><li>list</li></ul>",
		},
	}
	test.AssertEquals(t, got, expect)
	got = parseBlock(`<details>
Hello, world!
</details>`)
	expect = []blockElement{
		{
			kind:               blockElementDetails,
			detailsSummary:     "",
			detailsContentHTML: "<p>Hello, world!</p>",
		},
	}
	test.AssertEquals(t, got, expect)
	got = parseBlock(`<details>
<summary>summary</summary>
Hello, world!`)
	expect = []blockElement{
		{
			kind:               blockElementDetails,
			detailsSummary:     "summary",
			detailsContentHTML: "<p>Hello, world!</p>",
		},
	}
	test.AssertEquals(t, got, expect)

	// トグル (カスタム)
	got = parseBlock(`:::details Summary
Hello, world!
- list
- list
:::`)
	expect = []blockElement{
		{
			kind:               blockElementDetails,
			detailsSummary:     "Summary",
			detailsContentHTML: "<p>Hello, world!</p><ul><li>list</li><li>list</li></ul>",
		},
	}
	test.AssertEquals(t, got, expect)
	got = parseBlock(`:::details summary
Hello, world!`)
	expect = []blockElement{
		{
			kind:               blockElementDetails,
			detailsSummary:     "summary",
			detailsContentHTML: "<p>Hello, world!</p>",
		},
	}
	test.AssertEquals(t, got, expect)
}
