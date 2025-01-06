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

const TYPE_FEATURE_COLLECTION = "FeatureCollection"

// The Collection object represents a collection of spatially bounded elements.
// It contains a sequence of features, as defined by the GeoJSON FeatureCollection standard.
// This construct is designed to support "foreign members" for improved exchange
// of geospatial data. The value of a "foreign member" is determined by the application.
//
// This library reuses the "properties" attribute, which acts as a foreign member
// within the context of collections.
//
// To provide type-safe handling of collection properties, this library avoids using
// a generic interface{} type. Instead, it employs a type-tagging (or embedding) technique.
//
// Example:
//
//	type MyCollection struct {
//	  geojson.Collection[MyType]
//	  Name string `json:"name,omitempty"`
//	}
type Collection[T interface{ BoundingBox() BoundingBox }] struct {
	Features []T `json:"-"`
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

// EncodeGeoJSON is a helper function to implement GeoJSON codec
//
//	func (x MyCollection) MarshalJSON() ([]byte, error) {
//		type tStruct MyCollection
//		return x.Features.EncodeGeoJSON(tStruct(x))
//	}
func (c Collection[T]) EncodeGeoJSON(props any) ([]byte, error) {
	properties, err := json.Marshal(props)
	if err != nil {
		return nil, err
	}

	val := struct {
		Type       string          `json:"type"`
		BBox       BoundingBox     `json:"bbox,omitempty"`
		Features   []T             `json:"features,omitempty"`
		Properties json.RawMessage `json:"properties,omitempty"`
	}{
		Type:       TYPE_FEATURE_COLLECTION,
		BBox:       c.BoundingBox(),
		Features:   c.Features,
		Properties: properties,
	}

	return json.Marshal(val)
}

// Decodes FeatureCollection GeoJSON
// func (c *Collection[T]) UnmarshalJSON(b []byte) error {
// 	type anyCollection struct {
// 		Type     string          `json:"type"`
// 		Features json.RawMessage `json:"features,omitempty"`
// 	}

// 	any := anyCollection{}

// 	if err := json.Unmarshal(b, &any); err != nil {
// 		return err
// 	}

// 	if any.Type != TYPE_FEATURE_COLLECTION {
// 		return ErrUnsupportedType
// 	}

// 	if any.Features != nil {
// 		if err := json.Unmarshal(any.Features, &c.Features); err != nil {
// 			return err
// 		}
// 	}

// 	return nil
// }

// DecodeGeoJSON is a helper function to implement GeoJSON codec
//
//	func (x *MyCollection) UnmarshalJSON(b []byte) error {
//		type tStruct *MyCollection
//		return x.Features.DecodeGeoJSON(b, tStruct(x))
//	}
func (c *Collection[T]) DecodeGeoJSON(bytes []byte, props interface{}) error {
	val := struct {
		Type       string          `json:"type"`
		BBox       BoundingBox     `json:"bbox,omitempty"`
		Features   json.RawMessage `json:"features,omitempty"`
		Properties json.RawMessage `json:"properties,omitempty"`
	}{}

	if err := json.Unmarshal(bytes, &val); err != nil {
		return err
	}

	if val.Type != TYPE_FEATURE_COLLECTION {
		return ErrUnsupportedType
	}

	if val.Features != nil {
		if err := json.Unmarshal(val.Features, &c.Features); err != nil {
			return err
		}
	}

	if val.Properties != nil {
		if err := json.Unmarshal(val.Properties, &props); err != nil {
			return err
		}
	}

	return nil
}
