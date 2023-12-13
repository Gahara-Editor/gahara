<script lang="ts">
  import VideoLayout from "./VideoLayout.svelte";
  import MainMenuLayout from "./MainMenuLayout.svelte";
  import { router } from "./stores";
</script>

<div class="h-screen">
  {#if $router}
    <VideoLayout />
  {:else}
    <MainMenuLayout />
  {/if}
</div>
