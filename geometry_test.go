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
	"github.com/fogfish/it"
)

const (
	geometryUnknown = `
	{
			"type": "Unknown",
			"coordinates": [100.0, 0.0]
	}
	`

	geometryPoint = `
	{
			"type": "Point",
			"coordinates": [100.0, 0.0]
	}
	`

	geometryMultiPoint = `
	{
			"type": "MultiPoint",
			"coordinates": [
					[100.0, 0.0],
					[101.0, 1.0]
			]
	}
	`

	geometryLineString = `
	{
			"type": "LineString",
			"coordinates": [
					[100.0, 0.0],
					[101.0, 1.0]
			]
	}
	`

	geometryMultiLineString = `
	{
			"type": "MultiLineString",
			"coordinates": [
					[
							[100.0, 0.0],
							[101.0, 1.0]
					],
					[
							[102.0, 2.0],
							[103.0, 3.0]
					]
			]
	}
	`

	geometryPolygon = `
	{
			"type": "Polygon",
			"coordinates": [
					[
							[100.0, 0.0],
							[101.0, 0.0],
							[101.0, 1.0],
							[100.0, 1.0],
							[100.0, 0.0]
					]
			]
	}
	`

	geometryPolygonWithHole = `
	{
			"type": "Polygon",
			"coordinates": [
					[
							[100.0, 0.0],
							[101.0, 0.0],
							[101.0, 1.0],
							[100.0, 1.0],
							[100.0, 0.0]
					],
					[
							[100.8, 0.8],
							[100.8, 0.2],
							[100.2, 0.2],
							[100.2, 0.8],
							[100.8, 0.8]
					]
			]
	}
	`

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

func TestGeometryUnknown(t *testing.T) {
	var geo *geojson.Point
	err := json.Unmarshal([]byte(geometryUnknown), geo)

	it.Ok(t).
		IfNotNil(err).
		IfNil(geo)
}

func TestGeometryPoint(t *testing.T) {
	var geo *geojson.Point
	err := json.Unmarshal([]byte(geometryPoint), &geo)

	it.Ok(t).
		IfNil(err).
		If(geo).Should().Be().Like(&geojson.Point{})

	it.Ok(t).If(geo.Coords).Equal(
		geojson.Coord{100.0, 0.0},
	)

	bbox := geo.BoundingBox()
	it.Ok(t).If(bbox).Equal(
		geojson.BoundingBox{100.0, 0, 100.0, 0},
	)
}

func TestGeometryMultiPoint(t *testing.T) {
	var geo *geojson.MultiPoint
	err := json.Unmarshal([]byte(geometryMultiPoint), &geo)

	it.Ok(t).
		IfNil(err).
		If(geo).Should().Be().Like(&geojson.MultiPoint{})

	it.Ok(t).If(geo.Coords).Equal(
		geojson.Curve{
			{100.0, 0.0},
			{101.0, 1.0},
		},
	)

	bbox := geo.BoundingBox()
	it.Ok(t).If(bbox).Equal(
		geojson.BoundingBox{100.0, 0, 101.0, 1.0},
	)
}

func TestGeometryLineString(t *testing.T) {
	var geo *geojson.LineString
	err := json.Unmarshal([]byte(geometryLineString), &geo)

	it.Ok(t).
		IfNil(err).
		If(geo).Should().Be().Like(&geojson.LineString{})

	it.Ok(t).If(geo.Coords).Equal(
		geojson.Curve{
			{100.0, 0.0},
			{101.0, 1.0},
		},
	)

	bbox := geo.BoundingBox()
	it.Ok(t).If(bbox).Equal(
		geojson.BoundingBox{100.0, 0, 101.0, 1.0},
	)
}

func TestGeometryMultiLineString(t *testing.T) {
	var geo *geojson.MultiLineString
	err := json.Unmarshal([]byte(geometryMultiLineString), &geo)

	it.Ok(t).
		IfNil(err).
		If(geo).Should().Be().Like(&geojson.MultiLineString{})

	it.Ok(t).If(geo.Coords).Equal(
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
	)

	bbox := geo.BoundingBox()
	it.Ok(t).If(bbox).Equal(
		geojson.BoundingBox{100.0, 0, 103.0, 3.0},
	)
}

func TestGeometryPolygon(t *testing.T) {
	var geo *geojson.Polygon
	err := json.Unmarshal([]byte(geometryPolygon), &geo)

	it.Ok(t).
		IfNil(err).
		If(geo).Should().Be().Like(&geojson.Polygon{})

	it.Ok(t).If(geo.Coords).Equal(
		geojson.Surface{
			{
				{100.0, 0.0},
				{101.0, 0.0},
				{101.0, 1.0},
				{100.0, 1.0},
				{100.0, 0.0},
			},
		},
	)

	bbox := geo.BoundingBox()
	it.Ok(t).If(bbox).Equal(
		geojson.BoundingBox{100.0, 0, 101.0, 1.0},
	)
}

func TestGeometryPolygonWithHole(t *testing.T) {
	var geo *geojson.Polygon
	err := json.Unmarshal([]byte(geometryPolygonWithHole), &geo)

	it.Ok(t).
		IfNil(err).
		If(geo).Should().Be().Like(&geojson.Polygon{})

	it.Ok(t).If(geo.Coords).Equal(
		geojson.Surface{
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
		},
	)

	bbox := geo.BoundingBox()
	it.Ok(t).If(bbox).Equal(
		geojson.BoundingBox{100.0, 0, 101.0, 1.0},
	)
}

func TestGeometryMultiPolygon(t *testing.T) {
	var geo *geojson.MultiPolygon
	err := json.Unmarshal([]byte(geometryMultiPolygon), &geo)

	it.Ok(t).
		IfNil(err).
		If(geo).Should().Be().Like(&geojson.MultiPolygon{})

	it.Ok(t).If(geo.Coords).Equal(
		[]geojson.Surface{
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
		},
	)

	bbox := geo.BoundingBox()
	it.Ok(t).If(bbox).Equal(
		geojson.BoundingBox{100.0, 0, 103.0, 3.0},
	)
}

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
