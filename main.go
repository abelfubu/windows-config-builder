package main

import (
	"embed"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"

	"github.com/charmbracelet/huh"
)

//go:embed templates/*
var templates embed.FS

var home = os.Getenv("USERPROFILE")
var config = filepath.Join(home, ".config")

func main() {
	// Step 1: Install packages
	selectedPackages := installWingetPkgs()

	// Step 2: Create initial configuration files
	if confirm("Do you want to create initial configuration files?") {
		createInitialConfiguration(selectedPackages)
	}

	// Step 3: Ask about symlink
	if confirm("Do you want to add a symlink to your PowerShell profile?") {
		addPowershellSymlink()
	}

	// Step 4: Add Neovim symlink if Neovim was installed
	if slices.Contains(selectedPackages, "Neovim.Neovim") {
		addNeovimSymlink()
	}
}

func installWingetPkgs() []string {
	pkgs := []string{
		"Git.Git",
		"GitHub.cli",
		"JesseDuffield.Lazydocker",
		"CoreyButler.NVMforWindows",
		"Neovim.Neovim",
		"Microsoft.PowerShell.Preview",
		"Microsoft.PowerToys",
		"SUSE.RancherDesktop",
		"BurntSushi.ripgrep.MSVC",
		"zig.zig",
		"sharkdp.bat",
		"Clement.bottom",
		"eza-community.eza",
		"sharkdp.fd",
		"junegunn.fzf",
		"Derailed.k9s",
		"LGUG2Z.komorebi",
		"JesseDuffield.lazygit",
		"Starship.Starship",
		"LGUG2Z.whkd",
		"ajeetdsouza.zoxide",
	}

	descriptions := map[string]string{
		"Git.Git":                      "Distributed version control system",
		"GitHub.cli":                   "GitHub's official command line tool",
		"JesseDuffield.Lazydocker":     "Simple terminal UI for docker commands",
		"CoreyButler.NVMforWindows":    "Node Version Manager for Windows",
		"Neovim.Neovim":                "Hyperextensible Vim-based text editor",
		"Microsoft.PowerShell.Preview": "Cross-platform automation and configuration tool",
		"Microsoft.PowerToys":          "Windows system utilities to maximize productivity",
		"SUSE.RancherDesktop":          "Container management and Kubernetes on the desktop",
		"BurntSushi.ripgrep.MSVC":      "Recursively searches directories for a regex pattern",
		"zig.zig":                      "General-purpose programming language and toolchain",
		"sharkdp.bat":                  "A cat clone with wings (syntax highlighting)",
		"Clement.bottom":               "Cross-platform graphical process/system monitor",
		"eza-community.eza":            "Modern replacement for 'ls' with colors and icons",
		"sharkdp.fd":                   "Simple, fast and user-friendly alternative to 'find'",
		"junegunn.fzf":                 "Command-line fuzzy finder",
		"Derailed.k9s":                 "Terminal UI to interact with Kubernetes clusters",
		"LGUG2Z.komorebi":              "Tiling window manager for Windows",
		"JesseDuffield.lazygit":        "Simple terminal UI for git commands",
		"Starship.Starship":            "Cross-shell prompt for astronauts",
		"LGUG2Z.whkd":                  "Windows hotkey daemon",
		"ajeetdsouza.zoxide":           "Smarter cd command inspired by z and autojump",
	}

	var selected []string
	var options []huh.Option[string]

	// Find the longest package name for alignment
	maxLen := 0
	for _, pkg := range pkgs {
		if len(pkg) > maxLen {
			maxLen = len(pkg)
		}
	}

	for _, pkg := range pkgs {
		desc := descriptions[pkg]
		padding := (maxLen + 5) - len(pkg)
		spaces := strings.Repeat(" ", padding)
		options = append(options, huh.NewOption(fmt.Sprintf("%s%s - %s", pkg, spaces, desc), pkg))
	}

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewMultiSelect[string]().
				Title("Select packages to install").
				Options(options...).
				Value(&selected),
		),
	)

	form.Run()

	for _, pkg := range selected {
		fmt.Printf("Installing %s...\n", pkg)
		cmd := exec.Command("winget", "install", "--id", pkg, "--accept-package-agreements", "--accept-source-agreements", "-h")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Run()
	}

	if len(selected) > 0 {
		fmt.Println("✅ Done installing selected packages!")
	}

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
function ls { eza -l --icons $args }
function la { eza -la --icons $args }
function lt { eza --icons -TL $args }
`)...)
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
