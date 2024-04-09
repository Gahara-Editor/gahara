<script lang="ts">
  import { writable, type Writable } from "svelte/store";
  import { EventsEmit } from "../../wailsjs/runtime/runtime";

  export let isOpen: Writable<boolean> = writable(false);
  export let close: () => void = () => {};

  function keydown(e: KeyboardEvent) {
    e.stopPropagation();
    if (e.key === "Escape") close();
    if (e.key === "Enter") EventsEmit("evt_rename_clip");
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
    node.focus();
    returnFn.push(() => {
      node.removeEventListener("keydown", keydown);
      node.removeEventListener("transitionend", transitionEnd);
    });
    return {
      destroy: () => returnFn.forEach((fn) => fn()),
    };
  }
</script>

{#if $isOpen}
  <div
    id="modal"
    use:modalAction
    class="fixed top-0 left-0 min-h-screen w-full z-10 flex justify-center items-center opacity-100"
    tabindex={0}
    autofocus
  >
    <div
      id="backdrop"
      class="absolute w-full h-full opacity-40"
      on:click={close}
    />
    <div
      id="content-wrap"
      class="z-10 max-w-[70vw] rounded-[0.3rem] bg-obsbg overflow-hidden md:max-w-[100vw] px-4 py-2 flex flex-col gap-4 border-white border-2"
    >
      <slot name="header">
        <!-- fallback -->
        <div>
          <h1>Main Menu</h1>
        </div>
      </slot>

      <div id="content" class="max-w-[50vh] overflow-hidden">
        <slot name="content" />
      </div>

      <slot name="footer">
        <!-- fallback -->
        <div>
          <button
            on:click={close}
            class="rounded-lg bg-red-500 font-semibold text-white inline-flex items-center px-4 py-1.5 hover:bg-red-700 transition ease-in-out duration-200"
            >close</button
          >
        </div>
      </slot>
    </div>
  </div>
{/if}
