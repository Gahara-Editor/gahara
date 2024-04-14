package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/k1nho/gahara/internal/video"
	"github.com/wailsapp/wails/v2/pkg/menu"
	"github.com/wailsapp/wails/v2/pkg/menu/keys"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	wruntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

const (
	Success = "success"
	Failed  = "failed"
)

type Config struct {
	// GaharaDir: the workspace directory for gahara
	GaharaDir string `json:"gaharadir"`
	//ProjectDir: the project directory for a video editing project
	ProjectDir string `json:"projectdir,omitempty"`
}

// App struct
type App struct {
	// ctx: app context
	ctx context.Context
	// config: gahara configuration
	config Config
	// Timeline: the project timeline
	Timeline video.Timeline `json:"timeline"`
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{Timeline: video.NewTimeline()}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	err := a.gaharaSetup()
	if err != nil {
		wruntime.LogFatal(ctx, err.Error())
	}
}

// FilePicker: opens the native file picker for the user
// this spawns a proxy file creation, if the file is valid
func (a *App) FilePicker() error {
	fileFilter := wruntime.FileFilter{
		DisplayName: "Video Files(*.mov, *.mp4, *.mkv)",
		Pattern:     "*.mov;*.mp4;*.mkv;*.avi;*.wmv;*.webm;*.avchd",
	}

	openDialogOpts := wruntime.OpenDialogOptions{
		Title:   "Select File",
		Filters: []wruntime.FileFilter{fileFilter},
	}

	filepath, err := wruntime.OpenFileDialog(a.ctx, openDialogOpts)
	if err != nil {
		wruntime.LogError(a.ctx, err.Error())
		return err
	}

	go a.createProxyFile(filepath)
	return nil

}

// createWorkspace: creates a workspace directory to store the video projects of a user locally
func (a *App) createWorkspace() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		wruntime.LogFatal(a.ctx, "user home directory does not exists, cannot create workspace for Gahara")
		return "", fmt.Errorf("user home directory does not exists, cannot create workspace for Gahara")
	}

	gaharaDir := path.Join(homeDir, ".gahara")

	_, err = os.Stat(gaharaDir)
	// check if the gahara directory does not exists
	if os.IsNotExist(err) {
		if err := os.MkdirAll(gaharaDir, os.ModePerm); err != nil {
			wruntime.LogError(a.ctx, err.Error())
			return "", err
		}
		wruntime.LogInfo(a.ctx, "Gahara workspace has been created!")
	} else if err != nil {
		wruntime.LogError(a.ctx, "could not create gahara workspace")
		return "", err
	}
	return gaharaDir, nil

}

// createProjectWorkspace: creates a project directory to store the videos related to a project locally
func (a *App) CreateProjectWorkspace(projectName string) (string, error) {
	// create project workspace
	projectDir := path.Join(a.config.GaharaDir, projectName)
	file, err := os.Stat(projectDir)
	if os.IsNotExist(err) {
		if err := os.MkdirAll(projectDir, os.ModePerm); err != nil {
			wruntime.LogError(a.ctx, err.Error())
			return Failed, err
		}
		wruntime.LogInfo(a.ctx, fmt.Sprintf("project %s workspace has been created", projectName))
	} else if err != nil {
		wruntime.LogError(a.ctx, fmt.Sprintf("could not create the %s workspace\n", projectName))
		return Failed, err
	}

	if file != nil {
		wruntime.LogError(a.ctx, fmt.Sprintf("project name (%s) already exists in gahara workspace", projectName))
		return Failed, fmt.Errorf("project name (%s) already exists in gahara workspace", projectName)
	}

	a.config.ProjectDir = projectDir
	return Success, nil
}

// SetProjectDirectory: sets the project directory (used with loading projects)
func (a *App) SetProjectDirectory(projectDir string) {
	a.config.ProjectDir = path.Join(a.config.GaharaDir, projectDir)
}

// ReadGaharaWorkspace: retrieve all the project workspaces
func (a *App) ReadGaharaWorkspace() ([]string, error) {
	gaharaDir, err := os.Open(a.config.GaharaDir)
	if err != nil {
		wruntime.LogError(a.ctx, "could not read the gahara workspace")
		return nil, err
	}
	defer gaharaDir.Close()

	projects, err := gaharaDir.Readdir(0)
	if err != nil {
		wruntime.LogError(a.ctx, "could not retrieve projects in the gahara workspace")
		return nil, err
	}

	projectsDirectories := []string{}
	for _, project := range projects {
		if project.IsDir() {
			projectsDirectories = append(projectsDirectories, project.Name())
		}
	}

	if len(projectsDirectories) <= 0 {
		wruntime.LogError(a.ctx, "Gahara workspace exists, but no project workspace was found")
		return nil, fmt.Errorf("gahara workspace exists, but no project workspace was found")
	}

	wruntime.LogInfo(a.ctx, "project directories loaded successfully")
	return projectsDirectories, nil
}

