import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

// Project uses a `public/` directory that contains `index.html` and app entry files.
// Set Vite root to that folder and output the build to ../dist so Docker can copy /app/dist.
export default defineConfig({
  root: 'public',
  plugins: [react()],
  build: {
    outDir: '../dist'
  }
})
