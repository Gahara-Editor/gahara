import { draggedVideo, trackStore, videoStore, toolingStore } from "../stores";
import type { main } from "../../wailsjs/go/models";
import { InsertInterval } from "../../wailsjs/go/main/App";
import { NODE_VIDEO } from "./utils";

export function draggable(node: HTMLDivElement, data: main.Video) {
  let state = data;

  node.draggable = true;
  node.style.cursor = "grab";

  function handleDragStart(
    _: DragEvent & {
      currentTarget: EventTarget & HTMLDivElement;
    },
  ) {
    draggedVideo.setDraggedVideo(data);
    videoStore.viewVideo(draggedVideo.value());
  }

  node.addEventListener("dragstart", handleDragStart);

  return {
    update(data: main.Video) {
      state = data;
    },
    destroy() {
      node.removeEventListener("dragstart", handleDragStart);
    },
  };
}

export function dropzone(node: HTMLDivElement, opts) {
  let state = {
    dropEffect: "move",
    dragOverClass: "droppable",
    ...opts,
  };

  function handleDragEnter(
    e: DragEvent & {
      currentTarget: EventTarget & HTMLDivElement;
    },
  ) {
    e.currentTarget.classList.add(state.dragOverClass);
  }

  function handleDragLeave(
    e: DragEvent & { currentTarget: EventTarget & HTMLDivElement },
  ) {
    e.currentTarget.classList.remove(state.dragOverClass);
  }

  function handleDragOver(
    e: DragEvent & { currentTarget: EventTarget & HTMLDivElement },
  ) {
    e.preventDefault();

    e.dataTransfer.dropEffect = state.dropEffect;
  }

  function handleDrop(
    e: DragEvent & { currentTarget: EventTarget & HTMLDivElement },
  ) {
    e.preventDefault();
    e.currentTarget.classList.remove(state.dragOverClass);
    const videoRID = `${draggedVideo.value().filepath}/${
      draggedVideo.value().name
    }${draggedVideo.value().extension}`;

    // TODO: handle insertions at an specific part of the timeline
    InsertInterval(
      0,
      0,
      NODE_VIDEO,
      videoRID,
      draggedVideo.value().name,
      0,
      draggedVideo.value().duration,
    )
      .then((tVideo) => {
        // TODO: add to different tracks dynamically for now 0
        trackStore.addVideoToTrack(0, 0, tVideo);
        toolingStore.setVideoNode(tVideo);
        videoStore.setCurrentTime(tVideo.start);
      })
      .catch(() =>
        toolingStore.setActionMsg(
          `could not insert ${draggedVideo.value().name}`,
        ),
      );
  }

  node.addEventListener("dragenter", handleDragEnter);
  node.addEventListener("dragleave", handleDragLeave);
  node.addEventListener("dragover", handleDragOver);
  node.addEventListener("drop", handleDrop);

  return {
    destroy() {
      node.removeEventListener("dragenter", handleDragEnter);
      node.removeEventListener("dragleave", handleDragLeave);
      node.removeEventListener("dragover", handleDragOver);
      node.removeEventListener("drop", handleDrop);
    },
  };
}
