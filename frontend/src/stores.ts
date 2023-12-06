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

function createVideoStore() {
  const duration = writable<number>(0);
  const currentTime = writable<number>(0.0);
  const volume = writable<number>(0.5);
  const paused = writable<boolean>(true);
  const ended = writable<boolean>(false);

  const { set: setDur } = duration;
  const { set: setCurT } = currentTime;
  const { set: setVol } = volume;

  function setDuration(val: number) {
    setDur(val);
  }

  function getDuration(): number {
    return get(duration);
  }

  function setCurrentTime(val: number) {
    setCurT(val);
  }

  function getCurrentTime(): number {
    return get(currentTime);
  }

  function getVolume(): number {
    return get(volume);
  }

  function setVolume(val: number) {
    setVol(val);
  }

  function reset() {
    setDur(0);
    setCurT(0.0);
    setVol(0.5);
  }

  return {
    duration,
    currentTime,
    paused,
    ended,
    setDuration,
    setCurrentTime,
    getDuration,
    getCurrentTime,
    getVolume,
    setVolume,
    reset,
  };
}

export const router = createTwoPageRouterStore();
export const videoFiles = createFilesytemStore();
export const trackFiles = createFilesytemStore();
export const projectName = writable("");
export const currentVideo = writable("");
export const selectedTrack = writable("");
export const draggedVideo = createVideoTransferStore();
export const videoStore = createVideoStore();
export const currenTime = writable(0.0);
