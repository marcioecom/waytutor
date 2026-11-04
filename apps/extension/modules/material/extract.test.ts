import { describe, expect, it } from 'vitest'
import { extractMaterial } from './extract'

function extract(html: string) {
  return extractMaterial(new DOMParser().parseFromString(`<body>${html}</body>`, 'text/html'))
}

function texts(html: string) {
  return extract(html).passages.map((passage) => passage.text)
}

describe('extractMaterial', () => {
  it('turns each paragraph into one passage, keeping headings first', () => {
    expect(texts('<h1>Banco de dados</h1><p>Primeiro parágrafo.</p><p>Segundo parágrafo.</p>')).toEqual([
      'Banco de dados',
      'Primeiro parágrafo.',
      'Segundo parágrafo.',
    ])
  })

  it('reads a whole list as a single passage', () => {
    expect(texts('<ul><li>Primeiro item</li><li>Segundo item</li></ul>')).toEqual([
      'Primeiro item; Segundo item.',
    ])
  })

  it('reads one table row per passage', () => {
    expect(
      texts(
        '<table><tr><th>Sigla</th><th>Nome</th></tr><tr><td>SGBD</td><td>Banco de dados</td></tr></table>',
      ),
    ).toEqual(['Sigla. Nome', 'SGBD. Banco de dados'])
  })

  it('marks the passage a figure follows', () => {
    const material = extract(
      '<p>Antes da figura.</p><figure><img alt="Diagrama"><figcaption>Figura 1 — Diagrama</figcaption></figure><p>Depois da figura.</p>',
    )

    expect(material.passages.map((passage) => passage.text)).toEqual([
      'Antes da figura.',
      'Depois da figura.',
    ])
    expect(material.figureAfter).toEqual([0])
  })

  it('drops noise: stray urls, sources, captions and image file names', () => {
    expect(
      texts(
        '<p>https://exemplo.com/arquivo</p><p>Fonte: IBGE</p><p>Figura 2 — Gráfico</p><p>tabela-1.png</p><p>Conteúdo que deve ser lido.</p>',
      ),
    ).toEqual(['Conteúdo que deve ser lido.'])
  })

  it('stops at the references heading', () => {
    expect(texts('<p>Conteúdo real.</p><h2>Referências</h2><p>Não deve aparecer.</p>')).toEqual([
      'Conteúdo real.',
    ])
  })

  it('drops the first occurrence of a heading repeated later', () => {
    const html = [
      '<h2>Sumário</h2>',
      '<h3>Introdução</h3>',
      '<h3>Conclusão</h3>',
      '<h3>Introdução</h3>',
      '<p>Texto da introdução.</p>',
      '<h3>Conclusão</h3>',
      '<p>Texto da conclusão.</p>',
    ].join('')

    expect(texts(html)).toEqual([
      'Sumário',
      'Introdução',
      'Texto da introdução.',
      'Conclusão',
      'Texto da conclusão.',
    ])
  })

  it('follows the section while collecting passages', () => {
    const material = extract('<h2>Modelo relacional</h2><p>Uma tabela.</p><p>Outra tabela.</p>')

    expect(material.passages.map((passage) => passage.section)).toEqual([
      'Modelo relacional',
      'Modelo relacional',
      'Modelo relacional',
    ])
  })
})
