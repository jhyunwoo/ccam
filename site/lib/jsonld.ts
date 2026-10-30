import { abs, REPO, SITE } from "./paths";

export const VERSION = "0.13.1"; // shown on the site; bump with releases
export const CLAUDE_CODE_VERSION = "2.1.302";

// "CCAM" also names other projects; the full name tells search engines which one this is.
const FULL_NAME = "Claude Code Account Manager";
const WEBSITE = { "@id": SITE + "/#website" };
const ORGANIZATION = { "@id": SITE + "/#organization" };
export const AUTHOR = { "@type": "Person", name: "jhyunwoo", url: "https://github.com/jhyunwoo" };

export function websiteLd(lang: string) {
  return [
    {
      "@context": "https://schema.org",
      "@type": "WebSite",
      ...WEBSITE,
      name: "ccam",
      alternateName: FULL_NAME,
      url: SITE + "/",
      inLanguage: lang,
      publisher: ORGANIZATION,
      copyrightYear: 2026,
      copyrightHolder: { "@type": "Organization", name: "Moveto Team" },
    },
    {
      "@context": "https://schema.org",
      "@type": "Organization",
      ...ORGANIZATION,
      name: "ccam",
      alternateName: FULL_NAME,
      url: SITE + "/",
      logo: abs("/icon-512.png"),
      sameAs: [REPO],
    },
  ];
}

export function softwareLd(description: string, lang: string) {
  return {
    "@context": "https://schema.org",
    "@type": "SoftwareApplication",
    name: "ccam",
    alternateName: FULL_NAME,
    description,
    author: AUTHOR,
    publisher: ORGANIZATION,
    inLanguage: lang,
    applicationCategory: "DeveloperApplication",
    operatingSystem: "Linux, macOS",
    softwareVersion: VERSION,
    downloadUrl: `${REPO}/releases/latest`,
    license: `${REPO}/blob/main/LICENSE`,
    offers: { "@type": "Offer", price: "0", priceCurrency: "USD" },
    url: SITE + "/",
  };
}

export function articleLd(opts: { type: "TechArticle" | "BlogPosting"; title: string; description: string; url: string; published: string; updated: string; lang: string; image: string; crumbs: { name: string; url: string }[] }) {
  return [
    {
      "@context": "https://schema.org",
      "@type": opts.type,
      headline: opts.title,
      description: opts.description,
      inLanguage: opts.lang,
      datePublished: opts.published,
      dateModified: opts.updated,
      url: opts.url,
      mainEntityOfPage: opts.url,
      image: opts.image,
      author: AUTHOR,
      publisher: ORGANIZATION,
      isPartOf: WEBSITE,
    },
    {
      "@context": "https://schema.org",
      "@type": "BreadcrumbList",
      itemListElement: opts.crumbs.map((c, i) => ({ "@type": "ListItem", position: i + 1, name: c.name, item: c.url })),
    },
  ];
}

export function webPageLd(opts: { title: string; description: string; url: string; updated: string; lang: string }) {
  return {
    "@context": "https://schema.org",
    "@type": "WebPage",
    name: opts.title,
    description: opts.description,
    url: opts.url,
    inLanguage: opts.lang,
    dateModified: opts.updated,
    isPartOf: WEBSITE,
  };
}

export function faqLd(entries: { question: string; answer: string }[], lang: string) {
  return {
    "@context": "https://schema.org",
    "@type": "FAQPage",
    inLanguage: lang,
    mainEntity: entries.map((e) => ({ "@type": "Question", name: e.question, acceptedAnswer: { "@type": "Answer", text: e.answer } })),
  };
}

/** Serializes JSON-LD safely for a <script> element. */
export function ldJson(data: unknown): string {
  return JSON.stringify(data).replace(/</g, "\\u003c");
}
