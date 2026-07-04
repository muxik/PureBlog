import vue from '@vitejs/plugin-vue'
import { defineConfig } from 'vite'

export default defineConfig({
  plugins: [vue()],
  server: { port: 5173 },
  build: {
    rollupOptions: {
      output: {
        manualChunks(id) {
          if (id.includes('/node_modules/@codemirror/') || id.includes('/node_modules/@lezer/')) {
            return 'codemirror'
          }
        },
      },
    },
  },
})
