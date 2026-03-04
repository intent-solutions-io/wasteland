import type { BrowseResponse } from "./types";

// Start fetching the default browse data immediately when the module loads,
// before React mounts. This runs in parallel with auth status check.
let prefetchPromise: Promise<BrowseResponse> | null = null;

// Only prefetch for the root path with no query params (default browse).
if (typeof window !== "undefined" && window.location.pathname === "/" && !window.location.search) {
  prefetchPromise = fetch("/api/wanted")
    .then((r) => {
      if (!r.ok) throw new Error(r.statusText);
      return r.json() as Promise<BrowseResponse>;
    })
    .catch(() => null as unknown as BrowseResponse);
}

/**
 * Consume the prefetched browse response. Returns null if not available
 * (wrong page, already consumed, or fetch failed). Can only be consumed once.
 */
export function consumePrefetch(): Promise<BrowseResponse> | null {
  const p = prefetchPromise;
  prefetchPromise = null;
  return p;
}
