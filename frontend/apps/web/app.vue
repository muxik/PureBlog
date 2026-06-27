<template>
  <div class="site">
    <header class="site-header">
      <NuxtLink to="/" class="brand">PureBlog</NuxtLink>
    </header>
    <main class="site-main">
      <NuxtPage />
    </main>
    <footer class="site-footer">黛 · PureBlog v3</footer>
  </div>
</template>

<style>
.site {
  max-width: 42rem;
  margin: 0 auto;
  padding: 2rem 1.25rem 4rem;
}
.site-header {
  padding: 1rem 0 2.5rem;
}
.brand {
  font-family: var(--font-serif, serif);
  font-size: 1.25rem;
  color: var(--accent, #235a73);
  text-decoration: none;
  letter-spacing: 0.04em;
}
.site-footer {
  margin-top: 4rem;
  padding-top: 1.5rem;
  border-top: 1px solid var(--border, #e7e2d6);
  color: var(--ink-3, #8a857a);
  font-size: 0.85rem;
}
</style>
