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

func TestGeometryPoint(t *testing.T) {
	var geometry geojson.Geometry
	err := json.Unmarshal([]byte(geometryPoint), &geometry)

	it.Ok(t).
		IfNil(err).
		If(geometry.Coords).Should().Be().Like(&geojson.Point{})

	switch v := geometry.Coords.(type) {
	case *geojson.Point:
		it.Ok(t).If(v.Coords).Equal(
			geojson.Position{100.0, 0.0},
		)
	default:
		t.Error("Invalid Coords Type")
	}
}

func TestGeometryMultiPoint(t *testing.T) {
	var geometry geojson.Geometry
	err := json.Unmarshal([]byte(geometryMultiPoint), &geometry)

	it.Ok(t).
		IfNil(err).
		If(geometry.Coords).Should().Be().Like(&geojson.MultiPoint{})

	switch v := geometry.Coords.(type) {
	case *geojson.MultiPoint:
		it.Ok(t).If(v.Coords).Equal(
			geojson.Sequence{
				{100.0, 0.0},
				{101.0, 1.0},
			},
		)
	default:
		t.Error("Invalid Coords Type")
	}
}

func TestGeometryLineString(t *testing.T) {
	var geometry geojson.Geometry
	err := json.Unmarshal([]byte(geometryLineString), &geometry)

	it.Ok(t).
		IfNil(err).
		If(geometry.Coords).Should().Be().Like(&geojson.LineString{})

	switch v := geometry.Coords.(type) {
	case *geojson.LineString:
		it.Ok(t).If(v.Coords).Equal(
			geojson.Sequence{
				{100.0, 0.0},
				{101.0, 1.0},
			},
		)
	default:
		t.Error("Invalid Coords Type")
	}
}

func TestGeometryMultiLineString(t *testing.T) {
	var geometry geojson.Geometry
	err := json.Unmarshal([]byte(geometryMultiLineString), &geometry)

	it.Ok(t).
		IfNil(err).
		If(geometry.Coords).Should().Be().Like(&geojson.MultiLineString{})

	switch v := geometry.Coords.(type) {
	case *geojson.MultiLineString:
		it.Ok(t).If(v.Coords).Equal(
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
	default:
		t.Error("Invalid Coords Type")
	}
}

func TestGeometryPolygon(t *testing.T) {
	var geometry geojson.Geometry
	err := json.Unmarshal([]byte(geometryPolygon), &geometry)

	it.Ok(t).
		IfNil(err).
		If(geometry.Coords).Should().Be().Like(&geojson.Polygon{})

	switch v := geometry.Coords.(type) {
	case *geojson.Polygon:
		it.Ok(t).If(v.Coords).Equal(
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
	default:
		t.Error("Invalid Coords Type")
	}
}

func TestGeometryPolygonWithHole(t *testing.T) {
	var geometry geojson.Geometry
	err := json.Unmarshal([]byte(geometryPolygonWithHole), &geometry)

	it.Ok(t).
		IfNil(err).
		If(geometry.Coords).Should().Be().Like(&geojson.Polygon{})

	switch v := geometry.Coords.(type) {
	case *geojson.Polygon:
		it.Ok(t).If(v.Coords).Equal(
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
	default:
		t.Error("Invalid Coords Type")
	}
}

func TestGeometryMultiPolygon(t *testing.T) {
	var geometry geojson.Geometry
	err := json.Unmarshal([]byte(geometryMultiPolygon), &geometry)

	it.Ok(t).
		IfNil(err).
		If(geometry.Coords).Should().Be().Like(&geojson.MultiPolygon{})

	switch v := geometry.Coords.(type) {
	case *geojson.MultiPolygon:
		it.Ok(t).If(v.Coords).Equal(
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
	default:
		t.Error("Invalid Coords Type")
	}
}
