package sources

import (
	"os"
	"path/filepath"
	"syscall"

	"github.com/lxc/distrobuilder/shared"
)

// Pacman represents the pacstrap-like downloader.
type Pacman struct{}

// NewPacman creates a new Pacman instance.
func NewPacman() *Pacman {
	return &Pacman{}
}

// Run runs pacman.
func (s *Pacman) Run(definition shared.Definition, rootfsDir string) error {
	var args []string

	os.RemoveAll(rootfsDir)

	if definition.Image.ArchitectureMapped != "" {
		args = append(args, "--arch", definition.Image.ArchitectureMapped)
	}

	args = append(args, "--noconfirm", "--root", rootfsDir, "-Sy")

	if len(definition.Source.EarlyPackages) > 0 {
		args = append(args, definition.Source.EarlyPackages...)
	} else {
		args = append(args, "base")
	}

	// Define some paths
	procDir := filepath.Join(rootfsDir, "proc")
	sysDir := filepath.Join(rootfsDir, "sys")
	efivarsDir := filepath.Join(sysDir, "firmware", "efi", "efivars")
	udevDir := filepath.Join(rootfsDir, "dev")
	devptsDir := filepath.Join(udevDir, "pts")
	shmDir := filepath.Join(udevDir, "shm")
	runDir := filepath.Join(rootfsDir, "run")
	tmpDir := filepath.Join(rootfsDir, "tmp")

	// Create required directories in root
	err := os.MkdirAll(filepath.Join(rootfsDir, "var", "cache", "pacman", "pkg"), 0755)
	if err != nil {
		return err
	}
	err = os.MkdirAll(filepath.Join(rootfsDir, "var", "lib", "pacman"), 0755)
	if err != nil {
		return err
	}
	err = os.MkdirAll(filepath.Join(rootfsDir, "var", "log"), 0755)
	if err != nil {
		return err
	}
	err = os.MkdirAll(filepath.Join(rootfsDir, "etc"), 0755)
	if err != nil {
		return err
	}

	err = os.MkdirAll(udevDir, 0755)
	if err != nil {
		return err
	}
	err = os.MkdirAll(runDir, 0755)
	if err != nil {
		return err
	}
	err = os.MkdirAll(tmpDir, 0777)
	if err != nil {
		return err
	}
	err = os.MkdirAll(sysDir, 0555)
	if err != nil {
		return err
	}
	err = os.MkdirAll(procDir, 0555)
	if err != nil {
		return err
	}

	// Mount usual chroot directories
	err = shared.RunCommand("mount", "proc", procDir, "-t", "proc", "-o", "nosuid,noexec,nodev")
	if err != nil {
		return err
	}
	defer syscall.Unmount(procDir, 0)
	err = shared.RunCommand("mount", "sys", sysDir, "-t", "sysfs", "-o", "nosuid,noexec,nodev,ro")
	if err != nil {
		return err
	}
	defer syscall.Unmount(sysDir, 0)

	var efivarsStat os.FileInfo
	if efivarsStat, err = os.Stat(efivarsDir); err == nil && efivarsStat.IsDir() {
		err = shared.RunCommand("mount", "efivarfs", efivarsDir, "-t", "efivarfs", "-o", "nosuid,noexec,nodev")
		if err == nil {
			defer syscall.Unmount(efivarsDir, 0)
		}
	}

	err = shared.RunCommand("mount", "udev", udevDir, "-t", "devtmpfs", "-o", "mode=0755,nosuid")
	if err != nil {
		return err
	}
	defer syscall.Unmount(udevDir, 0)
	err = shared.RunCommand("mount", "devpts", devptsDir, "-t", "devpts", "-o", "mode=0620,gid=5,nosuid,noexec")
	if err != nil {
		return err
	}
	defer syscall.Unmount(devptsDir, 0)
	err = shared.RunCommand("mount", "shm", shmDir, "-t", "tmpfs", "-o", "mode=1777,nosuid,nodev")
	if err != nil {
		return err
	}
	defer syscall.Unmount(shmDir, 0)
	err = shared.RunCommand("mount", "run", runDir, "-t", "tmpfs", "-o", "nosuid,nodev,mode=0755")
	if err != nil {
		return err
	}
	defer syscall.Unmount(runDir, 0)
	err = shared.RunCommand("mount", "tmp", tmpDir, "-t", "tmpfs", "-o", "mode=1777,strictatime,nodev,nosuid")
	if err != nil {
		return err
	}
	defer syscall.Unmount(tmpDir, 0)

	// Run pacman
	err = shared.RunCommand("pacman", args...)
	if err != nil {
		return err
	}

	return nil
}
