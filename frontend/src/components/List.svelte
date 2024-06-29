<script lang="ts">
  import { isVideo, isVideoNode, type ListType } from "../lib/utils";
  import { searchListstore } from "../stores";
  import VideoNodeItem from "./VideoNodeItem.svelte";
  import VideoItem from "./VideoItem.svelte";
  const { activeList, searchIdx } = searchListstore;
</script>

<div
  id="content-wrap"
  class="z-10 bg-gblue0 min-h-[40vh] max-h-[40vh] min-w-[60vw] max-w-[60vw] rounded border-white border-2 p-2 overflow-y-scroll"
>
  <div id="video-clip-list">
    <h1 class="text-center font-semibold text-lg">Results</h1>
    <ul class="divide-y divide-gray-200 text-white">
      {#each $activeList as item, idx}
        {#if isVideo(item)}
          <VideoItem {item} {idx} selected={idx === $searchIdx} />
        {:else if isVideoNode(item)}
          <VideoNodeItem {item} {idx} selected={idx === $searchIdx} />
        {/if}
      {/each}
    </ul>
  </div>
</div>
