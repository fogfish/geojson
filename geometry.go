//
// Copyright (C) 2021 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/fogfish/geojson
//

package geojson

import "encoding/json"

type geometryType string

const (
	typePoint           = geometryType("Point")
	typeMultiPoint      = geometryType("MultiPoint")
	typeLineString      = geometryType("LineString")
	typeMultiLineString = geometryType("MultiLineString")
	typePolygon         = geometryType("Polygon")
	typeMultiPolygon    = geometryType("MultiPolygon")
)

// Geometry Object represents points, curves, and surfaces in coordinate space.
// It MUST be one of the seven geometry types.
type Geometry interface {
	MarshalGeoJSON() ([]byte, error)
	UnmarshalGeoJSON(b []byte) error
}

// UnmarshalJSON decodes Geometry from GeoJSON
func decodeGeometry(b []byte) (Geometry, error) {
	var gen struct {
		Type   geometryType    `json:"type"`
		Coords json.RawMessage `json:"coordinates"`
	}
	if err := json.Unmarshal(b, &gen); err != nil {
		return nil, err
	}

	var geo Geometry

	switch gen.Type {
	case typePoint:
		geo = &Point{}
	case typeMultiPoint:
		geo = &MultiPoint{}
	case typeLineString:
		geo = &LineString{}
	case typeMultiLineString:
		geo = &MultiLineString{}
	case typePolygon:
		geo = &Polygon{}
	case typeMultiPolygon:
		geo = &MultiPolygon{}
	default:
		return nil, ErrorUnsupportedType
	}

	err := geo.UnmarshalGeoJSON(gen.Coords)
	return geo, err
}

// Point type, the "coordinates" member is a single position.
type Point struct{ Coords Coord }

// MarshalGeoJSON encodes geometry type to GeoJSON
func (geo *Point) MarshalGeoJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"type":        typePoint,
		"coordinates": geo.Coords,
	})
}

// UnmarshalGeoJSON decodes geometry type from GeoJSON
func (geo *Point) UnmarshalGeoJSON(b []byte) error {
	if err := json.Unmarshal(b, &geo.Coords); err != nil {
		return err
	}
	return nil
}

// MultiPoint type, the "coordinates" member is an array of positions.
type MultiPoint struct{ Coords Curve }

// MarshalGeoJSON encodes geometry type to GeoJSON
func (geo *MultiPoint) MarshalGeoJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"type":        typeMultiPoint,
		"coordinates": geo.Coords,
	})
}

// UnmarshalGeoJSON decodes geometry type from GeoJSON
func (geo *MultiPoint) UnmarshalGeoJSON(b []byte) error {
	if err := json.Unmarshal(b, &geo.Coords); err != nil {
		return err
	}
	return nil
}

// LineString type, the "coordinates" member is an array of two or
// more positions.
type LineString struct{ Coords Curve }

// MarshalGeoJSON encodes geometry type to GeoJSON
func (geo *LineString) MarshalGeoJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"type":        typeLineString,
		"coordinates": geo.Coords,
	})
}

// UnmarshalGeoJSON decodes geometry type from GeoJSON
func (geo *LineString) UnmarshalGeoJSON(b []byte) error {
	if err := json.Unmarshal(b, &geo.Coords); err != nil {
		return err
	}
	return nil
}

// MultiLineString type, the "coordinates" member is an array of
// LineString coordinate arrays.
type MultiLineString struct{ Coords Surface }

// MarshalGeoJSON encodes geometry type to GeoJSON
func (geo *MultiLineString) MarshalGeoJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"type":        typeMultiLineString,
		"coordinates": geo.Coords,
	})
}

// UnmarshalGeoJSON decodes geometry type from GeoJSON
func (geo *MultiLineString) UnmarshalGeoJSON(b []byte) error {
	if err := json.Unmarshal(b, &geo.Coords); err != nil {
		return err
	}
	return nil
}

// Polygon is combinaton of exterior and interior linear rings,
// the first element is exterior ring, and others are interior rings.
//
// A linear ring is a closed LineString with four or more positions.
// The first and last positions are equivalent, and they MUST contain
// identical values; their representation SHOULD also be identical.
type Polygon struct{ Coords Surface }

// MarshalGeoJSON encodes geometry type to GeoJSON
func (geo *Polygon) MarshalGeoJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"type":        typePolygon,
		"coordinates": geo.Coords,
	})
}

// UnmarshalGeoJSON decodes geometry type from GeoJSON
func (geo *Polygon) UnmarshalGeoJSON(b []byte) error {
	if err := json.Unmarshal(b, &geo.Coords); err != nil {
		return err
	}
	return nil
}

// MultiPolygon type, the "coordinates" member is an array of
// Polygon coordinate arrays.
type MultiPolygon struct{ Coords []Surface }

// MarshalGeoJSON encodes geometry type to GeoJSON
func (geo *MultiPolygon) MarshalGeoJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"type":        typeMultiPolygon,
		"coordinates": geo.Coords,
	})
}

// UnmarshalGeoJSON decodes geometry type from GeoJSON
func (geo *MultiPolygon) UnmarshalGeoJSON(b []byte) error {
	if err := json.Unmarshal(b, &geo.Coords); err != nil {
		return err
	}
	return nil
}
