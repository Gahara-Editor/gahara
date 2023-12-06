<script lang="ts">
  import VideoLayout from "./components/VideoLayout.svelte";
  import MainMenuLayout from "./components/MainMenuLayout.svelte";
  import { router } from "./stores";
</script>

<div class="h-screen">
  {#if $router}
    <VideoLayout />
  {:else}
    <MainMenuLayout />
  {/if}
</div>
