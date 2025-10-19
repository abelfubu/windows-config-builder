package winget

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type PackageInstaller struct {
	packages string
}

func NewPackageInstaller() *PackageInstaller {
	installer := new(PackageInstaller)
	installer.loadInstalledPackages()

	return installer
}

func (p *PackageInstaller) loadInstalledPackages() {
	cmd := exec.Command("winget", "list")

	output, error := cmd.Output()

	if error != nil {
		fmt.Println("Error running 'winget list'", error)
		return
	}

	p.packages = string(output)
}

func (p *PackageInstaller) hasPackage(pkg string) bool {
	return strings.Contains(p.packages, pkg)
}

func (p *PackageInstaller) Install(pkgs []string) {
	var toInstall []string

	for _, pkg := range pkgs {
		if p.hasPackage(pkg) {
			fmt.Printf("📦 Package %s is already installed\n", pkg)
			continue
		}

		toInstall = append(toInstall, pkg)
	}

	if len(toInstall) == 0 {
		fmt.Println("✅ All packages already installed")
		return
	}

	base := []string{"install", "--silent", "--accept-package-agreements"}
	args := append(base, toInstall...)

	fmt.Printf("💾 Installing %v...\n", toInstall)

	cmd := exec.Command("winget", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}
