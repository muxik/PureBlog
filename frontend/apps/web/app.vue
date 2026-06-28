<template>
  <div class="site">
    <header class="site-header">
      <NuxtLink to="/" class="brand">PureBlog</NuxtLink>
      <nav class="site-nav">
        <NuxtLink to="/categories" class="nav-link">分类</NuxtLink>
        <NuxtLink to="/tags" class="nav-link">标签</NuxtLink>
      </nav>
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
  display: flex;
  align-items: baseline;
  gap: 1.5rem;
}
.brand {
  font-family: var(--font-serif, serif);
  font-size: 1.25rem;
  color: var(--accent, #235a73);
  text-decoration: none;
  letter-spacing: 0.04em;
}
.site-nav {
  display: flex;
  gap: 1rem;
}
.nav-link {
  font-size: 0.9rem;
  color: var(--ink-2, #6b655c);
  text-decoration: none;
}
.nav-link:hover {
  color: var(--accent, #235a73);
}
.site-footer {
  margin-top: 4rem;
  padding-top: 1.5rem;
  border-top: 1px solid var(--border, #e7e2d6);
  color: var(--ink-3, #8a857a);
  font-size: 0.85rem;
}
</style>
