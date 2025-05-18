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
		ID     string          `json:"id"`
		Entity json.RawMessage `json:"entity"`
	}

	Delete struct {
		ID string `json:"id"`
	}

	Search struct {
		Skip    int                          `json:"skip"`
		Take    int                          `json:"take"`
		Where   map[string]string            `json:"where,omitempty"`
		Sort    map[string]string            `json:"sort,omitempty"`
		Preload map[string]string            `json:"preload,omitempty"`
		Filters map[string]map[string]string `json:"filters,omitempty"`
	}

	Patch struct {
		ID      string            `json:"id"`
		Data    map[string]any    `json:"data"`
		Preload map[string]string `json:"preload,omitempty"`
	}
)
