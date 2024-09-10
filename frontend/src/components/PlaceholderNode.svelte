<script lang="ts">
  import type { TimelineNode } from "src/lib/timeline";
  import { formatSecondsToHMS } from "../lib/utils";
  import { toolingStore } from "../stores";
  import type { placeholder } from "wailsjs/go/models";

  const { editMode, timelineNode, boxLeftBound, trackZoom, cursorColor } =
    toolingStore;

  export let handleBoxRender: (
    e: MouseEvent & {
      currentTarget: EventTarget & HTMLDivElement;
    },
    pos: number,
    node: TimelineNode,
  ) => void;

  export let cutRangeBox: HTMLDivElement;
  export let node: placeholder.PlaceholderNode;
  export let pos: number;
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
  ></div>
{/if}
<div
  id={`node-${node.id}`}
  class="h-full bg-obsbg border-white border-2 rounded-md cursor-pointer select-none overflow-hidden flex flex-col justify-start"
  style={`width: ${
    (node.end - node.start) * $trackZoom < 120
      ? 120
      : (node.end - node.start) * $trackZoom
  }px; border-color: ${
    $timelineNode && $timelineNode.id === node.id ? $cursorColor : "#ffffff"
  }; `}
  on:click={(e) => handleBoxRender(e, pos, node)}
>
  <p class="text-sm font-semibold pl-1 pt-1">
    {node.name}
  </p>
  <p class="text-sm pl-1">
    {formatSecondsToHMS(node.end - node.start)}
  </p>
</div>
