import type { Material, Passage } from './material'

const SKIP_TAGS = new Set([
  'SCRIPT',
  'STYLE',
  'NOSCRIPT',
  'IFRAME',
  'SVG',
  'CANVAS',
  'NAV',
  'HEADER',
  'FOOTER',
  'ASIDE',
  'FORM',
  'BUTTON',
  'INPUT',
  'SELECT',
  'TEXTAREA',
  'VIDEO',
  'AUDIO',
])

const BLOCK_TAGS = new Set([
  'P',
  'LI',
  'BLOCKQUOTE',
  'H1',
  'H2',
  'H3',
  'H4',
  'H5',
  'H6',
  'TD',
  'TH',
  'DD',
  'DT',
])

const MAIN_SELECTORS = [
  'article',
  'main',
  '[role="main"]',
  '.d2l-html-block',
  '.course-content',
  '.lesson-content',
  '.page-content',
  '.content',
]

const URL_ONLY = /^(https?:\/\/|www\.)\S+$/i
const SOURCE_LINE = /^(fonte|source)\s*:/i
const FIGURE_CAPTION = /^(figura|tabela|figure|table)\s+\d+/i
const IMAGE_FILE = /\.(png|jpe?g|gif|webp|svg|bmp)\s*$/i
const REFERENCES = /^(referências|referencias|bibliografia)\b/i
const BOILERPLATE = /^na prática$/i
const RESUMES_READING = /^(saiba mais|conclusão|conclusao|considerações finais)\b/i
const NAVIGATION_ROLE = /^(navigation|menu|menubar|tablist|toolbar|search)$/i
const NAVIGATION_NAME = /(^|[\s_-])(nav|navigation|menu|toc|breadcrumb|sidebar|iterator|outline)([\s_-]|$)/

type Draft = {
  text: string
  section: string
  heading: boolean
}

function normalizeText(value: string): string {
  return value
    .replace(/\u00a0/g, ' ')
    .replace(/[\u200b-\u200d\ufeff]/g, '')
    .replace(/\s+/g, ' ')
    .trim()
}

function isNoise(text: string): boolean {
  if (text.length < 3) return true
  if (URL_ONLY.test(text)) return true
  if (SOURCE_LINE.test(text)) return true
  if (FIGURE_CAPTION.test(text)) return true
  if (IMAGE_FILE.test(text) && text.length < 120) return true
  if (/^[\p{P}\p{S}\s]+$/u.test(text)) return true
  if (/^bom trabalho!?$/i.test(text)) return true
  return false
}

// innerText inserts the line breaks that <br> and block elements imply, which
// textContent drops; test environments have no innerText.
function innerTextOf(element: Element): string {
  const rendered = (element as HTMLElement).innerText
  return typeof rendered === 'string' ? rendered : (element.textContent ?? '')
}

function textOf(element: Element): string {
  return normalizeText(innerTextOf(element))
}

function hasBlockChildren(element: HTMLElement): boolean {
  return element.querySelector('p, ul, ol, table, div, blockquote') !== null
}

function isVisible(element: HTMLElement): boolean {
  if (element.hidden) return false
  const style = getComputedStyle(element)
  return style.display !== 'none' && style.visibility !== 'hidden'
}

function linkDensity(element: HTMLElement): number {
  const text = textOf(element)
  if (!text) return 1

  const linkText = Array.from(element.querySelectorAll('a')).reduce(
    (total, link) => total + textOf(link).length,
    0,
  )
  return linkText / text.length
}

// The content frame ships a side menu, a breadcrumb and a section list around
// the Material. Without this the first Passages become a list of titles.
function isNavigationLike(element: HTMLElement): boolean {
  if (['NAV', 'ASIDE', 'HEADER', 'FOOTER'].includes(element.tagName)) return true

  const role = element.getAttribute('role')
  if (role && NAVIGATION_ROLE.test(role)) return true

  const name = `${element.id} ${typeof element.className === 'string' ? element.className : ''}`
  if (NAVIGATION_NAME.test(name.toLowerCase())) return true

  return element.querySelectorAll('a').length >= 3 && linkDensity(element) > 0.5
}

function findMainContainer(doc: Document): HTMLElement {
  let best: HTMLElement | null = null
  let bestScore = Number.NEGATIVE_INFINITY

  for (const selector of MAIN_SELECTORS) {
    for (const element of Array.from(doc.querySelectorAll<HTMLElement>(selector))) {
      if (!isVisible(element)) continue

      const text = textOf(element)
      if (text.length < 400) continue

      const density = linkDensity(element)
      if (density > 0.5 && element.querySelectorAll('a').length >= 3) continue

      const paragraphs = element.querySelectorAll('p').length
      const semantic = element.tagName === 'ARTICLE' || element.tagName === 'MAIN' ? 800 : 0
      const score = text.length + paragraphs * 40 - density * text.length + semantic

      if (score > bestScore) {
        bestScore = score
        best = element
      }
    }
  }

  return best ?? doc.body
}

