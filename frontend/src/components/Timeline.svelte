<script lang="ts">
  import { dropzone } from "../lib/dnd";
  import { onDestroy } from "svelte";
  import type { video } from "wailsjs/go/models";
  import {
    videoStore,
    toolingStore,
    trackStore,
    createBooleanStore,
  } from "../stores";
  import Playhead from "../icons/Playhead.svelte";
  import { slide } from "svelte/transition";
  import { cubicOut } from "svelte/easing";
  import { flip } from "svelte/animate";
  import { EventsOff, EventsOn } from "../../wailsjs/runtime/runtime";
  import EventModal from "./EventModal.svelte";
  import {
    InsertInterval,
    RenameVideoNode,
    ToggleLossless,
    MarkAllLossless,
    UnmarkAllLossless,
  } from "../../wailsjs/go/main/App";
  import RenameIcon from "../icons/RenameIcon.svelte";
  import SearchList from "../components/SearchList.svelte";
  import { formatSecondsToHMS } from "../lib/utils";

  const { isOpen, close, open } = createBooleanStore(false);
  const { setVideoSrc, currentTime, setCurrentTime } = videoStore;
  const {
    vimMode,
    cutStart,
    cutEnd,
    editMode,
    videoNode,
    videoNodePos,
    videoNodeWidth,
    videoNodeName,
    trackZoom,
    adjustTrackZoom,
    playheadPos,
    isMovingPlayhead,
    isMovingCutRangeBox,
    boxLeftBound,
    boxRightBound,
    clipCursorIdx,
    clipRegister,
    setActionMsg,
    setIsOpenSearchList,
    setClipCursorIdx,
    setVimMode,
    moveClipCursor,
    setClipRegister,
    moveCutRangeBox,
    movePlayhead,
    setCutEnd,
    setPlayheadPos,
    setVideoNode,
    setVideoNodeName,
    setVideoNodePos,
    setVideoNodeWidth,
    setClipStart,
    setClipEnd,
    setBoxLeftBound,
    setBoxRightBound,
    resetToolingStore,
  } = toolingStore;
  const {
    trackDuration,
    addVideoToTrack,
    setTrackTime,
    renameClipInTrack,
    toggleLosslessMarkofClip,
    markAllLossless,
    unmarkAllLossless,
    resetTrackStore,
  } = trackStore;

  let selectedID = 0;
  let trackNode: HTMLDivElement;
  let timelineNode: HTMLDivElement;
  let cutRangeBox: HTMLDivElement;
  let cutRangeSide: "left" | "right" | "middle" | "none";
  let cursorColor = "#ffffff";

  $: {
    cursorColor =
      $editMode === "remove"
        ? "#f7768e"
        : $vimMode && $editMode === "select"
        ? "#1abc9c"
        : "#ffffff";
  }

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
    setPlayheadPos(Math.min(e.currentTarget.offsetLeft, getTrackWidth()));
    setVideoNodeWidth(e.currentTarget.getBoundingClientRect().width);
    setCurrentTime(video.start);
    setClipStart(video.start);
    setClipEnd(video.end);
    setVideoNode(video);
    setVideoNodePos(pos);
    setVideoSrc(video.rid);
    setClipCursorIdx(pos);
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
    setClipCursorIdx(pos);
    setBoxLeftBound(e.currentTarget.offsetLeft);
    setBoxRightBound(e.currentTarget.offsetLeft + e.currentTarget.clientWidth);
  }

  function handleVideoNodeRename() {
    if ($videoNodeName === "") return;
    RenameVideoNode($videoNodePos, $videoNodeName)
      .then(() => {
        renameClipInTrack(0, $videoNodePos, $videoNodeName);
        setVideoNodeName("");
      })
      .catch(() => setActionMsg("could not rename clip"));
    setVimMode(true);
    close();
  }

  function handleKeybindTrackClipMove() {
    if (!$vimMode) return;
    const videoNodeDiv = document
      .getElementById(`track-${selectedID}`)
      ?.querySelector(`div:nth-child(${$clipCursorIdx + 1})`)
      ?.querySelector("div");
    if (videoNodeDiv) {
      videoNodeDiv.click();
      scrollToNode(videoNodeDiv);
    }
  }

  function scrollToNode(node: HTMLDivElement) {
    const timelineContainer = document.getElementById("timeline");
    const timelineRect = timelineContainer.getBoundingClientRect();
    const nodeRect = node.getBoundingClientRect();

    const isNodeVisible =
      nodeRect.left >= timelineRect.left &&
      nodeRect.right <= timelineRect.right;

    if (!isNodeVisible) {
      const scrollX =
        nodeRect.left - timelineRect.left + timelineContainer.scrollLeft;
      timelineContainer.scrollTo({
        left: scrollX,
        behavior: "smooth",
      });
    }
  }

  EventsOn("evt_open_rename_clip_modal", () => {
    if ($videoNode) {
      setVimMode(false);
      open();
    }
  });
  EventsOn("evt_rename_clip", () => {
    setVimMode(false);
    handleVideoNodeRename();
  });
  EventsOn("evt_track_move", (inc: number) => {
    moveClipCursor(inc);
    handleKeybindTrackClipMove();
  });
  EventsOn("evt_open_search_list", () => {
    if ($vimMode) {
      setVimMode(false);
      setIsOpenSearchList(true);
    }
  });
  EventsOn("evt_yank_clip", () => {
    if ($videoNode) {
      setClipRegister($videoNode);
      setActionMsg(`YANKED: ${$videoNode.name}`);
    }
  });
  EventsOn("evt_insertclip_edit", () => {
    if ($clipRegister) {
      InsertInterval(
        $clipRegister.rid,
        $clipRegister.name,
        $clipRegister.start,
        $clipRegister.end,
        $videoNodePos,
      )
        .then((tVideo) => {
          addVideoToTrack(0, tVideo, $videoNodePos);
          setVideoNode(tVideo);
          setVideoSrc(tVideo.rid);
          setCurrentTime(tVideo.start);
          setActionMsg(`PASTED: ${$videoNode.name}`);
        })
        .catch(() =>
          toolingStore.setActionMsg(
            `could not insert ${$clipRegister.name} from clip register`,
          ),
        );
    }
  });
  EventsOn("evt_zoom_timeline", (dir: string) => {
    adjustTrackZoom(dir);
    handleKeybindTrackClipMove();
  });

  EventsOn("evt_saved_timeline", (msg: string) => {
    setActionMsg(msg);
  });

  EventsOn("evt_toggle_lossless", () => {
    if ($videoNode) {
      ToggleLossless($videoNodePos)
        .then(() => {
          toggleLosslessMarkofClip(0, $videoNodePos);
        })
        .catch((err) => setActionMsg(err));
    }
  });

  EventsOn("evt_mark_all_lossless", () => {
    MarkAllLossless()
      .then(() => {
        markAllLossless();
        setActionMsg("-- MARKED CLIPS --");
      })
      .catch((err) => setActionMsg(err));
  });

  EventsOn("evt_unmark_all_lossless", () => {
    UnmarkAllLossless()
      .then(() => {
        unmarkAllLossless();
        setActionMsg("-- UNMARKED CLIPS --");
      })
      .catch((err) => setActionMsg(err));
  });

  onDestroy(() => {
    EventsOff(
      "evt_open_rename_clip_modal",
      "evt_rename_clip",
      "evt_track_move",
      "evt_open_search_list",
      "evt_yank_clip",
      "evt_insertclip_edit",
      "evt_zoom_timeline",
      "evt_saved_timeline",
      "evt_toggle_lossless",
      "evt_mark_all_lossless",
      "evt_unmark_all_lossless",
    );
    resetTrackStore();
    resetToolingStore();
  });
