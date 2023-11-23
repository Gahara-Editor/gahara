<script lang="ts">
  import TimelineArrow from "../icons/TimelineArrow.svelte";
  import { dropzone } from "../lib/dnd";
  import { trackFiles } from "../stores";
  import { onDestroy } from "svelte";

  let hoverPos = 0;
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

  onDestroy(() => {
    trackFiles.reset();
    hoverPos = 0;
  });
</script>

<div
  class="timeline h-full w-full bg-gdark border-t-2 border-t-white flex flex-col gap-4 pt-4 pb-4 relative"
  on:mousemove={(e) => (hoverPos = e.clientX)}
  use:dropzone={{}}
>
  <div
    class="absolute top-0 left-0 cursor-pointer"
    style={`left: ${hoverPos}px`}
  >
    <TimelineArrow />
  </div>
  <!-- VIDEO TRACKS -->
  {#each $trackFiles as track (track.filepath)}
    <div class="w-full h-28 bg-gred1">
      {track.name}
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
