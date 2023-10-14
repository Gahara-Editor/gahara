<script lang="ts">
  export let videoSrc: string;
  export let thumbnail: string;

  let duration: number;
  let currentTime: number;
  let playbackRate: number;
  let volume: number;
  let paused: boolean;
  let ended: boolean;
  let muted: boolean;
  let seeking: boolean;
  let buffered: TimeRanges;
  let played: TimeRanges;
  let seekable: TimeRanges;
  let videoWidth: number;
  let videoHeight: number;

  let videoContainer: HTMLElement;
  let video: HTMLVideoElement;
  let progress: HTMLProgressElement;

  function setProgressMax() {
    progress.max = duration;
  }

  function progressBarUpdate() {
    if (progress.max === 0) {
      progress.max = duration;
    }
    progress.value = currentTime;
    progress.style.width = `${Math.floor((currentTime * 100) / duration)}%`;
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
      document.exitFullscreen();
      videoContainer.setAttribute("data-fullscreen", "false");
    } else {
      videoContainer.requestFullscreen();
      videoContainer.setAttribute("data-fullscreen", "true");
    }
  }
</script>

{#if videoSrc}
  <figure bind:this={videoContainer}>
    <video
      id="video"
      preload="metadata"
      src={videoSrc}
      poster={thumbnail}
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
    <!-- Video Controls -->
    <ul id="video-controls" class="controls">
      <li>
        <button id="playpause" type="button" on:click={() => handlePlayPause()}
          >Play/Pause</button
        >
      </li>
      <li>
        <button id="stop" type="button" on:click={() => handleStop()}
          >Stop</button
        >
      </li>
      <li class="progress">
        <progress
          id="progress"
          value="0"
          bind:this={progress}
          on:click={(e) => progressSkipAhead(e)}
        >
          <span id="progress-bar"></span>
        </progress>
      </li>
      <li>
        <button id="mute" type="button" on:click={() => handleMute()}
          >Mute/Unmute</button
        >
      </li>
      <li>
        <button id="volinc" type="button" on:click={() => handleVolume("+")}
          >Vol+</button
        >
      </li>
      <li>
        <button id="voldec" type="button" on:click={() => handleVolume("-")}
          >Vol-</button
        >
      </li>
      <li>
        <button id="fs" type="button" on:click={() => handleFullScreen()}
          >Fullscreen</button
        >
      </li>
    </ul>
  </figure>
{/if}
