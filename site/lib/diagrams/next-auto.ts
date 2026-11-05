import { box, figure, glow, label, line, packet, reveal } from "../motion.ts";

// `ccam next --auto`: ask every logged-in account for its usage at once,
// then usage.Policy.Choose() takes, below auto_max_usage (overflowing to
// below 100), the least weekly usage within 3 points of the lowest 5-hour
// usage.
export default function nextAuto(): string {
  const h = 44, colW = 120, cols = [10, 140, 270], cx = cols.map((x) => x + colW / 2);
  const yTop = 10, yFetch = 70, yAcct = 148, yBest = 226, yRun = 296;
  const accts = [
    { name: "personal", sub: "5h 72%" },
    { name: "work", sub: "5h 8%" },
    { name: "side", sub: "5h 100%" },
  ];
  const parts: string[] = [
    line([[200, yTop + h], [200, yFetch]]),
    ...cx.map((x) => line([[200, yFetch + h], [200, 132], [x, 132], [x, yAcct]])),
    line([[cx[0], yAcct + h], [cx[0], 208], [200, 208], [200, yBest]]),
    line([[cx[1], yAcct + h], [cx[1], yBest]]),
    line([[200, yBest + h], [200, yRun]]),
    box(60, yTop, 280, h, "ccam next --auto", { sub: "pickAuto()" }),
    box(60, yFetch, 280, h, "FetchAll()", { sub: "usage.Options{}" }),
    ...accts.map((a, i) => box(cols[i], yAcct, colW, h, a.name, { sub: a.sub })),
    box(60, yBest, 280, h, "Policy.Choose()", { sub: "5h ≤ min+3 → min(week) · < max · < 100" }),
    box(60, yRun, 280, h, "LaunchWith(work)", { sub: "syscall.Exec(claude)" }),
  ];
  parts.push(
    glow(60, yTop, 280, h, 0),
    packet([200, yTop + h], [200, yFetch], 0.5),
    glow(60, yFetch, 280, h, 1.1),
    ...cx.map((x) => packet([x, 132], [x, yAcct], 1.6)),
    ...cols.map((x) => glow(x, yAcct, colW, h, 2.2)),
    reveal(3.2, label(cx[2], yAcct + h + 18, "✗", { anchor: "middle" })),
    packet([200, 208], [200, yBest], 3.4),
    glow(60, yBest, 280, h, 4.0),
    packet([200, yBest + h], [200, yRun], 4.6),
    glow(60, yRun, 280, h, 5.2),
    reveal(5.2, label(200, 364, "→ personal → work", { anchor: "middle" })),
  );
  return figure("next-auto", 400, 374, 8.4, parts);
}
