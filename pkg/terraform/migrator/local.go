package migrator

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
)

func (m *Migrator) initBackendLocal(ctx context.Context) error {
	if m.config.Backend == nil {
		return fmt.Errorf("backend config is nil")
	}

	path, err := m.config.Backend.Props.GetString("path")

	if err != nil {
		return err
	}

	f, err := os.OpenFile(path, os.O_RDONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	m.backendStateBytes, err = io.ReadAll(f)
	if err != nil {
		return err
	}

	if err = json.Unmarshal(m.backendStateBytes, &m.backendState); err != nil {
		return err
	}
	return nil
}
