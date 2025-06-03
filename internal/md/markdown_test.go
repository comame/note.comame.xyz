package md

import (
	_ "embed"
	"fmt"
	"testing"

	"github.com/comame/note.comame.xyz/internal/test"
)

func TestParseBlock(t *testing.T) {
	var expect []blockElement
	var got []blockElement

	inline := inlineElement{
		kind: inlineElementKindRoot,
		children: []inlineElement{
			{
				kind: inlineElementKindText,
				s:    "inline",
			},
		},
	}

	// 改行1つだけのパラグラフ
	got = parseBlock(`inline
inline`)
	expect = []blockElement{
		{
			kind: blockElementKindParagraph,
			children: inlineElement{
				kind: inlineElementKindRoot,
				children: []inlineElement{
					{
						kind: inlineElementKindText,
						s:    "inline",
					},
					{
						kind: inlineElementKindBreak,
					},
					{
						kind: inlineElementKindText,
						s:    "inline",
					},
				},
			},
		},
	}
	test.AssertEquals(t, got, expect)

	// 空行を挟んだパラグラフ
	got = parseBlock(`inline

inline`)
	expect = []blockElement{
		{
			kind:     blockElementKindParagraph,
			children: inline,
		},
		{
			kind: blockElementKindEmpty,
		},
		{
			kind:     blockElementKindParagraph,
			children: inline,
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
	got = parseBlock(`- [ ] inline
- [x] inline`)
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
			kind:     blockElementKindCodeBlock,
			codeName: "file",
			codeText: "source code\n",
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
			kind:           blockElementDetails,
			detailsSummary: "Summary",
			detailsContent: []blockElement{
				{
					kind: blockElementKindParagraph,
					children: inlineElement{
						kind: inlineElementKindRoot,
						children: []inlineElement{
							{
								kind: inlineElementKindText,
								s:    "Hello, world!",
							},
						},
					},
				},
				{
					kind:      blockElementKindList,
					children:  inlineElement{kind: inlineElementKindRoot, children: []inlineElement{{kind: inlineElementKindText, s: "list"}}},
					listLevel: 1,
				},
				{
					kind:      blockElementKindList,
					children:  inlineElement{kind: inlineElementKindRoot, children: []inlineElement{{kind: inlineElementKindText, s: "list"}}},
					listLevel: 1,
				},
			},
		},
	}
	test.AssertEquals(t, got, expect)
	got = parseBlock(`:::details summary
Hello, world!`)
	expect = []blockElement{
		{
			kind: blockElementKindParagraph,
			children: inlineElement{
				kind: inlineElementKindRoot,
				children: []inlineElement{
					{
						kind: inlineElementKindText,
						s:    ":::details summary",
					},
					{
						kind: inlineElementKindBreak,
					},
					{
						kind: inlineElementKindText,
						s:    "Hello, world!",
					},
				},
			},
		},
	}
	test.AssertEquals(t, got, expect)

	got = parseBlock(`; 定義1
: 説明1
; 定義2
: 説明2`)
	expect = []blockElement{
		{
			kind:                   blockElementDescriptionListItem,
			descriptionTerm:        "定義1",
			descriptionDescription: "説明1",
		},
		{
			kind:                   blockElementDescriptionListItem,
			descriptionTerm:        "定義2",
			descriptionDescription: "説明2",
		},
	}
	test.AssertEquals(t, got, expect)

	got = parseBlock(`; 定義
説明ではない`)
	expect = []blockElement{
		{
			kind: blockElementKindParagraph,
			children: inlineElement{
				kind: inlineElementKindRoot,
				children: []inlineElement{
					{
						kind: inlineElementKindText,
						s:    "; 定義",
					},
					{
						kind: inlineElementKindBreak,
					},
					{
						kind: inlineElementKindText,
						s:    "説明ではない",
					},
				},
			},
		},
	}
	test.AssertEquals(t, got, expect)
}

//go:embed spec.md
var specMD string

func ExampleToHTML() {
	html := ToHTML(specMD)

	fmt.Println(html)
	// Output:
	// <p>これはなに: Markdown とは微妙に互換性のないそれらしきもの</p><h1>見出し1</h1><h2>見出し2</h2><h3>見出し3</h3><h2>リスト</h2><ul><li>リスト</li><ul><li>リストアイテム</li><ul><li>リストアイテム</li></ul></ul><li>リストアイテム</li></ul><ul><li>空行を挟むと別のリストとして扱われる</li></ul><ul><li><input type='checkbox' inert>チェックボックス</li><li><input type='checkbox' checked inert>チェックボックス</li></ul><h2>段落</h2><p>通常の Markdown とは異なり、<br>行末の改行は無視されない。</p><p>空行をあけると別の段落として扱われる。</p><h2>コードブロック</h2><pre><code>Lorem ipsum
	// - コードブロック内では Markdown として解釈されない</code></pre><h2>インライン要素</h2><p><b>太字</b> <code>インラインコード</code> <a href="https://example.com">link</a> <a href="https://example.com">https://example.com</a></p><p><a href="https://example.com">l<code>i</code><b>n</b>k</a> のように、組み合わせることもできる</p><h2>拡張構文</h2><details><summary>サマリ</summary><p>内容は格納される</p></details><dl><dt>定義リスト</dt><dd>説明はここに入る</dd><dt>定義2</dt><dd>説明2</dd></dl><h2>非網羅的な既知の非互換</h2><h3>出力結果に違いが出る</h3><ul><li>インライン要素は左からマッチする</li><ul><li>例: <a href="https://example.com/inner">outer?[inner?</a>](https://example.com/outer)</li></ul><li>文末の改行は無視されず、<code>&lt;br&gt;</code> タグが挿入される</li><li>リスト間に空行を挟むと、別のリストとして認識される</li></ul><h3>制約</h3><ul><li>リストの開始記号はハイフン (<code>-</code>) のみ</li><li>リストのインデントは常に半角スペース2つ</li><li>画像はインラインに追加できない</li></ul>
}
