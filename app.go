package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/k1nho/gahara/internal/utils"
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
	ctx context.Context
	// config: gahara configuration
	config Config
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
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
func (a *App) FilePicker() (Video, error) {
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
		return Video{}, err
	}

	fileName := utils.GetFilename(filepath)
	name, ext, err := utils.GetNameAndExtension(fileName)
	if err != nil {
		wruntime.LogError(a.ctx, err.Error())
		return Video{}, err
	}

	if !utils.IsValidExtension(ext) {
		wruntime.LogError(a.ctx, "invalid file extension")
		return Video{}, fmt.Errorf("invalid file extension")
	}

	go a.createProxyFile(filepath, name)
	return Video{Name: fileName, FilePath: path.Join(a.config.ProjectDir, fileName)}, nil

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
		return nil, fmt.Errorf("Gahara workspace exists, but no project workspace was found")
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
			if filepath.Ext(project.Name()) == ".png" {
				continue
			}
			projectFiles = append(projectFiles, Video{Name: project.Name(), FilePath: path.Join(a.config.ProjectDir, project.Name())})
		}
	}

	if len(projectFiles) <= 0 {
		wruntime.LogError(a.ctx, "Project workspace exists, but no files were found")
		return nil, fmt.Errorf("Project workspace exists, but no files were found")
	}

	wruntime.LogInfo(a.ctx, "project directories loaded successfully")
	return projectFiles, nil
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
