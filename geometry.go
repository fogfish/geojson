//
// Copyright (C) 2021 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/fogfish/geojson
//

package geojson

import "encoding/json"

/*

Geometry Object represents points, curves, and surfaces in
coordinate space.
*/
type Geometry struct{ Coords }

/*

Coords is a member of a Geometry Object. It MUST be one of the
seven geometry types.
*/
type Coords interface {
	MarshalGeoJSON() ([]byte, error)
	UnmarshalGeoJSON(b []byte) error
}

/*

MarshalJSON encodes Geometry to GeoJSON
*/
func (geo Geometry) MarshalJSON() ([]byte, error) {
	return geo.Coords.MarshalGeoJSON()
}

/*

UnmarshalJSON decodes Geometry from GeoJSON
*/
func (geo *Geometry) UnmarshalJSON(b []byte) error {
	var gen struct {
		Type   string          `json:"type"`
		Coords json.RawMessage `json:"coordinates"`
	}
	if err := json.Unmarshal(b, &gen); err != nil {
		return err
	}

	switch gen.Type {
	case "Point":
		geo.Coords = &Point{}
		return geo.Coords.UnmarshalGeoJSON(gen.Coords)
	case "MultiPoint":
		geo.Coords = &MultiPoint{}
		return geo.Coords.UnmarshalGeoJSON(gen.Coords)
	case "LineString":
		geo.Coords = &LineString{}
		return geo.Coords.UnmarshalGeoJSON(gen.Coords)
	case "MultiLineString":
		geo.Coords = &MultiLineString{}
		return geo.Coords.UnmarshalGeoJSON(gen.Coords)
	case "Polygon":
		geo.Coords = &Polygon{}
		return geo.Coords.UnmarshalGeoJSON(gen.Coords)
	case "MultiPolygon":
		geo.Coords = &MultiPolygon{}
		return geo.Coords.UnmarshalGeoJSON(gen.Coords)
	default:
		return ErrorUnsupportedType
	}
}

/*

Point type, the "coordinates" member is a single position.
*/
type Point struct {
	Coords Position
}

func (geo *Point) MarshalGeoJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"type":        "Point",
		"coordinates": geo.Coords,
	})
}

func (geo *Point) UnmarshalGeoJSON(b []byte) error {
	if err := json.Unmarshal(b, &geo.Coords); err != nil {
		return err
	}
	return nil
}

/*

MultiPoint type, the "coordinates" member is an array of positions.
*/
type MultiPoint struct {
	Coords Sequence
}

func (geo *MultiPoint) MarshalGeoJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"type":        "MultiPoint",
		"coordinates": geo.Coords,
	})
}

func (geo *MultiPoint) UnmarshalGeoJSON(b []byte) error {
	if err := json.Unmarshal(b, &geo.Coords); err != nil {
		return err
	}
	return nil
}

/*

LineString type, the "coordinates" member is an array of two or
more positions.
*/
type LineString struct {
	Coords Sequence
}

func (geo *LineString) MarshalGeoJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"type":        "LineString",
		"coordinates": geo.Coords,
	})
}

func (geo *LineString) UnmarshalGeoJSON(b []byte) error {
	if err := json.Unmarshal(b, &geo.Coords); err != nil {
		return err
	}
	return nil
}

/*

MultiLineString type, the "coordinates" member is an array of
LineString coordinate arrays.
*/
type MultiLineString struct {
	Coords Surface
}

func (geo *MultiLineString) MarshalGeoJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"type":        "MultiLineString",
		"coordinates": geo.Coords,
	})
}

func (geo *MultiLineString) UnmarshalGeoJSON(b []byte) error {
	if err := json.Unmarshal(b, &geo.Coords); err != nil {
		return err
	}
	return nil
}

/*

Polygon is combinaton of exterior and interior linear rings,
the first element is exterior ring, and others are interior rings.

A linear ring is a closed LineString with four or more positions.
The first and last positions are equivalent, and they MUST contain
identical values; their representation SHOULD also be identical.
*/
type Polygon struct {
	Coords Surface
}

func (geo *Polygon) MarshalGeoJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"type":        "Polygon",
		"coordinates": geo.Coords,
	})
}

func (geo *Polygon) UnmarshalGeoJSON(b []byte) error {
	if err := json.Unmarshal(b, &geo.Coords); err != nil {
		return err
	}
	return nil
}

/*

MultiPolygon type, the "coordinates" member is an array of
Polygon coordinate arrays.
*/
type MultiPolygon struct {
	Coords []Surface
}

func (geo *MultiPolygon) MarshalGeoJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"type":        "MultiPolygon",
		"coordinates": geo.Coords,
	})
}

func (geo *MultiPolygon) UnmarshalGeoJSON(b []byte) error {
	if err := json.Unmarshal(b, &geo.Coords); err != nil {
		return err
	}
	return nil
}
