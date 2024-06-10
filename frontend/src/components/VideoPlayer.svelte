<script lang="ts">
  import {
    PlayIcon,
    PauseIcon,
    StopIcon,
  } from "@rgossiaux/svelte-heroicons/solid";
  import SpeakerIcon from "../icons/SpeakerIcon.svelte";
  import MutedIcon from "../icons/MutedIcon.svelte";
  import { onDestroy, onMount } from "svelte";
  import { videoStore, toolingStore } from "../stores";
  import ForwardIcon from "../icons/ForwardIcon.svelte";
  import BackwardIcon from "../icons/BackwardIcon.svelte";
  import {
    EventsEmit,
    EventsOff,
    EventsOn,
  } from "../../wailsjs/runtime/runtime";

  let volume: number;
  let muted: boolean;
  let videoWidth: number;
  let videoHeight: number;
  let videoContainer: HTMLElement;
  let video: HTMLVideoElement | null = null;

  let {
    source,
    duration,
    currentTime,
    paused,
    ended,
    playbackRate,
    handlePlaybackRate,
    setCurrentTime,
    resetVideo,
  } = videoStore;
  let { videoNode, editMode, cutEnd, isTrackPlaying, setIsTrackPlaying } =
    toolingStore;

  onMount(() => {
    setVideoPlayerDefaults();
  });

  $: muted = volume <= 0;
  $: {
    if ($editMode === "intervalCut") {
      if ($currentTime >= $cutEnd) {
        video.pause();
      }
    }
  }

  function setVideoPlayerDefaults() {
    volume = 0.5;
    paused.set(true);
    ended.set(false);
    muted = false;
  }

  function handlePlayPause() {
    if (video && video.readyState >= HTMLMediaElement.HAVE_CURRENT_DATA) {
      if ($paused || $ended) {
        video.play();
      } else {
        video.pause();
      }
      setIsTrackPlaying($paused);
    }
  }

  function handleStop() {
    if (video) video.pause();
    $currentTime = $videoNode ? $videoNode.start : 0;
  }

  function handleMute() {
    muted = !muted;
  }

  function handleTimeupdate() {
    if ($videoNode && $currentTime >= $videoNode.end)
      EventsEmit("evt_track_move", 1);
  }

  function seekVideo() {
    if ($videoNode) setCurrentTime($videoNode.start);
    if ($isTrackPlaying) video.play();
  }

  EventsOn("evt_play_track", () => {
    handlePlayPause();
  });

  onDestroy(() => {
    EventsOff("evt_play_track");
    resetVideo();
  });
</script>

<figure
  bind:this={videoContainer}
  class="h-full w-full overflow-hidden bg-gblue0 p-4 flex flex-col rounded-md border-white border-2"
>
  <div id="video-player" class="h-5/6">
    {#if $source}
      <video
        id="video"
        class="block h-full w-full object-contain bg-gprimary
        rounded-md"
        src={$videoNode
          ? $source + `#t=${$videoNode.start},${$videoNode.end}`
          : $source}
        bind:this={video}
        bind:duration={$duration}
        bind:ended={$ended}
        bind:currentTime={$currentTime}
        bind:playbackRate={$playbackRate}
        bind:paused={$paused}
        bind:volume
        bind:muted
        bind:videoWidth
        bind:videoHeight
        on:loadeddata={() => seekVideo()}
        on:timeupdate={() => handleTimeupdate()}
      >
        <track kind="captions" />
      </video>
    {:else}
      <div class="block h-full w-full bg-gprimary rounded-md"></div>
    {/if}
  </div>
  <!-- Video Controls -->
  {#if $source}
    {#if $videoNode}
      <div class="text-center text-lg font-semibold">
        {$videoNode.name}
      </div>
    {/if}
    <div
      id="video-controls"
      class="h-1/6 flex gap-2 items-center justify-center"
    >
      <button
        id="playbackrate"
        type="button"
        on:click={() => handlePlaybackRate("down")}
      >
        <BackwardIcon class="h-8 w-8 text-white" />
      </button>

      <button id="playpause" type="button" on:click={() => handlePlayPause()}>
        {#if $paused || $ended}
          <PlayIcon class="h-8 w-8 text-white" />
        {:else}
          <PauseIcon class="h-8 w-8 text-white" />
        {/if}
      </button>
      <button
        id="playbackrate"
        type="button"
        on:click={() => handlePlaybackRate("up")}
      >
        <ForwardIcon class="h-8 w-8 text-white" />
      </button>

      <button id="stop" type="button" on:click={() => handleStop()}>
        <StopIcon class="h-8 w-8 text-white" />
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
      <div class="px-2 rounded-md border-white border-2">
        x{$playbackRate}
      </div>
    </div>
  {/if}
</figure>
