import { derived, get, writable } from "svelte/store";
import type { main, video } from "../wailsjs/go/models";
import { InsertInterval } from "../wailsjs/go/main/App";
import {
  isPlaceholderItem,
  isVideo,
  isVideoItem,
  NODE_PLACEHOLDER,
  NODE_VIDEO,
  placeholderItem,
  type ListType,
} from "./lib/utils";
import { isNodeVideo } from "./lib/timeline";
import { handleKeybindTrackClipMove, type TimelineNode } from "./lib/timeline";

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

function createMainMenuStore() {
  const projects = writable<string[]>([]);
  const carouselIdx = writable<number>(0);

  const { set: setProjects, update: updateProjects } = projects;
  const { set: setCarouselIdx, update: updateCarouselIdx } = carouselIdx;

  function removeProject(projectName: string) {
    updateProjects((projects) =>
      projects.filter((project) => project !== projectName),
    );
  }

  function addProject(projectName: string) {
    updateProjects((projects) => {
      projects = [...projects, projectName];
      return projects;
    });
  }

  function resetMainMenuStore() {
    setCarouselIdx(0);
    setProjects([]);
  }

  return {
    projects,
    setProjects,
    addProject,
    removeProject,
    carouselIdx,
    setCarouselIdx,
    updateCarouselIdx,
    resetMainMenuStore,
  };
}

