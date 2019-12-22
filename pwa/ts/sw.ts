declare var self: ServiceWorkerGlobalScope;
export {};

const cacheName = "money-v1";
const DEBUG = false;

// FIXME: should we even use that? it's a pain to maintain a list manually...
// it's not even used by the fetch handler
const filesToCache = ["/", "index.html"];

function log(...args: any) {
  if (DEBUG) {
    console.log("[service worker]", ...args);
  }
}

self.addEventListener("install", (e: ExtendableEvent) => {});

self.addEventListener("activate", (event: ExtendableEvent) => {});

self.addEventListener("fetch", (event: FetchEvent) => {});
