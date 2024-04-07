<script lang="ts">
  import { router, projectName, mainMenuStore } from "../stores";
  import {
    GetProjectThumbnail,
    SetProjectDirectory,
    DeleteProject,
  } from "../../wailsjs/go/main/App";
  import { WindowSetTitle } from "../../wailsjs/runtime/runtime";
  import { onMount } from "svelte";
  import GaharaIcon from "../icons/GaharaIcon.svelte";
  import WarningIcon from "../icons/WarningIcon.svelte";
  import Modal from "../components/Modal.svelte";
  import TrashIcon from "../icons/TrashIcon.svelte";
  export let project = "";
  export let idx = 0;

  const { setRoute } = router;
  const { removeProject } = mainMenuStore;
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

  function deleteProject() {
    DeleteProject(project)
      .then(() => {
        removeProject(project);
      })
      .catch(console.log);
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
      <div class="flex w-full items-center justify-center">
        <GaharaIcon {idx} />
      </div>
    {/if}
  </div>
  <div class="flex-col justify-center items-center py-2">
    <h1 class="font-bold text-lg text-center py-2">{project}</h1>
    <div class="flex items-center justify-center py-2 gap-2">
      <button
        class="rounded-lg bg-indigo-500 font-semibold text-white px-4 py-1.5 hover:bg-indigo-700 transition ease-in-out duration-200 border-2 border-white"
        on:click={() => handleProjectSelected()}
      >
        Load
      </button>
      <Modal>
        <div slot="trigger" let:open>
          <button
            class="rounded-lg bg-gred1 font-semibold text-white px-4 py-1.5 hover:bg-gred transition ease-in-out duration-200 border-2 border-white"
            on:click={open}
          >
            Delete
          </button>
        </div>
        <div
          slot="header"
          class="font-semibold flex items-center justify-center gap-2"
        >
          <WarningIcon class="h-5 w-5 text-gyellow" />
          <h1>Delete File</h1>
        </div>
        <div slot="content">
          <p class="text-center font-semibold">
            Do you want to delete {project}?
          </p>
          <p>All of the files in this project will be permanently deleted</p>
        </div>
        <div
          slot="footer"
          let:store={{ close }}
          class="flex items-center justify-center gap-2"
        >
          <button
            class="flex items-center justify-center rounded-lg bg-gred1 font-semibold text-white px-4 py-1.5 hover:bg-gred transition ease-in-out duration-200 border-2 border-white gap-2"
            on:click={close}>Back</button
          >
          <button
            class="flex items-center justify-center rounded-lg bg-gblue0 font-semibold text-white px-4 py-1.5 hover:bg-gblue transition ease-in-out duration-200 border-2 border-white gap-2"
            on:click={() => deleteProject()}
          >
            <TrashIcon class="h-5 w-5 text-white" />
            <span>Delete Project</span>
          </button>
        </div>
      </Modal>
    </div>
  </div>
</div>
