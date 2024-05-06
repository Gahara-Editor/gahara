<script lang="ts">
  import { videoStore, toolingStore, videoFiles, trackStore } from "../stores";
  import type { main } from "../../wailsjs/go/models";
  import { InsertInterval } from "../../wailsjs/go/main/App";

  const {
    videoNodePos,
    setVimMode,
    isOpenSearchList,
    setIsOpenSearchList,
    setVideoNode,
  } = toolingStore;
  const { addVideoToTrack } = trackStore;
  const { searchFiles } = videoFiles;
  const { source, setVideoSrc, setCurrentTime, getDuration, viewVideo } =
    videoStore;

  let searchTerm = "";
  let searchIdx = -1;
  let searchList: main.Video[] = [];

  function handleSearch(
    e: Event & {
      currentTarget: EventTarget & HTMLInputElement;
    },
  ) {
    e.stopPropagation();
    searchList = searchFiles(searchTerm);
    if (searchList && searchList.length > 0) searchIdx = 0;
  }

  function keydown(e: KeyboardEvent) {
    e.stopPropagation();
    if (e.key === "Escape") {
      searchIdx = -1;
      searchTerm = "";
      searchList = [];
      setIsOpenSearchList(false);
      setVimMode(true);
    }
    if (e.key === "ArrowDown" || (e.key === "n" && e.ctrlKey)) moveSearchIdx(1);
    if (e.key === "ArrowUp" || (e.key === "p" && e.ctrlKey)) moveSearchIdx(-1);
    if (e.key === "Enter") {
      if (searchIdx >= 0 && searchIdx < searchList.length) {
        viewVideo(searchList[searchIdx]);

        const videoDuration = getDuration();
        if (videoDuration === 0) return;

        InsertInterval(
          $source,
          searchList[searchIdx].name,
          0,
          videoDuration,
          $videoNodePos,
        )
          .then((tVideo) => {
            addVideoToTrack(0, tVideo, $videoNodePos);
            setVideoNode(tVideo);
            setVideoSrc(tVideo.rid);
            setCurrentTime(tVideo.start);
          })
          .catch(() =>
            toolingStore.setActionMsg(
              `could not insert ${searchList[searchIdx].name}`,
            ),
          );
      }
    }
  }

  function moveSearchIdx(inc: number) {
    if (searchIdx === -1) return;
    if (searchIdx + inc < 0) searchIdx = searchList.length - 1;
    else if (searchIdx + inc >= searchList.length) searchIdx = 0;
    else searchIdx = searchIdx + inc;
  }

  function transitionEnd(e: TransitionEvent) {
    const node = e.target as HTMLElement;
    node.focus();
  }

  function modalAction(node: HTMLElement) {
    const returnFn = [];
    // for accessibility
    if (document.body.style.overflow !== "hidden") {
      const original = document.body.style.overflow;
      document.body.style.overflow = "hidden";
      returnFn.push(() => {
        document.body.style.overflow = original;
      });
    }
    node.addEventListener("keydown", keydown);
    node.addEventListener("transitionend", transitionEnd);

    const input = node.querySelector("input");
    if (input) input.focus();

    returnFn.push(() => {
      node.removeEventListener("keydown", keydown);
      node.removeEventListener("transitionend", transitionEnd);
    });
    return {
      destroy: () => returnFn.forEach((fn) => fn()),
    };
  }
</script>

{#if $isOpenSearchList}
  <div
    id="modal"
    use:modalAction
    class="fixed top-0 left-0 min-h-screen w-full z-10 flex flex-col justify-center items-center opacity-100 gap-2"
    tabindex={0}
  >
    <div
      id="backdrop"
      class="absolute w-full h-full opacity-40"
      on:click={() => setIsOpenSearchList(false)}
    />
    <div
      id="content-wrap"
      class="z-10 bg-gblue0 min-h-[40vh] max-h-[40vh] min-w-[60vw] max-w-[60vw] rounded border-white border-2 p-2 overflow-y-scroll"
    >
      <div id="video-clip-list">
        <h1 class="text-center font-semibold text-lg">Results</h1>
        <ul class="divide-y divide-gray-200 text-white">
          {#each searchList as video, i}
            {#if i === searchIdx}
              <li class=" bg-obsternary p-2 truncate rounded-md">
                <span class="text-teal font-semibold">></span>
                {video.name}
              </li>
            {:else}
              <li
                class=" bg-obsbg p-2 truncate rounded-md"
                on:click={() => (searchIdx = i)}
              >
                {video.name}
              </li>
            {/if}
          {/each}
        </ul>
      </div>
    </div>
    <div
      class="flex flex-col items-center bg-obsbg min-w-[60vw] rounded border-white border-2 py-2 gap-2"
    >
      <h1 class="font-semibold">Find Clips</h1>
      <input
        type="text"
        bind:value={searchTerm}
        on:input={handleSearch}
        class="p-1 rounded-sm text-black"
        autocorrect="off"
        autocomplete="off"
      />
    </div>
  </div>
{/if}
