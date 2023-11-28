import { get, writable } from "svelte/store";
import type { main } from "../wailsjs/go/models";

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

function createTwoPageRouterStore() {
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
    // TODO: handle duplicated keys
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

function createVideoTransferStore() {
  const video = writable<main.Video>(null);
  const { set } = video;

  function value(): main.Video {
    return get(video);
  }

  function setDraggedVideo(vid: main.Video) {
    set(vid);
  }

  return {
    setDraggedVideo,
    value,
  };
}

export const video = createTwoPageRouterStore();
export const videoFiles = createFilesytemStore();
export const trackFiles = createFilesytemStore();
export const projectName = writable("");
export const currentVideo = writable("");
export const selectedTrack = writable("");
export const draggedVideo = createVideoTransferStore();
