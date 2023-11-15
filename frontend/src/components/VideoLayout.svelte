<script lang="ts">
  import { FilePicker } from "../../wailsjs/go/main/App";
  import {
    XIcon,
    PlusCircleIcon,
    FilmIcon,
  } from "@rgossiaux/svelte-heroicons/solid";
  import { videoFiles, projectName } from "../stores";
  import { ReadProjectWorkspace } from "../../wailsjs/go/main/App";
  import { video } from "../stores";
  import { onDestroy } from "svelte";
  import Modal from "./Modal.svelte";
  import VideoPlayer from "./VideoPlayer.svelte";

  let fileUploadError = "";
  let videoSrc = "";
  let thumbnail = "";

  function loadProjectFiles() {
    ReadProjectWorkspace()
      .then((files) => videoFiles.addVideos(files))
      .catch(() => (fileUploadError = "No files in this project"));
  }
  loadProjectFiles();

  function selectFile() {
    FilePicker()
      .then((video) => {
        fileUploadError = "";
        videoFiles.addVideos([video]);
      })
      .catch(() => (fileUploadError = "No files selected"));
  }

  function viewVideo(path: string) {
    videoSrc = path;
  }

  onDestroy(() => {
    videoFiles.reset();
  });
</script>

<div class="h-full text-white rounded-md bg-red-500">
  <div class="flex items-start gap-2 p-4">
    <!-- Video Uploading -->
    <div
      class="flex flex-col text-white rounded-md bg-gblue0 border-white border-2 p-2 gap-2 w-1/4 h-fit"
    >
      <div class="flex items-center gap-2">
        <Modal>
          <div slot="trigger" let:open>
            <button
              class="bg-gdark hover:bg-gprimary rounded-lg px-4 py-2 border-2 border-white"
              on:click={open}
            >
              <FilmIcon class="h-6 w-6 text-white" />
            </button>
          </div>
          <div slot="header" class="font-semibold">
            <h1>Switch Project</h1>
          </div>
          <div slot="content">
            <p>Do you want to switch projects?</p>
          </div>
          <div slot="footer" let:store={{ close }} class="flex gap-2">
            <button
              class="rounded-lg bg-red-500 font-semibold text-white inline-flex items-center px-2 py-0.5 hover:bg-red-700 transition ease-in-out duration-200 border-2 border-white"
              on:click={close}>Back</button
            >
            <button
              class="bg-gdark rounded-lg px-2 py-0.5 border-2 border-white flex items-center gap-2 hover:bg-gblue0 transition ease-in-out duration-200"
              on:click={() => video.setMainMenuView()}
            >
              <FilmIcon class="h-6 w-6 text-white" />
              <span>Switch Project</span>
            </button>
          </div>
        </Modal>
        <h2 class="font-semibold text-md">{$projectName}</h2>
      </div>
      <div>
        <button
          class="bg-teal hover:bg-green2 px-2 py-1 rounded-full flex items-center gap-1 border-2 border-white transition ease-in-out duration-500"
          on:click={() => selectFile()}
        >
          <PlusCircleIcon class="h-5 w-5 text-white" />
          <span class="text-white font-medium">Video</span>
        </button>
      </div>
      {#if videoFiles}
        <div class="flex flex-col gap-2">
          {#each $videoFiles as video (video.name)}
            <div
              class="flex items-center bg-gprimary hover:bg-stone-700 rounded-lg p-2 cursor-grab transition ease-in-out duration-500 gap-2"
              on:click={() => viewVideo(video.filepath)}
            >
              <button
                class="bg-red-500 hover:bg-red-400 px-1 py-1 rounded-full"
                on:click={() => videoFiles.removeVideo(video.name)}
              >
                <XIcon class="h-3 w-3 text-white" />
              </button>
              <p class="text-sm">{video.name}</p>
            </div>
          {/each}
        </div>
      {/if}
    </div>
    {#if fileUploadError}
      {fileUploadError}
    {/if}
    <!-- Video Player -->
    <div id="video-player" class="h-[36rem] w-3/4">
      <VideoPlayer {videoSrc} {thumbnail} />
    </div>
  </div>
  <!-- Timeline -->
  <div>Timeline</div>
</div>
