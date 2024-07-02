import type { main, video } from "../../wailsjs/go/models";

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

export type ListType = main.Video | video.VideoNode;

export function isVideoNode(unit: ListType): unit is video.VideoNode {
  return (unit as video.VideoNode).losslessexport !== undefined;
}

export function isVideo(unit: ListType): unit is main.Video {
  return (unit as main.Video).duration !== undefined;
}
