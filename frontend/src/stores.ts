import { get, writable } from "svelte/store";
import type { main, video } from "../wailsjs/go/models";

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

function createRouterStore() {
  const route = writable<string>("main");
  const { set: setRoute } = route;

  const resetRouter = () => {
    setRoute("main");
  };

  return {
    setRoute,
    route,
    resetRouter,
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

  const resetVideoFiles = () => {
    set([]);
  };

  return {
    subscribe,
    addVideos,
    removeVideo,
    resetVideoFiles,
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

  function resetVideoTransfer() {
    set(null);
  }

  return {
    setDraggedVideo,
    resetVideoTransfer,
    value,
  };
}

function createTracksStore() {
  const tracks = writable<video.VideoNode[][]>([]);
  const trackTime = writable<number>(0.0);
  const trackDuration = writable<number>(0.0);

  const { subscribe, set, update } = tracks;
  const { set: setTrackDuration, update: updateTrackDuration } = trackDuration;
  const { set: setTrackTime } = trackTime;

  const addVideoToTrack = (id: number, video: video.VideoNode) => {
    // TODO: handle duplicated keys
    update((tracks) => {
      if (tracks.length === 0 || id > tracks.length) {
        tracks.push([video]);
      } else if (id >= 0 && id < tracks.length) {
        tracks[id] = [...tracks[id], video];
      }
      return tracks;
    });

    updateTrackDuration((tDuration) => (tDuration += video.end - video.start));
  };

  const removeAndAddIntervalToTrack = (
    id: number,
    pos: number,
    videoNodes: video.VideoNode[],
  ) => {
    update((tracks) => {
      if (pos < 0 || pos > tracks[0].length) {
        return tracks;
      }
      tracks[id].splice(pos, 1, ...videoNodes);
      return tracks;
    });
  };

  const removeVideoFromTrack = (id: number, videoNode: video.VideoNode) => {
    update((tracks) => {
      tracks[id] = tracks[0].filter((v) => v.id !== videoNode.id);
      return tracks;
    });
    updateTrackDuration(
      (tDuration) => (tDuration -= videoNode.end - videoNode.start),
    );
  };

  const resetTrackStore = () => {
    set([]);
    setTrackTime(0);
    setTrackDuration(0);
  };

  return {
    subscribe,
    trackTime,
    setTrackTime,
    addVideoToTrack,
    removeVideoFromTrack,
    removeAndAddIntervalToTrack,
    trackDuration,
    resetTrackStore,
  };
}

function createVideoStore() {
  const source = writable<string>("");
  const duration = writable<number>(0);
  const currentTime = writable<number>(0.0);
  const volume = writable<number>(0.5);
  const paused = writable<boolean>(true);
  const ended = writable<boolean>(false);

  const { set: setDuration } = duration;
  const { set: setCurrentTime } = currentTime;
  const { set: setVolume } = volume;
  const { set: setVideoSrc } = source;
  const { set: setPaused } = paused;
  const { set: setEnded } = ended;

  function viewVideo(video: main.Video) {
    setVideoSrc(`${video.filepath}/${video.name}${video.extension}`);
  }

  function getDuration(): number {
    return get(duration);
  }

  function resetVideo() {
    setDuration(0);
    setCurrentTime(0.0);
    setVolume(0.5);
    setVideoSrc("");
    setPaused(true);
    setEnded(false);
  }

  return {
    source,
    duration,
    getDuration,
    currentTime,
    paused,
    ended,
    viewVideo,
    setVideoSrc,
    setDuration,
    setCurrentTime,
    setVolume,
    resetVideo,
  };
}

function createVideoToolingStore() {
  // Edit modes
  const editMode = writable<string>("select");

  // Selected video information
  const videoNode = writable<video.VideoNode>(null);
  const videoNodePos = writable<number>(0);
  const videoNodeWidth = writable<number>(1);
  const videoNodeLeft = writable<number>(0);
  const { set: setVideoNode } = videoNode;
  const { set: setVideoNodePos } = videoNodePos;
  const { set: setVideoNodeWidth } = videoNodeWidth;
  const { set: setVideoNodeLeft } = videoNodeLeft;

  // Cut and range box operations
  const cutStart = writable<number>(0.0);
  const cutEnd = writable<number>(0.0);
  const clipStart = writable<number>(0.0);
  const clipEnd = writable<number>(0.0);
  const isMovingCutRangeBox = writable<boolean>(false);
  const boxLeftBound = writable<number>(0);
  const boxRightBound = writable<number>(0);
  const { set: setCutStart } = cutStart;
  const { set: setCutEnd } = cutEnd;
  const { set: setEditMode } = editMode;
  const { set: setClipStart } = clipStart;
  const { set: setClipEnd } = clipEnd;
  const { set: moveCutRangeBox } = isMovingCutRangeBox;
  const { set: setBoxLeftBound } = boxLeftBound;
  const { set: setBoxRightBound } = boxRightBound;

  // Playhead
  const playheadPos = writable<number>(0.0);
  const isMovingPlayhead = writable<boolean>(false);
  const { set: movePlayhead } = isMovingPlayhead;
  const { set: setPlayheadPos, update: updatePlayheadPos } = playheadPos;

  function resetToolingStore() {
    setVideoNode(null);
    setVideoNodePos(0);
    setCutStart(0.0);
    setCutEnd(0.0);
    setClipStart(0.0);
    setClipEnd(0.0);
    moveCutRangeBox(false);
    setBoxLeftBound(0);
    setBoxRightBound(0);
    setPlayheadPos(0);
    movePlayhead(false);
    setEditMode("select");
  }

  return {
    editMode,
    cutStart,
    setCutStart,
    cutEnd,
    setCutEnd,
    videoNode,
    setVideoNode,
    videoNodePos,
    setVideoNodePos,
    videoNodeWidth,
    setVideoNodeWidth,
    videoNodeLeft,
    setVideoNodeLeft,
    clipStart,
    setClipStart,
    clipEnd,
    setClipEnd,
    boxLeftBound,
    boxRightBound,
    setBoxLeftBound,
    setBoxRightBound,
    movePlayhead,
    playheadPos,
    setPlayheadPos,
    updatePlayheadPos,
    isMovingPlayhead,
    isMovingCutRangeBox,
    moveCutRangeBox,
    resetToolingStore,
  };
}

export const router = createRouterStore();
export const videoFiles = createFilesytemStore();
export const trackStore = createTracksStore();
export const projectName = writable("");
export const draggedVideo = createVideoTransferStore();
export const videoStore = createVideoStore();
export const toolingStore = createVideoToolingStore();
