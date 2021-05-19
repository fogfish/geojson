//
// Copyright (C) 2021 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/fogfish/geojson
//

package geojson

/*

Position is the fundamental geometry construct.
A position is an array of numbers that defines projected
coordinates. The first two elements are Euclidian coordinates
at plane in x, y order; easting, northing for projected
coordinates, longitude, and latitude for geographic coordinates.

One Position in the case of a Point geometry (0-dimensional point)
*/
type Position []float64

/*

Sequence of positions in the case of a LineString
or MultiPoint geometry (1-dimensional curve)
*/
type Sequence []Position

/*

Surface is an array of LineString or linear ring coordinates
in the case of a Polygon or MultiLineString geometry
(2-dimensional surface)
*/
type Surface []Sequence
