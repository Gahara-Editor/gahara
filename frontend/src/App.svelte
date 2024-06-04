<script lang="ts">
  import VideoLayout from "./VideoLayout.svelte";
  import MainMenuLayout from "./MainMenuLayout.svelte";
  import { exportOptionsStore, router } from "./stores";
  import ExportMenuLayout from "./ExportMenuLayout.svelte";
  import { EventsOff, EventsOn } from "../wailsjs/runtime/runtime";
  import { onDestroy } from "svelte";
  const { route, setRoute } = router;
  const { isProcessingVid } = exportOptionsStore;

  EventsOn("evt_change_route", (to: string) => {
    switch ($route) {
      case "video":
        setRoute(to);
        break;
      case "export":
        if (!$isProcessingVid) setRoute(to);
    }
  });

  onDestroy(() => {
    EventsOff("evt_change_route");
  });
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
