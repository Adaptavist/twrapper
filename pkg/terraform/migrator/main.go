package migrator

import (
	"context"
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"os"

	"github.com/adaptavist/terraform-wrapper/v1/cmd/twrapper/config"
	"github.com/adaptavist/terraform-wrapper/v1/pkg/terraform/model/state"
	"github.com/adaptavist/terraform-wrapper/v1/pkg/terraform/tfe-helper"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	tf "github.com/hashicorp/go-tfe"
	"go.uber.org/zap"
)

var reason = "Twrapper migration from backend"

type Migrator struct {
	config            *config.Config
	helper            *tfe.ClientHelper
	backendState      *state.State
	backendStateBytes []byte
	workspace         *tf.Workspace
	workspaceState    *tf.StateVersion
	project           *tf.Project
	s3                *s3.Client
}

// ACCESSORS

func (m *Migrator) Config() *config.Config {
	return m.config
}

func (m *Migrator) Helper() *tfe.ClientHelper {
	return m.helper
}

func (m *Migrator) Workspace() *tf.Workspace {
	return m.workspace
}

func (m *Migrator) WorkspaceState() *tf.StateVersion {
	return m.workspaceState
}

func (m *Migrator) Project() *tf.Project {
	return m.project
}

// CONSTRUCTOR

// New migration struct
func New(c *config.Config) (*Migrator, error) {
	helper, err := tfe.New()

	if err != nil {
		return nil, err
	}

	m := &Migrator{
		config: c,
		helper: helper,
	}

	return m, nil
}

// METHODS

// Init initializes the origin backend, and
func (m *Migrator) initBackend() error {
	switch m.config.Backend.Type {
	case "s3":
		return m.initBackendS3(context.TODO())
	case "local":
		return m.initBackendLocal(context.TODO())
	default:
		return fmt.Errorf("backend type %s is not supported", m.config.Backend.Type)
	}
}

// Initalises the terraform enterpise workspace and project we're migrating to
func (m *Migrator) initCloud() error {
	if m.config.Cloud.Workspace == "" {
		return fmt.Errorf("workspace is not set")
	}
	w, err := m.helper.CreateWorkspaceAndProjectIfNotExists(
		&m.config.Cloud.Organization,
		&m.config.Cloud.Project,
		&m.config.Cloud.Workspace,
	)

	if err != nil {
		return err
	}

	m.workspace = w

	// We've already created the workspace and project, so we don't need to do it again
	// So lets just get the project
	if m.config.Cloud.Project != "" {
		m.project, _ = m.helper.GetProjectByName(m.config.Cloud.Organization, m.config.Cloud.Project)
	}

	if m.workspace.CurrentStateVersion != nil {
		v, err := m.helper.Client.StateVersions.Read(m.helper.Context, m.workspace.CurrentStateVersion.ID)
		if err != nil {
			return err
		}
		m.workspaceState = v
	}

	return nil
}

// Init initializes the migrator by fetching backend and cloud state and information up front.
func (m *Migrator) Init() error {
	if err := m.initBackend(); err != nil {
		return err
	}
	if err := m.initCloud(); err != nil {
		return err
	}
	return nil
}

// Lock the cloud backend
func (m *Migrator) Lock() error {
	_, err := m.helper.Client.Workspaces.Lock(m.helper.Context, m.workspace.ID, tf.WorkspaceLockOptions{
		Reason: &reason,
	})
	if err != nil {
		m.config.Logger.Error("failed to lock workspace", zap.Error(err))
	}
	return err
}

// Unlock the cloud backend
func (m *Migrator) Unlock() error {
	_, err := m.helper.Client.Workspaces.Unlock(m.helper.Context, m.workspace.ID)
	if err != nil {
		m.config.Logger.Error("failed to unlock workspace", zap.Error(err))
	}
	return err
}

// Migrate the terraform backend
func (m *Migrator) Migrate() error {
	if m.backendState == nil {
		m.config.Logger.Info("no backend state to migrate")
		return nil
	}

	// we need to make sure the backend isn't newer than the cloud
	if m.workspaceState != nil && m.backendState.Serial > m.workspaceState.Serial {
		return fmt.Errorf("backend state newer than cloud state")
	}

	// encode backend state
	hash := fmt.Sprintf("%x", md5.Sum(m.backendStateBytes))
	state := base64.StdEncoding.EncodeToString(m.backendStateBytes)

	// lock cloud state
	if err := m.Lock(); err != nil {
		return err
	}

	// Ensure state is unlocked
	defer m.Unlock()

	// Make new version
	s, err := m.helper.Client.StateVersions.Create(m.helper.Context, m.workspace.ID, tf.StateVersionCreateOptions{
		Lineage: &m.backendState.Lineage,
		Serial:  &m.backendState.Serial,
		MD5:     &hash,
		State:   &state,
	})

	m.config.Logger.Debug("state version created", zap.String("state_version_id", s.ID))

	m.workspaceState = s

	return err
}

func (m *Migrator) Cleanup() error {
	switch m.config.Backend.Type {
	case "s3":
		return m.cleanupS3()
	case "local":
		path, err := m.config.Backend.Props.GetString("path")
		if err != nil {
			return err
		}
		return os.Remove(path)
	default:
		return fmt.Errorf("backend type %s is not supported", m.config.Backend.Type)
	}
}

// Migrate initialized the migration, migrates and cleans up after itself
func Migrate(c *config.Config) error {
	m, err := New(c)

	if err != nil {
		return err
	}

	if err := m.Init(); err != nil {
		return err
	}

	if err := m.Migrate(); err != nil {
		return err
	}

	if err := m.Cleanup(); err != nil {
		return err
	}

	return nil
}
