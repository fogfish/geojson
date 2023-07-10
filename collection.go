//
// Copyright (C) 2021 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/fogfish/geojson
//

package geojson

import (
	"encoding/json"
)

// Features Collection Type
type Collection[T interface{ BoundingBox() BoundingBox }] struct {
	Features []T
}

// BoundingBox of the features collection
func (c Collection[T]) BoundingBox() BoundingBox {
	if len(c.Features) == 0 {
		return nil
	}

	bbox := c.Features[0].BoundingBox()
	for i := 1; i < len(c.Features); i++ {
		bbox.Join(c.Features[i].BoundingBox())
	}

	return bbox
}

// Encodes FeatureCollection GeoJSON
func (c Collection[T]) MarshalJSON() ([]byte, error) {
	features, err := json.Marshal(c.Features)
	if err != nil {
		return nil, err
	}

	any := struct {
		Type     string          `json:"type"`
		BBox     BoundingBox     `json:"bbox,omitempty"`
		Features json.RawMessage `json:"features,omitempty"`
	}{
		Type:     "FeatureCollection",
		BBox:     c.BoundingBox(),
		Features: features,
	}

	return json.Marshal(any)
}

// Decodes FeatureCollection GeoJSON
func (c *Collection[T]) UnmarshalJSON(b []byte) error {
	type anyCollection struct {
		Type     string          `json:"type"`
		Features json.RawMessage `json:"features,omitempty"`
	}

	any := anyCollection{}

	if err := json.Unmarshal(b, &any); err != nil {
		return err
	}

	if any.Type != "FeatureCollection" {
		return ErrorUnsupportedType
	}

	if any.Features != nil {
		if err := json.Unmarshal(any.Features, &c.Features); err != nil {
			return err
		}
	}

	return nil
}
