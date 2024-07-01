import { toolingStore } from "../stores";

export function handleKeybindTrackClipMove() {
  const videoNodeDiv = document
    .getElementById(`track-0`)
    ?.querySelector(`div:nth-child(${toolingStore.getCursorIdx() + 1})`)
    ?.querySelector("div");
  if (videoNodeDiv) {
    videoNodeDiv.click();
    scrollToNode(videoNodeDiv);
  }
}

function scrollToNode(node: HTMLDivElement) {
  const timelineContainer = document.getElementById("timeline");
  const timelineRect = timelineContainer.getBoundingClientRect();
  const nodeRect = node.getBoundingClientRect();

  const isNodeVisible =
    nodeRect.left >= timelineRect.left && nodeRect.right <= timelineRect.right;

  if (!isNodeVisible) {
    const scrollX =
      nodeRect.left - timelineRect.left + timelineContainer.scrollLeft;
    timelineContainer.scrollTo({
      left: scrollX,
      behavior: "smooth",
    });
  }
}
