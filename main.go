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

type Symlink struct {
	Source string `json:"source"`
	Target string `json:"target"`
}

type Package struct {
	ConfigFolder *string   `json:"configFolder,omitempty"`
	Icon         string    `json:"icon"`
	Id           string    `json:"id"`
	Description  string    `json:"description"`
	Profile      []string  `json:"profile"`
	Symlinks     []Symlink `json:"symlinks,omitempty"`
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

func createInitialConfiguration(selectedPackages []string, pkgs []Package) {
	os.MkdirAll(config, os.ModePerm)

	profile := getFileContent("templates/profile.ps1")

	for _, pkg := range pkgs {
		if !slices.Contains(selectedPackages, pkg.Id) {
			continue
		}

		if len(pkg.Profile) > 0 {
			profile = fmt.Appendf(profile, "# %s\n", pkg.Id)
			for _, env := range pkg.Profile {
				profile = fmt.Appendf(profile, "%s\n", env)
			}
			profile = fmt.Appendf(profile, "\n")
		}

		if pkg.ConfigFolder != nil {
			os.CopyFS(filepath.Join(config, *pkg.ConfigFolder), os.DirFS(fmt.Sprintf("templates/%s", *pkg.ConfigFolder)))
		}

		if len(pkg.Symlinks) > 0 {
			fmt.Println("Creating symlinks for", pkg.Id)
			for _, link := range pkg.Symlinks {
				target := filepath.Join(os.Getenv(link.Target), link.Source)

				os.Remove(target)
				source := filepath.Join(filepath.Join(home, ".config"), link.Source)
				error := os.Symlink(source, target)
				if error != nil {
					fmt.Println("❌ Failed to create symlink:", error)
				}
			}
		}
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
