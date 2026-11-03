# 0001 — One synthesis per Passage

Each Passage is synthesised and played on its own, instead of generating the
audio of the whole Material in one call. Audio is cached per Passage, the next
one is synthesised while the current one plays, and speed is applied in the
player (`playbackRate`), never in the synthesis.

Providers cap input per request (5k to 40k characters depending on the model),
resuming from a saved position and skipping a paragraph are only cheap with
independent units, and a fix in the Pronunciation dictionary invalidates one
Passage instead of the whole Material.

**Consequences:** prosody can shift at the boundary between Passages. Keeping
voice, model and style instructions constant and synthesising whole paragraphs
(never lone sentences) hides most of it. Sentence highlight is estimated from
character offset, not from provider timestamps, which keeps the voice choice
free.

**Rejected:** one call per Material (hits the input cap, delays the start and
loses everything when a paragraph is skipped). Streaming speech: the Gateway
only exposes it as experimental, and it is not worth the dependency.
