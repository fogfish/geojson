//
// Copyright (C) 2021 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/fogfish/geojson
//

package geojson

// Coord is the fundamental geometry construct.
// A position is an array of numbers that defines projected
// coordinates. The first two elements are Euclidian coordinates
// at plane in x, y order; easting, northing for projected
// coordinates, longitude, and latitude for geographic coordinates.
//
// One Position in the case of a Point geometry (0-dimensional point)
type Coord []float64

// LatLng coordinates of the position
func (coords Coord) LatLng() (float64, float64) { return coords[1], coords[0] }
func (coords Coord) Lat() float64               { return coords[1] }
func (coords Coord) Lng() float64               { return coords[0] }

// FMap applies a function to each coords pair
func (coords Coord) FMap(f func(Coord)) { f(coords) }

// Sequence of positions in the case of a LineString
// or MultiPoint geometry (1-dimensional curve)
type Curve []Coord

// FMap applies a function to each coords pair
func (seq Curve) FMap(f func(Coord)) {
	for _, x := range seq {
		f(x)
	}
}

// Surface is an array of LineString or linear ring coordinates
// in the case of a Polygon or MultiLineString geometry
// (2-dimensional surface)
type Surface []Curve

// FMap applies a function to each coords pair
func (seq Surface) FMap(f func(Coord)) {
	for _, x := range seq {
		x.FMap(f)
	}
}
