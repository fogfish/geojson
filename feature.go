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

	"github.com/fogfish/curie/v2"
)

// Feature object represents a spatially bounded thing.
// This object contains geometry, a common identifier, and properties.
// The value of the properties is any JSON object, typically defined by
// an application.
//
// The library uses a type safe notation for the feature's property
// definition instead of generic interface{} type. It uses type tagging
// technique (or embedding):
//
//	type MyType struct {
//	  geojson.Feature
//	  Name      string `json:"name,omitempty"`
//	}
type Feature struct {
	ID       curie.IRI `json:"-"`
	Geometry Geometry  `json:"-"`
}

func (fea Feature) BoundingBox() BoundingBox { return fea.Geometry.BoundingBox() }

// EncodeGeoJSON is a helper function to implement GeoJSON codec
//
//	func (x MyType) MarshalJSON() ([]byte, error) {
//		type tStruct MyType
//		return x.Feature.EncodeGeoJSON(tStruct(x))
//	}
func (fea Feature) EncodeGeoJSON(props any) ([]byte, error) {
	properties, err := json.Marshal(props)
	if err != nil {
		return nil, err
	}

	geo := fea.Geometry
	if geo == nil {
		geo = &Point{Coords: Coord{}}
	}

	val := struct {
		Type       string          `json:"type"`
		BBox       BoundingBox     `json:"bbox,omitempty"`
		ID         curie.IRI       `json:"id,omitempty"`
		Geometry   Geometry        `json:"geometry,omitempty"`
		Properties json.RawMessage `json:"properties,omitempty"`
	}{
		ID:         fea.ID,
		Type:       "Feature",
		BBox:       geo.BoundingBox(),
		Geometry:   geo,
		Properties: properties,
	}

	return json.Marshal(val)
}

// anyGeoJSON is an internal type used for decode of GeoJSON
type anyGeoJSON struct {
	Type       string          `json:"type"`
	ID         curie.IRI       `json:"id,omitempty"`
	Geometry   json.RawMessage `json:"geometry,omitempty"`
	Properties json.RawMessage `json:"properties,omitempty"`
}

// DecodeGeoJSON is a helper function to implement GeoJSON codec
//
//	func (x *MyType) UnmarshalJSON(b []byte) error {
//		type tStruct *MyType
//		return x.Feature.DecodeGeoJSON(b, tStruct(x))
//	}
func (fea *Feature) DecodeGeoJSON(bytes []byte, props interface{}) error {
	any := anyGeoJSON{}

	if err := json.Unmarshal(bytes, &any); err != nil {
		return err
	}

	if any.Type != "Feature" {
		return ErrorUnsupportedType
	}

	return fea.decodeAnyGeoJSON(&any, props)
}

func (fea *Feature) decodeAnyGeoJSON(any *anyGeoJSON, props interface{}) error {
	if any.Geometry != nil {
		geo, err := decodeGeometry(any.Geometry)
		if err != nil {
			return err
		}
		fea.Geometry = geo
	}

	if any.Properties != nil {
		if err := json.Unmarshal(any.Properties, &props); err != nil {
			return err
		}
	}

	fea.ID = any.ID
	return nil
}

// New Feature from Geometry
func New(id curie.IRI, geometry Geometry) Feature {
	return Feature{ID: id, Geometry: geometry}
}

// NewPoint ⟼ Feature[Point]
func NewPoint(id curie.IRI, coords Coord) Feature {
	return Feature{
		ID:       id,
		Geometry: &Point{Coords: coords},
	}
}

// NewMultiPoint ⟼ Feature[MultiPoint]
func NewMultiPoint(id curie.IRI, coords Curve) Feature {
	return Feature{
		ID:       id,
		Geometry: &MultiPoint{Coords: coords},
	}
}

// NewLineString ⟼ Feature[LineString]
func NewLineString(id curie.IRI, coords Curve) Feature {
	return Feature{
		ID:       id,
		Geometry: &LineString{Coords: coords},
	}
}

// NewMultiLineString ⟼ Feature[MultiLineString]
func NewMultiLineString(id curie.IRI, coords Surface) Feature {
	return Feature{
		ID:       id,
		Geometry: &MultiLineString{Coords: coords},
	}
}

// NewPolygon ⟼ Feature[Polygon]
func NewPolygon(id curie.IRI, coords Surface) Feature {
	return Feature{
		ID:       id,
		Geometry: &Polygon{Coords: coords},
	}
}

// NewMultiPolygon ⟼ Feature[MultiPolygon]
func NewMultiPolygon(id curie.IRI, coords ...Surface) Feature {
	return Feature{
		ID:       id,
		Geometry: &MultiPolygon{Coords: coords},
	}
}
