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

	"github.com/fogfish/geojson"
	"github.com/fogfish/it/v2"
)

// fixtures for unit testing
var (
	coordPoint = geojson.Coord{100.0, 0.0}

	coordMultiPoint = geojson.Curve{
		{100.0, 0.0},
		{101.0, 1.0},
	}

	coordLineString = geojson.Curve{
		{100.0, 0.0},
		{101.0, 1.0},
	}

	coordMultiLineString = geojson.Surface{
		{
			{100.0, 0.0},
			{101.0, 1.0},
		},
		{
			{102.0, 2.0},
			{103.0, 3.0},
		},
	}

	coordPolygon = geojson.Surface{
		{
			{100.0, 0.0},
			{101.0, 0.0},
			{101.0, 1.0},
			{100.0, 1.0},
			{100.0, 0.0},
		},
	}

	coordPolygonWithHole = geojson.Surface{
		{
			{100.0, 0.0},
			{101.0, 0.0},
			{101.0, 1.0},
			{100.0, 1.0},
			{100.0, 0.0},
		},
		{
			{100.8, 0.8},
			{100.8, 0.2},
			{100.2, 0.2},
			{100.2, 0.8},
			{100.8, 0.8},
		},
	}

	coordMultiPolygon = geojson.Surfaces{
		{
			{
				{102.0, 2.0},
				{103.0, 2.0},
				{103.0, 3.0},
				{102.0, 3.0},
				{102.0, 2.0},
			},
		},
		{
			{
				{100.0, 0.0},
				{101.0, 0.0},
				{101.0, 1.0},
				{100.0, 1.0},
				{100.0, 0.0},
			},
			{
				{100.2, 0.2},
				{100.2, 0.8},
				{100.8, 0.8},
				{100.8, 0.2},
				{100.2, 0.2},
			},
		},
	}
)

func genGeoJSON(t string, shape geojson.Shape) []byte {
	b, _ := json.Marshal(map[string]any{
		"type":        t,
		"coordinates": shape,
	})
	return b
}

const (
	geometryCorrupted = `
	{
			type": "Unknown",
			"coordinates: [100.0, 0.0
	`

	// geometryUnknown = `
	// {
	// 		"type": "Unknown",
	// 		"coordinates": [100.0, 0.0]
	// }
	// `

	// geometryPoint = `
	// {
	// 		"type": "Point",
	// 		"coordinates": [100.0, 0.0]
	// }
	// `

	// geometryMultiPoint = `
	// {
	// 		"type": "MultiPoint",
	// 		"coordinates": [
	// 				[100.0, 0.0],
	// 				[101.0, 1.0]
	// 		]
	// }
	// `

	// geometryMultiPointX = `
	// {
	// 		"type": "MultiPointX",
	// 		"coordinates": [
	// 				[100.0, 0.0],
	// 				[101.0, 1.0]
	// 		]
	// }
	// `

	// geometryLineString = `
	// {
	// 		"type": "LineString",
	// 		"coordinates": [
	// 				[100.0, 0.0],
	// 				[101.0, 1.0]
	// 		]
	// }
	// `

	// geometryMultiLineString = `
	// {
	// 		"type": "MultiLineString",
	// 		"coordinates": [
	// 				[
	// 						[100.0, 0.0],
	// 						[101.0, 1.0]
	// 				],
	// 				[
	// 						[102.0, 2.0],
	// 						[103.0, 3.0]
	// 				]
	// 		]
	// }
	// `

	// geometryPolygon = `
	// {
	// 		"type": "Polygon",
	// 		"coordinates": [
	// 				[
	// 						[100.0, 0.0],
	// 						[101.0, 0.0],
	// 						[101.0, 1.0],
	// 						[100.0, 1.0],
	// 						[100.0, 0.0]
	// 				]
	// 		]
	// }
	// `

	// geometryPolygonWithHole = `
	// {
	// 		"type": "Polygon",
	// 		"coordinates": [
	// 				[
	// 						[100.0, 0.0],
	// 						[101.0, 0.0],
	// 						[101.0, 1.0],
	// 						[100.0, 1.0],
	// 						[100.0, 0.0]
	// 				],
	// 				[
	// 						[100.8, 0.8],
	// 						[100.8, 0.2],
	// 						[100.2, 0.2],
	// 						[100.2, 0.8],
	// 						[100.8, 0.8]
	// 				]
	// 		]
	// }
	// `

	geometryMultiPolygon = `
	{
			"type": "MultiPolygon",
			"coordinates": [
					[
							[
									[102.0, 2.0],
									[103.0, 2.0],
									[103.0, 3.0],
									[102.0, 3.0],
									[102.0, 2.0]
							]
					],
					[
							[
									[100.0, 0.0],
									[101.0, 0.0],
									[101.0, 1.0],
									[100.0, 1.0],
									[100.0, 0.0]
							],
							[
									[100.2, 0.2],
									[100.2, 0.8],
									[100.8, 0.8],
									[100.8, 0.2],
									[100.2, 0.2]
							]
					]
			]
	}
	`
)

