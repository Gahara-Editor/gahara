<script lang="ts">
  import { dropzone } from "../lib/dnd";
  import { onDestroy } from "svelte";
  import type { video } from "wailsjs/go/models";
  import { videoStore, toolingStore, trackStore } from "../stores";
  import Playhead from "../icons/Playhead.svelte";
  import { slide } from "svelte/transition";
  import { cubicIn, cubicOut } from "svelte/easing";
  import { flip } from "svelte/animate";

  const { setVideoSrc, currentTime, setCurrentTime } = videoStore;
  const {
    cutStart,
    cutEnd,
    editMode,
    videoNode,
    videoNodeWidth,
    playheadPos,
    isMovingPlayhead,
    isMovingCutRangeBox,
    boxLeftBound,
    boxRightBound,
    moveCutRangeBox,
    movePlayhead,
    setCutEnd,
    setPlayheadPos,
    setVideoNode,
    setVideoNodePos,
    setVideoNodeWidth,
    setClipStart,
    setClipEnd,
    setBoxLeftBound,
    setBoxRightBound,
    resetToolingStore,
  } = toolingStore;
  const { trackDuration, setTrackTime, resetTrackStore } = trackStore;

  let selectedID = 0;
  let trackNode: HTMLDivElement;
  let timelineNode: HTMLDivElement;
  let cutRangeBox: HTMLDivElement;
  let cutRangeSide: "left" | "right" | "middle" | "none";

  function getTrackWidth() {
    const track = document.getElementById(`track-${selectedID}`);
    return track.clientWidth;
  }

  function getTrackLeft() {
    const trackNode = document.getElementById(`track-${selectedID}`);
    return trackNode.getBoundingClientRect().left;
  }

  function handleEditModeMouseDown(
    e: MouseEvent & {
      currentTarget: EventTarget & HTMLDivElement;
    },
  ) {
    switch ($editMode) {
      case "timeline":
        movePlayhead(true);
        break;
      case "intervalCut":
        setCutSideDragged(e);
        break;
      default:
    }
  }

  function handleEditModeMouseMove(
    e: MouseEvent & {
      currentTarget: EventTarget & HTMLDivElement;
    },
    pos: number,
    tVideo: video.VideoNode,
  ) {
    switch ($editMode) {
      case "timeline":
        handleTimelineMove(e, pos, tVideo);
        break;
      case "intervalCut":
        adjustCutRange(e);
        break;
      default:
    }
  }

  function handleEditModeMouseUp() {
    switch ($editMode) {
      case "timeline":
        movePlayhead(false);
        break;
      case "intervalCut":
        moveCutRangeBox(false);
        break;
      default:
    }
  }

  function handleTimelineMove(
    e: MouseEvent & {
      currentTarget: EventTarget & HTMLDivElement;
    },
    pos: number,
    tVideo: video.VideoNode,
  ) {
    if ($isMovingPlayhead) {
      setPlayheadPos(Math.min(e.clientX - getTrackLeft(), getTrackWidth()));
      handleVideoNode(e, pos, tVideo);
    }
  }

  function setCutSideDragged(
    e: MouseEvent & {
      currentTarget: EventTarget & HTMLDivElement;
    },
  ) {
    e.preventDefault();
    e.stopPropagation();
    const { clientX } = e;
    const { left, width } = cutRangeBox.getBoundingClientRect();

    if (clientX <= left + width / 3) {
      cutRangeSide = "left";
    } else if (clientX >= left + (2 * width) / 3) {
      cutRangeSide = "right";
    } else {
      cutRangeSide = "middle";
    }
    moveCutRangeBox(true);
  }

  function adjustCutRange(
    e: MouseEvent & {
      currentTarget: EventTarget & HTMLDivElement;
    },
  ) {
    e.preventDefault();
    e.stopPropagation();

    if ($isMovingCutRangeBox) {
      const cutRangeBoxStyle = getComputedStyle(cutRangeBox);
      const boxLeft = parseFloat(cutRangeBoxStyle.left);
      const boxWidth = parseFloat(cutRangeBoxStyle.width);
      let [newBoxLeft, newBoxWidth, trackTime] = [-1, -1, 0];
      const adjustedX =
        e.clientX +
        timelineNode.scrollLeft -
        timelineNode.getBoundingClientRect().left;

      const mousePos = Math.max(
        $boxLeftBound,
        Math.min(adjustedX, $boxRightBound),
      );

      setCurrentTime(
        $videoNode.start +
          ($videoNode.end - $videoNode.start) *
            ((mousePos - $boxLeftBound) / $videoNodeWidth),
      );

      switch (cutRangeSide) {
        case "left":
          newBoxWidth =
            adjustedX <= $boxLeftBound
              ? boxWidth
              : boxWidth + (boxLeft - adjustedX);
          newBoxLeft = Math.max(
            $boxLeftBound,
            Math.min(adjustedX, $boxRightBound - boxWidth),
          );
          trackTime = (newBoxLeft / getTrackWidth()) * $trackDuration;
          cutStart.set($currentTime);
          break;
        case "right":
          newBoxWidth = Math.min(adjustedX - boxLeft, $boxRightBound - boxLeft);
          trackTime =
            ((boxLeft + newBoxWidth) / getTrackWidth()) * $trackDuration;
          cutEnd.set($currentTime);
          break;
        case "middle":
          const cutRangeBoxPos = boxLeft + (e.movementX || 0);
          newBoxLeft = Math.max(
            $boxLeftBound,
            Math.min(cutRangeBoxPos, $boxRightBound - boxWidth),
          );
          trackTime = (newBoxLeft / trackNode.clientWidth) * $trackDuration;
          cutStart.set($currentTime);
          break;
        default:
      }

      setTrackTime(trackTime);
      cutRangeBox.style.left = `${newBoxLeft === -1 ? boxLeft : newBoxLeft}px`;
      cutRangeBox.style.width = `${
        newBoxWidth === -1 ? boxWidth : newBoxWidth
      }px`;
    }
  }

  function handleBoxRender(
    e: MouseEvent & {
      currentTarget: EventTarget & HTMLDivElement;
    },
    pos: number,
    video: video.VideoNode,
  ) {
    setBoxLeftBound(e.currentTarget.offsetLeft);
    setBoxRightBound(e.currentTarget.offsetLeft + e.currentTarget.clientWidth);
    setVideoNodeWidth(e.currentTarget.getBoundingClientRect().width);
    setCurrentTime(video.start);
    setClipStart(video.start);
    setClipEnd(video.end);
    setVideoNode(video);
    setVideoNodePos(pos);
    setVideoSrc(video.rid);
  }

  function handleVideoNode(
    e: MouseEvent & {
      currentTarget: EventTarget & HTMLDivElement;
    },
    pos: number,
    video: video.VideoNode,
  ) {
    const clipWidth = e.currentTarget.getBoundingClientRect().width;
    const mousePos = e.clientX - e.currentTarget.getBoundingClientRect().left;
    const time =
      video.start + (video.end - video.start) * (mousePos / clipWidth);
    const trackTime =
      ((e.clientX + timelineNode.scrollLeft) / getTrackWidth()) *
      $trackDuration;
    setVideoSrc(video.rid);
    setVideoNodePos(pos);
    setVideoNode(video);
    setCurrentTime(time);
    setTrackTime(trackTime);
    setCutEnd(time);
    setBoxLeftBound(e.currentTarget.offsetLeft);
    setBoxRightBound(e.currentTarget.offsetLeft + e.currentTarget.clientWidth);
  }

  onDestroy(() => {
    resetTrackStore();
    resetToolingStore();
  });
