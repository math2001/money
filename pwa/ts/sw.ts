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

self.addEventListener("install", (e: ExtendableEvent) => {
  log("installing");
  e.waitUntil(
    caches.open(cacheName).then(cache => {
      return cache.addAll(filesToCache);
    })
  );
  log("done installing");
});

self.addEventListener("activate", (event: ExtendableEvent) => {
  log("activating");
  event.waitUntil(self.clients.claim());
  log("done activating");
});

self.addEventListener("fetch", (event: FetchEvent) => {
  log("fetching");
  event.respondWith(
    fetch(event.request)
      .then((resp: Response) => {
        const clone = resp.clone();
        self.caches.open(cacheName).then(cache => {
          cache.put(event.request, clone);
        });
        return resp;
      })
      .catch((err: Error) => {
        return caches
          .match(event.request, { ignoreSearch: true })
          .catch(resp => {
            return resp || err;
          });
      })
  );
  log("done fetching");
});
