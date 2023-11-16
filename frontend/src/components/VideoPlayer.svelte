<script lang="ts">
  import {
    PlayIcon,
    PauseIcon,
    StopIcon,
    DesktopComputerIcon,
  } from "@rgossiaux/svelte-heroicons/solid";
  import SpeakerIcon from "../icons/SpeakerIcon.svelte";
  import MutedIcon from "../icons/MutedIcon.svelte";
  export let videoSrc: string;

  let duration: number;
  let currentTime: number;
  let playbackRate: number;
  let volume: number;
  let paused: boolean = true;
  let ended: boolean = false;
  let muted: boolean = false;
  let seeking: boolean;
  let buffered: TimeRanges;
  let played: TimeRanges;
  let seekable: TimeRanges;
  let videoWidth: number;
  let videoHeight: number;

  let videoContainer: HTMLElement;
  let video: HTMLVideoElement;
  let progress: HTMLProgressElement;
  let fullscreen: HTMLButtonElement;

  function setProgressMax() {
    progress.max = duration;
  }

  function progressBarUpdate() {
    if (progress.max === 0) {
      progress.max = duration;
    }
    progress.value = currentTime;
  }

  function progressSkipAhead(
    e: MouseEvent & {
      currentTarget: EventTarget & HTMLProgressElement;
    },
  ) {
    const rect = progress.getBoundingClientRect();
    const pos = (e.pageX - rect.left) / progress.offsetWidth;
    currentTime = pos * duration;
  }

  function handlePlayPause() {
    if (paused || ended) {
      video.play();
    } else {
      video.pause();
    }
  }

  function handleStop() {
    video.pause();
    currentTime = 0;
    progress.value = 0;
  }

  function handleMute() {
    muted = !muted;
    console.log("muted ", muted);
  }

  function handleVolume(dir: string) {
    if (dir === "+" && volume < 1) {
      volume += 0.1;
    } else if (dir === "-" && volume > 0) {
      volume -= 0.1;
    }
  }

  function handleFullScreen() {
    if (document.fullscreenElement !== null) {
      document.exitFullscreen;
      videoContainer.setAttribute("data-fullscreen", "false");
    } else {
      videoContainer.requestFullscreen();
      videoContainer.setAttribute("data-fullscreen", "true");
    }
  }
</script>

{#if videoSrc}
  <figure
    bind:this={videoContainer}
    class="h-full w-full overflow-hidden bg-gblue0 p-4 flex flex-col rounded-md border-white border-2"
  >
    <div id="video-player" class="h-5/6">
      <video
        id="video"
        class="block h-full w-full object-contain bg-gprimary rounded-md"
        src={videoSrc}
        bind:this={video}
        bind:duration
        bind:buffered
        bind:played
        bind:seekable
        bind:seeking
        bind:ended
        bind:currentTime
        bind:playbackRate
        bind:paused
        bind:volume
        bind:muted
        bind:videoWidth
        bind:videoHeight
        on:loadedmetadata={() => setProgressMax()}
        on:timeupdate={() => progressBarUpdate()}
      >
        <track kind="captions" />
      </video>
    </div>
    <!-- Video Controls -->
    <div
      id="video-controls"
      class="h-1/6 flex gap-2 items-center justify-center"
    >
      <button id="playpause" type="button" on:click={() => handlePlayPause()}>
        {#if paused || ended}
          <PlayIcon class="h-8 w-8 text-white duration-200 ease-in-out" />
        {:else}
          <PauseIcon class="h-8 w-8 text-white" />
        {/if}
      </button>
      <button id="stop" type="button" on:click={() => handleStop()}>
        <StopIcon class="h-8 w-8 text-white" />
      </button>
      <div class="cursor-pointer h-6 w-full">
        <progress
          class="h-full w-full [&::-webkit-progress-bar]:rounded-sm [&::-webkit-progress-value]:rounded-sm [&::-webkit-progress-bar]:bg-slate-300 [&::-webkit-progress-value]:bg-gred1 [&::-moz-progress-bar]:bg-gred1"
          id="progress"
          value="0"
          bind:this={progress}
          on:click={(e) => progressSkipAhead(e)}
        >
          <span id="progress-bar"></span>
        </progress>
      </div>
      <button id="mute" type="button" on:click={() => handleMute()}>
        {#if muted}
          <MutedIcon />
        {:else}
          <SpeakerIcon />
        {/if}
      </button>
      <button id="volinc" type="button" on:click={() => handleVolume("+")}
        >Vol+</button
      >
      <button id="voldec" type="button" on:click={() => handleVolume("-")}
        >Vol-</button
      >
      <button
        id="fs"
        type="button"
        bind:this={fullscreen}
        on:click={() => handleFullScreen()}
      >
        <DesktopComputerIcon class="h-8 w-8 text-white" />
      </button>
    </div>
  </figure>
{/if}