</script>

<div
  class="timeline h-full w-full bg-gdark border-t-2 border-t-white flex flex-col gap-4 pt-4 pb-4 px-1 relative overflow-x-scroll overflow-y-hidden"
  bind:this={timelineNode}
  use:dropzone={{}}
  on:mouseup={() => handleEditModeMouseUp()}
>
  {#if $trackStore.length <= 0}
    <div class="flex justify-center items-center">
      <p class="text-white text-4xl font-semibold select-none">
        Drag And Drop Video Clips
      </p>
    </div>
  {:else if $editMode === "timeline"}
    <div
      class="absolute top-0 left-0 h-full w-3 z-10"
      style={`left: ${$playheadPos}px`}
    >
      <Playhead />
    </div>
  {/if}

  <!-- VIDEO TRACKS -->
  <!-- TODO Create an actual id for a track -->
  {#each $trackStore as track, id (id)}
    <div
      bind:this={trackNode}
      class="h-28 flex relative bg-gprimary gap-1 p-2 w-max"
      id={`track-${id}`}
    >
      <!-- Video Track -->
      {#each track as tVideo, pos (tVideo.id)}
        <div
          animate:flip={{ duration: 100 }}
          out:slide={{ axis: "x", duration: 100, easing: cubicOut }}
        >
          {#if $editMode === "intervalCut" && $videoNode && $videoNode.id === tVideo.id}
            <div
              class="absolute border-yellow-500 border-2 h-24 cursor-grab"
              style={`width: ${
                (tVideo.end - tVideo.start) * 20
              }px; left: ${$boxLeftBound}px`}
              bind:this={cutRangeBox}
              id="cut-range"
              on:mousemove={(e) => {
                handleEditModeMouseMove(e, pos, tVideo);
              }}
              on:mousedown={(e) => {
                handleEditModeMouseDown(e);
              }}
            ></div>
          {/if}
          <!-- Video Nodes of this track -->
          <div
            id={`videoNode-${tVideo.id}`}
            class="h-full bg-gblue0 border-white border-2 cursor-pointer select-none overflow-hidden"
            style={`width: ${
              (tVideo.end - tVideo.start) * 20
            }px; border-color: ${
              $editMode === "remove" &&
              $videoNode &&
              $videoNode.id === tVideo.id
                ? "#f7768e"
                : "#ffffff"
            }; `}
            on:click={(e) => handleBoxRender(e, pos, tVideo)}
            on:mousemove={(e) => {
              handleEditModeMouseMove(e, pos, tVideo);
            }}
            on:mousedown={(e) => {
              handleEditModeMouseDown(e);
            }}
          >
            Node
            {tVideo.end - tVideo.start}
          </div>
        </div>
      {/each}
    </div>
  {/each}
</div>

<style>
  .timeline:global(.droppable) {
    border-width: 2px;
    border-color: rgb(122, 162, 247);
  }

  .timeline:global(.droppable) * {
    pointer-events: none;
  }
</style>
