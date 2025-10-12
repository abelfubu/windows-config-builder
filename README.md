# Windows Config Builder

A powerful tool that automates the setup of a Windows development environment using the winget package manager. Quickly install essential development tools and configure them with sensible defaults.

## Features

- **Interactive Package Selection**: Choose from a curated list of popular development tools
- **Automated Installation**: Uses winget to install selected packages
- **Smart Configuration**: Creates configuration files for installed tools
- **Profile Management**: Sets up PowerShell profile with useful aliases and functions
- **Cross-tool Integration**: Configures tools to work well together

## Supported Packages

| Package | Description |
|---------|-------------|
| **Git** | Distributed version control system |
| **GitHub CLI** | GitHub's official command line tool |
| **Lazydocker** | Simple terminal UI for docker commands |
| **NVM for Windows** | Node Version Manager for Windows |
| **Neovim** | Hyperextensible Vim-based text editor |
| **PowerShell Preview** | Cross-platform automation and configuration tool |
| **PowerToys** | Windows system utilities to maximize productivity |
| **Rancher Desktop** | Container management and Kubernetes on the desktop |
| **ripgrep** | Recursively searches directories for a regex pattern |
| **Zig** | General-purpose programming language and toolchain |
| **bat** | A cat clone with wings (syntax highlighting) |
| **bottom** | Cross-platform graphical process/system monitor |
| **eza** | Modern replacement for 'ls' with colors and icons |
| **fd** | Simple, fast and user-friendly alternative to 'find' |
| **fzf** | Command-line fuzzy finder |
| **k9s** | Terminal UI to interact with Kubernetes clusters |
| **komorebi** | Tiling window manager for Windows |
| **lazygit** | Simple terminal UI for git commands |
| **Starship** | Cross-shell prompt for astronauts |
| **whkd** | Windows hotkey daemon |
| **zoxide** | Smarter cd command inspired by z and autojump |

## Installation

### Option 1: Download from Releases
1. Go to the [Releases page](../../releases)
2. Download the latest `windows-config-builder.exe`
3. Run the executable

### Option 2: Install via winget (when available)
```powershell
winget install Abelfubu.WindowsConfigBuilder
```

### Option 3: Build from source
```bash
git clone https://github.com/abelfubu/windows-config-builder.git
cd windows-config-builder
go build -o windows-config-builder.exe
```

## Usage

1. **Run the tool**: Execute `windows-config-builder.exe`
2. **Select packages**: Use the interactive interface to choose which packages to install
3. **Confirm installation**: The tool will install selected packages using winget
4. **Configure tools**: Choose to create initial configuration files
5. **Set up profiles**: Optionally create PowerShell profile symlinks

## What it configures

### PowerShell Profile
- **Starship**: Beautiful cross-shell prompt with custom theme
- **Zoxide**: Smart directory jumping
- **eza aliases**: Modern `ls` replacement with icons
- **fzf**: Fuzzy finder with VS Code dark theme colors
- **Environment variables**: Proper config directories for bat and komorebi

### Tool Configurations
- **Neovim**: Basic configuration in `~/.config/nvim/`
- **Starship**: Custom prompt configuration
- **Symlinks**: Automatically creates symlinks to keep configs in `~/.config/`

## Requirements

- Windows 10/11
- winget package manager (usually pre-installed)
- PowerShell (for profile configuration)

## Configuration Location

All configuration files are stored in `%USERPROFILE%\.config\` following XDG Base Directory specification:

```
~/.config/
├── profile.ps1          # PowerShell profile
├── nvim/                # Neovim configuration
│   └── init.lua
├── starship/            # Starship prompt configuration
│   └── starship.toml
├── bat/                 # bat configuration directory
└── komorebi/            # komorebi configuration directory
```

## Development

### Prerequisites
- Go 1.21 or later

### Building
```bash
go build -o windows-config-builder.exe
```

### Contributing
1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Author

**Abel de la Fuente** - [@abelfubu](https://github.com/abelfubu)

---

*Made with ❤️ for the Windows development community*
