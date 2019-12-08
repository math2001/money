const cacheName = 'money'

const filesToCache = [
  '/',
  'index.html',
]

self.addEventListener('install', e => {
  console.log("[service worker] installing...")
  e.waitUntil(
    caches.open(cacheName).then(cache => {
      console.log("[service worker] got cache")
      return cache.addAll(filesToCache)
    }).then(done => {
      console.log("[service worker] done fetching resources", done)
    })
  )
})


self.addEventListener('activate', e => {
  e.waitUntil(self.clients.Claim())
})

self.addEventListener('fetch', event => {

  event.respondWith(
    fetch(event.request).catch(err => {
      return caches.match(event.request, {ignoreSearch: true})
    })
  )
})
