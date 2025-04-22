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
	"fmt"
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
	Geometry() Shape
	BoundingBox() BoundingBox
	unmarshalGeoJSON(b []byte) error
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
	case "":
		return nil, nil
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
		return nil, fmt.Errorf("type %s is not supported as GeoJSON %s", gen.Type, "Geometry")
	}

	err := geo.unmarshalGeoJSON(gen.Coords)
	return geo, err
}

// Point type, the "coordinates" member is a single position.
type Point struct {
	Coords Coord `json:"coordinates"`
}

func (geo *Point) Geometry() Shape { return geo.Coords }

// BoundingBox around the point
func (geo *Point) BoundingBox() BoundingBox {
	if len(geo.Coords) == 0 {
		return nil
	}

	return boundingBox(geo.Coords, geo.Coords)
}

// Encode Point Geometry to GeoJSON format
func (geo *Point) MarshalJSON() ([]byte, error) {
	type Struct Point
	return json.Marshal(&struct {
		Type geometryType `json:"type"`
		*Struct
	}{
		Type:   typePoint,
		Struct: (*Struct)(geo),
	})
}

// Decode Point Geometry from GeoJSON format
func (geo *Point) UnmarshalJSON(b []byte) error {
	type Struct Point
	var bag struct {
		Type geometryType `json:"type"`
		*Struct
	}

	if err := json.Unmarshal(b, &bag); err != nil {
		return err
	}

	if bag.Type != typePoint {
		return fmt.Errorf("type %s is not supported as GeoJSON %s", bag.Type, typePoint)
	}

	*geo = (Point)(*bag.Struct)
	return nil
}

// UnmarshalGeoJSON decodes geometry type from GeoJSON
func (geo *Point) unmarshalGeoJSON(b []byte) error {
	if err := json.Unmarshal(b, &geo.Coords); err != nil {
		return err
	}
	return nil
}

// MultiPoint type, the "coordinates" member is an array of positions.
type MultiPoint struct {
	Coords Curve `json:"coordinates"`
}

func (geo *MultiPoint) Geometry() Shape { return geo.Coords }

// BoundingBox around MultiPoint
func (geo *MultiPoint) BoundingBox() BoundingBox {
	if len(geo.Coords) == 0 {
		return nil
	}

	return boundingBox(geo.Coords[0], geo.Coords)
}

// Encode MultiPoint Geometry to GeoJSON format
func (geo *MultiPoint) MarshalJSON() ([]byte, error) {
	type Struct MultiPoint
	return json.Marshal(&struct {
		Type geometryType `json:"type"`
		*Struct
	}{
		Type:   typeMultiPoint,
		Struct: (*Struct)(geo),
	})
}

// Decode MultiPoint Geometry from GeoJSON format
func (geo *MultiPoint) UnmarshalJSON(b []byte) error {
	type Struct MultiPoint
	var bag struct {
		Type geometryType `json:"type"`
		*Struct
	}

	if err := json.Unmarshal(b, &bag); err != nil {
		return err
	}

	if bag.Type != typeMultiPoint {
		return fmt.Errorf("type %s is not supported as GeoJSON %s", bag.Type, typePoint)
	}

	*geo = (MultiPoint)(*bag.Struct)
	return nil
}

// UnmarshalGeoJSON decodes geometry type from GeoJSON
func (geo *MultiPoint) unmarshalGeoJSON(b []byte) error {
	if err := json.Unmarshal(b, &geo.Coords); err != nil {
		return err
	}
	return nil
}

// LineString type, the "coordinates" member is an array of two or
// more positions.
type LineString struct {
	Coords Curve `json:"coordinates"`
}

func (geo *LineString) Geometry() Shape { return geo.Coords }

// BoundingBox around LineString
func (geo *LineString) BoundingBox() BoundingBox {
	if len(geo.Coords) == 0 {
		return nil
	}

	return boundingBox(geo.Coords[0], geo.Coords)
}

// Encode Point Geometry to GeoJSON format
func (geo *LineString) MarshalJSON() ([]byte, error) {
	type Struct LineString
	return json.Marshal(&struct {
		Type geometryType `json:"type"`
		*Struct
	}{
		Type:   typeLineString,
		Struct: (*Struct)(geo),
	})
}

// Decode Point Geometry from GeoJSON format
func (geo *LineString) UnmarshalJSON(b []byte) error {
	type Struct LineString
	var bag struct {
		Type geometryType `json:"type"`
		*Struct
	}

	if err := json.Unmarshal(b, &bag); err != nil {
		return err
	}

	if bag.Type != typeLineString {
		return fmt.Errorf("type %s is not supported as GeoJSON %s", bag.Type, typePoint)
	}

	*geo = (LineString)(*bag.Struct)
	return nil
}

