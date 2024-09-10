import type { placeholder, audio, main, video } from "../../wailsjs/go/models";

export const NODE_VIDEO = "NODE_VIDEO";
export const NODE_AUDIO = "NODE_AUDIO";
export const NODE_PLACEHOLDER = "NODE_PLACEHOLDER";
export const NODE_TEXT = "NODE_TEXT";

export type ListType =
  | main.Video
  | video.VideoNode
  | audio.AudioNode
  | placeholder.PlaceholderNode;

export function isVideoItem(unit: ListType): unit is video.VideoNode {
  return unit && (unit as video.VideoNode).losslessexport !== undefined;
}

export function isVideo(unit: ListType): unit is main.Video {
  return unit && (unit as main.Video).duration !== undefined;
}

export function isPlaceholderItem(
  unit: ListType,
): unit is placeholder.PlaceholderNode {
  return (
    unit && (unit as placeholder.PlaceholderNode).type === NODE_PLACEHOLDER
  );
}

export const placeholderItem: placeholder.PlaceholderNode = {
  type: NODE_PLACEHOLDER,
  id: "",
  rid: "",
  start: 0,
  end: 300,
  name: "placeholder",
};

export function formatSecondsToHMS(seconds: number): string {
  const hours = Math.floor(seconds / 3600);
  const minutes = Math.floor((seconds % 3600) / 60);
  const remainingSeconds = Math.floor(seconds % 60);

  const formattedHours = hours < 10 ? `0${hours}` : `${hours}`;
  const formattedMinutes = minutes < 10 ? `0${minutes}` : `${minutes}`;
  const formattedSeconds =
    remainingSeconds < 10 ? `0${remainingSeconds}` : `${remainingSeconds}`;

  return `${formattedHours}:${formattedMinutes}:${formattedSeconds}`;
}

export function scrollVertical(node: HTMLLIElement) {
  const listContainer = document.getElementById("content-wrap");
  const listRect = listContainer.getBoundingClientRect();
  const nodeRect = node.getBoundingClientRect();

  const isNodeVisible =
    nodeRect.top >= listRect.top && nodeRect.bottom <= listRect.bottom;

  if (!isNodeVisible) {
    const scrollY = nodeRect.top - listRect.top + listContainer.scrollTop;
    listContainer.scrollTo({
      top: scrollY,
      behavior: "smooth",
    });
  }
}
