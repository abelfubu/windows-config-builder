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
	Icon        string `json:"icon"`
	Id          string `json:"id"`
	Description string `json:"description"`
}

//go:embed templates/*
var templates embed.FS

var home = os.Getenv("USERPROFILE")
var config = filepath.Join(home, ".config")

func main() {
	// Step 1: Install packages
	selected := selectPackages()
	installer := winget.NewPackageInstaller()
	installer.Install(selected)

	// Step 2: Create initial configuration files
	if confirm("Do you want to create initial configuration files?") {
		createInitialConfiguration(selected)
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

func selectPackages() []string {
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

func createInitialConfiguration(selectedPackages []string) {
	os.MkdirAll(config, os.ModePerm)

	profile := getFileContent("templates/profile.ps1")

	if slices.Contains(selectedPackages, "Starship.Starship") {
		os.MkdirAll(filepath.Join(config, "starship"), os.ModePerm)
		starshipToml := getFileContent("templates/starship.toml")
		os.WriteFile(filepath.Join(config, "starship", "starship.toml"), starshipToml, os.ModePerm)
		profile = append(profile, []byte(`
# Starship
$Env:STARSHIP_CONFIG="$HOME/.config/starship/starship.toml"
Invoke-Expression (&starship init powershell)
`)...)
	}

	if slices.Contains(selectedPackages, "ajeetdsouza.zoxide") {
		profile = append(profile, []byte(`
# Zoxide
Invoke-Expression (& { (zoxide init powershell | Out-String) })
`)...)
	}

	if slices.Contains(selectedPackages, "Neovim.Neovim") {
		os.MkdirAll(filepath.Join(config, "nvim"), os.ModePerm)
		initVim := getFileContent("templates/init.lua")
		os.WriteFile(filepath.Join(config, "nvim", "init.lua"), initVim, os.ModePerm)
	}

	if slices.Contains(selectedPackages, "sharkdp.bat") {
		profile = append(profile, []byte(`
# Configure bat config directory
$Env:BAT_CONFIG_DIR="$HOME\.config\bat"
`)...)
	}

	if slices.Contains(selectedPackages, "junegunn.fzf") {
		profile = append(profile, []byte(`
# FZF
$Env:FZF_DEFAULT_OPTS=@"
--preview='bat --color=always {}'
--bind ctrl-u:preview-up,ctrl-d:preview-down,ctrl-p:toggle-preview
--color=bg+:#264f78,spinner:#569cd6,hl:#dcdcaa
--color=fg:#d4d4d4,header:#4ec9b0,info:#d4d4d4,pointer:#569cd6
--color=marker:#264f78,fg+:#ffffff,prompt:#4fc1ff,hl+:#f44747
--color=selected-bg:#264f78
--multi
"@
`)...)
	}

	if slices.Contains(selectedPackages, "LGUG2Z.komorebi") {
		profile = append(profile, []byte(`
# KOMOREBI
$Env:KOMOREBI_CONFIG_HOME = "$HOME\.config\komorebi"
`)...)
	}

	if slices.Contains(selectedPackages, "eza-community.eza") {
		profile = append(profile, []byte(`
# EZA
if (Test-Path alias:ls) { Remove-Item alias:ls }
function ls { eza -l --icons $args }
function la { eza -la --icons $args }
function lt { eza --icons -TL $args }
`)...)
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
