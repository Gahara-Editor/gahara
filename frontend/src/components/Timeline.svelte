<script lang="ts">
  import { dropzone } from "../lib/dnd";
  import { onDestroy } from "svelte";
  import type { video } from "../../wailsjs/go/models";
  import {
    videoStore,
    toolingStore,
    trackStore,
    createBooleanStore,
    searchListstore,
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
  import { NODE_AUDIO, NODE_PLACEHOLDER } from "../lib/utils";
  import VideoNode from "./VideoNode.svelte";
  import AudioNode from "./AudioNode.svelte";
  import PlaceholderNode from "./PlaceholderNode.svelte";
  import {
    scrollToNode,
    isNodeVideo,
    type TimelineNode,
    scrollToTrack,
  } from "../lib/timeline";

  const { isOpen, close, open } = createBooleanStore(false);
  const { setVideoSrc, currentTime, setCurrentTime } = videoStore;
  const {
    vimMode,
    cutStart,
    cutEnd,
    editMode,
    timelineNode,
    timelineNodePos,
    timelineNodeWidth,
    timelineNodeName,
    adjustTrackZoom,
    playheadPos,
    isMovingPlayhead,
    isMovingCutRangeBox,
    boxLeftBound,
    boxRightBound,
    cursorColor,
    trackCursorIdx,
    clipCursorIdx,
    clipRegister,
    setActionMsg,
    setIsOpenSearchList,
    setClipCursorIdx,
    setVimMode,
    moveTrackCursor,
    moveClipCursor,
    setClipRegister,
    moveCutRangeBox,
    movePlayhead,
    setCutEnd,
    setPlayheadPos,
    setTimelineNode,
    setTimelineNodeName,
    setTimelineNodePos,
    setTimelineNodeWidth,
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
  } = trackStore;
  const { setSearchTerm } = searchListstore;

  let trackNode: HTMLDivElement;
  let currentNode: HTMLDivElement;
  let cutRangeBox: HTMLDivElement;
  let cutRangeSide: "left" | "right" | "middle" | "none";

  $: {
    $cursorColor =
      $editMode === "remove"
        ? "#f7768e"
        : $vimMode && $editMode === "select"
        ? "#1abc9c"
        : "#ffffff";
  }

  function getTrackWidth() {
    const track = document.getElementById(`track-${$trackCursorIdx}`);
    return track.clientWidth;
  }

  function getTrackLeft() {
    const trackNode = document.getElementById(`track-${$trackCursorIdx}`);
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
      handletimelineNode(e, pos, tVideo);
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
        currentNode.scrollLeft -
        currentNode.getBoundingClientRect().left;

      const mousePos = Math.max(
        $boxLeftBound,
        Math.min(adjustedX, $boxRightBound),
      );

      setCurrentTime(
        $timelineNode.start +
          ($timelineNode.end - $timelineNode.start) *
            ((mousePos - $boxLeftBound) / $timelineNodeWidth),
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
    node: TimelineNode,
  ) {
    setBoxLeftBound(e.currentTarget.offsetLeft);
    setBoxRightBound(e.currentTarget.offsetLeft + e.currentTarget.clientWidth);
    setPlayheadPos(Math.min(e.currentTarget.offsetLeft, getTrackWidth()));
    setTimelineNodeWidth(e.currentTarget.getBoundingClientRect().width);
    setTimelineNodePos(pos);
    setClipCursorIdx(pos);
    setVideoSrc(node.rid);
    setTimelineNode(node);

    if (isNodeVideo(node)) {
      setCurrentTime(node.start);
      setClipStart(node.start);
      setClipEnd(node.end);
    }
  }

  function handletimelineNode(
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
      ((e.clientX + currentNode.scrollLeft) / getTrackWidth()) * $trackDuration;
    setVideoSrc(video.rid);
    setTimelineNodePos(pos);
    setTimelineNode(video);
    setCurrentTime(time);
    setTrackTime(trackTime);
    setCutEnd(time);
    setClipCursorIdx(pos);
    setBoxLeftBound(e.currentTarget.offsetLeft);
    setBoxRightBound(e.currentTarget.offsetLeft + e.currentTarget.clientWidth);
  }

  function handletimelineNodeRename() {
    if ($timelineNodeName === "") return;
    RenameVideoNode($trackCursorIdx, $timelineNodePos, $timelineNodeName)
      .then(() => {
        renameClipInTrack($trackCursorIdx, $timelineNodePos, $timelineNodeName);
        setTimelineNodeName("");
      })
      .catch(() => setActionMsg("could not rename clip"));
    setVimMode(true);
    close();
  }

  function handleKeybindTrackClipMove() {
    if (!$vimMode) return;
    const timelineNodeDiv = document
      .getElementById(`track-${$trackCursorIdx}`)
      ?.querySelector(`div:nth-child(${$clipCursorIdx + 1})`)
      ?.querySelector("div");
    if (timelineNodeDiv) {
      timelineNodeDiv.click();
      scrollToNode(timelineNodeDiv);
    }
  }

  function handleKeybindTrackMove() {
    const trackNode = document.getElementById(`track-${$trackCursorIdx}`);
    if (trackNode) scrollToTrack(trackNode);
  }

  EventsOn("evt_open_rename_clip_modal", () => {
    if ($timelineNode) {
      setVimMode(false);
      open();
    }
  });
  EventsOn("evt_rename_clip", () => {
    setVimMode(false);
    handletimelineNodeRename();
  });
  EventsOn("evt_track_move", (inc: number) => {
    moveTrackCursor(inc);
    moveClipCursor(0);
    handleKeybindTrackClipMove();
    handleKeybindTrackMove();
  });
  EventsOn("evt_clip_move", (inc: number) => {
    moveClipCursor(inc);
    handleKeybindTrackClipMove();
  });
  EventsOn("evt_open_search_list", () => {
    if ($vimMode) {
      setVimMode(false);
      setIsOpenSearchList(true);
    }
  });
  EventsOn("evt_search_timeline_clip", () => {
    if ($vimMode) {
      setVimMode(false);
      setSearchTerm("/x ");
      setIsOpenSearchList(true);
    }
  });
  EventsOn("evt_search_placeholder", () => {
    if ($vimMode) {
      setVimMode(false);
      setSearchTerm("/p ");
      setIsOpenSearchList(true);
    }
  });
  EventsOn("evt_yank_clip", () => {
    if ($timelineNode) {
      setClipRegister($timelineNode);
      setActionMsg(`YANKED: ${$timelineNode.name}`);
    }
  });
  EventsOn("evt_insertclip_edit", () => {
    if ($clipRegister) {
      InsertInterval(
        $trackCursorIdx,
        $timelineNodePos,
        $timelineNode.type,
        $clipRegister.rid,
        $clipRegister.name,
        $clipRegister.start,
        $clipRegister.end,
      )
        .then((tVideo) => {
          addVideoToTrack($trackCursorIdx, $timelineNodePos, tVideo);
          setTimelineNode(tVideo);
          setVideoSrc(tVideo.rid);
          setCurrentTime(tVideo.start);
          setActionMsg(`PASTED: ${$timelineNode.name}`);
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
    if ($timelineNode) {
      ToggleLossless(0, $timelineNodePos)
        .then(() => {
          toggleLosslessMarkofClip(0, $timelineNodePos);
        })
        .catch((err) => setActionMsg(err));
    }
  });

  EventsOn("evt_mark_all_lossless", () => {
    MarkAllLossless(0)
      .then(() => {
        markAllLossless();
        setActionMsg("-- MARKED CLIPS --");
      })
      .catch((err) => setActionMsg(err));
  });

  EventsOn("evt_unmark_all_lossless", () => {
    UnmarkAllLossless(0)
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
      "evt_clip_move",
      "evt_open_search_list",
      "evt_yank_clip",
      "evt_insertclip_edit",
      "evt_zoom_timeline",
      "evt_saved_timeline",
      "evt_toggle_lossless",
      "evt_mark_all_lossless",
      "evt_unmark_all_lossless",
      "evt_search_timeline_clip",
      "evt_search_placeholder",
    );
    resetToolingStore();
  });
</script>

<div
  class="timeline h-full w-full bg-gdark border-t-2 border-t-white flex flex-col gap-4 pt-4 pb-4 px-1 relative overflow-x-scroll overflow-y-hidden"
  id="timeline"
  bind:this={currentNode}
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
          rename video clip ({$timelineNode.name})
        </p>
        <input
          type="text"
          bind:value={$timelineNodeName}
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
          handletimelineNodeRename();
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

  <!-- TRACKS -->
  <!-- TODO Create an actual id for a track -->
  {#each $trackStore as track, id (id)}
    <div
      bind:this={trackNode}
      class={`h-36 flex relative ${
        $trackCursorIdx === id
          ? "border-gyellow border-[3px]"
          : "border-white border-0"
      } bg-gblue0 gap-1 p-2 w-max rounded-md`}
      id={`track-${id}`}
    >
      {#if track.length > 0}
        <!-- Nodes -->
        {#each track as node, pos (node.id)}
          <div
            animate:flip={{ duration: 100 }}
            out:slide={{ axis: "x", duration: 100, easing: cubicOut }}
          >
            {#if isNodeVideo(node)}
              <VideoNode
                {node}
                {pos}
                {cutRangeBox}
                {handleBoxRender}
                {handleEditModeMouseDown}
                {handleEditModeMouseMove}
              />
            {:else if node.type === NODE_AUDIO}
              <AudioNode {node} {pos} {cutRangeBox} />
            {:else if node.type === NODE_PLACEHOLDER}
              <PlaceholderNode {node} {pos} {cutRangeBox} {handleBoxRender} />
            {/if}
          </div>
        {/each}
      {:else}
        <span>T</span>
      {/if}
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
