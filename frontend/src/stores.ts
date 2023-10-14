import { writable } from "svelte/store";
import { main } from "../wailsjs/go/models";

export function createBooleanStore(initial: boolean) {
  const isOpen = writable(initial);
  const { set, update } = isOpen;
  return {
    isOpen,
    open: () => set(true),
    close: () => set(false),
    toggle: () => update((n) => !n),
  };
}

function createVideoStore() {
  const { subscribe, set, update } = writable(false);

  const setVideoLayoutView = () => {
    set(true);
  };

  const setMainMenuView = () => {
    set(false);
  };

  return {
    subscribe,
    update,
    setVideoLayoutView,
    setMainMenuView,
  };
}

function createFilesytemStore() {
  const { subscribe, set, update } = writable<main.Video[]>([]);

  const addVideos = (videos: main.Video[]) => {
    update((projectFiles) => (projectFiles = [...projectFiles, ...videos]));
  };

  const removeVideo = (fileName: string) => {
    update(
      (projectFiles) =>
        (projectFiles = projectFiles.filter(
          (video) => video.name !== fileName,
        )),
    );
  };

  const reset = () => {
    set([]);
  };

  return {
    subscribe,
    addVideos,
    removeVideo,
    reset,
  };
}

export const video = createVideoStore();
export const videoFiles = createFilesytemStore();
export const projectName = writable("");
