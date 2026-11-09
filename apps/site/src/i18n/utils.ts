import { defaultLang, htmlLang, ui, type Lang, type UiKey } from './ui';

export function getLangFromUrl(url: URL): Lang {
  const [, lang] = url.pathname.split('/');
  if (lang in ui) return lang as Lang;
  return defaultLang;
}

export function getHtmlLang(url: URL): string {
  return htmlLang[getLangFromUrl(url)];
}

export function useTranslations(lang: Lang) {
  return function t(key: UiKey): string {
    return ui[lang][key] ?? ui[defaultLang][key];
  };
}

// Default language lives at the root, so only non-default paths get a prefix.
export function useTranslatedPath(lang: Lang) {
  return function translatePath(path: string, target: Lang = lang): string {
    return target === defaultLang ? path : `/${target}${path}`;
  };
}
