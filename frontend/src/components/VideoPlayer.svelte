<script lang="ts">
  import {
    PlayIcon,
    PauseIcon,
    StopIcon,
    DesktopComputerIcon,
  } from "@rgossiaux/svelte-heroicons/solid";
  import SpeakerIcon from "../icons/SpeakerIcon.svelte";
  import MutedIcon from "../icons/MutedIcon.svelte";
  import { onDestroy, onMount } from "svelte";
  import { videoStore, toolingStore } from "../stores";
  import ForwardIcon from "../icons/ForwardIcon.svelte";
  import BackwardIcon from "../icons/BackwardIcon.svelte";

  const { videoNode } = toolingStore;

  let playbackRate: number;
  let volume: number;
  let muted: boolean;
  let seeking: boolean;
  let buffered: TimeRanges;
  let played: TimeRanges;
  let seekable: TimeRanges;
  let videoWidth: number;
  let videoHeight: number;

  let videoContainer: HTMLElement;
  let video: HTMLVideoElement;
  let videoSrc: string;

  let { source, duration, currentTime, paused, ended } = videoStore;
  let { editMode, cutStart, cutEnd } = toolingStore;

  onMount(() => {
    videoSrc = $source;
    setVideoPlayerDefaults();
  });

  $: muted = volume <= 0;
  $: {
    if (videoSrc !== $source) {
      videoSrc = $source;
      setVideoPlayerDefaults();
    }
  }
  $: {
    if ($editMode === "cut") {
      if ($currentTime >= $cutEnd) {
        video.pause();
      }
    }
  }

  function setVideoPlayerDefaults() {
    playbackRate = 1;
    volume = 0.5;
    paused.set(true);
    ended.set(false);
    muted = false;
  }

  function handlePlayPause() {
    if ($editMode === "cut") {
      currentTime.set($cutStart);
    }

    if ($paused || $ended) {
      video.play();
    } else {
      video.pause();
    }
  }

  function handleStop() {
    video.pause();
    $currentTime = $videoNode ? $videoNode.start : 0;
  }

  function handleMute() {
    muted = !muted;
  }

  function handlePlaybackRate(dir: string) {
    switch (dir) {
      case "down":
        playbackRate =
          playbackRate / 2 > 0.125 ? playbackRate / 2 : playbackRate;
        break;
      default:
        playbackRate = playbackRate * 2 < 8 ? playbackRate * 2 : playbackRate;
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

  onDestroy(() => {
    videoSrc = "";
  });
</script>

<figure
  bind:this={videoContainer}
  class="h-full w-full overflow-hidden bg-gblue0 p-4 flex flex-col rounded-md border-white border-2"
>
  <div id="video-player" class="h-5/6">
    {#if videoSrc}
      <video
        id="video"
        class="block h-full w-full object-contain bg-gprimary
        rounded-md"
        src={videoSrc}
        bind:this={video}
        bind:duration={$duration}
        bind:buffered
        bind:played
        bind:seekable
        bind:seeking
        bind:ended={$ended}
        bind:currentTime={$currentTime}
        bind:playbackRate
        bind:paused={$paused}
        bind:volume
        bind:muted
        bind:videoWidth
        bind:videoHeight
      >
        <track kind="captions" />
      </video>
    {:else}
      <div class="block h-full w-full bg-gprimary rounded-md"></div>
    {/if}
  </div>
  <!-- Video Controls -->
  {#if videoSrc}
    <div
      id="video-controls"
      class="h-1/6 flex gap-2 items-center justify-center"
    >
      <button id="playpause" type="button" on:click={() => handlePlayPause()}>
        {#if $paused || $ended}
          <PlayIcon class="h-8 w-8 text-white" />
        {:else}
          <PauseIcon class="h-8 w-8 text-white" />
        {/if}
      </button>
      <button id="stop" type="button" on:click={() => handleStop()}>
        <StopIcon class="h-8 w-8 text-white" />
      </button>
      <button
        id="playbackrate"
        type="button"
        on:click={() => handlePlaybackRate("down")}
      >
        <BackwardIcon class="h-8 w-8 text-white" />
      </button>

      <button
        id="playbackrate"
        type="button"
        on:click={() => handlePlaybackRate("up")}
      >
        <ForwardIcon class="h-8 w-8 text-white" />
      </button>

      <button id="mute" type="button" on:click={() => handleMute()}>
        {#if muted}
          <MutedIcon />
        {:else}
          <SpeakerIcon />
        {/if}
      </button>
      <div class="flex flex-col items-center">
        <input
          type="range"
          min="0"
          max="1"
          step="0.01"
          bind:value={volume}
          class="w-20 cursor-pointer [&::-webkit-slider-runnable-track]:bg-slate-300 [&::-webkit-slider-runnable-track]:rounded-xl"
        />
      </div>
      <button id="fs" type="button" on:click={() => handleFullScreen()}>
        <DesktopComputerIcon class="h-8 w-8 text-white" />
      </button>
      <div class="px-2 rounded-md border-white border-2">
        x{playbackRate}
      </div>
    </div>
  {/if}
</figure>
