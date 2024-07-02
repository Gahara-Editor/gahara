<script lang="ts">
  import { searchListstore } from "../stores";
  import type { main } from "../../wailsjs/go/models";
  import { scrollVertical } from "../lib/searchlist";
  export let idx: number;
  export let item: main.Video;
  export let selected: boolean;

  let el: HTMLLIElement = null;
  const { setSearchIdx } = searchListstore;
  $: if (selected && el) {
    scrollVertical(el);
  }
</script>

{#if selected}
  <li class=" bg-obsternary p-2 truncate rounded-md" bind:this={el}>
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
