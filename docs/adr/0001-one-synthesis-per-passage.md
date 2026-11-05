# 0001. One synthesis per Passage

Each Passage is synthesised and played on its own, instead of generating the
audio of the whole Material in one call. Audio is cached per Passage, the next
one is synthesised while the current one plays, and the player applies speed
through `playbackRate`, never in the synthesis.

Providers cap input per request at 5k to 40k characters depending on the model.
Resuming from a saved position and skipping a paragraph are cheap only with
independent units. A fix in the Pronunciation dictionary invalidates one Passage
instead of the whole Material.

**Consequences:** prosody can shift at the boundary between Passages. Most of it
goes away when voice, model and style instructions stay constant and synthesis
works on whole paragraphs, never lone sentences. Sentence highlight comes from
character offset, not from provider timestamps, which keeps the voice choice
free.

**Rejected:** one call per Material, which hits the input cap, delays the start
and loses everything when a paragraph is skipped. Also rejected: streaming
speech, because the Gateway exposes it only as experimental and the dependency
is not worth it.
