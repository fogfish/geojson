//
// Copyright (C) 2021 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/fogfish/geojson
//

package geojson_test

import (
	"encoding/json"
	"testing"

	"github.com/fogfish/curie"
	"github.com/fogfish/geojson"
	"github.com/fogfish/it"
)

const (
	featureInvalid = `
		{
			"type": "Unknown",
			"geometry": {
				"type": "Point",
				"coordinates": [102.0, 0.5]
			},
			"properties": {
				"name": "Helsinki"
			}
		}
	`

	featurePoint = `
		{
			"type": "Feature",
			"geometry": {
				"type": "Point",
				"coordinates": [102.0, 0.5]
			},
			"properties": {
				"name": "Helsinki"
			}
		}
	`
)

type City struct {
	Name string `json:"name,omitempty"`
}

type GeoJsonCity struct {
	geojson.Feature
	City
}

func (x GeoJsonCity) MarshalJSON() ([]byte, error) {
	type tStruct GeoJsonCity
	return x.Feature.EncodeGeoJSON(tStruct(x))
}

func (x *GeoJsonCity) UnmarshalJSON(b []byte) error {
	type tStruct *GeoJsonCity
	return x.Feature.DecodeGeoJSON(b, tStruct(x))
}

func TestFeatureDecode(t *testing.T) {
	var city GeoJsonCity
	err := json.Unmarshal([]byte(featurePoint), &city)

	it.Ok(t).
		IfNil(err).
		If(city.Name).Equal("Helsinki").
		If(city.Geometry).Should().Be().Like(geojson.Point{})

	switch v := city.Geometry.(type) {
	case *geojson.Point:
		it.Ok(t).If(v.Coords).Equal(geojson.Coord{102.0, 0.5})
	default:
		t.Errorf("Invaid Coords Type")
	}
}

func TestFeatureEncodePoint(t *testing.T) {
	city := GeoJsonCity{
		Feature: geojson.NewPoint(geojson.Coord{100.0, 0.0}),
		City:    City{Name: "Helsinki"},
	}

	data, err := json.Marshal(city)
	it.Ok(t).IfNil(err)

	var c GeoJsonCity
	err = json.Unmarshal([]byte(data), &c)

	it.Ok(t).
		IfNil(err).
		IfNil(c.ID).
		If(c.Name).Equal(city.Name).
		If(c.Geometry).Should().Be().Like(geojson.Point{})
}

func TestFeatureEncodeMultiPoint(t *testing.T) {
	city := GeoJsonCity{
		Feature: geojson.NewMultiPoint(
			geojson.Curve{
				{100.0, 0.0},
				{101.0, 1.0},
			},
		),
		City: City{Name: "Helsinki"},
	}

	data, err := json.Marshal(city)
	it.Ok(t).IfNil(err)

	var c GeoJsonCity
	err = json.Unmarshal([]byte(data), &c)

	it.Ok(t).
		IfNil(err).
		IfNil(c.ID).
		If(c.Name).Equal(city.Name).
		If(c.Geometry).Should().Be().Like(geojson.MultiPoint{})
}

func TestFeatureEncodeLineString(t *testing.T) {
	city := GeoJsonCity{
		Feature: geojson.NewLineString(
			geojson.Curve{
				{100.0, 0.0},
				{101.0, 1.0},
			},
		),
		City: City{Name: "Helsinki"},
	}

	data, err := json.Marshal(city)
	it.Ok(t).IfNil(err)

	var c GeoJsonCity
	err = json.Unmarshal([]byte(data), &c)

	it.Ok(t).
		IfNil(err).
		IfNil(c.ID).
		If(c.Name).Equal(city.Name).
		If(c.Geometry).Should().Be().Like(geojson.LineString{})
}

func TestFeatureEncodeMultiLineString(t *testing.T) {
	city := GeoJsonCity{
		Feature: geojson.NewMultiLineString(
			geojson.Surface{
				{
					{100.0, 0.0},
					{101.0, 1.0},
				},
				{
					{102.0, 2.0},
					{103.0, 3.0},
				},
			},
		),
		City: City{Name: "Helsinki"},
	}

	data, err := json.Marshal(city)
	it.Ok(t).IfNil(err)

	var c GeoJsonCity
	err = json.Unmarshal([]byte(data), &c)

	it.Ok(t).
		IfNil(err).
		IfNil(c.ID).
		If(c.Name).Equal(city.Name).
		If(c.Geometry).Should().Be().Like(geojson.MultiLineString{})
}

func TestFeatureEncodePolygon(t *testing.T) {
	city := GeoJsonCity{
		Feature: geojson.NewPolygon(
			geojson.Surface{
				{
					{100.0, 0.0},
					{101.0, 0.0},
					{101.0, 1.0},
					{100.0, 1.0},
					{100.0, 0.0},
				},
			},
		),
		City: City{Name: "Helsinki"},
	}

	data, err := json.Marshal(city)
	it.Ok(t).IfNil(err)

	var c GeoJsonCity
	err = json.Unmarshal([]byte(data), &c)

	it.Ok(t).
		IfNil(err).
		IfNil(c.ID).
		If(c.Name).Equal(city.Name).
		If(c.Geometry).Should().Be().Like(geojson.Polygon{})
}

func TestFeatureEncodeMultiPolygon(t *testing.T) {
	city := GeoJsonCity{
		Feature: geojson.NewMultiPolygon(
			geojson.Surface{
				{
					{102.0, 2.0},
					{103.0, 2.0},
					{103.0, 3.0},
					{102.0, 3.0},
					{102.0, 2.0},
				},
			},
		),
		City: City{Name: "Helsinki"},
	}

	data, err := json.Marshal(city)
	it.Ok(t).IfNil(err)

	var c GeoJsonCity
	err = json.Unmarshal([]byte(data), &c)

	it.Ok(t).
		IfNil(err).
		IfNil(c.ID).
		If(c.Name).Equal(city.Name).
		If(c.Geometry).Should().Be().Like(geojson.MultiPolygon{})
}

func TestFeatureWithID(t *testing.T) {
	city := GeoJsonCity{
		Feature: geojson.NewPoint(geojson.Coord{100.0, 0.0}).WithID("city:helsinki"),
		City:    City{Name: "Helsinki"},
	}

	data, err := json.Marshal(city)
	it.Ok(t).IfNil(err)

	var c GeoJsonCity
	err = json.Unmarshal([]byte(data), &c)

	it.Ok(t).
		IfNil(err).
		IfNotNil(c.ID).
		If(*c.ID).Equal(curie.New("city:helsinki"))
}

func TestFeatureWithIRI(t *testing.T) {
	iri := curie.New("city:helsinki")
	city := GeoJsonCity{
		Feature: geojson.NewPoint(geojson.Coord{100.0, 0.0}).WithIRI(iri),
		City:    City{Name: "Helsinki"},
	}

	data, err := json.Marshal(city)
	it.Ok(t).IfNil(err)

	var c GeoJsonCity
	err = json.Unmarshal([]byte(data), &c)

	it.Ok(t).
		IfNil(err).
		IfNotNil(c.ID).
		If(*c.ID).Equal(iri)
}

func TestFeatureInvalidDecode(t *testing.T) {
	var city GeoJsonCity
	err := json.Unmarshal([]byte(featureInvalid), &city)

	it.Ok(t).
		IfNotNil(err).
		If(city.Name).Equal("").
		IfNil(city.Geometry)
}
