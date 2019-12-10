declare var self: ServiceWorkerGlobalScope;
export {};

console.log("hello world from service worker!!");

const cacheName = "money-v1";
const DEBUG = true;

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
  log("sw fetching");
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
        console.log(err);
        console.log(typeof err);
        return caches.match(event.request, { ignoreSearch: true }) || err;
      })
  );
  // event.respondWith(
  //   fetch(event.request)
  //     .then((resp: Response) => {
  //       let clone = resp.clone();
  //       self.caches.open(cacheName).then(cache => {
  //         cache.put(event.request, clone);
  //       });
  //       return resp;
  //     })
  //     .catch(err => {
  //       // try to get from the cach
  //       return caches.match(event.request, { ignoreSearch: true }) || err;
  //     })
  //     .then(resp => {
  //       return resp;
  //     })
  // );
});
