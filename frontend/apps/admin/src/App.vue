<template>
  <div class="admin">
    <RouterView />
  </div>
</template>

<style>
.admin {
  max-width: 48rem;
  margin: 0 auto;
  padding: 2rem 1.25rem;
}
.admin input,
.admin button {
  font: inherit;
}
.admin input {
  display: block;
  width: 100%;
  padding: 0.5rem 0.7rem;
  margin: 0.4rem 0;
  border: 1px solid var(--border, #e7e2d6);
  border-radius: 4px;
  background: var(--paper, #faf8f2);
  color: inherit;
}
.admin button {
  padding: 0.45rem 1rem;
  margin-right: 0.5rem;
  border: 1px solid var(--accent, #235a73);
  border-radius: 4px;
  background: var(--accent, #235a73);
  color: #fff;
  cursor: pointer;
}
.admin button.ghost {
  background: transparent;
  color: var(--accent, #235a73);
}
.err {
  color: #b23b3b;
}
</style>
