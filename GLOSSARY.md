# WayTutor

WayTutor is a browser extension that reads the student portal content aloud. It
pronounces terms and acronyms correctly, pauses so the listener can look at
figures, and ends with a spoken summary.

## Terms

| Term | Definition | Code name | Avoid |
| --- | --- | --- | --- |
| **Topic** | A unit of content in the portal, opened through the sequence viewer. It may hold a Material, a quiz, a PDF or a video. | `Topic` | page, lesson, activity |
| **Material** | The textual content that can be read aloud, extracted from a Topic: headings, paragraphs, lists, tables and figures. | `Material` | document, article, PDF |
| **Passage** | The smallest continuous speakable unit of a Material, usually one paragraph or one list item. | `Passage` | chunk, block, snippet |
| **Figure** | An image inside the Material. Its caption and credits are not part of the Reading. | `Figure` | picture, graph |
| **Noise** | A Passage that exists in the Material but is not read: figure captions, source credits, stray URLs, image file names, bibliography and institutional boilerplate. | `Noise` | junk, irrelevant content |
| **Reading** | Playing a Material aloud from the first Passage to the last. | `Reading` | narration, playback, audio |
| **Figure pause** | The silence inserted after the Passage that precedes a Figure, so the listener can look at the image. | `figurePauseSeconds` | image pause, break |
| **Reading progress** | The stored state of one Reading: Material, current Passage, voice, speed, finished. | `ReadingProgress` | player state, bookmark |
| **Voice** | A selectable voice for the Reading. Changing it applies from the next Passage. | `Voice` | timbre, speaker |
| **Summary** | Text and audio generated from a Material that was already read, strictly grounded in it. It introduces nothing external to the Material. | `Summary` | explanation, free summary |
| **Pronunciation dictionary** | A persistent set of rules describing how a term is spoken, including acronyms, proper nouns and foreign terms. | `Pronunciation` | lexicon |
| **Account** | A person identified by a Google account, with status `pending`, `active` or `paused`. | `Account` | user record, profile |
| **Access session** | An opaque token that authenticates one extension installation. | `AccessSession` | login, token |

There is one Voice: the Summary is read with the same voice as the Reading.

The word **session** alone is never used. It is either Reading progress or
Access session.
