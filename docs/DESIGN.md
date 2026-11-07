# Design

## Theme

**The Green Study.** The panel is a study desk at night: a dark forest room, the
pool of a green banker's lamp, a stack of books, a pocket watch and a capybara
professor reading with you. It is not "dark mode"; it is a warm place next to
the white portal.

The world lives in three objects: the **capybara** (emotional state), the
**banker's lamp** (the light that says what is happening) and the **pocket
watch** (progress, in place of a progress bar). Everything else stays a tool:
title, controls and summary with a clean hierarchy.

The panel is part of Waypoint Labs and borrows its instrument language: **Geom**
as the only typeface, a chartreuse **signal** for live state and a small
ordered-dither texture in the lamp's light. The study room stays warm; the
instrument stays precise.

Behavior references: SwapGun (world plus mascot in the layout plus an object in
place of the spinner), Duolingo (mascot as a state machine), Lofi Girl (study as
a place), Waypoint Labs (navigation instrument, dither, chartreuse signal).

## Colors

**Restrained on a dark surface**: two green-tinted neutrals plus cream ink, one
sage accent for action, one chartreuse signal for live state and the lamp's
half-light for the figure pause. The lamp bulb is the only warm light in the
interface, because it is a physical object in the room, not a UI accent.

| Token | Value | Use |
| --- | --- | --- |
| `--color-night` | `#0f150e` | panel background |
| `--color-night-2` | `#151c13` | control surfaces |
| `--color-night-3` | `#1b2418` | fields, hover, disabled |
| `--color-night-4` | `#232e1e` | watch track, skeleton |
| `--color-line` | `#2a3524` | rules and borders |
| `--color-line-2` | `#3d4a33` | control borders |
| `--color-ink` | `#f3eedc` | primary text (15.9:1) |
| `--color-ink-muted` | `#c6cbb4` | secondary text (11.1:1) |
| `--color-ink-faint` | `#98a08a` | meta and signature (6.8:1) |
| `--color-moss` | `#a9cc8f` | primary action, selection, current passage, wordmark accent (10.3:1) |
| `--color-moss-bright` | `#c1dda8` | hover |
| `--color-moss-deep` | `#547a44` | lamp shade highlight, list hover |
| `--color-signal` | `oklch(0.79 0.15 115)` | live progress arc, live status dot, focus ring |
| `--color-bulb` | `#ffd79a` | lamp bulb on |
| `--color-bulb-dim` | `#c9a06a` | lamp bulb paused, figure annotation |
| `--color-bulb-off` | `#6b5a43` | lamp bulb idle or error |
| `--color-brass` | `#b98a4f` | lamp arm and base |
| `--color-leaf` | `#466b4f` | capybara cardigan |
| `--color-capy` | `#9a6a45` | fur |
| `--color-capy-2` | `#a4714b` | muzzle |
| `--color-capy-ear` | `#7a5236` | ears |
| `--color-capy-ink` | `#241b15` | eyes, glasses, nose |
| `--color-danger` | `#f2a48e` | error (8.3:1 on error background) |

Rules:

- Text on moss or signal is always `--color-night`. Never cream on either.
- Moss is action and selection; signal is **only** live reading; the lamp's
  half-light is **only** the figure pause. A color may not do two jobs.
- The body carries a soft green radial glow at the top; the lamp adds a warm
  pool locally. The room is green, the lamp is warm.
- The scene is the only place where more than two hues meet.

## Typography

One family, as on the Waypoint site: **Geom** 400–700, OFL, packaged in
`public/fonts` (28 KB). It carries the wordmark, titles, UI and data. The
geometric forms read as an instrument, not as a generic startup sans.

| Role | Size | Weight | Use |
| --- | --- | --- | --- |
| Wordmark | 19 px | 600 | `WayTutor` |
| Material title | 15 px | 600 | at most 3 lines, 1.35 line-height |
| Screen titles | 18–20 px | 600 | `Resumo`, login, empty state |
| Interface | 11.5–13 px | 400–500 | buttons, body, passage list |
| Data | 10–11 px | 400–600 | meta, passage number, speed, status, account |

- `tabular-nums` on everything that changes with the reading.
- Uppercase only in the status badge (tracking 0.14 em), never as a section
  label.
- Body prose maxes out at 65–75ch; the panel is narrower than that by design.

## Components