</script>

<div
  class="timeline h-full w-full bg-gdark border-t-2 border-t-white flex flex-col gap-4 pt-4 pb-4 px-1 relative overflow-x-scroll overflow-y-hidden"
  id="timeline"
  bind:this={timelineNode}
  use:dropzone={{}}
  on:mouseup={() => handleEditModeMouseUp()}
>
  <SearchList />
  <EventModal {isOpen} {close}>
    <div
      slot="header"
      class="font-semibold flex items-center justify-center gap-2"
    >
      <RenameIcon class="h-5 w-5 text-gyellow" />
      <h1>Rename Clip</h1>
    </div>
    <div slot="content">
      <div class="flex flex-col items-center justify-center gap-2">
        <p class="text-center font-semibold">
          rename video clip ({$videoNode.name})
        </p>
        <input
          type="text"
          bind:value={$videoNodeName}
          class="p-1 rounded-sm text-black"
          autocorrect="off"
          autocomplete="off"
        />
      </div>
    </div>
    <div slot="footer" class="flex items-center justify-center gap-2">
      <button
        class="flex items-center justify-center rounded-lg bg-gred1 font-semibold text-white px-4 py-1.5 hover:bg-gred transition ease-in-out duration-200 border-2 border-white gap-2"
        on:click={close}>Back</button
      >
      <button
        class="flex items-center justify-center rounded-lg bg-gblue0 font-semibold text-white px-4 py-1.5 hover:bg-gblue transition ease-in-out duration-200 border-2 border-white gap-2"
        on:click={() => {
          handleVideoNodeRename();
        }}
      >
        <span>Rename</span>
      </button>
    </div>
  </EventModal>

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
                (tVideo.end - tVideo.start) * $trackZoom < 120
                  ? 120
                  : (tVideo.end - tVideo.start) * $trackZoom
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
            class="h-full bg-obsbg border-white border-2 cursor-pointer select-none overflow-hidden flex flex-col justify-start"
            style={`width: ${
              (tVideo.end - tVideo.start) * $trackZoom < 120
                ? 120
                : (tVideo.end - tVideo.start) * $trackZoom
            }px; border-color: ${
              $videoNode && $videoNode.id === tVideo.id
                ? cursorColor
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
            <p>
              {tVideo.name}
            </p>
            <p>
              {formatSecondsToHMS(tVideo.end - tVideo.start)}
            </p>
            {#if tVideo.losslessexport}
              <span class="w-6 font-bold text-lg text-center text-gyellow">
                M
              </span>
            {/if}
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
