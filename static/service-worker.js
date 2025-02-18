const CACHE_NAME = "cs2-browser-v1";
const ASSETS = [
  "/",
  "/static/style.css",
  "/static/manifest.json",
  "/static/icon-192.png",
  "/static/icon-512.png",
  "https://unpkg.com/htmx.org@1.9.10",
  "https://unpkg.com/alpinejs@3.x.x/dist/cdn.min.js",
];

self.addEventListener("install", (event) => {
  console.log("Service Worker: Installing...");
  event.waitUntil(
    caches.open(CACHE_NAME).then((cache) => {
      console.log("Service Worker: Caching files");
      return cache.addAll(ASSETS);
    }),
  );
});

self.addEventListener("activate", (event) => {
  console.log("Service Worker: Activated");
  event.waitUntil(
    caches.keys().then((cacheNames) => {
      return Promise.all(
        cacheNames.map((cache) => {
          if (cache !== CACHE_NAME) {
            return caches.delete(cache);
          }
        }),
      );
    }),
  );
});

self.addEventListener("fetch", (event) => {
  event.respondWith(
    caches
      .match(event.request)
      .then((response) => response || fetch(event.request)),
  );
});
