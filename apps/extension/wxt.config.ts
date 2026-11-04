import tailwindcss from '@tailwindcss/vite'
import { defineConfig } from 'wxt'

export default defineConfig({
  modules: ['@wxt-dev/module-react'],
  vite: () => ({ plugins: [tailwindcss()] }),
  manifest: {
    name: 'WayTutor',
    description: 'Lê em voz alta o material do portal e resume o que foi lido.',
    action: { default_title: 'WayTutor' },
    permissions: ['storage', 'sidePanel', 'identity'],
    host_permissions: [
      'https://ava.catolica.edu.br/*',
      'https://conteudo.catolica.edu.br/*',
      'https://*.sequences.api.brightspace.com/*',
      'https://waytutor-api.marcio.run/*',
      'http://localhost:8080/*',
    ],
  },
})
