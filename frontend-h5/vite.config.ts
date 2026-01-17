import { defineConfig } from 'vite'
import fs from 'node:fs'
import path from 'node:path'

function brandLoginPage() {
  return {
    name: 'brand-login-page',
    configureServer(server: any) {
      server.middlewares.use((req: any, res: any, next: any) => {
        const method = (req.method || 'GET').toUpperCase()
        if (method !== 'GET' && method !== 'HEAD') return next()

        const url = req.url || ''
        const pathname = url.split('?')[0]
        if (pathname !== '/brand/login' && pathname !== '/brand/login/') return next()

        const htmlPath = path.resolve(__dirname, 'brand.html')
        const html = fs.readFileSync(htmlPath, 'utf-8')
        res.statusCode = 200
        res.setHeader('Content-Type', 'text/html; charset=utf-8')
        res.end(html)
      })
    },
  }
}

export default defineConfig({
  plugins: [brandLoginPage()],
  server: {
    host: '0.0.0.0',
    port: 3100,
    proxy: {
      '/api': {
        target: 'http://localhost:8889',
        changeOrigin: true
      }
    }
  }
})