type GeoJSON struct {
	geojson.Feature
}

func (x GeoJSON) MarshalJSON() ([]byte, error) {
	type tStruct GeoJSON
	return x.Feature.EncodeGeoJSON(tStruct(x))
}

func (x *GeoJSON) UnmarshalJSON(b []byte) error {
	type tStruct *GeoJSON
	return x.Feature.DecodeGeoJSON(b, tStruct(x))
}

func testGeometry[T geojson.Geometry](
	t *testing.T,
	typeOf string,
	coord geojson.Shape,
	bbox geojson.BoundingBox,
) {
	t.Helper()

	var geo T

	t.Run("Success", func(t *testing.T) {
		err := json.Unmarshal(genGeoJSON(typeOf, coord), &geo)

		it.Then(t).Should(
			it.Nil(err),
			it.Equiv(geo.BoundingBox(), bbox),
			it.Equiv(geo.Geometry(), coord),
		)
	})

	t.Run("Not Supported", func(t *testing.T) {
		it.Then(t).Should(
			it.Fail(
				func() error {
					return json.Unmarshal(genGeoJSON("Unknown", coord), geo)
				},
			).Contain("type Unknown is not supported"),
		)
	})

	t.Run("Corrupted", func(t *testing.T) {
		it.Then(t).Should(
			it.Fail(
				func() error {
					return json.Unmarshal([]byte(geometryCorrupted), geo)
				},
			).Contain("invalid character"),
		)
	})
}

func TestGeometryPoint(t *testing.T) {
	testGeometry[*geojson.Point](t, "Point", coordPoint,
		geojson.BoundingBox{100.0, 0, 100.0, 0},
	)
}

func TestGeometryMultiPoint(t *testing.T) {
	testGeometry[*geojson.MultiPoint](t, "MultiPoint", coordMultiPoint,
		geojson.BoundingBox{100.0, 0, 101.0, 1.0},
	)
}

func TestGeometryLineString(t *testing.T) {
	testGeometry[*geojson.LineString](t, "LineString", coordLineString,
		geojson.BoundingBox{100.0, 0, 101.0, 1.0},
	)
}

func TestGeometryMultiLineString(t *testing.T) {
	testGeometry[*geojson.MultiLineString](t, "MultiLineString", coordMultiLineString,
		geojson.BoundingBox{100.0, 0, 103.0, 3.0},
	)
}

func TestGeometryPolygon(t *testing.T) {
	testGeometry[*geojson.Polygon](t, "Polygon", coordPolygon,
		geojson.BoundingBox{100.0, 0, 101.0, 1.0},
	)
}

func TestGeometryPolygonWithHole(t *testing.T) {
	testGeometry[*geojson.Polygon](t, "Polygon", coordPolygonWithHole,
		geojson.BoundingBox{100.0, 0, 101.0, 1.0},
	)
}

func TestGeometryMultiPolygon(t *testing.T) {
	testGeometry[*geojson.MultiPolygon](t, "MultiPolygon", coordMultiPolygon,
		geojson.BoundingBox{100.0, 0, 103.0, 3.0},
	)
}

/*
func TestEmptyGeometry(t *testing.T) {
	it.Ok(t).
		IfTrue(geojson.NewPoint("", nil).BoundingBox() == nil).
		IfTrue(geojson.NewPoint("", geojson.Coord{}).BoundingBox() == nil).
		IfTrue(geojson.NewMultiPoint("", nil).BoundingBox() == nil).
		IfTrue(geojson.NewMultiPoint("", geojson.Curve{}).BoundingBox() == nil).
		IfTrue(geojson.NewLineString("", nil).BoundingBox() == nil).
		IfTrue(geojson.NewLineString("", geojson.Curve{}).BoundingBox() == nil).
		IfTrue(geojson.NewMultiLineString("", nil).BoundingBox() == nil).
		IfTrue(geojson.NewMultiLineString("", geojson.Surface{}).BoundingBox() == nil).
		IfTrue(geojson.NewPolygon("", nil).BoundingBox() == nil).
		IfTrue(geojson.NewPolygon("", geojson.Surface{}).BoundingBox() == nil).
		IfTrue(geojson.NewMultiPolygon("", nil).BoundingBox() == nil).
		IfTrue(geojson.NewMultiPolygon("", geojson.Surface{}).BoundingBox() == nil)
}
*/
