import { scanMessageType } from '../modules/material/messages'
import { extractMaterial } from '../modules/material/extract'

export default defineContentScript({
  matches: [
    'https://ava.catolica.edu.br/*',
    'https://conteudo.catolica.edu.br/*',
    'https://*.sequences.api.brightspace.com/*',
  ],
  allFrames: true,
  matchAboutBlank: true,
  runAt: 'document_idle',
  main() {
    chrome.runtime.onMessage.addListener((message: unknown, _sender, sendResponse) => {
      if (!message || typeof message !== 'object') return
      if ((message as { type?: unknown }).type !== scanMessageType) return

      sendResponse(extractMaterial(document))
    })
  },
})
