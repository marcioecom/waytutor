export const languages = {
  pt: 'Português',
  en: 'English'
} as const;

export type Lang = keyof typeof languages;

export const defaultLang: Lang = 'pt';

/** Valid `html lang` values, keyed like `languages`. */
export const htmlLang: Record<Lang, string> = {
  pt: 'pt-BR',
  en: 'en'
};

/** UI strings shared by the layout, header and footer of every page. */
export const ui = {
  pt: {
    'a.skip': 'Pular para o conteúdo',
    'cta.install': 'Instalar',
    'footer.project': 'WayTutor, projeto da Waypoint Labs.',
    'footer.contact': 'Contato:',
    'footer.privacy': 'Política de privacidade'
  },
  en: {
    'a.skip': 'Skip to content',
    'cta.install': 'Install',
    'footer.project': 'WayTutor, a Waypoint Labs project.',
    'footer.contact': 'Contact:',
    'footer.privacy': 'Privacy policy'
  }
} as const;

export type UiKey = keyof (typeof ui)[typeof defaultLang];
