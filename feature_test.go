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

	"github.com/fogfish/curie/v2"
	"github.com/fogfish/geojson"
	"github.com/fogfish/it/v2"
)

const (
	city_helsinki = curie.IRI("city:helsinki")

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

	featurePointEmpty = `
	{
		"type": "Feature",
		"geometry": null,
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

	it.Then(t).Should(
		it.Nil(err),
		it.Equal(city.Name, "Helsinki"),
		it.Like(city.Geometry, &geojson.Point{geojson.Coord{102.0, 0.5}}),
	)
}

func TestFeatureDecodeEmpty(t *testing.T) {
	var city GeoJsonCity
	err := json.Unmarshal([]byte(featurePointEmpty), &city)

	it.Then(t).Should(
		it.Nil(err),
		it.Equal(city.Name, "Helsinki"),
		it.Nil(city.Geometry),
	)
}

func TestFeatureEncodePoint(t *testing.T) {
	city := GeoJsonCity{
		Feature: geojson.NewPoint(city_helsinki, geojson.Coord{100.0, 0.0}),
		City:    City{Name: "Helsinki"},
	}

	data, err := json.Marshal(city)
	it.Then(t).Should(it.Nil(err))

	var c GeoJsonCity
	err = json.Unmarshal([]byte(data), &c)

	it.Then(t).Should(
		it.Nil(err),
		it.Equal(c.ID, city_helsinki),
		it.Equal(c.Name, city.Name),
		it.Like(c.Geometry, &geojson.Point{geojson.Coord{100.0, 0.0}}),
	)
}

func TestFeatureEncodePointEmpty(t *testing.T) {
	city := GeoJsonCity{
		Feature: geojson.NewPoint(city_helsinki, geojson.Coord{}),
		City:    City{Name: "Helsinki"},
	}

	data, err := json.Marshal(city)
	it.Then(t).Should(it.Nil(err))

	var c GeoJsonCity
	err = json.Unmarshal([]byte(data), &c)

	it.Then(t).Should(
		it.Nil(err),
		it.Equal(c.ID, city_helsinki),
		it.Equal(c.Name, city.Name),
		it.Like(c.Geometry, &geojson.Point{geojson.Coord{}}),
	)
}

func TestFeatureEncodeUndefined(t *testing.T) {
	city := GeoJsonCity{
		City: City{Name: "Helsinki"},
	}

	data, err := json.Marshal(city)
	it.Then(t).Should(it.Nil(err))

	var c GeoJsonCity
	err = json.Unmarshal(data, &c)

	it.Then(t).Should(
		it.Nil(err),
		it.Equal(c.ID, ""),
		it.Equal(c.Name, city.Name),
		it.Nil(c.Geometry),
	)
}

func TestFeatureEncodeMultiPoint(t *testing.T) {
	city := GeoJsonCity{
		Feature: geojson.NewMultiPoint(city_helsinki,
			geojson.Curve{
				{100.0, 0.0},
				{101.0, 1.0},
			},
		),
		City: City{Name: "Helsinki"},
	}

	data, err := json.Marshal(city)
	it.Then(t).Should(it.Nil(err))

	var c GeoJsonCity
	err = json.Unmarshal([]byte(data), &c)

	it.Then(t).Should(
		it.Nil(err),

		it.Equal(c.ID, city_helsinki),
		it.Equal(c.Name, city.Name),
		it.TypeOf[*geojson.MultiPoint](c.Geometry),
	)
}

func TestFeatureEncodeLineString(t *testing.T) {
	city := GeoJsonCity{
		Feature: geojson.NewLineString(city_helsinki,
			geojson.Curve{
				{100.0, 0.0},
				{101.0, 1.0},
			},
		),
		City: City{Name: "Helsinki"},
	}

	data, err := json.Marshal(city)
	it.Then(t).Should(it.Nil(err))

	var c GeoJsonCity
	err = json.Unmarshal([]byte(data), &c)

	it.Then(t).Should(
		it.Nil(err),
		it.Equal(c.ID, city_helsinki),
		it.Equal(c.Name, city.Name),
		it.TypeOf[*geojson.LineString](c.Geometry),
	)
}

func TestFeatureEncodeMultiLineString(t *testing.T) {
	city := GeoJsonCity{
		Feature: geojson.NewMultiLineString(city_helsinki,
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
	it.Then(t).Should(it.Nil(err))

	var c GeoJsonCity
	err = json.Unmarshal([]byte(data), &c)

	it.Then(t).Should(
		it.Nil(err),
		it.Equal(c.ID, city_helsinki),
		it.Equal(c.Name, city.Name),
		it.TypeOf[*geojson.MultiLineString](c.Geometry),
	)
}

func TestFeatureEncodePolygon(t *testing.T) {
	city := GeoJsonCity{
		Feature: geojson.NewPolygon(city_helsinki,
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
	it.Then(t).Should(it.Nil(err))

	var c GeoJsonCity
	err = json.Unmarshal([]byte(data), &c)

	it.Then(t).Should(
		it.Nil(err),
		it.Equal(c.ID, city_helsinki),
		it.Equal(c.Name, city.Name),
		it.TypeOf[*geojson.Polygon](c.Geometry),
	)
}

func TestFeatureEncodeMultiPolygon(t *testing.T) {
	city := GeoJsonCity{
		Feature: geojson.NewMultiPolygon(city_helsinki,
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
	it.Then(t).Should(it.Nil(err))

	var c GeoJsonCity
	err = json.Unmarshal([]byte(data), &c)

	it.Then(t).Should(
		it.Nil(err),
		it.Equal(c.ID, city_helsinki),
		it.Equal(c.Name, city.Name),
		it.TypeOf[*geojson.MultiPolygon](c.Geometry),
	)
}

func TestFeatureInvalidDecode(t *testing.T) {
	var city GeoJsonCity
	err := json.Unmarshal([]byte(featureInvalid), &city)

	it.Then(t).ShouldNot(
		it.Nil(err),
	).Should(
		it.Equal(city.Name, ""),
		it.Nil(city.Geometry),
	)
}
