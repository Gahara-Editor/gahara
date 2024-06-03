<div align="center">
  <p><img src="./public/GaharaGithubIcon.svg" width="144" alt="GaharaVideoEditorIcon"/></p>
  <h1>🎬 Gahara 🎬</h1>
  <strong>A simple video editor</strong>
  <br>
</div>
<br>
<p align="center">
  <img src="https://img.shields.io/github/languages/code-size/Gahara-Editor/gahara" alt="GitHub code size in bytes">
  <a href="https://github.com/Gahara-Editor/gahara/issues" style="text-decoration:none;">
    <img src="https://img.shields.io/github/issues/Gahara-Editor/gahara" alt="GitHub issues">
  </a>
 <img src="https://img.shields.io/github/commit-activity/w/Gahara-Editor/gahara" alt="GitHub commit activity">
 <img src="https://github.com/Gahara-Editor/gahara/actions/workflows/tests.yml/badge.svg" alt="code coverage badge">
</p>

https://github.com/Gahara-Editor/gahara/assets/59541661/c9f07d7a-9e8e-4bb2-96b1-002f13764724

Gahara is a lightweight video editor app powered by the multimedia tool FFmpeg and built with the Wails framework. It provides a GUI to manage video projects and operate on multiple video clips with functions such as adding, removing, and trimming segments. Additionally, it has exporting features such as concatenating video clips (lossy), and support for converting videos to various formats supported by FFmpeg, such as .mp4, .avi, .wmv, and more.

## 🚀 Features

- **Add**, **remove**, and **cut** video clips
- **Video concatenation:** merge multiple video clips together (Lossy)
- **Video clip extraction:** cut and extract smaller segments from a larger video clip (Lossless)
- **Video format conversion:** transform the current format to another during export (.mp4, .avi, .wmv, etc)
- **Manage projects:** ability to create, and delete multiple video projects
- **Vim inspired keybinds:** delete, yank, paste, reorder, and move through the project timeline with Vim keybinds
- **Video clip labeling:** ability to rename video clips

## 📜 Requirements

- [FFmpeg](https://ffmpeg.org/download.html) >= 6.0
- [Node.js](https://nodejs.org/en/download)
- [Go](https://go.dev/dl/)
- [Wails](https://wails.io/docs/gettingstarted/installation)

## 📦 Installation

### Releases

**_Coming Soon_**

### Build from source

Make sure to have the [Requirements](#-requirements) installed.

1. Clone the repo

```bash
git clone git@github.com:Gahara-Editor/gahara.git
```

or

```bash
git clone https://github.com/Gahara-Editor/gahara.git
```

2. Move to the repo

```bash
cd gahara
```

3. Build the binary

```bash
wails build
```

4. Application will be under `build/bin/`
