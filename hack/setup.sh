#!/bin/bash

PLATFORM=$1

if [[ "$PLATFORM" == "darwin/universal" ]]; then
  echo "Setting up FFmpeg for macOS"
  curl -O https://evermeet.cx/ffmpeg/ffmpeg-115960-g3a5202d026.zip
  unzip ffmpeg-115960-g3a5202d026.zip
  mkdir -p resources/darwin/
  mv ffmpeg resources/darwin/ffmpeg
  chmod +x resources/darwin/ffmpeg
  rm -rf ffmpeg-115960-g3a5202d026.zip
elif [[ "$PLATFORM" == "linux/amd64" ]]; then
  echo "Setting up FFmpeg for Linux"
elif [[ "$PLATFORM" == "windows" ]]; then
  echo "Setting up FFmpeg for Windows"
else
  echo "Unsupported platform: $PLATFORM"
  exit 1
fi
