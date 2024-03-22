<script lang="ts">
  import { ChevronDownIcon } from "@rgossiaux/svelte-heroicons/solid";
  import { ExportVideo, GetOutputFileSavePath } from "../wailsjs/go/main/App";
  import { exportOptionsStore, router } from "./stores";
  import { onDestroy } from "svelte";
  import FolderOpenIcon from "./icons/FolderOpenIcon.svelte";
  import { EventsOff, EventsOn } from "../wailsjs/runtime/runtime";

  const { setRoute } = router;
  const {
    filename,
    outputPath,
    preset,
    codec,
    setCodec,
    videoFormat,
    resolution,
    videoFormats,
    resolutionOpts,
    presetOpts,
    isProcessingVid,
    processingMsg,
    progressPercentage,
    getCompatibleCodecs,
    setProgressPercentage,
    setProcessingMsg,
    setIsProcessingVid,
    getExportOptions,
    setOutputPath,
    resetExportOptionsStore,
  } = exportOptionsStore;

  $: {
    switch ($videoFormat) {
      case ".webm":
        setCodec("libvpx-vp9");
        break;
      case ".ogv":
        setCodec("libtheora");
        break;
      default:
        setCodec("libx264");
    }
  }

  function handleOutputPath() {
    GetOutputFileSavePath()
      .then((outpath) => {
        setOutputPath(outpath);
      })
      .catch(console.log);
  }

  function handleExport() {
    setIsProcessingVid(true);
    setProcessingMsg("");
    ExportVideo(getExportOptions())
      .then((msg) => setProcessingMsg(msg))
      .catch((msg) => setProcessingMsg(msg))
      .finally(() => {
        setIsProcessingVid(false);
        setProgressPercentage(0.0);
      });
  }

  EventsOn("evt_encoding_progress", (progress: number) => {
    setProgressPercentage(progress);
  });

  onDestroy(() => {
    EventsOff("evt_encoding_progress");
    resetExportOptionsStore();
  });
</script>

