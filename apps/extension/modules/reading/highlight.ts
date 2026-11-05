import { highlightMessageType } from '@/modules/material/messages'

export function highlightPassage(
  tabId: number,
  frameId: number,
  index: number,
): Promise<unknown> {
  return chrome.tabs
    .sendMessage(tabId, { type: highlightMessageType, index }, { frameId })
    .catch(() => undefined)
}
