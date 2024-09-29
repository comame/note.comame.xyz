package md

import (
	"testing"

	"github.com/comame/note.comame.xyz/internal/test"
)

func TestBlockElementsToHTML(t *testing.T) {
	inline := inlineElement{
		kind: inlineElementKindRoot,
		children: []inlineElement{
			{
				kind: inlineElementKindText,
				s:    "str",
			},
		},
	}

	expect := "<p>str</p>"
	got := blockElementsToHTML([]blockElement{
		{
			kind:     blockElementKindParagraph,
			children: inline,
		},
	})
	test.AssertSame(t, got, expect)

	expect = "<ul><li>str</li><li>str</li><ul><li>str</li><li>str</li></ul><li>str</li></ul>"
	got = blockElementsToHTML([]blockElement{
		{
			kind:      blockElementKindList,
			children:  inline,
			listLevel: 1,
		},
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
			listLevel: 2,
		},
		{
			kind:      blockElementKindList,
			children:  inline,
			listLevel: 1,
		},
	})
	test.AssertSame(t, got, expect)

	expect = "<ul><li>str</li></ul><ul><li>str</li></ul>"
	got = blockElementsToHTML([]blockElement{
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
	})
	test.AssertSame(t, got, expect)

	expect = "<ul><li><input type='checkbox' checked inert>str</li><li><input type='checkbox' inert>str</li></ul>"
	got = blockElementsToHTML([]blockElement{
		{
			kind:              blockElementKindList,
			children:          inline,
			listLevel:         1,
			checkboxList:      true,
			checkboxIsChecked: true,
		},
		{
			kind:              blockElementKindList,
			children:          inline,
			listLevel:         1,
			checkboxList:      true,
			checkboxIsChecked: false,
		},
	})
	test.AssertSame(t, got, expect)

	expect = "<ul><li>str</li><ul><li>str</li><ul><li>str</li></ul></ul></ul>"
	got = blockElementsToHTML([]blockElement{
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
	})
	test.AssertSame(t, got, expect)

	expect = "<ul><li>str</li><ul><li>str</li><ul><li>str</li></ul></ul></ul><p>str</p>"
	got = blockElementsToHTML([]blockElement{
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
			kind:     blockElementKindParagraph,
			children: inline,
		},
	})
	test.AssertSame(t, got, expect)

	expect = "<ul><li>str</li><ul><li>str</li><ul><li>str</li></ul></ul><li>str</li></ul>"
	got = blockElementsToHTML([]blockElement{
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
	})
	test.AssertSame(t, got, expect)

	expect = "<figure><img src=\"https://example.com/example.png\" alt=\"image\"><figcaption>image</figcaption></figure>"
	got = blockElementsToHTML([]blockElement{
		{
			kind:         blockElementKindImage,
			imageSrc:     "https://example.com/example.png",
			imageCaption: "image",
		},
	})
	test.AssertSame(t, got, expect)

	expect = "<h1>str</h1><h2>str</h2><h3>str</h3>"
	got = blockElementsToHTML([]blockElement{
		{
			kind:     blockElementKindHeading1,
			children: inline,
		}, {
			kind:     blockElementKindHeading2,
			children: inline,
		}, {
			kind:     blockElementKindHeading3,
			children: inline,
		},
	})
	test.AssertSame(t, got, expect)

	expect = "<p>str</p><pre><code>source code</code></pre><p>str</p>"
	got = blockElementsToHTML([]blockElement{
		{
			kind:     blockElementKindParagraph,
			children: inline,
		},
		{
			kind:     blockElementKindCodeBlock,
			codeText: "source code",
		},
		{
			kind:     blockElementKindParagraph,
			children: inline,
		},
	})
	test.AssertSame(t, got, expect)
}

func TestInlineElementTreeToHTML(t *testing.T) {
	// ただの文字列
	test.AssertSame(
		t,
		inlineElementToHTML(inlineElement{
			kind: inlineElementKindRoot,
			children: []inlineElement{
				{
					kind: inlineElementKindText,
					s:    "Hello, world!",
				},
			},
		}),
		"Hello, world!",
	)

	// 太字
	test.AssertSame(
		t,
		inlineElementToHTML(inlineElement{
			kind: inlineElementKindRoot,
			children: []inlineElement{
				{
					kind: inlineElementKindBold,
					children: []inlineElement{
						{
							kind: inlineElementKindText,
							s:    "Hello, world!",
						},
					},
				},
			},
		}),
		"<b>Hello, world!</b>",
	)

	// 太字リンク
	test.AssertSame(
		t,
		inlineElementToHTML(inlineElement{
			kind: inlineElementKindRoot,
			children: []inlineElement{
				{
					kind:     inlineElementKindLink,
					linkHref: "https://example.com/example.html",
					children: []inlineElement{
						{
							kind: inlineElementKindBold,
							children: []inlineElement{
								{
									kind: inlineElementKindText,
									s:    "Hello, world!",
								},
							},
						},
					},
				},
			},
		}),
		"<a href=\"https://example.com/example.html\"><b>Hello, world!</b></a>",
	)

	// 複数
	test.AssertSame(
		t,
		inlineElementToHTML(inlineElement{
			kind: inlineElementKindRoot,
			children: []inlineElement{
				{
					kind: inlineElementKindText,
					s:    "Hello, world!",
				},
				{
					kind: inlineElementKindBold,
					children: []inlineElement{
						{
							kind: inlineElementKindText,
							s:    "Hello, world!",
						},
					},
				},
				{
					kind:     inlineElementKindLink,
					linkHref: "https://example.com/example.html",
					children: []inlineElement{
						{
							kind: inlineElementKindBold,
							children: []inlineElement{
								{
									kind: inlineElementKindText,
									s:    "Hello, world!",
								},
							},
						},
					},
				},
			},
		}),
		"Hello, world!<b>Hello, world!</b><a href=\"https://example.com/example.html\"><b>Hello, world!</b></a>",
	)
}
