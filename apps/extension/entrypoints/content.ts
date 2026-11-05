import { highlightMessageType, scanMessageType } from '@/modules/material/messages'
import { extractMaterial } from '@/modules/material/extract'

const HIGHLIGHT_CLASS = 'waytutor-current'

// This runs inside the portal, where our tokens do not exist.
const HIGHLIGHT_STYLE = `
.${HIGHLIGHT_CLASS} {
  outline: 2px solid #f0a63c;
  outline-offset: 4px;
  border-radius: 4px;
  background-color: rgb(240 166 60 / 0.12);
  transition: background-color 150ms ease, outline-color 150ms ease;
}
`

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
    let elements: (HTMLElement | null)[] = []

    injectStyle()

    chrome.runtime.onMessage.addListener((message: unknown, _sender, sendResponse) => {
      if (!message || typeof message !== 'object') return

      const command = message as { type?: string; index?: number }

      if (command.type === scanMessageType) {
        const extracted = extractMaterial(document)
        elements = extracted.elements
        sendResponse(extracted.material)
        return
      }

      if (command.type === highlightMessageType && typeof command.index === 'number') {
        highlight(command.index)
      }
    })

    function highlight(index: number) {
      clearHighlight()
      const element = elements[index]
      if (!element) return

      element.classList.add(HIGHLIGHT_CLASS)
      element.scrollIntoView({ block: 'center', behavior: 'smooth' })
    }

    function clearHighlight() {
      for (const element of document.querySelectorAll(`.${HIGHLIGHT_CLASS}`)) {
        element.classList.remove(HIGHLIGHT_CLASS)
      }
    }

    function injectStyle() {
      if (document.getElementById('waytutor-highlight')) return

      const style = document.createElement('style')
      style.id = 'waytutor-highlight'
      style.textContent = HIGHLIGHT_STYLE
      document.documentElement.append(style)
    }
  },
})
