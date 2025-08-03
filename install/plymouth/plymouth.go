package plymouth

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"archriot-installer/logger"
)

type PlymouthManager struct {
	homeDir   string
	configDir string
}

func NewPlymouthManager() (*PlymouthManager, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	return &PlymouthManager{
		homeDir:   homeDir,
		configDir: filepath.Join(homeDir, ".local/share/archriot/config/default/plymouth"),
	}, nil
}

func (pm *PlymouthManager) InstallPlymouth() error {
	logger.Log("Progress", "System", "Plymouth", "Checking configuration")

	if !pm.needsReinstall() {
		logger.Log("Complete", "System", "Plymouth", "Already configured")
		return nil
	}

	logger.Log("Progress", "System", "Plymouth", "Installing")

	// Install Plymouth package
	logger.Log("Progress", "Package", "Plymouth", "Installing package")
	if err := pm.runCommand("yay", "-S", "--noconfirm", "plymouth"); err != nil {
		return fmt.Errorf("failed to install Plymouth: %v", err)
	}
	logger.Log("Success", "Package", "Plymouth", "Package installed")

	// Backup and configure mkinitcpio
	logger.Log("Progress", "System", "Plymouth", "Configuring mkinitcpio")
	if err := pm.configureMkinitcpio(); err != nil {
		return err
	}
	logger.Log("Success", "System", "Plymouth", "mkinitcpio configured")

	// Configure bootloader
	logger.Log("Progress", "System", "Plymouth", "Configuring bootloader")
	if err := pm.configureBootloader(); err != nil {
		return err
	}
	logger.Log("Success", "System", "Plymouth", "Bootloader configured")

	// Install theme
	logger.Log("Progress", "System", "Plymouth", "Installing theme")
	if err := pm.installTheme(); err != nil {
		return err
	}
	logger.Log("Success", "System", "Plymouth", "Theme installed")

	// Set default theme
	logger.Log("Progress", "System", "Plymouth", "Setting default theme")
	if err := pm.runCommand("sudo", "plymouth-set-default-theme", "-R", "archriot"); err != nil {
		return fmt.Errorf("failed to set default theme: %v", err)
	}
	logger.Log("Success", "System", "Plymouth", "Theme activated")

	logger.Log("Success", "System", "Plymouth", "Installation complete")
	return nil
}

func (pm *PlymouthManager) needsReinstall() bool {
	// Check if Plymouth installed
	if _, err := exec.LookPath("plymouth"); err != nil {
		return true
	}

	// Check if theme set
	cmd := exec.Command("sudo", "plymouth-set-default-theme")
	output, err := cmd.Output()
	if err != nil || strings.TrimSpace(string(output)) != "archriot" {
		return true
	}

	// Check if theme files exist
	if _, err := os.Stat("/usr/share/plymouth/themes/archriot/archriot.plymouth"); err != nil {
		return true
	}

	// Check logo content
	newLogo := filepath.Join(pm.configDir, "logo.png")
	currentLogo := "/usr/share/plymouth/themes/archriot/logo.png"
	if pm.filesChanged(newLogo, currentLogo) {
		return true
	}

	return false
}

func (pm *PlymouthManager) filesChanged(file1, file2 string) bool {
	data1, err1 := os.ReadFile(file1)
	data2, err2 := os.ReadFile(file2)
	if err1 != nil || err2 != nil {
		return true
	}
	return string(data1) != string(data2)
}

func (pm *PlymouthManager) configureMkinitcpio() error {
	// Backup
	pm.runCommand("sudo", "cp", "/etc/mkinitcpio.conf", "/etc/mkinitcpio.conf.bak")

	// Check if plymouth already in hooks
	content, err := os.ReadFile("/etc/mkinitcpio.conf")
	if err != nil {
		return err
	}

	if strings.Contains(string(content), "plymouth") {
		return nil
	}

	// Add plymouth hook - EXACT shell script logic
	lines := strings.Split(string(content), "\n")
	for i, line := range lines {
		if strings.HasPrefix(line, "HOOKS=") {
			if strings.Contains(line, "encrypt") {
				// LUKS: plymouth must come AFTER consolefont but BEFORE block
				// Matches: sudo sed -i '/^HOOKS=/s/\(consolefont.*\)\(block\)/\1 plymouth \2/' /etc/mkinitcpio.conf
				lines[i] = strings.Replace(line, " block", " plymouth block", 1)
			} else if strings.Contains(line, "base systemd") {
				// systemd initramfs
				lines[i] = strings.Replace(line, "base systemd", "base systemd plymouth", 1)
			} else if strings.Contains(line, "base udev") {
				// Standard udev initramfs - same LUKS logic
				lines[i] = strings.Replace(line, " block", " plymouth block", 1)
			} else {
				// Fallback: add after base udev
				lines[i] = strings.Replace(line, "base udev", "base udev plymouth", 1)
			}
			break
		}
	}

	// Write file
	newContent := strings.Join(lines, "\n")
	if err := pm.writeFile("/etc/mkinitcpio.conf", newContent); err != nil {
		return err
	}

	// Regenerate initramfs
	logger.Log("Progress", "System", "Plymouth", "Regenerating initramfs (slow)")
	err = pm.runCommand("sudo", "mkinitcpio", "-P")
	if err == nil {
		logger.Log("Success", "System", "Plymouth", "Initramfs regenerated")
	}
	return err
}

