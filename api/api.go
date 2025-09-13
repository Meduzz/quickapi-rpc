package api

import "encoding/json"

type (
	Create struct {
		Entity json.RawMessage `json:"entity"`
	}

	Read struct {
		ID      string            `json:"id"`
		Preload map[string]string `json:"preload,omitempty"`
	}

	Update struct {
		ID     string                       `json:"id"`
		Entity json.RawMessage              `json:"entity"`
		Hooks  map[string]map[string]string `json:"hooks,omitempty"`
	}

	Delete struct {
		ID    string                       `json:"id"`
		Hooks map[string]map[string]string `json:"hooks,omitempty"`
	}

	Search struct {
		Skip    int                          `json:"skip"`
		Take    int                          `json:"take"`
		Where   map[string]string            `json:"where,omitempty"`
		Sort    map[string]string            `json:"sort,omitempty"`
		Preload map[string]string            `json:"preload,omitempty"`
		Hooks   map[string]map[string]string `json:"hooks,omitempty"`
	}

	Patch struct {
		ID      string                       `json:"id"`
		Data    map[string]any               `json:"data"`
		Preload map[string]string            `json:"preload,omitempty"`
		Hooks   map[string]map[string]string `json:"hooks,omitempty"`
	}
)
