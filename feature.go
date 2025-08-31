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
// definition instead of generic any type. It uses type tagging
// technique (or embedding):
//
//	type MyType struct {
//	  geojson.Feature
//	  Name      string `json:"name,omitempty"`
//	}
type Feature struct {
	Geometry Geometry `json:"-"`
}

func (fea Feature) unapply() Geometry { return fea.Geometry }
func (fea *Feature) apply(g Geometry) { fea.Geometry = g }

// Return the bounding box of the feature's geometry.
func (fea Feature) BoundingBox() BoundingBox {
	if fea.Geometry == nil {
		return nil
	}

	return fea.Geometry.BoundingBox()
}

// IFeature is a phantom type of the Feature itself, used for type-safe marshaling and unmarshaling.
type IFeature interface {
	BoundingBox() BoundingBox
	unapply() Geometry
	apply(Geometry)
}

//
// Encoder
//

// EncodeGeoJSON is a helper function to implement GeoJSON codec
//
//	func (x MyType) MarshalJSON() ([]byte, error) {
//	  type tStruct MyType
//	  return x.Feature.EncodeGeoJSON(x.ID, tStruct(x))
//	}
func (fea Feature) EncodeGeoJSON(props any) ([]byte, error) {
	return fencoder(&fea, props)
}

// Encodes object as GeoJSON
func Marshal[T IFeature](obj T) ([]byte, error) {
	return fencoder(obj, obj)
}

func fencoder(fea IFeature, obj any) ([]byte, error) {
	geometry := fea.unapply()
	properties, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}

	// Note: skip bounding box for the point.
	var bbox BoundingBox
	switch geometry.(type) {
	case nil:
		bbox = nil
	case *Point:
		bbox = nil
	default:
		bbox = geometry.BoundingBox()
	}

	// ID is optional if the feature has a unique identifier.
	var id string
	if identity, has := any(obj).(interface{ GeoJsonID() string }); has {
		id = identity.GeoJsonID()
	}

	val := struct {
		Type       string          `json:"type"`
		ID         string          `json:"id,omitempty"`
		BBox       BoundingBox     `json:"bbox,omitempty"`
		Geometry   Geometry        `json:"geometry"`
		Properties json.RawMessage `json:"properties,omitempty"`
	}{
		Type:       TYPE_FEATURE,
		ID:         id,
		BBox:       bbox,
		Geometry:   geometry,
		Properties: properties,
	}

	return json.Marshal(val)
}

//
// Decoder
//

// DecodeGeoJSON is a helper function to implement GeoJSON codec
//
//	func (x *MyType) UnmarshalJSON(b []byte) (err error) {
//	  type tStruct *MyType
//	  x.ID, err = x.Feature.DecodeGeoJSON(b, tStruct(x))
//	  return
//	}
func (fea *Feature) DecodeGeoJSON(bytes []byte, props any) error {
	return fdecode(bytes, fea, props)
}

// Decodes GeoJSON object
func Unmarshal[T IFeature](bytes []byte, obj T) error {
	return fdecode(bytes, obj, obj)
}

func fdecode(bytes []byte, fe IFeature, obj any) error {
	var fea struct {
		Type       string          `json:"type"`
		ID         string          `json:"id,omitempty"`
		Geometry   json.RawMessage `json:"geometry"`
		Properties json.RawMessage `json:"properties,omitempty"`
	}

	if err := json.Unmarshal(bytes, &fea); err != nil {
		return err
	}

	if fea.Type != TYPE_FEATURE {
		return ErrUnsupportedType
	}

	if fea.Properties != nil {
		if err := json.Unmarshal(fea.Properties, &obj); err != nil {
			return err
		}

		// ID is optional if the feature has a unique identifier.
		if identity, has := any(obj).(interface{ SetGeoJsonID(string) }); has {
			identity.SetGeoJsonID(fea.ID)
		}
	}

	if fea.Geometry != nil {
		geo, err := decodeGeometry(fea.Geometry)
		if err != nil {
			return err
		}
		fe.apply(geo)
	}

	return nil
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
