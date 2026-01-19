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
        if (pathname === '/app.js') {
          const jsPath = path.resolve(__dirname, 'app.js')
          const js = fs.readFileSync(jsPath, 'utf-8')
          res.statusCode = 200
          res.setHeader('Content-Type', 'application/javascript; charset=utf-8')
          res.setHeader('Cache-Control', 'no-store')
          res.end(js)
          return
        }

        const isBrandEntry =
          pathname === '/brand' ||
          pathname === '/brand/' ||
          pathname === '/brand/login' ||
          pathname === '/brand/login/'
        if (!isBrandEntry) return next()

        const htmlPath = path.resolve(__dirname, 'brand.html')
        const html = fs.readFileSync(htmlPath, 'utf-8')
        res.statusCode = 200
        res.setHeader('Content-Type', 'text/html; charset=utf-8')
        res.setHeader('Cache-Control', 'no-store')
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
