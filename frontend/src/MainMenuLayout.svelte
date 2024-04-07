<script lang="ts">
  import {
    CreateProjectWorkspace,
    ReadGaharaWorkspace,
  } from "../wailsjs/go/main/App";
  import {
    ArrowLeftIcon,
    ArrowRightIcon,
  } from "@rgossiaux/svelte-heroicons/solid";
  import { router, projectName, mainMenuStore } from "./stores";
  import { WindowSetTitle } from "../wailsjs/runtime/runtime";
  import ProjectCard from "./components/ProjectCard.svelte";
  import { onDestroy, onMount } from "svelte";
  import { flip } from "svelte/animate";
  import { slide } from "svelte/transition";

  const { setRoute } = router;
  const {
    projects,
    carouselIdx,
    setProjects,
    updateCarouselIdx,
    resetMainMenuStore,
  } = mainMenuStore;

  let createProjectView: boolean;
  let loadProjectView: boolean;
  let mainMenuError: string = "";

  onMount(() => {
    loadProjects();
  });

  function createProject() {
    CreateProjectWorkspace($projectName)
      .then(() => {
        mainMenuError = "";
        WindowSetTitle($projectName);
        setRoute("video");
      })
      .catch((msg) => (mainMenuError = msg));
  }

  function loadProjects() {
    ReadGaharaWorkspace()
      .then((res) => setProjects(res))
      .catch(console.log);
  }

  function handleNext() {
    if ($carouselIdx < $projects.length - 1) {
      updateCarouselIdx((idx) => idx + 1);
    }
  }

  function handlePrevious() {
    if ($carouselIdx > 0) {
      updateCarouselIdx((idx) => idx - 1);
    }
  }
  function MainMenuView() {
    createProjectView = false;
    loadProjectView = false;
    mainMenuError = "";
  }

  function CreateProjectView() {
    createProjectView = true;
    loadProjectView = false;
    mainMenuError = "";
  }

  function LoadProjectView() {
    createProjectView = false;
    loadProjectView = true;
    mainMenuError = "";
  }

  onDestroy(() => {
    resetMainMenuStore();
  });
</script>

<div class="h-full flex flex-col items-center justify-center gap-4">
  <div class="flex items-center gap-2 bg-indigo-500 rounded-full">
    <div class="rounded-full bg-indigo-600 p-2">
      <img src="images/iconmin.png" alt="Gahara Logo" class="w-20 p-1" />
    </div>
    <h1 class="text-6xl text-white font-bold p-2">Gahara</h1>
  </div>
  <div class="flex gap-4">
    {#if !createProjectView}
      <button
        class="rounded-lg bg-indigo-500 border-white border-2 font-semibold text-white inline-flex items-center px-4 py-1.5 hover:bg-indigo-700 transition ease-in-out duration-200"
        on:click={() => {
          $projectName = "";
          CreateProjectView();
        }}>Create Project</button
      >
    {:else}
      <button
        class="rounded-lg bg-red-500 border-white border-2 font-semibold text-white inline-flex items-center px-4 py-1.5 hover:bg-red-700 transition ease-in-out duration-200"
        on:click={() => MainMenuView()}>Back</button
      >
    {/if}
    {#if !loadProjectView}
      <button
        class="rounded-lg bg-purple-500 border-white border-2 font-semibold text-white inline-flex items-center px-4 py-1.5 hover:bg-purple-700 transition ease-in-out duration-200"
        on:click={() => LoadProjectView()}>Load Project</button
      >
    {:else}
      <button
        class="rounded-lg bg-red-500 border-white border-2 font-semibold text-white inline-flex items-center px-4 py-1.5 hover:bg-red-700 transition ease-in-out duration-200"
        on:click={() => MainMenuView()}>Back</button
      >
    {/if}
  </div>
  {#if createProjectView}
    <div class="flex flex-col gap-3">
      <label for="project-name" class=" flex self-center text-white font-bold"
        >Project Name</label
      >
      <input
        type="text"
        autocorrect="off"
        autocomplete="off"
        bind:value={$projectName}
        id="project-name"
        class="border-solid border-indigo-300 border-2 rounded-md focus:border-indigo-500 focus:outline-none"
      />
      <button
        class="rounded-full bg-indigo-500 font-semibold text-white px-4 py-1.5 hover:bg-indigo-700 transition ease-in-out duration-200 self-center"
        on:click={() => createProject()}>Start</button
      >
    </div>
  {/if}
  {#if loadProjectView}
    <div class="flex flex-col items-center justify-center">
      {#if $projects.length > 0}
        <div class="flex items-center justify-center gap-4">
          <button
            on:click={handlePrevious}
            disabled={$carouselIdx === 0}
            class={$carouselIdx === 0 ? "text-gray-400" : "text-white"}
          >
            <ArrowLeftIcon class="h-6 w-6" />
          </button>
          <div class="flex items-center justify-center p-2 gap-2">
            {#each $projects as project, idx (project)}
              <div animate:flip={{ duration: 400 }} in:slide>
                {#if idx >= $carouselIdx && idx < $carouselIdx + 3}
                  <ProjectCard {project} {idx} />
                {/if}
              </div>
            {/each}
          </div>
          <button
            on:click={handleNext}
            disabled={$carouselIdx >= $projects.length - 3}
            class={$carouselIdx >= $projects.length - 3
              ? "text-gray-400"
              : "text-white"}
          >
            <ArrowRightIcon class="h-6 w-6" />
          </button>
        </div>
      {:else}
        <p class="text-white text-lg">No projects found</p>
      {/if}
    </div>
  {/if}
  {#if mainMenuError}
    <p class="rounded-md bg-gred1 text-white text-md p-2">{mainMenuError}</p>
  {/if}
</div>
