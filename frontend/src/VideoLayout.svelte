<script lang="ts">
  import {
    FilePicker,
    ReadProjectWorkspace,
    LoadTimeline,
    SaveTimeline,
    ResetTimeline,
  } from "../wailsjs/go/main/App";
  import {
    XIcon,
    FilmIcon,
    ArrowSmDownIcon,
  } from "@rgossiaux/svelte-heroicons/solid";
  import { router, videoStore, videoFiles, trackStore } from "./stores";
  import { onMount } from "svelte";
  import Modal from "./components/Modal.svelte";
  import VideoPlayer from "./components/VideoPlayer.svelte";
  import Timeline from "./components/Timeline.svelte";
  import { draggable } from "./lib/dnd";
  import { WindowSetTitle } from "../wailsjs/runtime/runtime";
  import ToolingLayout from "./ToolingLayout.svelte";
  import FloppyDisk from "./icons/FloppyDisk.svelte";
  import FolderOpenIcon from "./icons/FolderOpenIcon.svelte";

  const { setVideoSrc, resetVideo } = videoStore;
  const { addVideoToTrack, trackTime, trackDuration } = trackStore;
  const { resetVideoFiles } = videoFiles;
  const { setRoute } = router;
  let fileUploadError = "";

  function loadProjectFiles() {
    if ($videoFiles.length === 0) {
      ReadProjectWorkspace()
        .then((files) => videoFiles.addVideos(files))
        .catch(() => (fileUploadError = "No files in this project"));
    }
  }

  function loadTimeline() {
    LoadTimeline()
      .then((timeline) => {
        timeline.video_nodes.forEach((videoNode) => {
          addVideoToTrack(0, videoNode);
        });
        setVideoSrc(timeline.video_nodes[0].rid);
      })
      .catch(console.log);
  }

  function saveTimeline() {
    SaveTimeline()
      .then(() => console.log("timeline was saved"))
      .catch(() => console.log("could not save timeline"));
  }

  function formatSecondsToHMS(seconds: number): string {
    const hours = Math.floor(seconds / 3600);
    const minutes = Math.floor((seconds % 3600) / 60);
    const remainingSeconds = Math.floor(seconds % 60);

    const formattedHours = hours < 10 ? `0${hours}` : `${hours}`;
    const formattedMinutes = minutes < 10 ? `0${minutes}` : `${minutes}`;
    const formattedSeconds =
      remainingSeconds < 10 ? `0${remainingSeconds}` : `${remainingSeconds}`;

    return `${formattedHours}:${formattedMinutes}:${formattedSeconds}`;
  }

  onMount(() => {
    loadProjectFiles();
    loadTimeline();
  });

  function selectFile() {
    FilePicker()
      .then((video) => {
        fileUploadError = "";
        videoFiles.addVideos([video]);
      })
      .catch(() => (fileUploadError = "No files selected"));
  }
</script>

<div class="h-full text-white rounded-md bg-gprimary flex flex-col">
  <div class="flex items-start gap-2 p-4">
    <!-- Video Uploading -->
    <div
      class="flex flex-col text-white rounded-md bg-gblue0 border-white border-2 p-2 gap-2 w-1/4 h-[28rem]"
    >
      <div class="flex items-center gap-2">
        <Modal>
          <div slot="trigger" let:open>
            <button
              class="bg-gdark px-2 py-1 rounded-md flex items-center gap-1 border-2 border-white"
              on:click={open}
            >
              <FilmIcon class="h-5 w-5 text-white" />
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
              class="bg-gdark rounded-lg px-2 py-1 border-2 border-white flex items-center gap-2 hover:bg-gblue0 transition ease-in-out duration-200"
              on:click={() => {
                resetVideoFiles();
                resetVideo();
                ResetTimeline();
                WindowSetTitle("Gahara");
                setRoute("main");
              }}
            >
              <FilmIcon class="h-5 w-5 text-white" />
              <span>Switch Project</span>
            </button>
          </div>
        </Modal>
        <button
          class="bg-gdark px-2 py-1 rounded-md flex items-center gap-1 border-2 border-white"
          on:click={() => selectFile()}
        >
          <FolderOpenIcon class="h-5 w-5 text-white" />
        </button>
        <button
          class="bg-gdark px-2 py-1 rounded-md flex items-center gap-1 border-2 border-white"
          on:click={() => saveTimeline()}
        >
          <FloppyDisk class="h-5 w-5 text-white" />
        </button>

        <button
          class="bg-gdark px-2 py-1 rounded-md flex items-center gap-1 border-2 border-white"
          on:click={() => setRoute("export")}
        >
          <ArrowSmDownIcon class="h-5 w-5 text-white" />
        </button>
      </div>
      <div class="flex items-center gap-2">
        {#if fileUploadError}
          <div>
            {fileUploadError}
          </div>
        {/if}
      </div>
      {#if videoFiles}
        <div class="flex flex-col gap-2 max-h-screen overflow-y-auto">
          {#each $videoFiles as video (video.name)}
            <div
              use:draggable={video}
              class="flex items-center bg-gprimary hover:bg-stone-700 rounded-lg p-2 cursor-grab transition ease-in-out duration-500 gap-2"
              on:click={() => {
                fileUploadError = "";
              }}
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
    <!-- Video Player -->
    <div id="video-player" class="h-[28rem] w-3/4">
      <VideoPlayer />
    </div>
  </div>
  <ToolingLayout />
  <div class="flex justify-center select-none">
    {formatSecondsToHMS($trackTime)}: {formatSecondsToHMS($trackDuration)}
  </div>
  <!-- Timeline -->
  <Timeline />
</div>