// ReadProjectWorkspace: retrieve all the files in the project workspace
func (a *App) ReadProjectWorkspace() ([]Video, error) {
	projectDir, err := os.Open(a.config.ProjectDir)
	if err != nil {
		wruntime.LogError(a.ctx, "could not read the gahara workspace")
		return nil, err
	}
	defer projectDir.Close()

	files, err := projectDir.Readdir(0)
	if err != nil {
		wruntime.LogError(a.ctx, "could not retrieve projects in the gahara workspace")
		return nil, err
	}

	projectFiles := []Video{}
	for _, project := range files {
		if !project.IsDir() {
			if !video.IsValidExtension(filepath.Ext(project.Name())) {
				continue
			}
			projectFiles = append(projectFiles, Video{Name: strings.Split(project.Name(), ".")[0], Extension: filepath.Ext(project.Name()), FilePath: a.config.ProjectDir})
		}
	}

	if len(projectFiles) <= 0 {
		wruntime.LogError(a.ctx, "Project workspace exists, but no files were found")
		return nil, fmt.Errorf("project workspace exists, but no files were found")
	}

	wruntime.LogInfo(a.ctx, "project files loaded successfully")
	return projectFiles, nil
}

// DeleteProject: deletes a video project and all of its related files
func (a *App) DeleteProject(name string) error {
	path := path.Join(a.config.GaharaDir, name)
	err := os.RemoveAll(path)
	if err != nil {
		return fmt.Errorf("could not delete the project")
	}
	return nil
}

// DeleteProjectFile: delete a project file (root id form)
func (a *App) DeleteProjectFile(rid string) error {
	_, err := os.Stat(rid)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("file %s  does not exists", video.GetFilename(rid))
		}
		return err
	}

	err = os.Remove(rid)
	if err != nil {
		return err
	}
	return nil
}

func (a *App) AppMenu(menuItems ...menu.MenuItem) *menu.Menu {
	return menu.NewMenuFromItems(menu.AppMenu(), menu.EditMenu(), menu.WindowMenu())
}

// SetDefaultAppMenu: it sets the default menu
func (a *App) SetDefaultAppMenu() {
	runtime.MenuSetApplicationMenu(a.ctx, a.AppMenu())
	runtime.MenuUpdateApplicationMenu(a.ctx)
}

