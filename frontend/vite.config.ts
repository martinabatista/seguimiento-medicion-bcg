import { defineConfig, type Plugin } from 'vite'
import react from '@vitejs/plugin-react'

// Los navegadores llevan cualquier nombre terminado en .localhost a la propia compu, así que la app
// también se abre en http://software-metrics.localhost:5173. Este plugin hace que Vite muestre ese link.
function linkConNombre(): Plugin {
  return {
    name: 'link-con-nombre',
    configureServer(server) {
      const mostrar = server.printUrls
      server.printUrls = () => {
        if (server.resolvedUrls) {
          server.resolvedUrls.local = server.resolvedUrls.local.map((url) =>
            url.replace('//localhost', '//software-metrics.localhost'),
          )
        }
        mostrar()
      }
    },
  }
}

// En desarrollo, Vite reenvía los pedidos a /api al backend (go run ./cmd/api, puerto 8080):
// así el frontend y la API se ven como un solo sitio y no hace falta tocar el backend.
export default defineConfig({
  plugins: [react(), linkConNombre()],
  server: {
    proxy: { '/api': 'http://localhost:8080' },
  },
})