// UnmarshalGeoJSON decodes geometry type from GeoJSON
func (geo *LineString) unmarshalGeoJSON(b []byte) error {
	if err := json.Unmarshal(b, &geo.Coords); err != nil {
		return err
	}
	return nil
}

// MultiLineString type, the "coordinates" member is an array of
// LineString coordinate arrays.
type MultiLineString struct {
	Coords Surface `json:"coordinates"`
}

func (geo *MultiLineString) Geometry() Shape { return geo.Coords }

// BoundingBox around MultiLineString
func (geo *MultiLineString) BoundingBox() BoundingBox {
	if len(geo.Coords) == 0 || len(geo.Coords[0]) == 0 {
		return nil
	}

	return boundingBox(geo.Coords[0][0], geo.Coords)
}

// Encode MultiLineString Geometry to GeoJSON format
func (geo *MultiLineString) MarshalJSON() ([]byte, error) {
	type Struct MultiLineString
	return json.Marshal(&struct {
		Type geometryType `json:"type"`
		*Struct
	}{
		Type:   typeMultiLineString,
		Struct: (*Struct)(geo),
	})
}

// Decode MultiLineString Geometry from GeoJSON format
func (geo *MultiLineString) UnmarshalJSON(b []byte) error {
	type Struct MultiLineString
	var bag struct {
		Type geometryType `json:"type"`
		*Struct
	}

	if err := json.Unmarshal(b, &bag); err != nil {
		return err
	}

	if bag.Type != typeMultiLineString {
		return fmt.Errorf("type %s is not supported as GeoJSON %s", bag.Type, typePoint)
	}

	*geo = (MultiLineString)(*bag.Struct)
	return nil
}

// UnmarshalGeoJSON decodes geometry type from GeoJSON
func (geo *MultiLineString) unmarshalGeoJSON(b []byte) error {
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
type Polygon struct {
	Coords Surface `json:"coordinates"`
}

func (geo *Polygon) Geometry() Shape { return geo.Coords }

// BoundingBox around Polygon
func (geo *Polygon) BoundingBox() BoundingBox {
	if len(geo.Coords) == 0 || len(geo.Coords[0]) == 0 {
		return nil
	}

	return boundingBox(geo.Coords[0][0], geo.Coords)
}

// Encode Point Geometry to GeoJSON format
func (geo *Polygon) MarshalJSON() ([]byte, error) {
	type Struct Polygon
	return json.Marshal(&struct {
		Type geometryType `json:"type"`
		*Struct
	}{
		Type:   typePolygon,
		Struct: (*Struct)(geo),
	})
}

// Decode Point Geometry from GeoJSON format
func (geo *Polygon) UnmarshalJSON(b []byte) error {
	type Struct Polygon
	var bag struct {
		Type geometryType `json:"type"`
		*Struct
	}

	if err := json.Unmarshal(b, &bag); err != nil {
		return err
	}

	if bag.Type != typePolygon {
		return fmt.Errorf("type %s is not supported as GeoJSON %s", bag.Type, typePoint)
	}

	*geo = (Polygon)(*bag.Struct)
	return nil
}

// UnmarshalGeoJSON decodes geometry type from GeoJSON
func (geo *Polygon) unmarshalGeoJSON(b []byte) error {
	if err := json.Unmarshal(b, &geo.Coords); err != nil {
		return err
	}
	return nil
}

// MultiPolygon type, the "coordinates" member is an array of
// Polygon coordinate arrays.
type MultiPolygon struct {
	Coords Surfaces `json:"coordinates"`
}

func (geo *MultiPolygon) Geometry() Shape { return geo.Coords }

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

// Encode Point Geometry to GeoJSON format
func (geo *MultiPolygon) MarshalJSON() ([]byte, error) {
	type Struct MultiPolygon
	return json.Marshal(&struct {
		Type geometryType `json:"type"`
		*Struct
	}{
		Type:   typeMultiPolygon,
		Struct: (*Struct)(geo),
	})
}

// Decode Point Geometry from GeoJSON format
func (geo *MultiPolygon) UnmarshalJSON(b []byte) error {
	type Struct MultiPolygon
	var bag struct {
		Type geometryType `json:"type"`
		*Struct
	}

	if err := json.Unmarshal(b, &bag); err != nil {
		return err
	}

	if bag.Type != typeMultiPolygon {
		return fmt.Errorf("type %s is not supported as GeoJSON %s", bag.Type, typePoint)
	}

	*geo = (MultiPolygon)(*bag.Struct)
	return nil
}

// UnmarshalGeoJSON decodes geometry type from GeoJSON
func (geo *MultiPolygon) unmarshalGeoJSON(b []byte) error {
	if err := json.Unmarshal(b, &geo.Coords); err != nil {
		return err
	}
	return nil
}
