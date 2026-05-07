import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

export default defineConfig({
  plugins: [vue()],
  css: {
    preprocessorOptions: {
      sass: {
        silenceDeprecations: ['import']
      }
    }
  },
  optimizeDeps: {
    include: ['@dennislee928/nothingx-react-components']
  }
})
