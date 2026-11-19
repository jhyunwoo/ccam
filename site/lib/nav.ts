// The docs sidebar, in reading order: start, the features ccam exists for,
// the features around them, task guides, then reference detail.
// Groups only decide where a page is listed; URLs keep their section folder.
// Every English page must appear here exactly once (scripts/check-links.ts).
export const nav = [
  { key: "start", pages: ["docs/install", "docs/quick-start"] },
  { key: "core", pages: ["features/accounts", "features/launch", "features/dashboard", "features/share", "features/usage", "features/next"] },
  { key: "more", pages: ["features/shortcuts", "features/folder", "features/pick", "features/agent", "features/agents", "features/api", "features/mcp", "features/ext", "features/migrate", "features/doctor", "features/update"] },
  { key: "guides", pages: ["guides/multiple-claude-code-accounts", "docs/sharing", "docs/usage-limits", "guides/usage-limit-workflow", "guides/remote-control"] },
  { key: "reference", pages: ["docs/commands", "guides/claude-config-dir", "compare/aliases", "compare/account-switchers", "docs/faq"] },
  { key: "blog", pages: ["blog/introducing-ccam"] },
] as const;

export type GroupKey = (typeof nav)[number]["key"];
