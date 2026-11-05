import tailwindcss from '@tailwindcss/vite'
import { defineConfig } from 'wxt'

export default defineConfig({
  modules: ['@wxt-dev/module-react'],
  vite: () => ({ plugins: [tailwindcss()] }),
  manifest: {
    name: 'WayTutor',
    description: 'Lê em voz alta o material do portal e resume o que foi lido.',
    icons: {
      16: 'icon/16.png',
      32: 'icon/32.png',
      48: 'icon/48.png',
      128: 'icon/128.png',
    },
    action: {
      default_title: 'WayTutor',
      default_icon: { 16: 'icon/16.png', 32: 'icon/32.png' },
    },
    permissions: ['storage', 'sidePanel', 'identity', 'webNavigation'],
    host_permissions: [
      'https://ava.catolica.edu.br/*',
      'https://conteudo.catolica.edu.br/*',
      'https://*.sequences.api.brightspace.com/*',
      'https://waytutor-api.marcio.run/*',
      'http://localhost:8080/*',
    ],
  },
  hooks: {
    // A published item should not ask for access to a development server.
    'build:manifestGenerated': (wxt, manifest) => {
      if (wxt.config.mode === 'development') return

      manifest.host_permissions = manifest.host_permissions?.filter(
        (host: string) => !host.includes('localhost'),
      )
    },
  },
})