func (pm *PlymouthManager) configureBootloader() error {
	// systemd-boot
	if _, err := os.Stat("/boot/loader/entries"); err == nil {
		entries, _ := filepath.Glob("/boot/loader/entries/*.conf")
		for _, entry := range entries {
			entryName := filepath.Base(entry)
			if strings.Contains(entryName, "fallback") {
				continue
			}
			content, _ := os.ReadFile(entry)
			if strings.Contains(string(content), "splash") {
				continue
			}
			lines := strings.Split(string(content), "\n")
			for i, line := range lines {
				if strings.HasPrefix(line, "options ") {
					lines[i] = line + " splash quiet"
					break
				}
			}
			pm.writeFile(entry, strings.Join(lines, "\n"))
		}
		return nil
	}

	// GRUB
	if _, err := os.Stat("/etc/default/grub"); err == nil {
		pm.runCommand("sudo", "cp", "/etc/default/grub", "/etc/default/grub.bak")

		content, _ := os.ReadFile("/etc/default/grub")
		if strings.Contains(string(content), "splash") {
			return nil
		}

		lines := strings.Split(string(content), "\n")
		for i, line := range lines {
			if strings.HasPrefix(line, "GRUB_CMDLINE_LINUX_DEFAULT=") {
				// Extract content between first and last quotes like shell script
				start := strings.Index(line, "\"")
				end := strings.LastIndex(line, "\"")
				if start != -1 && end != -1 && start != end {
					cmdline := line[start+1 : end]
					if !strings.Contains(cmdline, "splash") {
						cmdline += " splash"
					}
					if !strings.Contains(cmdline, "quiet") {
						cmdline += " quiet"
					}
					lines[i] = fmt.Sprintf("GRUB_CMDLINE_LINUX_DEFAULT=\"%s\"", strings.TrimSpace(cmdline))
				}
				break
			}
		}

		if err := pm.writeFile("/etc/default/grub", strings.Join(lines, "\n")); err != nil {
			return err
		}
		logger.Log("Progress", "System", "Plymouth", "Regenerating GRUB config")
		return pm.runCommand("sudo", "grub-mkconfig", "-o", "/boot/grub/grub.cfg")
	}

	// UKI
	if _, err := os.Stat("/etc/cmdline.d"); err == nil {
		// Check existing files for splash/quiet like shell script does
		files, _ := filepath.Glob("/etc/cmdline.d/*.conf")
		hasSplash := false
		hasQuiet := false

		for _, file := range files {
			content, err := os.ReadFile(file)
			if err != nil {
				continue
			}
			if strings.Contains(string(content), "splash") {
				hasSplash = true
			}
			if strings.Contains(string(content), "quiet") {
				hasQuiet = true
			}
		}

		// Only add what's missing, append like shell script
		if !hasSplash {
			pm.runCommand("sh", "-c", "echo 'splash' | sudo tee -a /etc/cmdline.d/archriot.conf")
		}
		if !hasQuiet {
			pm.runCommand("sh", "-c", "echo 'quiet' | sudo tee -a /etc/cmdline.d/archriot.conf")
		}
		return nil
	}
	if _, err := os.Stat("/etc/kernel/cmdline"); err == nil {
		pm.runCommand("sudo", "cp", "/etc/kernel/cmdline", "/etc/kernel/cmdline.bak")
		content, _ := os.ReadFile("/etc/kernel/cmdline")
		cmdline := strings.TrimSpace(string(content))
		if !strings.Contains(cmdline, "splash") {
			cmdline += " splash"
		}
		if !strings.Contains(cmdline, "quiet") {
			cmdline += " quiet"
		}
		return pm.writeFile("/etc/kernel/cmdline", strings.TrimSpace(cmdline))
	}

	return nil
}

func (pm *PlymouthManager) installTheme() error {
	// Check if theme files exist first
	if _, err := os.Stat(pm.configDir); err != nil {
		return fmt.Errorf("Plymouth theme files not found at %s", pm.configDir)
	}

	logoFile := filepath.Join(pm.configDir, "logo.png")
	if _, err := os.Stat(logoFile); err != nil {
		return fmt.Errorf("logo.png not found at %s", logoFile)
	}

	// Cleanup
	pm.runCommand("sudo", "rm", "-rf", "/usr/share/plymouth/themes/archriot")
	pm.runCommand("sudo", "rm", "-f", "/var/lib/plymouth/themes/default.plymouth")

	// Create directory
	if err := pm.runCommand("sudo", "mkdir", "-p", "/usr/share/plymouth/themes/archriot"); err != nil {
		return err
	}

	// Copy files
	return pm.runCommand("sudo", "cp", "-r", pm.configDir+"/.", "/usr/share/plymouth/themes/archriot/")
}

func (pm *PlymouthManager) runCommand(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	return cmd.Run()
}

func (pm *PlymouthManager) writeFile(path, content string) error {
	cmd := exec.Command("sudo", "tee", path)
	cmd.Stdin = strings.NewReader(content)
	return cmd.Run()
}