<div class="h-full rounded-md bg-gprimary flex justify-center items-center">
  <div class="flex flex-col items-center justify-center gap-2">
    <div
      class="flex flex-col gap-4 bg-obsbg p-4 border-2 border-white rounded-md"
    >
      <!-- SAVE PATH -->
      <div class="flex items-center justify-between gap-1">
        <label for="export-path" class="text-white">Location</label>
        <input
          id="export-path"
          disabled
          type="text"
          bind:value={$outputPath}
          class="p-1 rounded-sm"
        />
        <button
          class="bg-gdark px-2 py-1 rounded-md flex items-center gap-1 border-2 border-white"
          on:click={() => handleOutputPath()}
        >
          <FolderOpenIcon class="h-5 w-5 text-white" />
        </button>
      </div>
      <!-- Filename -->
      <div class="flex items-center justify-between">
        <label for="filename-path" class="text-white">File Name:</label>
        <input
          id="filename"
          type="text"
          autocorrect="off"
          autocomplete="off"
          bind:value={$filename}
          class="p-1 rounded-sm"
        />
      </div>
      <!-- Video Format -->
      <div class="flex items-center justify-between">
        <label for="videoFormats" class="text-white">Format (.mp4, .mov):</label
        >
        <div class="relative inline-flex">
          <select
            id="videoFormats"
            bind:value={$videoFormat}
            class="block appearance-none w-full bg-white border border-indigo-500 hover:border-gray-500 px-4 py-2 pr-8 rounded leading-tight focus:outline-none focus:bg-white focus:border-indigo-600"
          >
            {#each videoFormats as vFormat (vFormat)}
              <option value={vFormat}>
                {vFormat}
              </option>
            {/each}
          </select>
          <div
            class="pointer-events-none absolute inset-y-0 right-0 flex items-center px-2 text-indigo-500"
          >
            <ChevronDownIcon class="h-6 w-6" />
          </div>
        </div>
      </div>

      <!-- Video Codec -->
      <div class="flex items-center justify-between">
        <label for="videoCodecs" class="text-white">Codec (H264, H265):</label>
        <div class="relative inline-flex">
          <select
            id="videoCodecs"
            bind:value={$codec}
            class="block appearance-none w-full bg-white border border-indigo-500 hover:border-gray-500 px-4 py-2 pr-8 rounded leading-tight focus:outline-none focus:bg-white focus:border-indigo-600"
          >
            {#each getCompatibleCodecs($videoFormat) as codec (codec[0])}
              <option value={codec[1]}>
                {codec[0]}
              </option>
            {/each}
          </select>
          <div
            class="pointer-events-none absolute inset-y-0 right-0 flex items-center px-2 text-indigo-500"
          >
            <ChevronDownIcon class="h-6 w-6" />
          </div>
        </div>
      </div>
      <!--  Resolution -->
      <div class="flex items-center justify-between">
        <label for="resolutionOpts" class="text-white">Resolution:</label>
        <div class="relative inline-flex">
          <select
            id="resolutionOpts"
            bind:value={$resolution}
            class="block appearance-none w-full bg-white border border-indigo-500 hover:border-gray-500 px-4 py-2 pr-8 rounded leading-tight focus:outline-none focus:bg-white focus:border-indigo-600"
          >
            {#each resolutionOpts as resolutionOpt (resolutionOpt)}
              <option value={resolutionOpt}>
                {resolutionOpt}
              </option>
            {/each}
          </select>
          <div
            class="pointer-events-none absolute inset-y-0 right-0 flex items-center px-2 text-indigo-500"
          >
            <ChevronDownIcon class="h-6 w-6" />
          </div>
        </div>
      </div>
      <!-- Preset -->
      <div class="flex items-center justify-between">
        <label for="presetOpts" class="text-white">Encoding Speed:</label>
        <div class="relative inline-flex">
          <select
            id="presetOpts"
            bind:value={$preset}
            class="block appearance-none w-full bg-white border border-indigo-500 hover:border-gray-500 px-4 py-2 pr-8 rounded leading-tight focus:outline-none focus:bg-white focus:border-indigo-600"
          >
            {#each presetOpts as presetOpt (presetOpt)}
              <option value={presetOpt}>
                {presetOpt}
              </option>
            {/each}
          </select>
          <div
            class="pointer-events-none absolute inset-y-0 right-0 flex items-center px-2 text-indigo-500"
          >
            <ChevronDownIcon class="h-6 w-6" />
          </div>
        </div>
      </div>
      <!-- Actions -->
      {#if !$isProcessingVid}
        <div class="flex justify-center items-center gap-2">
          <button
            class="rounded-lg bg-gred1 border-white border-2 font-semibold text-white inline-flex items-center px-4 py-1.5 hover:bg-gred transition ease-in-out duration-200"
            on:click={() => setRoute("video")}>Back</button
          >
          <button
            class="rounded-lg bg-gblue0 border-white border-2 font-semibold text-white inline-flex items-center px-4 py-1.5 hover:bg-gblue transition ease-in-out duration-200"
            on:click={() => handleExport()}>Export</button
          >
        </div>
      {:else}
        <div
          class="rounded-lg bg-gblue0 border-white border-2 font-semibold text-white inline-flex items-center justify-center px-4 py-1.5 gap-2"
        >
          <svg
            class="animate-spin -ml-1 mr-3 h-5 w-5 text-white"
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
          <span> Processing video... </span>
        </div>
        <!-- Progress Bar -->
        <div
          class="relative h-6 bg-obsbg border-white border-2 rounded-full overflow-hidden"
        >
          <div
            class="h-full bg-gblue0"
            style={`width: ${$progressPercentage}%; transition: width 0.3s ease-in-out;`}
          ></div>
          <div
            class="absolute inset-0 flex items-center justify-center text-white font-semibold"
          >
            {$progressPercentage}%
          </div>
        </div>
      {/if}
    </div>
    {#if $processingMsg}
      <div class="rounded-md border-2 border-white bg-obsbg text-white p-2">
        {$processingMsg}
      </div>
    {/if}
  </div>
</div>
