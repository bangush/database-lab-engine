/*
Provision wrapper

2019 © Postgres.ai
*/

package provision

import (
	"fmt"
	"strings"
	"time"

	"github.com/pkg/errors"

	"gitlab.com/postgres-ai/database-lab/pkg/log"
)

const (
	// ModeZfs defines provisioning for ZFS.
	ModeZfs = "zfs"
)

// NoRoomError defines a specific error type.
type NoRoomError struct {
	msg string
}

type State struct {
	InstanceID        string
	InstanceIP        string
	DockerContainerID string
	SessionID         string
}

type Session struct {
	ID   string
	Name string

	// Database.
	Host     string
	Port     uint
	User     string
	Password string

	// For user-defined username and password.
	ephemeralUser     string
	ephemeralPassword string
}

type Config struct {
	Mode string `yaml:"mode"`

	ModeZfs ModeZfsConfig `yaml:"zfs"`

	// Postgres options.
	PgVersion    string `yaml:"pgVersion"`
	PgBindir     string `yaml:"pgBindir"`
	PgDataSubdir string `yaml:"pgDataSubdir"`

	// Database user will be created with the specified credentials.
	DbUsername string
	DbPassword string
}

// TODO(anatoly): Merge with disk from models?
type Disk struct {
	Size     uint64
	Free     uint64
	DataSize uint64
}

type Snapshot struct {
	ID          string
	CreatedAt   time.Time
	DataStateAt time.Time
}

type SessionState struct {
	CloneSize uint64
}

type Provision interface {
	Init() error
	Reinit() error

	StartSession(string, string, ...string) (*Session, error)
	StopSession(*Session) error
	ResetSession(*Session, ...string) error

	CreateSnapshot(string) error
	GetSnapshots() ([]*Snapshot, error)

	RunPsql(*Session, string) (string, error)

	GetDiskState() (*Disk, error)
	GetSessionState(*Session) (*SessionState, error)
}

type provision struct {
	config Config
}

func NewProvision(config Config) (Provision, error) {
	// nolint
	switch config.Mode {
	case ModeZfs:
		log.Dbg("Using ZFS mode.")
		return NewProvisionModeZfs(config)
	}

	return nil, errors.New("unsupported mode specified")
}

// Check validity of a configuration and show a message for each violation.
func IsValidConfig(c Config) bool {
	result := true

	if len(c.PgVersion) == 0 && len(c.PgBindir) == 0 {
		log.Err("Either pgVersion or pgBindir should be set.")
		result = false
	}

	if len(c.PgBindir) > 0 && strings.HasSuffix(c.PgBindir, "/") {
		log.Err("Remove tailing slash from pgBindir.")
	}

	switch c.Mode {
	case ModeZfs:
		result = result && isValidConfigModeZfs(c)
	default:
		log.Err("Unsupported mode specified.")
		result = false
	}

	return result
}

func (s *Session) GetConnStr(dbname string) string {
	connStr := "sslmode=disable"

	if len(s.Host) > 0 {
		connStr += " host=" + s.Host
	}

	if s.Port > 0 {
		connStr += fmt.Sprintf(" port=%d", s.Port)
	}

	if len(s.User) > 0 {
		connStr += " user=" + s.User
	}

	if len(s.Password) > 0 {
		connStr += " password=" + s.Password
	}

	if len(dbname) > 0 {
		connStr += " dbname=" + dbname
	}

	return connStr
}

// NewNoRoomError instances a new NoRoomError.
func NewNoRoomError(errorMessage string) error {
	return &NoRoomError{msg: errorMessage}
}

func (e *NoRoomError) Error() string {
	// TODO(anatoly): Change message.
	return "session cannot be started because there is no room: " + e.msg
}