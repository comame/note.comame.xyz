async function awaitWasm() {
  return new Promise((resolve) => {
    (function loop() {
      setTimeout(() => {
        console.log("loop");

        if (!window.go_parseMarkdown) {
          loop();
          return;
        }
        resolve();
      }, 100);
    })();
  });
}

(async function () {
  await awaitWasm();

  document.getElementById("input").value = `👇👇👇編集してみてね👇👇👇

# Heading 1
## Heading 2
### Heading 3

Lorem **ipsum** [dolor **sit amet** consectetur](https://example.com) adipisicing elit. \`Nostrum assumenda fuga enim\` ullam impedit quibusdam necessitatibus excepturi earum? Animi placeat porro, quis veniam numquam cum provident dolore eum fugiat maxime.

- list
- list
    - [link](https://example.com)
    - <https://example.com>
- list
    - **list**
- list

\`\`\`main.go
package main

import log

func main() {
	log.Println("Hello, world!")
}
\`\`\`
`;
  document.getElementById("output").innerHTML = go_parseMarkdown(
    document.getElementById("input").value
  );
})();
