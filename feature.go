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

	"github.com/fogfish/curie"
)

/*

Feature object represents a spatially bounded thing.
This object contains geometry, a common identifier, and properties.
The value of the properties is any JSON object, typically defined by
an application.

The library uses a type safe notation for the feature's property
definition instead of generic interface{} type. It uses type tagging
technique (or embeding):

  type MyType struct {
    geojson.Feature
    Name      string `json:"name,omitempty"`
  }

*/
type Feature struct {
	ID       *curie.IRI `json:"-"`
	Geometry *Geometry  `json:"-"`
}

/*

WithID sets feature ID from string
*/
func (fea Feature) WithID(iri string) Feature {
	fea.ID = curie.New(iri).This()
	return fea
}

/*

WithIRI sets feature ID from CURIE (compact IRI type)
*/
func (fea Feature) WithIRI(iri curie.IRI) Feature {
	fea.ID = &iri
	return fea
}

/*

EncodeGeoJSON is a helper function to implement GeoJSON codec

  func (x MyType) MarshalJSON() ([]byte, error) {
	  type tStruct MyType
	  return x.Feature.EncodeGeoJSON(tStruct(x))
  }
*/
func (fea Feature) EncodeGeoJSON(props interface{}) ([]byte, error) {
	properties, err := json.Marshal(props)
	if err != nil {
		return nil, err
	}

	any := struct {
		Type       string          `json:"type"`
		ID         *curie.IRI      `json:"id,omitempty"`
		Geometry   *Geometry       `json:"geometry,omitempty"`
		Properties json.RawMessage `json:"properties,omitempty"`
	}{
		ID:         fea.ID,
		Type:       "Feature",
		Geometry:   fea.Geometry,
		Properties: properties,
	}

	return json.Marshal(any)
}

/*

DecodeGeoJSON is a helper function to implement GeoJSON codec

  func (x *MyType) UnmarshalJSON(b []byte) error {
	  type tStruct *MyType
	  return x.Feature.DecodeGeoJSON(b, tStruct(x))
  }
*/
func (fea *Feature) DecodeGeoJSON(bytes []byte, props interface{}) error {
	any := struct {
		Type       string          `json:"type"`
		ID         *curie.IRI      `json:"id,omitempty"`
		Geometry   json.RawMessage `json:"geometry,omitempty"`
		Properties json.RawMessage `json:"properties,omitempty"`
	}{}

	if err := json.Unmarshal(bytes, &any); err != nil {
		return err
	}

	if any.Type != "Feature" {
		return ErrorUnsupportedType
	}

	if any.Geometry != nil {
		geo := Geometry{}
		if err := json.Unmarshal(any.Geometry, &geo); err != nil {
			return err
		}
		fea.Geometry = &geo
	}

	if any.Properties != nil {
		if err := json.Unmarshal(any.Properties, &props); err != nil {
			return err
		}
	}

	fea.ID = any.ID
	return nil
}

/*

NewPoint ⟼ Feature[Point]
*/
func NewPoint(coords ...float64) Feature {
	return Feature{
		Geometry: &Geometry{
			Coords: &Point{Coords: coords},
		},
	}
}

/*

NewMultiPoint ⟼ Feature[MultiPoint]
*/
func NewMultiPoint(coords Sequence) Feature {
	return Feature{
		Geometry: &Geometry{
			Coords: &MultiPoint{Coords: coords},
		},
	}
}

/*

NewLineString ⟼ Feature[LineString]
*/
func NewLineString(coords Sequence) Feature {
	return Feature{
		Geometry: &Geometry{
			Coords: &LineString{Coords: coords},
		},
	}
}

/*

NewMultiLineString ⟼ Feature[MultiLineString]
*/
func NewMultiLineString(coords Surface) Feature {
	return Feature{
		Geometry: &Geometry{
			Coords: &MultiLineString{Coords: coords},
		},
	}
}

/*

NewPolygon ⟼ Feature[Polygon]
*/
func NewPolygon(coords Surface) Feature {
	return Feature{
		Geometry: &Geometry{
			Coords: &Polygon{Coords: coords},
		},
	}
}

/*

NewMultiPolygon ⟼ Feature[MultiPolygon]
*/
func NewMultiPolygon(coords ...Surface) Feature {
	return Feature{
		Geometry: &Geometry{
			Coords: &MultiPolygon{Coords: coords},
		},
	}
}