export type ExtractedMaterial = {
  material: Material
  elements: (HTMLElement | null)[]
}

export function extractMaterial(doc: Document = document): ExtractedMaterial {
  const drafts: Draft[] = []
  const elements: (HTMLElement | null)[] = []
  const figures: number[] = []

  let skipRest = false
  let skipBoilerplate = false
  let section = ''
  let title = ''
  let lastWasFigure = false

  const pushText = (raw: string, element: HTMLElement | null, heading = false) => {
    const text = normalizeText(raw)
    if (isNoise(text)) return
    if (drafts.at(-1)?.text === text) return

    drafts.push({ text, section, heading })
    elements.push(element)
    lastWasFigure = false
  }

  const pushFigure = () => {
    const last = drafts.length - 1
    if (last < 0 || lastWasFigure) return

    figures.push(last)
    lastWasFigure = true
  }

  const walk = (node: Node) => {
    if (!(node instanceof HTMLElement)) return
    if (SKIP_TAGS.has(node.tagName)) return
    if (node.getAttribute('aria-hidden') === 'true') return
    if (node.matches('#sidebar, .d2l-sequence-viewer-header')) return
    if (isNavigationLike(node)) return

    if (node.tagName === 'FIGURE') {
      pushFigure()
      return
    }
    if (node.tagName === 'FIGCAPTION') return
    if (node.tagName === 'IMG') {
      const image = node as HTMLImageElement
      const decorative = image.getAttribute('alt') === '' || image.role === 'presentation'
      if (!decorative && image.width > 80) pushFigure()
      return
    }

    if (/^H[1-6]$/.test(node.tagName)) {
      const text = textOf(node)
      if (!text) return

      if (REFERENCES.test(text)) {
        skipRest = true
        return
      }
      if (BOILERPLATE.test(text)) {
        skipBoilerplate = true
        return
      }
      if (RESUMES_READING.test(text)) skipBoilerplate = false
      if (skipRest) return

      if (!title) title = text
      section = text
      pushText(text, node, true)
      return
    }

    if (node.tagName === 'UL' || node.tagName === 'OL') {
      if (skipRest || skipBoilerplate) return
      if (!isVisible(node)) return

      const items = Array.from(node.querySelectorAll(':scope > li'))
        .map((item) => textOf(item).replace(/[.;]+$/, ''))
        .filter((item) => item && !isNoise(item))

      // A list is one speakable unit: one request, one highlight, one Passage.
      if (items.length) pushText(`${items.join('; ')}.`, node)
      return
    }

    if (node.tagName === 'TABLE') {
      for (const row of Array.from(node.querySelectorAll('tr'))) {
        const cells = Array.from(row.querySelectorAll('th, td'))
          .map((cell) => textOf(cell))
          .filter(Boolean)

        if (cells.length) pushText(cells.join('. '), row as HTMLElement)
      }
      return
    }

    if (BLOCK_TAGS.has(node.tagName) && !hasBlockChildren(node)) {
      if (skipRest || skipBoilerplate) return
      if (!isVisible(node)) return

      const text = textOf(node)
      if (/^prezado\(a\) estudante/i.test(text)) return
      pushText(text, node)
      return
    }

    for (const child of Array.from(node.childNodes)) walk(child)
  }

  walk(findMainContainer(doc))

  if (!drafts.length && !skipRest) pushText(innerTextOf(doc.body), null)

  // A summary often repeats the section headings. When a heading shows up again
  // later, the first occurrence is the summary and should not be read.
  const lastHeading = new Map<string, number>()
  drafts.forEach((draft, position) => {
    if (draft.heading) lastHeading.set(draft.text.toLowerCase(), position)
  })

  const newIndex = new Map<number, number>()
  const passages: Passage[] = []
  const keptElements: (HTMLElement | null)[] = []

  drafts.forEach((draft, position) => {
    if (draft.heading && lastHeading.get(draft.text.toLowerCase()) !== position) return

    newIndex.set(position, passages.length)
    passages.push({ text: draft.text, section: draft.section })
    keptElements.push(elements[position] ?? null)
  })

  const figureAfter = [
    ...new Set(
      figures.flatMap((position) => {
        const moved = newIndex.get(position)
        return moved === undefined ? [] : [moved]
      }),
    ),
  ]

  const firstHeading = doc.querySelector('h1, h2')
  const heading = firstHeading ? textOf(firstHeading) : ''

  return {
    material: {
      title: title || heading || 'Material sem título',
      url: doc.location?.href ?? '',
      passages,
      figureAfter,
    },
    elements: keptElements,
  }
}
