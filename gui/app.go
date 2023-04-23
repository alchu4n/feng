package main

import (
	"context"

	"github.com/ac0d3r/feng/pkg/appvuln"
	"github.com/ac0d3r/feng/pkg/macho"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx     context.Context
	cancel  func()
	appvuln *appvuln.AppScanx
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{
		appvuln: appvuln.New(),
	}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	ctx, cancel := context.WithCancel(ctx)
	a.ctx = ctx
	a.cancel = cancel
}

func (a *App) shutdown(_ context.Context) {
	a.cancel()
}

func (a *App) SelectMacho() string {
	selection, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Title:                      "Select MachO File",
		ShowHiddenFiles:            true,
		TreatPackagesAsDirectories: true,
	})
	if err != nil {
		return ""
	}
	return selection
}

func (a *App) ParseMacho(path string) []*macho.MachOInfo {
	info, err := macho.Parse(path)
	if err != nil {
		return nil
	}

	return info
}

func (a *App) AppScan() []*appvuln.Info {
	vulns, err := a.appvuln.Scan()
	if err != nil {
		return nil
	}

	return vulns
}
