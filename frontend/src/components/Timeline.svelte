<script lang="ts">
  import TimelineArrow from "../icons/TimelineArrow.svelte";
  import { dropzone } from "../lib/dnd";
  import { trackFiles } from "../stores";
  import { onDestroy } from "svelte";
  import type { main } from "wailsjs/go/models";
  import { GenerateThumbnail } from "../../wailsjs/go/main/App";

  let hoverPos = 0;
  let timestamp = new Date().getTime();
  let duration: number;
  let currentTime: number;
  let playbackRate: number;
  let volume: number;
  let paused: boolean;
  let ended: boolean;
  let muted: boolean;
  let seeking: boolean;
  let buffered: TimeRanges;
  let played: TimeRanges;
  let seekable: TimeRanges;
  let videoWidth: number;
  let videoHeight: number;

  function loadThumbnail(track: main.Video) {
    return `${track.filepath}/${track.name}.png`;
  }

  function createThumbnail(track: main.Video) {
    GenerateThumbnail(`${track.filepath}/${track.name}${track.extension}`)
      .then(() => (timestamp = new Date().getTime()))
      .catch((e) => console.log(e));
  }

  function calculateMaxTrackWidth() {
    const trackElements = document.querySelectorAll(".video-track");
    let maxWidth = 0;

    trackElements.forEach((trackElement) => {
      const trackWidth = trackElement.clientWidth;
      if (trackWidth > maxWidth) {
        maxWidth = trackWidth;
      }
    });

    return maxWidth;
  }

  onDestroy(() => {
    trackFiles.reset();
    hoverPos = 0;
  });
</script>

<div
  class="timeline h-full w-full bg-gdark border-t-2 border-t-white flex flex-col gap-4 pt-4 pb-4 relative"
  on:mousemove={(e) =>
    (hoverPos = Math.min(e.clientX, calculateMaxTrackWidth()))}
  use:dropzone={{}}
>
  <div
    class="absolute top-0 left-0 cursor-pointer"
    style={`left: ${hoverPos}px`}
  >
    <TimelineArrow />
  </div>
  <!-- VIDEO TRACKS -->
  {#each $trackFiles as track (track.filepath + track.name)}
    <div
      class="w-1/4 h-28 bg-teal rounded-md border-gblue border-2 video-track"
    >
      <img
        src={loadThumbnail(track) + `?${timestamp}`}
        alt={`video: ${track.name}`}
        class="h-full w-full"
        on:error={() => createThumbnail(track)}
      />
    </div>
  {/each}
</div>

<style>
  .timeline:global(.droppable) {
    border-width: 2px;
    border-color: rgb(122, 162, 247);
  }

  .timeline:global(.droppable) * {
    pointer-events: none;
  }
</style>
