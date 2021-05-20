//
// Copyright (C) 2021 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/fogfish/geojson
//

package geojson

// Index of the coordinate fraction
const (
	X = 0
	Y = 1
)

/*

Position is the fundamental geometry construct.
A position is an array of numbers that defines projected
coordinates. The first two elements are Euclidian coordinates
at plane in x, y order; easting, northing for projected
coordinates, longitude, and latitude for geographic coordinates.

One Position in the case of a Point geometry (0-dimensional point)
*/
type Position []float64

// LatLng coordinates of the position
func (coords Position) LatLng() (float64, float64) {
	return coords[1], coords[0]
}

// FMap applies a function to each coords pair
func (coords Position) FMap(f func(Position)) {
	f(coords)
}

/*

Sequence of positions in the case of a LineString
or MultiPoint geometry (1-dimensional curve)
*/
type Sequence []Position

// FMap applies a function to each coords pair
func (seq Sequence) FMap(f func(Position)) {
	for _, x := range seq {
		f(x)
	}
}

/*

Surface is an array of LineString or linear ring coordinates
in the case of a Polygon or MultiLineString geometry
(2-dimensional surface)
*/
type Surface []Sequence

// FMap applies a function to each coords pair
func (seq Surface) FMap(f func(Position)) {
	for _, x := range seq {
		x.FMap(f)
	}
}
