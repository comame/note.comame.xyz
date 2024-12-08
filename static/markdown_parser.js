const worker = new Worker("/static/markdown_worker.js", {
  type: "module",
});

let nextMarkdownInput = null;
let previousMarkdownInput = null;
let processing = false;

const interval = () => {
  if (processing || nextMarkdownInput === previousMarkdownInput) {
    return;
  }
  processing = true;
  previousMarkdownInput = nextMarkdownInput;

  worker.postMessage({
    payload: nextMarkdownInput,
    type: "parse_markdown",
  });
};
interval();
setInterval(interval, 100);

/**
 * @param {string} md
 */
export function onInput(md) {
  nextMarkdownInput = md;
}

/**
 * @param {(html: string) => void} f
 */
export function setCallback(f) {
  worker.addEventListener("message", ({ data }) => {
    processing = false;

    const { payload } = data;
    f(payload);
  });
}
