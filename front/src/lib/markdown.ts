export async function parse(md: string): Promise<string> {
  const { parse: parseOriginal } = await import("./markdown_wasm");

  return parseOriginal(md);
}
