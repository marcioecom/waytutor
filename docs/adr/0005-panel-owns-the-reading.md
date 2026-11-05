# 0005. The panel owns the Reading

In v1 the service worker owned the Reading state machine, because the worker is
the only context every surface can reach. Chrome also terminates the worker
whenever it needs to, so the state machine had to recover itself from
`chrome.storage` on every wake-up.

In v2 the side panel owns the Reading. The panel is a page, so it lives as long
as it is open, it can call `chrome.tabs.sendMessage`, `chrome.webNavigation` and
`chrome.identity.launchWebAuthFlow` on its own, and react-query holds the state
without a persistence layer of its own. The service worker keeps one job,
`setPanelBehavior`, because that has to be registered before a panel exists. The
panel also owns the audio element. There is no offscreen document, because the
context that decides what to play is the same one that can play it.

**Consequences:** closing the panel ends the Reading. For a side panel product
that is the expected behaviour, and it removes the whole class of "the worker
died mid-Reading" bugs. Preferences, Reading progress and the Access session are
the only things written to `chrome.storage`.

**Rejected:** keeping the state in the worker, which fights the platform;
putting the state in the content script, a page we do not control that navigates
away; and an offscreen document for audio. The offscreen document came from v1,
where the worker owned the state and could not play sound. With the panel owning
the Reading it would add a second context, a message protocol and one more
permission for no benefit.

The client keeps a look-ahead buffer of three Passages and nothing more. The
durable cache is the server's, keyed by everything that changes the audio, so
re-reading a Material never reaches the model again.
