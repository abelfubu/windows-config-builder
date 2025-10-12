package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"

	"github.com/charmbracelet/huh"
)

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
		"MSIX\\Abelfubu.RaindropCommandPaletteExtension",
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

	var selected []string
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewMultiSelect[string]().
				Title("Select packages to install").
				Options(huh.NewOptions(pkgs...)...).
				Value(&selected),
		),
	)

	form.Run()

	for _, pkg := range selected {
		fmt.Printf("Installing %s...\n", pkg)
		cmd := exec.Command("winget", "install", "--id", pkg, "--accept-package-agreements", "--accept-source-agreements", "-h")
		cmd.Run()
	}

	fmt.Println("✅ Done installing selected packages!")

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

	profile := getFileContent("profile.ps1")

	if slices.Contains(selectedPackages, "Starship.Starship") {
		fmt.Println("Creating initial Starship configuration...")
		os.MkdirAll(filepath.Join(config, "starship"), os.ModePerm)
		starshipToml := getFileContent("starship.toml")
		os.WriteFile(filepath.Join(config, "starship", "starship.toml"), starshipToml, os.ModePerm)
		profile = append(profile, []byte("\n# Initialize Starship\n$env:STARSHIP_CONFIG=\"$HOME/.config/starship/starship.toml\"\nInvoke-Expression (&starship init powershell)\n")...)
	}

	if slices.Contains(selectedPackages, "ajeetdsouza.zoxide") {
		fmt.Println("Creating initial Zoxide configuration...")
		profile = append(profile, []byte("\n# Initialize Zoxide\nInvoke-Expression (& { (zoxide init powershell | Out-String) })\n")...)
	}

	if slices.Contains(selectedPackages, "Neovim.Neovim") {
		fmt.Println("Creating initial Neovim configuration...")
		os.MkdirAll(filepath.Join(config, "nvim"), os.ModePerm)
		initVim := getFileContent("init.lua")
		os.WriteFile(filepath.Join(config, "nvim", "init.lua"), initVim, os.ModePerm)
	}

	if slices.Contains(selectedPackages, "eza-community.eza") {
		profile = append(profile, []byte("\n# Alias ls to eza\nfunction ls { eza -l --icons $args }\nfunction la { eza -la --icons $args }\nfunction lt { eza --icons -TL $args }\n")...)
	}

	os.WriteFile(filepath.Join(config, "profile.ps1"), profile, os.ModePerm)
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
	content, error := os.ReadFile(path)

	if error != nil {
		fmt.Println("❌ Failed to read file:", error)
		return []byte{}
	}

	return content
}
