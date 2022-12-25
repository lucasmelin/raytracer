//go:build mage

package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path"
	"runtime"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

type Install mg.Namespace

// Runs `go mod download` and then builds the `raytracer` binary.
func Build() error {
	if err := sh.Run("go", "mod", "download"); err != nil {
		return err
	}
	return sh.Run("go", "build", "-o", "raytracer", "./cmd")
}

// Runs the `raytracer` binary, building it first if necessary.
func Run() error {
	mg.Deps(Build)
	return sh.Run("./raytracer")
}

// Displays the generated image, generating it first if necessary.
func View() error {
	pngImage := "image.png"
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}
	if _, err := os.Stat(path.Join(cwd, pngImage)); errors.Is(err, os.ErrNotExist) {
		fmt.Printf("%s not found, running raytracer ✨\n", pngImage)
		mg.Deps(Build, Run)
	}
	return sh.Run("open", "-a", "Preview", pngImage)
}

// Removes the generated PNG image from disk.
func Clean() error {
	pngImage := "image.png"
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}
	fullImagePath := path.Join(cwd, pngImage)
	if _, err := os.Stat(fullImagePath); err == nil {
		fmt.Printf("%s found, removing from disk ✨\n", pngImage)
		return os.Remove(fullImagePath)
	}
	fmt.Printf("%s not found, nothing to remove ✨\n", pngImage)
	return nil
}

// Installs all system and Go dependencies.
func (Install) Deps() error {
	var sdlCommand []string
	if runtime.GOOS == "linux" {
		sdlCommand = []string{"sudo", "apt-get", "update", "&&", "sudo", "apt-get", "install", "libsdl2{,-image,-mixer,-ttf,-gfx}-dev"}
	} else if runtime.GOOS == "darwin" {
		sdlCommand = []string{"brew", "install", "sdl2{,_image,_mixer,_ttf,_gfx}", "pkg-config"}
	} else {
		return errors.New("unknown OS")
	}
	if err := sh.Run(sdlCommand[0], sdlCommand[1:]...); err != nil {
		return err
	}
	return sh.Run("go", "mod", "download")
}

// Runs the unit tests.
func Test() error {
	mg.Deps(Install.Deps)
	if _, err := exec.LookPath("gotestsum"); err == nil {
		output, err := sh.Output("gotestsum", "--no-color=false")
		fmt.Printf("%s\n", output)
		return err
	}
	output, err := sh.Output("go", "test", "-v", "./...")
	fmt.Printf("%s\n", output)
	return err
}
