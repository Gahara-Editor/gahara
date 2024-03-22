<script lang="ts">
  import VideoLayout from "./VideoLayout.svelte";
  import MainMenuLayout from "./MainMenuLayout.svelte";
  import { router } from "./stores";

  const { route } = router;
</script>

<div class="h-screen">
  {#if $route === "video"}
    <VideoLayout />
  {:else if $route === "export"}
    <div></div>
  {:else}
    <MainMenuLayout />
  {/if}
</div>
