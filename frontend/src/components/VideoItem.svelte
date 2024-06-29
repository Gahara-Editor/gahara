<script lang="ts">
  import { searchListstore } from "../stores";
  import type { main } from "../../wailsjs/go/models";
  export let idx: number;
  export let item: main.Video;
  export let selected: boolean;

  const { setSearchIdx } = searchListstore;
</script>

{#if selected}
  <li class=" bg-obsternary p-2 truncate rounded-md">
    <span class="text-teal font-semibold">></span>
    <span class="font-semibold text-gyellow">[CLIP]</span>
    {item.name}
  </li>
{:else}
  <li
    class=" bg-obsbg p-2 truncate rounded-md"
    on:click={() => setSearchIdx(idx)}
  >
    <span class="font-semibold text-gyellow">[CLIP]</span>
    {item.name}
  </li>
{/if}
