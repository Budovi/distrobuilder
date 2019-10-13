package sources

import (
	"github.com/lxc/distrobuilder/shared"
	"os"
)

// ManjaroBootstrap represents the manjaro-bootstrap downloader.
type ManjaroBootstrap struct{}

// NewManjaroBootstrap creates a new ManjaroBootstrap instance.
func NewManjaroBootstrap() *ManjaroBootstrap {
	return &ManjaroBootstrap{}
}

// Run runs manjaro-bootstrap.
func (s *ManjaroBootstrap) Run(definition shared.Definition, rootfsDir string) error {
	var args []string

	os.RemoveAll(rootfsDir)

	if definition.Image.ArchitectureMapped != "" {
		args = append(args, "-a", definition.Image.ArchitectureMapped)
	}

	if definition.Source.URL != "" {
		args = append(args, "-r", definition.Source.URL)
	}

	args = append(args, rootfsDir)

	err := shared.RunCommand("manjaro-bootstrap", args...)
	if err != nil {
		return err
	}

	return nil
}
