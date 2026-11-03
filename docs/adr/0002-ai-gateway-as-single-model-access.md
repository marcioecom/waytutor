# 0002 — Vercel AI Gateway as the single model access

Every model call — speech and the Summary LLM — goes through the Vercel AI
Gateway with one key and a monthly budget. The choice is commercial, not
technical: the credits already exist, the Gateway charges the provider list
price without markup, and it exposes Gemini, ElevenLabs, OpenAI and others in
one catalogue, so changing voice or model is configuration instead of a new
account.

The Gateway exposes speech in beta and does not implement the OpenAI-compatible
`/v1/audio/speech` route. The call is `POST /v4/ai/speech-model`, which returns
base64 audio in JSON. Streaming speech exists only as
`experimental_streamSpeech`, which is why the product is not streaming: audio is
chunked per Passage instead.

**Consequences:** speech depends on a beta endpoint with a protocol version
header. If the model id or the protocol changes, the fix is confined to
`internal/speech`. The Gateway key lives only on the server.

**Rejected:** calling Google directly (a second invoice and a second quota for
the same capability), and keeping the key in the extension, which cannot work
with more than one account.
