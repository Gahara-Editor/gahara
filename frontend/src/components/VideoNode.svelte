<script lang="ts">
  import { formatSecondsToHMS } from "../lib/utils";
  import { toolingStore } from "../stores";
  import type { video } from "wailsjs/go/models";

  const { editMode, timelineNode, boxLeftBound, trackZoom, cursorColor } =
    toolingStore;

  export let cutRangeBox: HTMLDivElement;
  export let node: video.VideoNode;
  export let pos: number;
  export let handleEditModeMouseMove: (
    e: MouseEvent & {
      currentTarget: EventTarget & HTMLDivElement;
    },
    pos: number,
    node: video.VideoNode,
  ) => void;
  export let handleEditModeMouseDown: (
    e: MouseEvent & {
      currentTarget: EventTarget & HTMLDivElement;
    },
  ) => void;
  export let handleBoxRender: (
    e: MouseEvent & {
      currentTarget: EventTarget & HTMLDivElement;
    },
    pos: number,
    video: video.VideoNode,
  ) => void;
</script>

{#if $editMode === "intervalCut" && $timelineNode && $timelineNode.id === node.id}
  <div
    class="absolute border-yellow-500 border-2 h-24 cursor-grab"
    style={`width: ${
      (node.end - node.start) * $trackZoom < 120
        ? 120
        : (node.end - node.start) * $trackZoom
    }px; left: ${$boxLeftBound}px`}
    bind:this={cutRangeBox}
    id="cut-range"
    on:mousemove={(e) => {
      handleEditModeMouseMove(e, pos, node);
    }}
    on:mousedown={(e) => {
      handleEditModeMouseDown(e);
    }}
  ></div>
{/if}
<div
  id={`videoNode-${node.id}`}
  class="h-full bg-obsbg border-white border-2 rounded-md cursor-pointer select-none overflow-hidden flex flex-col justify-start"
  style={`width: ${
    (node.end - node.start) * $trackZoom < 120
      ? 120
      : (node.end - node.start) * $trackZoom
  }px; border-color: ${
    $timelineNode && $timelineNode.id === node.id ? $cursorColor : "#ffffff"
  }; `}
  on:click={(e) => handleBoxRender(e, pos, node)}
  on:mousemove={(e) => {
    handleEditModeMouseMove(e, pos, node);
  }}
  on:mousedown={(e) => {
    handleEditModeMouseDown(e);
  }}
>
  <p class="text-sm font-semibold pl-1 pt-1">
    {node.name}
  </p>
  <p class="text-sm pl-1">
    {formatSecondsToHMS(node.end - node.start)}
  </p>
  {#if node.losslessexport}
    <span class="w-6 font-bold text-lg text-center text-gyellow"> M </span>
  {/if}
</div>
