import type { audio, placeholder, video } from "../../wailsjs/go/models";
import { toolingStore } from "../stores";
import { NODE_AUDIO, NODE_PLACEHOLDER, NODE_VIDEO } from "./utils";

export type TimelineNode =
  | video.VideoNode
  | audio.AudioNode
  | placeholder.PlaceholderNode;

export function isNodeVideo(unit: TimelineNode): unit is video.VideoNode {
  return unit && (unit as video.VideoNode).type === NODE_VIDEO;
}

export function isNodeAudio(unit: TimelineNode): unit is audio.AudioNode {
  return unit && (unit as audio.AudioNode).type === NODE_AUDIO;
}

export function isNodePlaceholder(
  unit: TimelineNode,
): unit is placeholder.PlaceholderNode {
  return (
    unit && (unit as placeholder.PlaceholderNode).type === NODE_PLACEHOLDER
  );
}

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

export function scrollToNode(node: HTMLDivElement) {
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

export function scrollToTrack(node: HTMLElement) {
  const timelineContainer = document.getElementById("timeline");
  const timelineRect = timelineContainer.getBoundingClientRect();
  const nodeRect = node.getBoundingClientRect();

  const isNodeVisible =
    nodeRect.top >= timelineRect.top && nodeRect.bottom <= timelineRect.bottom;

  if (!isNodeVisible) {
    const scrollY =
      nodeRect.top - timelineRect.top + timelineContainer.scrollTop;
    timelineContainer.scrollTo({
      top: scrollY,
      behavior: "smooth",
    });
  }
}
