<script lang="ts">
  import {
    CheckIcon,
    VideoCameraIcon,
    HandIcon,
    XIcon,
  } from "@rgossiaux/svelte-heroicons/solid";
  import { toolingStore, trackStore } from "./stores";
  import {
    AddTrack,
    RemoveInterval,
    RemoveTrack,
    SplitInterval,
  } from "../wailsjs/go/main/App";
  import TrashIcon from "./icons/TrashIcon.svelte";
  import SliceIntervalIcon from "./icons/SliceIntervalIcon.svelte";
  import BoxIntevalIcon from "./icons/BoxIntevalIcon.svelte";
  import { EventsEmit, EventsOff, EventsOn } from "../wailsjs/runtime/runtime";
  import { onDestroy } from "svelte";
  import VimIcon from "./icons/VimIcon.svelte";

  const {
    vimMode,
    updateVimMode,
    editMode,
    cutStart,
    cutEnd,
    timelineNodePos,
    timelineNode,
    trackCursorIdx,
    setEditMode,
    setClipRegister,
    setActionMsg,
  } = toolingStore;
  const {
    addTrack,
    removeTrack,
    removeVideoFromTrack,
    removeAndAddIntervalToTrack,
  } = trackStore;

  let executeEdit = false;

  $: {
    executeEdit = $editMode !== "select" ? true : false;
  }

  function handleAddTrack() {
    AddTrack().then(() => {
      addTrack();
    });
  }

  function handleRemoveTrack() {
    const tid = $trackCursorIdx;
    RemoveTrack(tid).then(() => {
      removeTrack(tid);
    });
  }

  function handleTwoCut() {
    if ($editMode === "timeline" && $timelineNode) {
      SplitInterval(
        $trackCursorIdx,
        $timelineNodePos,
        "evt_slice_cut",
        $timelineNode.start,
        $cutEnd,
      )
        .then((nodes) => {
          if (nodes.length > 0) {
            removeAndAddIntervalToTrack(0, $timelineNodePos, nodes);
            setTimeout(() => {
              EventsEmit("evt_clip_move", 1);
            }, 120);
          }
        })
        .catch(() => setActionMsg(`could not cut ${$timelineNode.name}`));
    }
  }

  function handleEditAction() {
    if ($timelineNode) {
      switch ($editMode) {
        case "intervalCut":
          SplitInterval(
            $trackCursorIdx,
            $timelineNodePos,
            "intervalCut",
            $cutStart,
            $cutEnd,
          )
            .then((nodes) => {
              removeAndAddIntervalToTrack(0, $timelineNodePos, nodes);
              setTimeout(() => {
                EventsEmit("evt_clip_move", 1);
              }, 120);
            })
            .catch(() => setActionMsg(`could not cut ${$timelineNode.name}`));
          break;
        case "remove":
          setClipRegister($timelineNode);
          RemoveInterval($trackCursorIdx, $timelineNodePos)
            .then(() => {
              removeVideoFromTrack($trackCursorIdx, $timelineNode);
              // timeout to catch up UI shifting on first element
              if ($timelineNodePos === 0) {
                setTimeout(() => {
                  EventsEmit("evt_clip_move", 0);
                }, 120);
              } else EventsEmit("evt_clip_move", -1);
            })
            .catch(() =>
              setActionMsg(`could not delete ${$timelineNode.name}`),
            );
          break;
      }
    }
    setEditMode("select");
    setActionMsg("-- SELECT --");
  }

  EventsOn("evt_toggle_vim_mode", () => {
    updateVimMode((mode) => !mode);
    if ($vimMode) setActionMsg("-- SELECT --");
    else setActionMsg("-- GAHARA --");
  });

  EventsOn("evt_change_vim_mode", (mode: string) => {
    setEditMode(mode);
    setActionMsg(`-- ${mode.toUpperCase()} --`);
  });
  EventsOn("evt_splitclip_edit", () => {
    handleTwoCut();
  });
  EventsOn("evt_execute_edit", () => {
    if ($vimMode) handleEditAction();
  });

  EventsOn("evt_add_track", () => {
    handleAddTrack();
  });
  EventsOn("evt_remove_track", () => {
    handleRemoveTrack();
  });
  onDestroy(() => {
    EventsOff(
      "evt_toggle_vim_mode",
      "evt_change_vim_mode",
      "evt_splitclip_edit",
      "evt_execute_edit",
      "evt_add_track",
      "evt_remove_track",
    );
  });
</script>

<div class="flex items-center p-2 justify-center gap-2">
  {#if executeEdit && $editMode !== "timeline"}
    <div
      class="flex items-center justify-center bg-gblue0 rounded-md border-white border-2 p-2 gap-2"
    >
      <button
        class="bg-ggreen px-2 py-1 rounded-md flex items-center border-2 border-white"
        on:click={() => handleEditAction()}
      >
        <CheckIcon class="h-5 w-5" />
      </button>
      <button
        class="bg-gred1 px-2 py-1 rounded-md flex items-center border-2 border-white"
        on:click={() => setEditMode("select")}
      >
        <XIcon class="h-5 w-5" />
      </button>
    </div>
  {/if}

  <div
    id="video-tooling"
    class="flex items-center justify-center bg-gblue0 rounded-md border-white border-2 p-2 gap-2"
  >
    <button
      class="bg-gdark px-2 py-1 rounded-md flex items-center gap-1 border-2 border-white"
      on:click={() => updateVimMode((mode) => !mode)}
    >
      <VimIcon class={$vimMode ? "h-5 w-5 text-teal" : "h-5 w-5"} />
    </button>
    <button
      class="bg-gdark px-3 py-0.5 rounded-md flex items-center gap-1 border-2 border-white"
      on:click={() => handleAddTrack()}
    >
      T
    </button>
    <button
      class="bg-gdark px-2 py-1 rounded-md flex items-center gap-1 border-2 border-white"
      style={$editMode === "select"
        ? "border-color: #eab308"
        : "border-color: white"}
      on:click={() => setEditMode("select")}
    >
      <HandIcon class="h-5 w-5 text-white" />
    </button>

    <button
      class="bg-gdark px-2 py-1 rounded-md flex items-center gap-1 border-2 border-white"
      disabled={$trackStore.length <= 0}
      style={$editMode === "timeline"
        ? "border-color: #eab308"
        : "border-color: white"}
      on:click={() => setEditMode("timeline")}
    >
      <VideoCameraIcon class="h-5 w-5 text-white" />
    </button>
    <button
      class="bg-gdark px-2 py-1 rounded-md flex items-center gap-1 border-2 border-white"
      disabled={$trackStore.length <= 0}
      on:click={() => handleTwoCut()}
    >
      <SliceIntervalIcon class="h-5 w-5 text-white" />
    </button>
    <button
      class="bg-gdark px-2 py-1 rounded-md flex items-center gap-1 border-2 border-white"
      style={$editMode === "intervalCut"
        ? "border-color: #eab308"
        : "border-color: white"}
      disabled={$trackStore.length <= 0}
      on:click={() => {
        setEditMode("intervalCut");
      }}
    >
      <BoxIntevalIcon class="h-5 w-5 text-white" />
    </button>

    <button
      class="bg-gdark px-2 py-1 rounded-md flex items-center gap-1 border-2 border-white"
      style={$editMode === "remove"
        ? "border-color: #eab308"
        : "border-color: white"}
      disabled={$trackStore.length <= 0}
      on:click={() => {
        setEditMode("remove");
      }}
    >
      <TrashIcon class="h-5 w-5 text-white" />
    </button>
  </div>
</div>
