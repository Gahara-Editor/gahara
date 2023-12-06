<script lang="ts">
  import TimelineArrow from "../icons/TimelineArrow.svelte";
  import { dropzone } from "../lib/dnd";
  import { trackFiles, selectedTrack } from "../stores";
  import { onDestroy } from "svelte";
  import type { main } from "wailsjs/go/models";
  import { GenerateThumbnail } from "../../wailsjs/go/main/App";
  import { videoStore } from "../stores";

  let hoverPos = 0;
  let moveTimeline = false;
  let timestamp = new Date().getTime();
  let trackNode: HTMLDivElement;

  const { getDuration, currentTime } = videoStore;

  $: {
    if (trackNode) {
      hoverPos = ($currentTime / getDuration()) * trackNode.clientWidth;
    }
  }

  function viewVideo(video: main.Video) {
    selectedTrack.set(`${video.filepath}/${video.name}${video.extension}`);
  }

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

  function moveTimelineBar() {
    moveTimeline = true;
  }

  function stopTimelineBar() {
    moveTimeline = false;
  }

  function handleTimelineMove(
    e: MouseEvent & {
      currentTarget: EventTarget & HTMLDivElement;
    },
  ) {
    if (moveTimeline) {
      hoverPos = Math.min(e.clientX, calculateMaxTrackWidth());
      currentTime.set((hoverPos / trackNode.clientWidth) * getDuration());
    }
  }

  onDestroy(() => {
    trackFiles.reset();
    hoverPos = 0;
  });
</script>

<div
  class="timeline h-full w-full bg-gdark border-t-2 border-t-white flex flex-col gap-4 pt-4 pb-4 px-1 relative"
  use:dropzone={{}}
  on:mousedown={() => moveTimelineBar()}
  on:mouseup={() => stopTimelineBar()}
  on:mousemove={(e) => handleTimelineMove(e)}
>
  {#if $trackFiles.length <= 0}
    <div class="flex justify-center items-center">
      <p class="text-white text-4xl font-semibold select-none">
        Drag And Drop Video Clips
      </p>
    </div>
  {:else}
    <div class="absolute top-0 left-0 h-full w-3" style={`left: ${hoverPos}px`}>
      <TimelineArrow />
    </div>
  {/if}

  <!-- VIDEO TRACKS -->
  {#each $trackFiles as track (track.filepath + track.name)}
    <div
      class="w-1/4 h-28 bg-teal rounded-md border-white border-2 video-track cursor-pointer"
      on:click={() => viewVideo(track)}
      bind:this={trackNode}
    >
      <img
        src={loadThumbnail(track) + `?${timestamp}`}
        alt={`video: ${track.name}`}
        class="h-full w-full select-none"
        draggable={false}
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
