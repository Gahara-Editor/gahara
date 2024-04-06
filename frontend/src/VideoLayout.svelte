<script lang="ts">
  import {
    FilePicker,
    ReadProjectWorkspace,
    LoadTimeline,
    SaveTimeline,
    ResetTimeline,
    DeleteRIDReferences,
    DeleteProjectFile,
  } from "../wailsjs/go/main/App";
  import {
    XIcon,
    FilmIcon,
    ArrowSmDownIcon,
  } from "@rgossiaux/svelte-heroicons/solid";
  import type { main } from "../wailsjs/go/models";
  import { router, videoStore, videoFiles, trackStore } from "./stores";
  import { onDestroy, onMount } from "svelte";
  import Modal from "./components/Modal.svelte";
  import VideoPlayer from "./components/VideoPlayer.svelte";
  import Timeline from "./components/Timeline.svelte";
  import { draggable } from "./lib/dnd";
  import {
    EventsOff,
    EventsOn,
    WindowSetTitle,
  } from "../wailsjs/runtime/runtime";
  import ToolingLayout from "./ToolingLayout.svelte";
  import FloppyDisk from "./icons/FloppyDisk.svelte";
  import FolderOpenIcon from "./icons/FolderOpenIcon.svelte";
  import TrashIcon from "./icons/TrashIcon.svelte";
  import WarningIcon from "./icons/WarningIcon.svelte";

  const { setVideoSrc, resetVideo } = videoStore;
  const {
    addVideoToTrack,
    removeRIDReferencesFromTrack,
    trackTime,
    trackDuration,
  } = trackStore;
  const {
    pipelineMessages,
    videoFilesError,
    removeVideoFile,
    addPipelineMsg,
    removePipelineMsg,
    addVideos,
    setVideoFilesError,
    resetVideoFiles,
  } = videoFiles;
  const { setRoute } = router;

  onMount(() => {
    loadProjectFiles();
    loadTimeline();
  });

  function loadProjectFiles() {
    if ($videoFiles.length === 0) {
      ReadProjectWorkspace()
        .then((files) => addVideos(files))
        .catch(() => setVideoFilesError("No files in this project"));
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

  function selectFile() {
    FilePicker()
      .then(() => {})
      .catch(() => setVideoFilesError("no file selected"));
  }

  async function deleteVideoFile(video: main.Video) {
    try {
      const rid = `${video.filepath}/${video.name}${video.extension}`;
      await Promise.all([DeleteRIDReferences(rid), DeleteProjectFile(rid)]);
      await SaveTimeline();
      removeRIDReferencesFromTrack(0, rid);
      removeVideoFile(video.name);
    } catch (err) {
      setVideoFilesError(err);
    }
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

  EventsOn("evt_proxy_file_created", (video: main.Video) => {
    setVideoFilesError("");
    addVideos([video]);
    removePipelineMsg();
  });
  EventsOn("evt_proxy_error_msg", (msg: string) => {
    setVideoFilesError(msg);
  });
  EventsOn("evt_proxy_pipeline_msg", (msg: string) => {
    addPipelineMsg(msg);
  });

  onDestroy(() => {
    EventsOff(
      "evt_proxy_file_created",
      "evt_error_msg",
      "evt_proxy_pipeline_msg",
    );
  });
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
              class="rounded-lg bg-gred1 font-semibold text-white inline-flex items-center px-2 py-0.5 hover:bg-gred transition ease-in-out duration-200 border-2 border-white"
              on:click={close}>Back</button
            >
            <button
              class="bg-gblue0 rounded-lg px-2 py-1 border-2 border-white flex items-center gap-2 hover:bg-gblue transition ease-in-out duration-200"
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
        {#if $videoFilesError}
          <div>
            {$videoFilesError}
          </div>
        {/if}
      </div>
      {#if $videoFiles}
        <div class="flex flex-col gap-2 max-h-screen overflow-y-auto">
          {#each $videoFiles as video (video.name)}
            <div
              use:draggable={video}
              class="flex items-center bg-gprimary hover:bg-stone-700 rounded-lg p-2 cursor-grab transition ease-in-out duration-500 gap-2 border-white border-2"
            >
              <Modal>
                <div slot="trigger" let:open>
                  <button
                    class="bg-red-500 hover:bg-red-400 px-1 py-1 rounded-full"
                    on:click={open}
                  >
                    <XIcon class="h-3 w-3 text-white" />
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
                    Do you want to delete {video.name}?
                  </p>
                  <p>
                    All of the clips in the timeline that reference this file
                    will also be deleted
                  </p>
                </div>
                <div
                  slot="footer"
                  let:store={{ close }}
                  class="flex items-center justify-center gap-2"
                >
                  <button
                    class="rounded-lg bg-gred1 font-semibold text-white inline-flex items-center px-2 py-0.5 hover:bg-gred transition ease-in-out duration-200 border-2 border-white"
                    on:click={close}>Back</button
                  >
                  <button
                    class="bg-gblue0 rounded-lg px-2 py-1 border-2 border-white flex items-center gap-2 hover:bg-gblue transition ease-in-out duration-200"
                    on:click={() => {
                      deleteVideoFile(video);
                    }}
                  >
                    <TrashIcon class="h-5 w-5 text-white" />
                    <span>Delete File</span>
                  </button>
                </div>
              </Modal>

              <p class="text-sm">{video.name}</p>
            </div>
          {/each}
          {#each $pipelineMessages as msg (msg)}
            <div
              class="flex items-center bg-obsbg hover:bg-obsternary rounded-lg p-2 cursor-grab transition ease-in-out duration-500 gap-2 border-white border-2 overflow-hidden whitespace-nowrap"
            >
              <svg
                class="animate-spin h-5 w-5 text-white min-w-[20px]"
                xmlns="http://www.w3.org/2000/svg"
                fill="none"
                viewBox="0 0 24 24"
              >
                <circle
                  class="opacity-25"
                  cx="12"
                  cy="12"
                  r="10"
                  stroke="currentColor"
                  stroke-width="4"
                ></circle>
                <path
                  class="opacity-75"
                  fill="currentColor"
                  d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
                ></path>
              </svg>

              <p class="text-sm font-bold">{msg}</p>
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
