package api

import "encoding/json"

type (
	Create struct {
		Entity json.RawMessage `json:"entity"`
	}

	Read struct {
		ID      string                       `json:"id"`
		Filters map[string]map[string]string `json:"filters,omitempty"`
	}

	Update struct {
		ID      string                       `json:"id"`
		Entity  json.RawMessage              `json:"entity"`
		Filters map[string]map[string]string `json:"filters,omitempty"`
	}

	Delete struct {
		ID      string                       `json:"id"`
		Filters map[string]map[string]string `json:"filters,omitempty"`
	}

	Search struct {
		Skip    int                          `json:"skip"`
		Take    int                          `json:"take"`
		Where   map[string]string            `json:"where,omitempty"`
		Filters map[string]map[string]string `json:"filters,omitempty"`
	}

	Patch struct {
		ID      string                       `json:"id"`
		Data    map[string]any               `json:"data"`
		Filters map[string]map[string]string `json:"filters,omitempty"`
	}
)
