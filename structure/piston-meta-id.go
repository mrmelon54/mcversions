package structure

import (
	"encoding/json"
	"github.com/Masterminds/semver/v3"
)

type PistonMetaId struct {
	normal  bool
	version *semver.Version
	other   string
}

func NewPistonMetaId(id string) (*PistonMetaId, error) {
	ver, err := semver.NewVersion(id)
	if err == nil {
		return &PistonMetaId{normal: true, version: ver, other: id}, nil
	}
	return &PistonMetaId{other: id}, nil
}

func (id PistonMetaId) String() string {
	return id.other
}

func (id *PistonMetaId) MarshalJSON() ([]byte, error) {
	return json.Marshal(id.other)
}

func (id *PistonMetaId) UnmarshalJSON(data []byte) error {
	id.normal = false
	err := json.Unmarshal(data, &id.other)
	if err != nil {
		return err
	}
	ver, err := semver.NewVersion(id.other)
	if err == nil {
		id.normal = true
		id.version = ver
	}
	return nil
}

func (id *PistonMetaId) Equal(other *PistonMetaId) bool {
	if id.normal && other.normal {
		return id.version.Equal(other.version)
	}
	return id.String() == other.String()
}

func PistonMetaIdCheckConstraints(id *PistonMetaId, con *semver.Constraints) bool {
	if id.normal {
		return con.Check(id.version)
	}
	return false
}
