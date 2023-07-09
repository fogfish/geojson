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
	"github.com/fogfish/it"
)

func TestPosition(t *testing.T) {
	p := geojson.Coord{100.0, 0.0}
	lat, lng := p.LatLng()

	seq := []float64{}
	p.FMap(func(x geojson.Coord) {
		lat, lng := x.LatLng()
		seq = append(seq, lng, lat)
	})

	it.Ok(t).
		If(lat).Equal(0.0).
		If(lng).Equal(100.0).
		If(p.Lat()).Equal(lat).
		If(p.Lng()).Equal(lng).
		If(seq).Equal([]float64{100.0, 0.0})
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

	it.Ok(t).
		If(seq).Equal([]float64{100.0, 0.0, 100.0, 0.0})
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

	it.Ok(t).
		If(seq).Equal([]float64{100.0, 0.0, 101.0, 0.0, 101.0, 1.0, 100.0, 1.0, 100.0, 0.0})
}
