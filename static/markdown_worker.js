import "/out/dist/wasm_exec.js";

self.onmessage = async (e) => {
  await instantiate;

  const type = e.data["type"];
  const payload = e.data["payload"];

  switch (type) {
    case "parse_markdown": {
      self.postMessage({
        type: "parse_markdown",
        payload: go_parseMarkdown(payload),
      });
    }
  }
};

const instantiate = new Promise(async (resolve) => {
  const go = new Go();

  const { instance } = await WebAssembly.instantiateStreaming(
    fetch("/out/dist/goapp.wasm"),
    go.importObject
  );

  go.run(instance);
  resolve();
});
