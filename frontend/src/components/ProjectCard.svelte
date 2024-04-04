<script lang="ts">
  import { router, projectName } from "../stores";
  import {
    GetProjectThumbnail,
    SetProjectDirectory,
  } from "../../wailsjs/go/main/App";
  import { WindowSetTitle } from "../../wailsjs/runtime/runtime";
  import { onMount } from "svelte";
  import GaharaIcon from "../icons/GaharaIcon.svelte";
  export let project = "";
  export let idx = 0;

  const { setRoute } = router;
  let imgpath = "";

  onMount(() => {
    GetProjectThumbnail(project)
      .then((path) => {
        imgpath = path;
      })
      .catch((_) => (imgpath = ""));
  });

  function handleProjectSelected() {
    $projectName = project;
    SetProjectDirectory($projectName).then(() => {
      WindowSetTitle($projectName);
      setRoute("video");
    });
  }
</script>

<div
  class="bg-obsbg rounded w-80 overflow-hidden border-white border-2 text-white"
>
  <div
    class="relative h-48 w-full flex items-center justify-center bg-gradient-to-r from-indigo-800 via-indigo-400 to-indigo-800 border-b-white border-b-2"
  >
    {#if imgpath}
      <img src={imgpath} alt={project} draggable={false} />
    {:else}
      <div class="flex w-full items-center justify-center z-10">
        <GaharaIcon {idx} />
      </div>
    {/if}
  </div>
  <div class="flex-col justify-center items-center py-2">
    <h1 class="font-bold text-lg text-center py-2">{project}</h1>
    <div class="flex items-center justify-center py-2">
      <button
        class="rounded-full bg-indigo-500 font-semibold text-white px-4 py-1.5 hover:bg-indigo-700 transition ease-in-out duration-200"
        on:click={() => handleProjectSelected()}
      >
        Load
      </button>
    </div>
  </div>
</div>