// EnableVideoMenus: It enables the menus for the video layout view
func (a *App) EnableVideoMenus() {
	timelineMenu := menu.NewMenu()
	timelineMenu.AddText("Open File", keys.CmdOrCtrl("o"), func(cd *menu.CallbackData) {
		wruntime.EventsEmit(a.ctx, video.EVT_UPLOAD_FILE)
	})

	timelineMenu.AddText("Play Track", keys.Shift("space"), func(cd *menu.CallbackData) {
		wruntime.EventsEmit(a.ctx, video.EVT_PLAY_TRACK)
	})
	timelineMenu.AddText("Save Timeline", keys.CmdOrCtrl("s"), func(cd *menu.CallbackData) {
		err := a.SaveTimeline()
		if err != nil {
			wruntime.LogError(a.ctx, "could not save timeline")
		}
		wruntime.EventsEmit(a.ctx, video.EVT_SAVED_TIMELINE, "-- SAVED --")
	})
	timelineMenu.AddText("Rename clip", keys.CmdOrCtrl("r"), func(cd *menu.CallbackData) {
		wruntime.EventsEmit(a.ctx, video.EVT_OPEN_RENAME_CLIP_MODAL)
	})
	timelineMenu.AddText("Toggle Vim Mode", keys.CmdOrCtrl("i"), func(cd *menu.CallbackData) {
		wruntime.EventsEmit(a.ctx, video.EVT_TOGGLE_VIM_MODE)
		wruntime.EventsEmit(a.ctx, video.EVT_TRACK_MOVE, 0)
	})
	vimCommandsMenu := timelineMenu.AddSubmenu("Vim Commands")
	vimCommandsMenu.AddText("Normal Mode", keys.Key("i"), func(cd *menu.CallbackData) {
		wruntime.EventsEmit(a.ctx, video.EVT_CHANGE_VIM_MODE, "select")
	})
	vimCommandsMenu.AddText("Delete Mode", keys.Key("d"), func(cd *menu.CallbackData) {
		wruntime.EventsEmit(a.ctx, video.EVT_CHANGE_VIM_MODE, "remove")
	})
	vimCommandsMenu.AddText("Timeline Mode", keys.Key("t"), func(cd *menu.CallbackData) {
		wruntime.EventsEmit(a.ctx, video.EVT_CHANGE_VIM_MODE, "timeline")
	})
	vimCommandsMenu.AddText("Split clip", keys.Key("x"), func(cd *menu.CallbackData) {
		wruntime.EventsEmit(a.ctx, video.EVT_SPLITCLIP_EDIT)
	})
	vimCommandsMenu.AddText("Yank clip", keys.Key("y"), func(cd *menu.CallbackData) {
		wruntime.EventsEmit(a.ctx, video.EVT_YANK_CLIP)
	})
	vimCommandsMenu.AddText("Paste clip", keys.Key("p"), func(cd *menu.CallbackData) {
		wruntime.EventsEmit(a.ctx, video.EVT_INSERTCLIP_EDIT)
	})
	vimCommandsMenu.AddText("Execute Edit", keys.Key("enter"), func(cd *menu.CallbackData) {
		wruntime.EventsEmit(a.ctx, video.EVT_EXECUTE_EDIT)
	})
	vimCommandsMenu.AddText("Move Track Left", keys.Key("h"), func(cd *menu.CallbackData) {
		wruntime.EventsEmit(a.ctx, video.EVT_TRACK_MOVE, -1)
	})
	vimCommandsMenu.AddText("Move Track Right", keys.Key("l"), func(cd *menu.CallbackData) {
		wruntime.EventsEmit(a.ctx, video.EVT_TRACK_MOVE, 1)
	})
	vimCommandsMenu.AddText("Move to Beginning of Track", keys.Key("0"), func(cd *menu.CallbackData) {
		wruntime.EventsEmit(a.ctx, video.EVT_TRACK_MOVE, -len(a.Timeline.VideoNodes))
	})
	vimCommandsMenu.AddText("Move to End of Track", keys.Key("$"), func(cd *menu.CallbackData) {
		wruntime.EventsEmit(a.ctx, video.EVT_TRACK_MOVE, len(a.Timeline.VideoNodes))
	})
	vimCommandsMenu.AddText("Zoom In Timeline", keys.Shift("+"), func(cd *menu.CallbackData) {
		wruntime.EventsEmit(a.ctx, video.EVT_ZOOM_TIMELINE, "in")
	})
	vimCommandsMenu.AddText("Zoom Out Timeline", keys.Shift("-"), func(cd *menu.CallbackData) {
		wruntime.EventsEmit(a.ctx, video.EVT_ZOOM_TIMELINE, "out")
	})
	vimCommandsMenu.AddText("Save Timeline", keys.Shift("w"), func(cd *menu.CallbackData) {
		err := a.SaveTimeline()
		if err != nil {
			wruntime.LogError(a.ctx, "could not save timeline")
		}
		wruntime.EventsEmit(a.ctx, video.EVT_SAVED_TIMELINE, "-- SAVED --")
	})
	vimCommandsMenu.AddText("Open Search List", keys.Key("/"), func(cd *menu.CallbackData) {
		wruntime.EventsEmit(a.ctx, video.EVT_OPEN_SEARCH_LIST)
	})

	appMenu := a.AppMenu()
	appMenu.Items = append(appMenu.Items, &menu.MenuItem{
		Label:   "Timeline",
		SubMenu: timelineMenu,
	})
	runtime.MenuSetApplicationMenu(a.ctx, appMenu)
	runtime.MenuUpdateApplicationMenu(a.ctx)
}

// GaharaSetup: setup of gahara on startup (workspace and config.json)
func (a *App) gaharaSetup() error {
	gaharaDir, err := a.createWorkspace()
	if err != nil {
		return err
	}
	wruntime.LogInfo(a.ctx, "Gahara workspace has been found!")

	// config.json
	configPath := path.Join(gaharaDir, "config.json")
	// check that exists first
	_, err = os.Stat(configPath)
	if os.IsNotExist(err) {
		// create config file
		file, err := os.Create(configPath)
		if err != nil {
			return err
		}
		defer file.Close()

		gaharaConfig := Config{
			GaharaDir: gaharaDir,
		}

		bytes, err := json.MarshalIndent(gaharaConfig, "", "\t")
		if err != nil {
			wruntime.LogError(a.ctx, "could not marshal Config struct")
			return err
		}

		_, err = file.Write(bytes)
		if err != nil {
			wruntime.LogError(a.ctx, "could not write bytes into json file")
			return err
		}

		wruntime.LogInfo(a.ctx, "config.json file for gahara has been created!")

	} else if err != nil {
		wruntime.LogError(a.ctx, fmt.Sprintf("could not setup gahara: %s\n", err.Error()))
		return err

	}

	// the file exists, read it into the struct
	bytes, err := os.ReadFile(configPath)
	if err != nil {
		wruntime.LogError(a.ctx, "could not read the config file path")
		return err
	}

	err = json.Unmarshal(bytes, &a.config)
	if err != nil {
		wruntime.LogError(a.ctx, "could not unmarshal the config")
		return err
	}

	wruntime.LogInfo(a.ctx, "config.json file has been found!")
	return nil
}
