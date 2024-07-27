package whisper

import (
	"encoding/json"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

//////////////////////////////////////////////////////////////////////////////
// TYPES

type Model struct {
	Id      string `json:"id" writer:",width:28,wrap"`
	Object  string `json:"object,omitempty" writer:"-"`
	Path    string `json:"path,omitempty" writer:",width:40,wrap"`
	Created int64  `json:"created,omitempty"`
	OwnedBy string `json:"owned_by,omitempty"`
}

type models struct {
	sync.RWMutex

	// Path to the models directory
	path string

	// list of all models
	models []*Model
}

//////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (m *Model) String() string {
	data, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return err.Error()
	}
	return string(data)
}

//////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

// Return a model by its Id
func (m *models) ById(name string) *Model {
	m.RLock()
	defer m.RUnlock()
	name = modelNameToId(name)
	for _, model := range m.models {
		if model.Id == name {
			return model
		}
	}
	return nil
}

// Return a model by name and path
func (m *models) ByPath(name, dest string) *Model {
	m.RLock()
	defer m.RUnlock()
	path := filepath.Join(dest, name)
	for _, model := range m.models {
		if model.Path == path {
			return model
		}
	}
	return nil
}

// Rescan models directory
func (m *models) Rescan() error {
	m.Lock()
	defer m.Unlock()
	if models, err := listModels(m.path); err != nil {
		return err
	} else {
		m.models = models
	}
	return nil
}

//////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func listModels(path string) ([]*Model, error) {
	result := make([]*Model, 0, 100)

	// Walk filesystem
	return result, fs.WalkDir(os.DirFS(path), ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}

		// Ignore hidden files or files without a .bin extension
		if strings.HasPrefix(d.Name(), ".") {
			return nil
		}
		if filepath.Ext(d.Name()) != extModel {
			return nil
		}

		// Ignore files we can't get information on
		info, err := d.Info()
		if err != nil {
			return nil
		}

		// Ignore non-regular files
		if !d.Type().IsRegular() {
			return nil
		}

		// Ignore files less than 8MB
		if info.Size() < 8*1024*1024 {
			return nil
		}

		// Get model information
		model := new(Model)
		model.Object = "model"
		model.Path = path
		model.Created = info.ModTime().Unix()

		// Generate an Id for the model
		model.Id = modelNameToId(filepath.Base(path))

		// Append to result
		result = append(result, model)

		// Continue walking
		return nil
	})
}

func modelNameToId(name string) string {
	// We replace all non-alphanumeric characters with underscores
	return strings.Map(func(r rune) rune {
		if r >= 'a' && r <= 'z' {
			return r
		}
		if r >= 'A' && r <= 'Z' {
			return r
		}
		if r >= '0' && r <= '9' {
			return r
		}
		if r == '.' || r == '-' {
			return r
		}
		return '_'
	}, name)
}
