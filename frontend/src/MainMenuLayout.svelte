<script lang="ts">
  import {
    CreateProjectWorkspace,
    ReadGaharaWorkspace,
    SetProjectDirectory,
  } from "../wailsjs/go/main/App";
  import { ChevronDownIcon } from "@rgossiaux/svelte-heroicons/solid";
  import { router, projectName } from "./stores";
  import { WindowSetTitle } from "../wailsjs/runtime/runtime";

  let createProjectView: boolean;
  let loadProjectView: boolean;
  let mainMenuError: string = "";
  let projects: string[];

  function createProject() {
    CreateProjectWorkspace($projectName)
      .then(() => {
        mainMenuError = "";
        WindowSetTitle($projectName);
        router.setVideoLayoutView();
      })
      .catch((msg) => (mainMenuError = msg));
  }

  function loadProjects() {
    ReadGaharaWorkspace()
      .then((res) => (projects = res))
      .catch(console.log);
  }
  loadProjects();

  function handleProjectSelected() {
    SetProjectDirectory($projectName).then(() => {
      WindowSetTitle($projectName);
      router.setVideoLayoutView();
    });
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
        class="rounded-lg bg-indigo-500 font-semibold text-white inline-flex items-center px-4 py-1.5 hover:bg-indigo-700 transition ease-in-out duration-200"
        on:click={() => {
          $projectName = "";
          CreateProjectView();
        }}>Create Project</button
      >
    {:else}
      <button
        class="rounded-lg bg-red-500 font-semibold text-white inline-flex items-center px-4 py-1.5 hover:bg-red-700 transition ease-in-out duration-200"
        on:click={() => MainMenuView()}>Back</button
      >
    {/if}
    {#if !loadProjectView}
      <button
        class="rounded-lg bg-purple-500 font-semibold text-white inline-flex items-center px-4 py-1.5 hover:bg-purple-700 transition ease-in-out duration-200"
        on:click={() => LoadProjectView()}>Load Project</button
      >
    {:else}
      <button
        class="rounded-lg bg-red-500 font-semibold text-white inline-flex items-center px-4 py-1.5 hover:bg-red-700 transition ease-in-out duration-200"
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
    <div class="flex flex-col gap-1">
      {#if projects}
        <div class="flex gap-2">
          <div class="relative inline-flex">
            <select
              id="projectNames"
              bind:value={$projectName}
              class="block appearance-none w-full bg-white border border-indigo-500 hover:border-gray-500 px-4 py-2 pr-8 rounded leading-tight focus:outline-none focus:bg-white focus:border-indigo-600"
            >
              {#each projects as project (project)}
                <option value={project}>
                  {project.length > 25
                    ? `${project.substring(0, 25)}...`
                    : project}
                </option>
              {/each}
            </select>
            <div
              class="pointer-events-none absolute inset-y-0 right-0 flex items-center px-2 text-indigo-500"
            >
              <ChevronDownIcon class="h-6 w-6" />
            </div>
          </div>
          <button
            disabled={$projectName == ""}
            class="rounded-full bg-indigo-500 font-semibold text-white px-4 py-1.5 hover:bg-indigo-700 transition ease-in-out duration-200"
            on:click={() => handleProjectSelected()}>Load</button
          >
        </div>
        <div class="text-white 2xl">project previes</div>
      {:else}
        <p class="text-white text-lg">No projects found</p>
      {/if}
    </div>
  {/if}
  {#if mainMenuError}
    <p class="rounded-md bg-gred1 text-white text-md p-2">{mainMenuError}</p>
  {/if}
</div>
