<script lang="ts">
  import { toolingStore, searchListstore } from "../stores";
  import List from "./List.svelte";

  const { setVimMode, isOpenSearchList, setIsOpenSearchList, setActionMsg } =
    toolingStore;
  const {
    searchTerm,
    moveSearchIdx,
    search,
    executeAction,
    resetSearchListStore,
  } = searchListstore;

  function handleSearch(
    e: Event & {
      currentTarget: EventTarget & HTMLInputElement;
    },
  ) {
    e.stopPropagation();
    search();
  }

  async function keydown(e: KeyboardEvent) {
    e.stopPropagation();
    if (e.key === "Escape") {
      resetSearchListStore();
      setIsOpenSearchList(false);
      setVimMode(true);
    }
    if (e.key === "ArrowDown" || (e.key === "n" && e.ctrlKey)) moveSearchIdx(1);
    if (e.key === "ArrowUp" || (e.key === "p" && e.ctrlKey)) moveSearchIdx(-1);
    if (e.key === "Enter") {
      await executeAction()
        .then()
        .catch((e) => setActionMsg(e));
    }
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
    <List />
    <div
      class="flex flex-col items-center bg-obsbg min-w-[60vw] rounded border-white border-2 py-2 gap-2"
    >
      <h1 class="font-semibold">Search</h1>
      <input
        type="text"
        bind:value={$searchTerm}
        on:input={handleSearch}
        class="p-1 rounded-sm text-black"
        autocorrect="off"
        autocomplete="off"
      />
    </div>
  </div>
{/if}
