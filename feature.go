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

const TYPE_FEATURE = "Feature"

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
	Geometry Geometry `json:"-"`
}

func (fea Feature) BoundingBox() BoundingBox {
	if fea.Geometry == nil {
		return nil
	}

	return fea.Geometry.BoundingBox()
}

// EncodeGeoJSON is a helper function to implement GeoJSON codec
//
//	func (x MyType) MarshalJSON() ([]byte, error) {
//	  type tStruct MyType
//	  return x.Feature.EncodeGeoJSON(x.ID, tStruct(x))
//	}
func (fea Feature) EncodeGeoJSON(id string, props any) ([]byte, error) {
	properties, err := json.Marshal(props)
	if err != nil {
		return nil, err
	}

	// Note: skip bounding box for the point.
	var bbox BoundingBox
	switch fea.Geometry.(type) {
	case nil:
		bbox = nil
	case *Point:
		bbox = nil
	default:
		bbox = fea.Geometry.BoundingBox()
	}

	val := struct {
		Type       string          `json:"type"`
		BBox       BoundingBox     `json:"bbox,omitempty"`
		ID         string          `json:"id,omitempty"`
		Geometry   Geometry        `json:"geometry"`
		Properties json.RawMessage `json:"properties,omitempty"`
	}{
		ID:         id,
		Type:       TYPE_FEATURE,
		BBox:       bbox,
		Geometry:   fea.Geometry,
		Properties: properties,
	}

	return json.Marshal(val)
}

// anyGeoJSON is an internal type used for decode of GeoJSON
type anyGeoJSON struct {
	Type       string          `json:"type"`
	ID         string          `json:"id,omitempty"`
	Geometry   json.RawMessage `json:"geometry"`
	Properties json.RawMessage `json:"properties,omitempty"`
}

// DecodeGeoJSON is a helper function to implement GeoJSON codec
//
//	func (x *MyType) UnmarshalJSON(b []byte) (err error) {
//	  type tStruct *MyType
//	  x.ID, err = x.Feature.DecodeGeoJSON(b, tStruct(x))
//	  return
//	}
func (fea *Feature) DecodeGeoJSON(bytes []byte, props any) (string, error) {
	obj := anyGeoJSON{}

	if err := json.Unmarshal(bytes, &obj); err != nil {
		return "", err
	}

	if obj.Type != TYPE_FEATURE {
		return "", ErrUnsupportedType
	}

	return fea.decodeAnyGeoJSON(&obj, props)
}

func (fea *Feature) decodeAnyGeoJSON(obj *anyGeoJSON, props any) (string, error) {
	if obj.Geometry != nil {
		geo, err := decodeGeometry(obj.Geometry)
		if err != nil {
			return "", err
		}
		fea.Geometry = geo
	}

	if obj.Properties != nil {
		if err := json.Unmarshal(obj.Properties, &props); err != nil {
			return "", err
		}
	}

	return obj.ID, nil
}

// New Feature from Geometry
func New(geometry Geometry) Feature {
	return Feature{Geometry: geometry}
}

// NewPoint ⟼ Feature[Point]
func NewPoint(coords Coord) Feature {
	return Feature{
		Geometry: &Point{Coords: coords},
	}
}

// NewMultiPoint ⟼ Feature[MultiPoint]
func NewMultiPoint(coords Curve) Feature {
	return Feature{
		Geometry: &MultiPoint{Coords: coords},
	}
}

// NewLineString ⟼ Feature[LineString]
func NewLineString(coords Curve) Feature {
	return Feature{
		Geometry: &LineString{Coords: coords},
	}
}

// NewMultiLineString ⟼ Feature[MultiLineString]
func NewMultiLineString(coords Surface) Feature {
	return Feature{
		Geometry: &MultiLineString{Coords: coords},
	}
}

// NewPolygon ⟼ Feature[Polygon]
func NewPolygon(coords Surface) Feature {
	return Feature{
		Geometry: &Polygon{Coords: coords},
	}
}

// NewMultiPolygon ⟼ Feature[MultiPolygon]
func NewMultiPolygon(coords ...Surface) Feature {
	return Feature{
		Geometry: &MultiPolygon{Coords: coords},
	}
}
