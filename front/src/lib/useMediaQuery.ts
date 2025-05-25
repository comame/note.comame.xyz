import { useSyncExternalStore } from "react";

class MediaQueryStore {
  list: MediaQueryList;

  constructor(query: string) {
    this.list = window.matchMedia(query);
  }

  subscribe(onStoreChange: () => void): () => void {
    this.list.addEventListener("change", onStoreChange);

    return () => {
      this.list.removeEventListener("change", onStoreChange);
    };
  }

  getSnapshot(): boolean {
    return this.list.matches;
  }
}

const stores: Map<string, MediaQueryStore> = new Map();

export function useMediaQuery(query: string): boolean {
  let store: MediaQueryStore;
  if (stores.has(query)) {
    store = stores.get(query)!;
  } else {
    store = new MediaQueryStore(query);
    stores.set(query, store);
  }

  const matches = useSyncExternalStore(
    store.subscribe.bind(store),
    store.getSnapshot.bind(store)
  );
  return matches;
}
