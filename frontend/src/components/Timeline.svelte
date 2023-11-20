<script lang="ts">
  import TimelineArrow from "../icons/TimelineArrow.svelte";

  export let videoSrc: string;
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
  let tracks = [videoSrc];

  function handleBarProgress(
    e: MouseEvent & {
      currentTarget: EventTarget & HTMLDivElement;
    },
  ) {}
</script>

<div
  class="h-full w-full bg-gdark border-t-2 border-t-white flex flex-col gap-4 pt-4 pb-4 relative"
  on:mousemove={(e) => (hoverPos = e.clientX)}
>
  <div class="absolute top-0 left-0" style={`left: ${hoverPos}px`}>
    <TimelineArrow />
  </div>
  <!-- VIDEO TRACK -->
  <div class="bg-red-500 h-28">vide</div>
  <div class="bg-green-500 h-28">vide</div>
</div>