function createFilesytemStore() {
  const { subscribe, set, update } = writable<main.Video[]>([]);
  const videoFilesError = writable<string>("");
  const pipelineMessages = writable<string[]>([]);

  const { set: setVideoFilesError } = videoFilesError;
  const { set: setPipelineMsgs, update: updatePipelineMsgs } = pipelineMessages;

  function addPipelineMsg(msg: string) {
    updatePipelineMsgs((msgs) => (msgs = [...msgs, msg]));
  }

  function removePipelineMsg() {
    updatePipelineMsgs((msgs) => (msgs = [...msgs.slice(1)]));
  }

  function addVideos(videos: main.Video[]) {
    // TODO: handle duplicated keys
    update((projectFiles) => (projectFiles = [...projectFiles, ...videos]));
  }

  function removeVideoFile(fileName: string) {
    update(
      (projectFiles) =>
        (projectFiles = projectFiles.filter(
          (video) => video.name !== fileName,
        )),
    );
  }

  function searchFiles(query: string): main.Video[] {
    return get(videoFiles).filter((video) => {
      return video.name.toLowerCase().includes(query);
    });
  }

  const resetVideoFiles = () => {
    set([]);
    setPipelineMsgs([]);
    setVideoFilesError("");
  };

  return {
    subscribe,
    addVideos,
    videoFilesError,
    setVideoFilesError,
    pipelineMessages,
    addPipelineMsg,
    removePipelineMsg,
    removeVideoFile,
    searchFiles,
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
  const tracks = writable<TimelineNode[][]>([]);
  const trackTime = writable<number>(0.0);
  const trackDuration = writable<number>(0.0);
  const { subscribe, set, update } = tracks;
  const { set: setTrackDuration, update: updateTrackDuration } = trackDuration;
  const { set: setTrackTime } = trackTime;

  function addTrack() {
    update((tracks) => {
      tracks.push([]);
      return tracks;
    });
  }

  function removeTrack(tid: number) {
    update((tracks) => {
      if (tid >= 0 && tid < tracks.length) {
        tracks.splice(tid, 1);
      }
      return tracks;
    });
  }

  function addVideoToTrack(
    tid: number,
    pos: number,
    node: TimelineNode,
    mode: string = "none",
  ) {
    // TODO: handle duplicated keys
    update((tracks) => {
      if (tracks.length === 0 || tid > tracks.length) {
        tracks.push([node]);
      } else if (tid >= 0 && tid < tracks.length) {
        if (mode === "append" || tracks[tid].length === 0)
          tracks[tid] = [...tracks[tid], node];
        else {
          if (pos >= 0 && pos < tracks[tid].length) {
            tracks[tid].splice(pos, 0, node);
          }
        }
      }
      return tracks;
    });
    updateTrackDuration((tDuration) => (tDuration += node.end - node.start));
  }

  function removeAndAddIntervalToTrack(
    id: number,
    pos: number,
    nodes: TimelineNode[],
  ) {
    update((tracks) => {
      if (pos < 0 || pos > tracks[0].length) {
        return tracks;
      }
      tracks[id].splice(pos, 1, ...nodes);
      return tracks;
    });
  }

  function removeVideoFromTrack(id: number, node: TimelineNode) {
    update((tracks) => {
      if (!tracks[id]) return tracks;
      tracks[id] = tracks[id].filter((v) => v.id !== node.id);
      return tracks;
    });
    updateTrackDuration((tDuration) => (tDuration -= node.end - node.start));
  }

  function removeRIDReferencesFromTrack(id: number, rid: string) {
    let durationRemoved = 0;
    update((tracks) => {
      if (!tracks[id]) return tracks;
      tracks[id] = tracks[id].filter((v) => {
        if (v.rid === rid) {
          durationRemoved += v.end - v.start;
        }
        return v.rid !== rid;
      });
      return tracks;
    });
    updateTrackDuration((tDuration) => (tDuration -= durationRemoved));
  }

  function renameClipInTrack(id: number, pos: number, name: string) {
    update((tracks) => {
      if (!tracks[id]) return tracks;
      if (pos < 0 || pos > tracks[0].length) return tracks;
      tracks[id][pos].name = name;
      return tracks;
    });
  }

  function toggleLosslessMarkofClip(id: number, pos: number) {
    update((tracks) => {
      if (!tracks[id]) return tracks;
      if (pos < 0 || pos > tracks[0].length) return tracks;
      const node = tracks[id][pos];
      if (isNodeVideo(node)) node.losslessexport = !node.losslessexport;
      return tracks;
    });
  }

  function markAllLossless() {
    update((tracks) => {
      if (!tracks[0]) return tracks;
      for (let track of tracks[0]) {
        if (isNodeVideo(track)) track.losslessexport = true;
      }
      return tracks;
    });
  }

  function unmarkAllLossless() {
    update((tracks) => {
      if (!tracks[0]) return tracks;
      for (let track of tracks[0]) {
        if (isNodeVideo(track)) track.losslessexport = false;
      }
      return tracks;
    });
  }

  function searchTracks(query: string): video.VideoNode[] {
    const sTracks = get(tracks);
    if (sTracks.length <= 0) return [];

    return sTracks.reduce<video.VideoNode[]>((acc, track) => {
      const filteredNodes = track.filter(
        (node): node is video.VideoNode =>
          isNodeVideo(node) && node.name.toLowerCase().includes(query),
      );
      return acc.concat(filteredNodes);
    }, []);
  }

  function getClipPosInTrack(id: string): number {
    let sTracks = get(tracks);
    if (sTracks.length <= 0) return -1;
    return sTracks[0].findIndex((track) => track.id === id);
  }

  function resetTrackStore() {
    set([]);
    setTrackTime(0);
    setTrackDuration(0);
  }

  return {
    subscribe,
    trackTime,
    getClipPosInTrack,
    setTrackTime,
    addTrack,
    removeTrack,
    addVideoToTrack,
    removeVideoFromTrack,
    removeAndAddIntervalToTrack,
    removeRIDReferencesFromTrack,
    renameClipInTrack,
    toggleLosslessMarkofClip,
    markAllLossless,
    unmarkAllLossless,
    trackDuration,
    searchTracks,
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
  const playbackRate = writable<number>(1);

  const { set: setDuration } = duration;
  const { set: setCurrentTime, update: updateCurrenTime } = currentTime;
  const { set: setVolume } = volume;
  const { set: setVideoSrc } = source;
  const { set: setPaused } = paused;
  const { set: setEnded } = ended;
  const { set: setPlaybackRate, update: updatePlaybackRate } = playbackRate;

  function viewVideo(video: main.Video) {
    setVideoSrc(`${video.filepath}/${video.name}${video.extension}`);
  }

  function getDuration(): number {
    return get(duration);
  }

  function handlePlaybackRate(dir: string) {
    switch (dir) {
      case "down":
        if (get(playbackRate) / 2 < 0.5) return;
        updatePlaybackRate((playbackRate) => playbackRate / 2);
        break;
      default:
        if (get(playbackRate) * 2 > 2) return;
        updatePlaybackRate((playbackRate) => playbackRate * 2);
    }
  }

  function resetVideo() {
    setPlaybackRate(1);
    setDuration(0);
    setCurrentTime(0.0);
    setVolume(0.5);
    setVideoSrc("");
    setPaused(true);
    setEnded(false);
  }

  return {
    source,
    playbackRate,
    handlePlaybackRate,
    duration,
    getDuration,
    currentTime,
    paused,
    ended,
    viewVideo,
    setPlaybackRate,
    setVideoSrc,
    setDuration,
    setCurrentTime,
    updateCurrenTime,
    setVolume,
    resetVideo,
  };
}

function createVideoToolingStore() {
  const numberOfClipsInTrack = derived(trackStore, ($trackStore) => {
    if ($trackStore[getTrackCursorIdx()])
      return $trackStore[getTrackCursorIdx()].length;
    return 0;
  });

  // Messages
  const actionMessage = writable<string>("-- GAHARA --");
  const { set: setActionMsg } = actionMessage;

  // Modals
  const isOpenSearchList = writable<boolean>(false);
  const { set: setIsOpenSearchList } = isOpenSearchList;

  // Vim
  const clipCursorIdx = writable<number>(0);
  const trackCursorIdx = writable<number>(0);
  const clipRegister = writable<TimelineNode>(null);
  const cursorColor = writable<string>("#ffffff");
  const { set: setTrackCursorIdx, update: updateTrackCursorIdx } =
    trackCursorIdx;
  const { set: setClipCursorIdx, update: updateClipCursorIdx } = clipCursorIdx;
  const { set: setClipRegister } = clipRegister;

  // Edit modes
  const editMode = writable<string>("select");
  const vimMode = writable<boolean>(false);
  const { set: setVimMode, update: updateVimMode } = vimMode;
  const { set: setEditMode } = editMode;

  // Selected timeline node information
  const timelineNode = writable<TimelineNode>(null);
  const timelineNodePos = writable<number>(0);
  const timelineNodeWidth = writable<number>(1);
  const timelineNodeLeft = writable<number>(0);
  const timelineNodeName = writable<string>("");
  const { set: setTimelineNode } = timelineNode;
  const { set: setTimelineNodePos } = timelineNodePos;
  const { set: setTimelineNodeWidth } = timelineNodeWidth;
  const { set: setTimelineNodeLeft } = timelineNodeLeft;
  const { set: setTimelineNodeName } = timelineNodeName;

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
  const { set: setClipStart } = clipStart;
  const { set: setClipEnd } = clipEnd;
  const { set: moveCutRangeBox } = isMovingCutRangeBox;
  const { set: setBoxLeftBound } = boxLeftBound;
  const { set: setBoxRightBound } = boxRightBound;

  // Playback
  const playheadPos = writable<number>(0.0);
  const isMovingPlayhead = writable<boolean>(false);
  const isTrackPlaying = writable<boolean>(false);
  const trackZoom = writable<number>(20);
  const { set: movePlayhead } = isMovingPlayhead;
  const { set: setPlayheadPos, update: updatePlayheadPos } = playheadPos;
  const { set: setIsTrackPlaying, update: updateIsTrackPlaying } =
    isTrackPlaying;
  const { set: setTrackZoom, update: updateTrackZoom } = trackZoom;

  function getVimMode(): boolean {
    return get(vimMode);
  }

  function getCursorIdx(): number {
    return get(clipCursorIdx);
  }

  function getTrackCursorIdx(): number {
    return get(trackCursorIdx);
  }

  function moveTrackCursor(inc: number) {
    if (!get(vimMode)) return;
    const numTracks = get(trackStore).length;
    updateTrackCursorIdx((trackIdx) => {
      if (trackIdx + inc >= numTracks) return numTracks - 1;
      if (trackIdx + inc < 0) return 0;
      return trackIdx + inc;
    });
  }

  function moveClipCursor(inc: number) {
    if (!get(vimMode)) return;
    const numClips = get(numberOfClipsInTrack);
    updateClipCursorIdx((cursorIdx) => {
      if (cursorIdx + inc >= numClips) return numClips - 1;
      if (cursorIdx + inc < 0) return 0;
      return cursorIdx + inc;
    });
  }
  function adjustTrackZoom(dir: string) {
    switch (dir) {
      case "in":
        if (get(trackZoom) * 2 > 40) return;
        updateTrackZoom((trackZoom) => trackZoom * 2);
        break;

      default:
        if (get(trackZoom) / 2 < 5) return;
        updateTrackZoom((trackZoom) => trackZoom / 2);
    }
  }

  function resetToolingStore() {
    setActionMsg("-- GAHARA --");
    setIsOpenSearchList(false);
    setVimMode(false);
    setClipCursorIdx(0);
    setTrackCursorIdx(0);
    setClipRegister(null);
    setTimelineNode(null);
    setTimelineNodePos(0);
    setCutStart(0.0);
    setCutEnd(0.0);
    setClipStart(0.0);
    setClipEnd(0.0);
    moveCutRangeBox(false);
    setBoxLeftBound(0);
    setBoxRightBound(0);
    setPlayheadPos(0);
    setIsTrackPlaying(false);
    movePlayhead(false);
    setTimelineNodeName("");
    setEditMode("select");
  }

  return {
    actionMessage,
    setActionMsg,
    isOpenSearchList,
    setIsOpenSearchList,
    vimMode,
    getVimMode,
    setVimMode,
    updateVimMode,
    cursorColor,
    trackCursorIdx,
    setTrackCursorIdx,
    moveTrackCursor,
    clipCursorIdx,
    getCursorIdx,
    getTrackCursorIdx,
    setClipCursorIdx,
    moveClipCursor,
    clipRegister,
    setClipRegister,
    editMode,
    setEditMode,
    cutStart,
    setCutStart,
    cutEnd,
    setCutEnd,
    timelineNode,
    setTimelineNode,
    timelineNodeName,
    setTimelineNodeName,
    timelineNodePos,
    setTimelineNodePos,
    timelineNodeWidth,
    setTimelineNodeWidth,
    timelineNodeLeft,
    setTimelineNodeLeft,
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
    isTrackPlaying,
    setIsTrackPlaying,
    updateIsTrackPlaying,
    trackZoom,
    adjustTrackZoom,
    setTrackZoom,
    setPlayheadPos,
    updatePlayheadPos,
    isMovingPlayhead,
    isMovingCutRangeBox,
    moveCutRangeBox,
    resetToolingStore,
  };
}

function createExportOptionsStore() {
  const videoFormats = [
    ".mp4",
    ".mov",
    ".mkv",
    ".m4v",
    ".avi",
    ".wmv",
    ".flv",
    ".h264",
    ".hevc",
    ".3gp",
    ".3g2",
    ".ogv",
    ".264",
    ".265",
    ".webm",
  ];

  const videoCodecs = [
    ["H.264", "libx264"],
    ["H.264rgb", "libx264rgb"],
    ["H.265 (HEVC)", "libx265"],
    ["Lossless", "copy"],
    ["VP9", "libvpx-vp9"],
    ["Theora", "libtheora"],
  ];

  const crfOpts = ["18", "19", "20", "21", "22", "23"];
  const resolutionOpts = [
    "640x480",
    "1280x720",
    "1920x1080",
    "2560x1440",
    "3840x2160",
  ];
  const presetOpts = ["slow", "medium", "fast"];

  const filename = writable<string>("myvideo");
  const resolution = writable<string>("1920x1080");
  const codec = writable<string>("libx264");
  const videoFormat = writable<string>(".mp4");
  const preset = writable<string>("medium");
  const crf = writable<string>("18");
  const outputPath = writable<string>("");
  const isProcessingVid = writable<boolean>(false);
  const processingMsg = writable<string>("");
  const progressPercentage = writable<number>(0.0);
  const videoProcessingResults = writable<Array<any>>([]);
  const exportTabIndex = writable<number>(1);

  const { set: setFilename } = filename;
  const { set: setResolution } = resolution;
  const { set: setCodec } = codec;
  const { set: setVideoFormat } = videoFormat;
  const { set: setPreset } = preset;
  const { set: setCrf } = crf;
  const { set: setOutputPath } = outputPath;
  const { set: setIsProcessingVid } = isProcessingVid;
  const { set: setProcessingMsg } = processingMsg;
  const { set: setProgressPercentage } = progressPercentage;
  const {
    set: setVideoProcessingResults,
    update: updateVideoProcessingResults,
  } = videoProcessingResults;
  const { set: setExportTabIndex } = exportTabIndex;

  function getCompatibleCodecs(selectedFormat: string) {
    switch (selectedFormat) {
      case ".webm":
        return videoCodecs.filter((c) => c[0] === "VP9");
      case ".ogv":
        return videoCodecs.filter((c) => c[0] === "Theora");
      default:
        return videoCodecs.filter((c) => c[0] !== "VP9" && c[0] !== "Theora");
    }
  }

  function getExportOptions(): video.ProcessingOpts {
    const exportOpts: video.ProcessingOpts = {
      input_path: "",
      output_path: get(outputPath),
      filename: get(filename),
      video_format: get(videoFormat),
      codec: get(codec),
      resolution: get(resolution),
      preset: get(preset),
      crf: get(crf),
    };
    return exportOpts;
  }

  function addProcessingResult(result: any) {
    updateVideoProcessingResults((results) => [...results, result]);
  }

  function resetExportOptionsStore() {
    setFilename("myvideo");
    setResolution("1920x1080");
    setCodec("libx264");
    setVideoFormat(".mp4");
    setPreset("medium");
    setCrf("18");
    setOutputPath("");
    setIsProcessingVid(false);
    setProcessingMsg("");
    setProgressPercentage(0.0);
    setVideoProcessingResults([]);
    setExportTabIndex(1);
  }

  return {
    filename,
    setFilename,
    resolution,
    resolutionOpts,
    codec,
    setCodec,
    getCompatibleCodecs,
    videoFormat,
    videoFormats,
    videoCodecs,
    preset,
    presetOpts,
    crf,
    crfOpts,
    outputPath,
    setOutputPath,
    progressPercentage,
    setProgressPercentage,
    getExportOptions,
    setIsProcessingVid,
    isProcessingVid,
    processingMsg,
    setProcessingMsg,
    videoProcessingResults,
    setVideoProcessingResults,
    addProcessingResult,
    exportTabIndex,
    setExportTabIndex,
    resetExportOptionsStore,
  };
}

function createSearchListStore() {
  // searching
  let searchTerm = writable<string>("");
  let searchIdx = writable<number>(-1);
  let activeList = writable<ListType[]>([]);

  // duration
  let hours = writable<number>(0);
  let minutes = writable<number>(1);
  let seconds = writable<number>(0);

  const { set: setSearchTerm } = searchTerm;
  const { set: setSearchIdx, update: updateSearchIdx } = searchIdx;
  const { set: setActiveList } = activeList;
  const { set: setHours } = hours;
  const { set: setMinutes } = minutes;
  const { set: setSeconds } = seconds;

  function durationInSeconds(): number {
    return get(hours) * 3600 + get(minutes) * 60 + get(seconds);
  }

  function moveSearchIdx(inc: number) {
    const N = get(activeList).length;
    if (get(searchIdx) === -1) return;
    if (get(searchIdx) + inc < 0) setSearchIdx(N - 1);
    else if (get(searchIdx) + inc >= N) setSearchIdx(0);
    else setSearchIdx(get(searchIdx) + inc);
  }

  function search() {
    const query = get(searchTerm).toLowerCase();
    const commandRgx = /^\/([a-z])\s(.*)$/;
    const match = query.match(commandRgx);

    if (match) {
      switch (match[1]) {
        case "x":
          setActiveList(trackStore.searchTracks(match[2]));
          break;
        case "p":
          setActiveList([placeholderItem]);
          break;
        default:
      }
    } else {
      setActiveList(videoFiles.searchFiles(query));
    }
    if (get(activeList).length > 0) setSearchIdx(0);
  }

  async function executeAction() {
    const idx = get(searchIdx);
    const aList = get(activeList);

    if (idx >= 0 && idx < aList.length) {
      const item = aList[idx];
      if (isVideoItem(item)) {
        let nodeidx = trackStore.getClipPosInTrack(item.id);
        if (nodeidx !== -1) {
          toolingStore.setClipCursorIdx(nodeidx);
          handleKeybindTrackClipMove();
          resetSearchListStore();
          toolingStore.setIsOpenSearchList(false);
          toolingStore.setVimMode(true);
        }
      } else if (isVideo(item)) {
        videoStore.viewVideo(item);
        InsertInterval(
          get(toolingStore.trackCursorIdx),
          get(toolingStore.timelineNodePos),
          NODE_VIDEO,
          get(videoStore.source),
          item.name,
          0,
          item.duration,
        )
          .then((tVideo) => {
            trackStore.addVideoToTrack(
              get(toolingStore.trackCursorIdx),
              get(toolingStore.timelineNodePos),
              tVideo,
            );
            toolingStore.setTimelineNode(tVideo);
            toolingStore.setActionMsg("-- CLIP ADDED --");
            videoStore.setVideoSrc(tVideo.rid);
            videoStore.setCurrentTime(tVideo.start);
          })
          .catch(() =>
            toolingStore.setActionMsg(`could not insert ${aList[idx].name}`),
          );
      } else if (isPlaceholderItem(item)) {
        InsertInterval(
          get(toolingStore.trackCursorIdx),
          get(toolingStore.timelineNodePos),
          NODE_PLACEHOLDER,
          "",
          "placeholder",
          0,
          durationInSeconds(),
        )
          .then((node) => {
            trackStore.addVideoToTrack(
              get(toolingStore.trackCursorIdx),
              get(toolingStore.timelineNodePos),
              node,
            );
            toolingStore.setTimelineNode(node);
            toolingStore.setActionMsg("-- PLACEHOLDER ADDED --");
          })
          .catch(() =>
            toolingStore.setActionMsg(`could not insert ${aList[idx].name}`),
          );
      }
    }
  }

  function resetSearchListStore() {
    setSearchTerm("");
    setActiveList([]);
    setSearchIdx(-1);
  }

  return {
    hours,
    minutes,
    seconds,
    setHours,
    setMinutes,
    setSeconds,
    durationInSeconds,
    activeList,
    search,
    searchTerm,
    setSearchTerm,
    searchIdx,
    moveSearchIdx,
    setSearchIdx,
    updateSearchIdx,
    executeAction,
    resetSearchListStore,
  };
}

export const router = createRouterStore();
export const videoFiles = createFilesytemStore();
export const mainMenuStore = createMainMenuStore();
export const trackStore = createTracksStore();
export const projectName = writable("");
export const draggedVideo = createVideoTransferStore();
export const videoStore = createVideoStore();
export const toolingStore = createVideoToolingStore();
export const exportOptionsStore = createExportOptionsStore();
export const searchListstore = createSearchListStore();
