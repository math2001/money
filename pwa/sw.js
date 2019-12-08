const cacheName = 'money-v1'

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
    fetch(event.request)
    .then(resp => {
      // save in the cache
      console.log('[service worker] got response, put into cache (should have next message)', resp, resp.clone)
      let clone = resp.clone()
      self.caches.open(cacheName).then(cache => {
        cache.put(event.request, clone)
        console.log('[service worker] saved response in the cache', clone)
      })
      return resp

    }).catch(err => {
      // try to get from the cach
      return caches.match(event.request, {ignoreSearch: true}) || err
    }).then(resp => {
      console.log('[service worker] just checking what we got', resp)
      return resp
    })
  )
})