- **Brand plus state** — wordmark `Way` cream plus `Tutor` moss; tagline in
  Geom lowercase at the right. Status is a dot plus uppercase Geom label
  (`PRONTO`, `LENDO`, `PAUSADO`, `PREPARANDO`); the dot is signal and pulses
  while reading.
- **Pocket watch** — 82 px button; 56 px face; progress ring (r=38, stroke 3.5)
  in signal with animated `stroke-dashoffset`; 9×6 moss crown at the top. Hover
  lights the face; paused, the icon stays moss on the dark face. The ring is the
  only place where progress is a live signal.
- **Capybara** — 104×96 SVG, one component with four moods: green cardigan,
  cream shirt, round glasses, wide muzzle, tiny ears, blinking eyes. The body
  tilts −1.2° when paused and +1.2° on error.
- **Banker's lamp** — 36×90 SVG: base, brass arm and a green glass shade. It
  sits at the left of the stage, on the same table line as the capybara. The
  bulb and halo change intensity by state; the lamp does not move.
- **Book stack** — 26×22 SVG: four books (slate blue, paper, wine, page block)
  resting on the table line between the lamp base and the capybara. Scenery,
  not state: it does not animate.
- **Light** — radial halo behind the capybara and the lamp; opacity by state
  (0.30 idle to 0.95 reading), a 6 second breathing during reading and a 1.6
  second warm-up during preparation.
- **Dither** — an 8×8 ordered Bayer stipple in signal, masked to a circle inside
  the lamp's light pool. It is a texture, not decoration: it brightens with the
  reading state (10% error to 35% reading) and drifts 8 px over 7 seconds. The
  light itself is a blurred ellipse, so neither layer ends on a rectangle. The
  WebGL field of the site is not used here; this is its quiet, static sibling.
- **Table line** — 1 px of `--color-line-2` with faded ends at the foot of the
  stage. It makes the lamp, the books and the capybara share one plane.
- **Passage list** — scrollable, 13 px, Geom numbers; the current passage gets
  a 2 px moss left border with a soft pulse; previous and next stay muted. A
  figure pause is marked with a small picture icon and `Figura a seguir` in the
  lamp's half-light (`--color-bulb-dim`).
- **Speed** — segmented control with six options, 30 px tall, Geom tabular
  numbers; active in moss with dark text.
- **Voice settings** — `<details>` of 40 px with the current voice at the right;
  inside, a custom voice picker and the figure pause. Rare adjustments stay
  folded. The picker groups **Favoritas** above **Todas as vozes**; each voice
  has a star, and the current voice keeps a signal check.
- **Summary** — a screen with a Geom title, scrollable text and a moss primary
  button.
- **Error** — text in `--color-danger` plus the capybara's tilting; a banner
  for blocking errors.
- **Loading** — a small spinner next to a plain sentence; the scan is short.

Every control has hover, focus-visible (signal ring), active and disabled.
Targets are at least 36 px (exception: the speed segment, 30 px tall with a wide
horizontal hit area).

## Mascot states

| App state | Mood | Eyes | Mouth | Lamp | Body |
| --- | --- | --- | --- | --- | --- |
| Playing | `reading` | down | straight line | bulb on, halo 0.95 plus breathing | straight |
| Preparing | `thinking` | closed | half smile | bulb on, halo pulses 0.55–0.92 | straight |
| Paused / ready | `waiting` | forward | half smile | bulb half lit, halo 0.45 | tilt −1.2° |
| Error | `error` | wide | wavy | bulb low, halo 0.30 | tilt +1.2° |

Rules: the mascot never covers a control or text; the blink is 0.25 s every
7.2 s; no state relies on color alone.

## Layout

- Useful width 320 to 440 px; 16 px side padding.
- Single column, 12 to 14 px gaps: brand, material head, stage, passage list,
  speed, voice settings, account footer.
- The material title is context; the stage is the hero.
- The passage list flexes and scrolls; the stage and controls stay fixed.
- Below 360 px the stage, watch and arrows shrink together; below 300 px the
  scene (lamp plus capybara) leaves and the transport remains.

## Motion

- Blink 7.2 s; light breathing 6 s; light warm-up 1.6 s (preparing); ring
  520 ms; halo 700 ms; tilt 420 ms; dither drift 7 s.
- Curve `cubic-bezier(0.22, 1, 0.36, 1)`.
- Nothing animates layout, nothing enters in a cascade, nothing covers text.
- `prefers-reduced-motion: reduce` turns everything off, including the blink,
  the warm-up and the dither drift.
