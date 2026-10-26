import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

// En desarrollo, Vite reenvía los pedidos a /api al backend (go run ./cmd/api, puerto 8080):
// así el frontend y la API se ven como un solo sitio y no hace falta tocar el backend.
export default defineConfig({
  plugins: [react()],
  server: {
    proxy: { '/api': 'http://localhost:8080' },
  },
})
