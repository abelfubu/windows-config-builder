package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"

	"github.com/abelfubu/windows-config-builder/pkg/winget"

	"github.com/charmbracelet/huh"
)

type Package struct {
	Profile     []string `json:"profile"`
	Icon        string   `json:"icon"`
	Id          string   `json:"id"`
	Description string   `json:"description"`
}

//go:embed templates/*
var templates embed.FS

var home = os.Getenv("USERPROFILE")
var config = filepath.Join(home, ".config-test")

func main() {
	// Step 1: Install packages
	packages := getPackages()
	selected := selectPackages(packages)
	installer := winget.NewPackageInstaller()
	installer.Install(selected)

	// Step 2: Create initial configuration files
	if confirm("Do you want to create initial configuration files?") {
		createInitialConfiguration(selected, packages)
	}

	// Step 3: Ask about symlink
	if confirm("Do you want to add a symlink to your PowerShell profile?") {
		addPowershellSymlink()
	}

	// Step 4: Add Neovim symlink if Neovim was installed
	if slices.Contains(selected, "Neovim.Neovim") {
		addNeovimSymlink()
	}
}

func getPackages() []Package {
	data, err := templates.ReadFile("templates/packages.json")
	if err != nil {
		fmt.Println("❌ Failed to read packages.json:", err)
		return nil
	}

	var pkgs []Package
	if err := json.Unmarshal(data, &pkgs); err != nil {
		fmt.Println("❌ Failed to parse packages.json:", err)
		return nil
	}

	return pkgs
}

func selectPackages(pkgs []Package) []string {
	var options []huh.Option[string]
	for _, pkg := range pkgs {
		options = append(options, huh.NewOption(fmt.Sprintf("%s %-30s %s", pkg.Icon, pkg.Id, pkg.Description), pkg.Id))
	}

	var selected []string
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewMultiSelect[string]().
				Title("Select packages to install").
				Options(options...).
				Value(&selected),
		),
	)

	form.Run()

	return selected
}

func addPowershellSymlink() {
	src := filepath.Join(config, "profile.ps1")
	cmd := exec.Command("pwsh", "-Command", "echo $profile")

	output, error := cmd.Output()

	if error != nil {
		fmt.Println("❌ Failed to execute PowerShell command:", error)
		return
	}

	profilePath := strings.TrimSpace(string(output))

	os.MkdirAll(filepath.Dir(profilePath), os.ModePerm)

	os.Remove(profilePath)

	err := os.Symlink(src, profilePath)
	if err != nil {
		fmt.Println("❌ Failed to create symlink:", err)
	} else {
		fmt.Println("✅ Symlink created at PowerShell profile!")
	}
}

func addNeovimSymlink() {
	local := os.Getenv("LOCALAPPDATA")
	src := filepath.Join(config, "nvim")
	nvimConfigDir := filepath.Join(local, "nvim")

	// Remove existing file or symlink if present
	os.Remove(nvimConfigDir)

	err := os.Symlink(src, nvimConfigDir)
	if err != nil {
		fmt.Println("❌ Failed to create symlink:", err)
	} else {
		fmt.Println("✅ Symlink created at PowerShell profile!")
	}
}

func createInitialConfiguration(selectedPackages []string, pkgs []Package) {
	os.MkdirAll(config, os.ModePerm)

	profile := getFileContent("templates/profile.ps1")

	for _, pkg := range pkgs {
		if !slices.Contains(selectedPackages, pkg.Id) || len(pkg.Profile) == 0 {
			continue
		}

		profile = fmt.Appendf(profile, "# %s\n", pkg.Id)
		for _, env := range pkg.Profile {
			profile = fmt.Appendf(profile, "%s\n", env)
		}
		profile = fmt.Appendf(profile, "\n")
	}

	if true {
		os.WriteFile(filepath.Join(config, "profile.ps1"), profile, os.ModePerm)
		return
	}

	if slices.Contains(selectedPackages, "Starship.Starship") {
		os.MkdirAll(filepath.Join(config, "starship"), os.ModePerm)
		starshipToml := getFileContent("templates/starship.toml")
		os.WriteFile(filepath.Join(config, "starship", "starship.toml"), starshipToml, os.ModePerm)
	}

	if slices.Contains(selectedPackages, "Neovim.Neovim") {
		os.MkdirAll(filepath.Join(config, "nvim"), os.ModePerm)
		initVim := getFileContent("templates/init.lua")
		os.WriteFile(filepath.Join(config, "nvim", "init.lua"), initVim, os.ModePerm)
	}

	if slices.Contains(selectedPackages, "Derailed.k9s") {
		os.CopyFS(filepath.Join(config, "k9s"), os.DirFS("./templates/k9s"))
		os.Symlink(filepath.Join(os.Getenv("LOCALAPPDATA"), "k9s"), filepath.Join(config, "k9s"))
	}

	os.WriteFile(filepath.Join(config, "profile.ps1"), profile, os.ModePerm)
	fmt.Print("✅ Initial configuration files created!\n")
}

func confirm(message string) bool {
	var result bool

	prompt := huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().Title(message).Value(&result),
		),
	)

	prompt.Run()

	return result
}

func getFileContent(path string) []byte {
	content, error := templates.ReadFile(path)

	if error != nil {
		fmt.Println("❌ Failed to read embedded file:", error)
		return []byte{}
	}

	return content
}

func Find[T any](slice []T, predicate func(T) bool) (T, bool) {
	var zero T
	for _, v := range slice {
		if predicate(v) {
			return v, true
		}
	}

	return zero, false
}
