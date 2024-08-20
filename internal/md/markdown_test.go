package md

import (
	"fmt"
)

func Example() {
	str := `# Title
Hello, world!

## Title 2
Lorem **ipsum** [dolor **sit amet** consectetur](https://example.com) adipisicing elit. Nostrum assumenda fuga enim ullam impedit quibusdam necessitatibus excepturi earum? Animi placeat porro, quis veniam numquam cum provident dolore eum fugiat maxime.

aaa
## 日本語だよ😄
- foo
  - bar
    - baz
  - foo
    - bar
- baz
- <https://example.com>
![caption](https://example.com/img.png)

日本語だよ😄`

	html := blockElementsToHTML(parseBlock(str))
	fmt.Println(html)

	// Output:
	// <h1>Title</h1><p>Hello, world!</p><h2>Title 2</h2><p>Lorem <b>ipsum</b> <a href="https://example.com">dolor <b>sit amet</b> consectetur</a> adipisicing elit. Nostrum assumenda fuga enim ullam impedit quibusdam necessitatibus excepturi earum? Animi placeat porro, quis veniam numquam cum provident dolore eum fugiat maxime.</p><p>aaa</p><h2>日本語だよ😄</h2><ul><li>foo</li><ul><li>bar</li><ul><li>baz</li></ul><li>foo</li><ul><li>bar</li></ul></ul><li>baz</li><li><a href="https://example.com">https://example.com</a></li></ul><figure><img src="https://example.com/img.png" alt="caption"><figcaption>caption</figcaption></figure><p>日本語だよ😄</p>
}
