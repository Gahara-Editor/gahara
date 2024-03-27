<script lang="ts">
  import VideoLayout from "./VideoLayout.svelte";
  import MainMenuLayout from "./MainMenuLayout.svelte";
  import { router } from "./stores";
  import ExportMenuLayout from "./ExportMenuLayout.svelte";

  const { route } = router;
</script>

<div class="h-screen">
  {#if $route === "video"}
    <VideoLayout />
  {:else if $route === "export"}
    <ExportMenuLayout />
  {:else}
    <MainMenuLayout />
  {/if}
</div>
