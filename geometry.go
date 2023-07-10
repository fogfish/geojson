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
	BoundingBox() BoundingBox
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

// BoundingBox around the point
func (geo *Point) BoundingBox() BoundingBox {
	if len(geo.Coords) == 0 {
		return nil
	}

	return boundingBox(geo.Coords, geo.Coords)
}

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

// BoundingBox around MultiPoint
func (geo *MultiPoint) BoundingBox() BoundingBox {
	if len(geo.Coords) == 0 {
		return nil
	}

	return boundingBox(geo.Coords[0], geo.Coords)
}

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

// BoundingBox around LineString
func (geo *LineString) BoundingBox() BoundingBox {
	if len(geo.Coords) == 0 {
		return nil
	}

	return boundingBox(geo.Coords[0], geo.Coords)
}

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

// BoundingBox around MultiLineString
func (geo *MultiLineString) BoundingBox() BoundingBox {
	if len(geo.Coords) == 0 || len(geo.Coords[0]) == 0 {
		return nil
	}

	return boundingBox(geo.Coords[0][0], geo.Coords)
}

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

// BoundingBox around Polygon
func (geo *Polygon) BoundingBox() BoundingBox {
	if len(geo.Coords) == 0 || len(geo.Coords[0]) == 0 {
		return nil
	}

	return boundingBox(geo.Coords[0][0], geo.Coords)
}

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

// BoundingBox around MultiPolygon
func (geo *MultiPolygon) BoundingBox() BoundingBox {
	if len(geo.Coords) == 0 || len(geo.Coords[0]) == 0 {
		return nil
	}

	bbox := boundingBox(geo.Coords[0][0][0], geo.Coords[0])
	for i := 0; i < len(geo.Coords); i++ {
		surface := geo.Coords[i]
		if len(surface) > 0 && len(surface[0]) > 0 {
			bbox.Join(boundingBox(surface[0][0], surface))
		}
	}

	return bbox
}

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
