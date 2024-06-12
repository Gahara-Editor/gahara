<script lang="ts">
  import VideoLayout from "./VideoLayout.svelte";
  import MainMenuLayout from "./MainMenuLayout.svelte";
  import {
    exportOptionsStore,
    router,
    videoFiles,
    videoStore,
    projectName,
  } from "./stores";
  import ExportMenuLayout from "./ExportMenuLayout.svelte";
  import {
    EventsOff,
    EventsOn,
    WindowSetTitle,
  } from "../wailsjs/runtime/runtime";
  import { ResetTimeline } from "../wailsjs/go/main/App";
  import { onDestroy } from "svelte";
  const { route, setRoute } = router;
  const { isProcessingVid } = exportOptionsStore;
  const { resetVideo } = videoStore;
  const { resetVideoFiles } = videoFiles;

  function cleanupProjectWorkspace() {
    resetVideoFiles();
    resetVideo();
    ResetTimeline();
    WindowSetTitle("Gahara");
  }

  EventsOn("evt_change_route", (to: string) => {
    switch ($route) {
      case "video":
        switch (to) {
          case "main":
            cleanupProjectWorkspace();
            break;
          case "export":
            WindowSetTitle(`Export - ${$projectName}`);
        }
        setRoute(to);
        break;
      case "export":
        if (!$isProcessingVid) {
          WindowSetTitle($projectName);
          setRoute(to);
        }
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
