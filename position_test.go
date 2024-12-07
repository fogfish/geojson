//
// Copyright (C) 2021 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/fogfish/geojson
//

package geojson_test

import (
	"testing"

	"github.com/fogfish/geojson"
	"github.com/fogfish/it/v2"
)

func TestPosition(t *testing.T) {
	p := geojson.Coord{100.0, 0.0}
	lat, lng := p.LatLng()

	seq := []float64{}
	p.FMap(func(x geojson.Coord) {
		lat, lng := x.LatLng()
		seq = append(seq, lng, lat)
	})

	it.Then(t).Should(
		it.Equal(lat, 0.0),
		it.Equal(lng, 100.0),
		it.Equal(p.Lat(), lat),
		it.Equal(p.Lng(), lng),
		it.Seq(seq).Equal(100.0, 0.0),
	)
}

func TestSequence(t *testing.T) {
	p := geojson.Curve{
		{100.0, 0.0},
		{100.0, 0.0},
	}

	seq := []float64{}
	p.FMap(func(x geojson.Coord) {
		lat, lng := x.LatLng()
		seq = append(seq, lng, lat)
	})

	it.Then(t).Should(
		it.Seq(seq).Equal(100.0, 0.0, 100.0, 0.0),
	)
}

func TestSurface(t *testing.T) {
	p := geojson.Surface{
		{
			{100.0, 0.0},
			{101.0, 0.0},
			{101.0, 1.0},
			{100.0, 1.0},
			{100.0, 0.0},
		},
	}

	seq := []float64{}
	p.FMap(func(x geojson.Coord) {
		lat, lng := x.LatLng()
		seq = append(seq, lng, lat)
	})

	it.Then(t).Should(
		it.Seq(seq).Equal(100.0, 0.0, 101.0, 0.0, 101.0, 1.0, 100.0, 1.0, 100.0, 0.0),
	)
}

func TestSurfaces(t *testing.T) {
	p := geojson.Surfaces{
		{
			{
				{100.0, 0.0},
				{101.0, 0.0},
				{101.0, 1.0},
				{100.0, 1.0},
				{100.0, 0.0},
			},
		},
	}

	seq := []float64{}
	p.FMap(func(x geojson.Coord) {
		lat, lng := x.LatLng()
		seq = append(seq, lng, lat)
	})

	it.Then(t).Should(
		it.Seq(seq).Equal(100.0, 0.0, 101.0, 0.0, 101.0, 1.0, 100.0, 1.0, 100.0, 0.0),
	)
}

func TestBBox(t *testing.T) {
	bbox := geojson.BoundingBox{-10.0, -20.0, +10.0, +20.0}

	it.Then(t).Should(
		it.Equal(bbox.SouthWest().Lng(), -10.0),
		it.Equal(bbox.SouthWest().Lat(), -20.0),
		it.Equal(bbox.NorthEast().Lng(), +10.0),
		it.Equal(bbox.NorthEast().Lat(), +20.0),
	)
}
