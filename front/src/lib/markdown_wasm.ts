import "../../../out/dist/wasm_exec";

// @ts-expect-error
const go = new Go();

const { instance } = await WebAssembly.instantiateStreaming(
  fetch("/out/dist/goapp.wasm"),
  go.importObject
);

go.run(instance);

console.log("wasm markdown parser loaded!");

export function parse(md: string): string {
  // @ts-expect-error
  return go_parseMarkdown(md);
}
