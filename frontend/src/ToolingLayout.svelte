<script lang="ts">
  import {
    CheckIcon,
    VideoCameraIcon,
    HandIcon,
    XIcon,
  } from "@rgossiaux/svelte-heroicons/solid";
  import { toolingStore, trackStore } from "./stores";
  import { RemoveInterval, SplitInterval } from "../wailsjs/go/main/App";
  import TrashIcon from "./icons/TrashIcon.svelte";
  import SliceIntervalIcon from "./icons/SliceIntervalIcon.svelte";
  import BoxIntevalIcon from "./icons/BoxIntevalIcon.svelte";

  const { editMode, cutStart, cutEnd, videoNodePos, videoNode } = toolingStore;
  const { removeVideoFromTrack, removeAndAddIntervalToTrack } = trackStore;

  let executeEdit = false;

  $: {
    executeEdit = $editMode !== "select" ? true : false;
  }

  function setEditMode(mode: string) {
    editMode.set(mode);
  }

  function handleTwoCut() {
    if ($editMode === "timeline" && $videoNode) {
      SplitInterval("evt_slice_cut", $videoNodePos, $videoNode.start, $cutEnd)
        .then((nodes) => {
          if (nodes.length > 0) {
            removeAndAddIntervalToTrack(0, $videoNodePos, nodes);
          }
        })
        .catch(console.log);
    }
  }

  function handleEditAction() {
    switch ($editMode) {
      case "intervalCut":
        if ($videoNode) {
          SplitInterval("intervalCut", $videoNodePos, $cutStart, $cutEnd)
            .then((nodes) => {
              removeAndAddIntervalToTrack(0, $videoNodePos, nodes);
            })
            .catch(console.log);
        }
        setEditMode("select");
        break;
      case "remove":
        if ($videoNode) {
          RemoveInterval($videoNodePos)
            .then(() => {
              removeVideoFromTrack(0, $videoNode);
            })
            .catch(console.log);
        }
        setEditMode("select");
        break;
      default:
    }
  }
</script>

<div class="w-full flex items-center p-2 justify-center">
  <div
    id="video-tooling"
    class="flex items-center justify-center bg-gblue0 rounded-md border-white border-2 p-2 gap-2"
  >
    {#if executeEdit && $editMode !== "timeline"}
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
    {/if}
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
