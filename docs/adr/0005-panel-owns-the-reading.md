# 0005 — The panel owns the Reading

In v1 the service worker owned the Reading state machine, because the worker is
the only context every surface can reach. The worker is also the context Chrome
kills whenever it wants, so the state machine had to recover itself from
`chrome.storage` on every wake-up.

In v2 the side panel owns the Reading. The panel is a page: it lives as long as
it is open, it can call `chrome.tabs.sendMessage`, `chrome.webNavigation`,
`chrome.identity.launchWebAuthFlow` and `chrome.runtime.sendMessage` on its own,
and react-query holds the state without a persistence layer of its own. The
service worker is left with one job — `setPanelBehavior` — because that has to
be registered before a panel exists. The offscreen document owns the audio
element and fetches the audio itself, so nothing travels as base64.

**Consequences:** closing the panel ends the Reading. For a side panel product
that is the expected behaviour, and it removes the whole class of "the worker
died mid-Reading" bugs. Preferences, Reading progress and the Access session are
the only things written to `chrome.storage`.

**Rejected:** keeping the state in the worker (fights the platform), and putting
the state in the content script (a page we do not control and that navigates
away).

The browser HTTP cache replaces v1's hand-written Cache Storage layer: audio
responses are immutable and carry an `ETag`, and regenerating one costs nothing
because the server caches the object in R2.
