//
// Copyright (C) 2021 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/fogfish/geojson
//

package geojson

// All position types implements shape interface,
// allowing map function over coordinates.
type Shape interface{ FMap(func(Coord)) }

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

type Surfaces []Surface

// FMap applies a function to each coords pair
func (seq Surfaces) FMap(f func(Coord)) {
	for _, x := range seq {
		x.FMap(f)
	}
}

// Bounding Box: The value of the bbox member MUST be an array of
// length 2*n where n is the number of dimensions represented in the
// contained geometries, with all axes of the most southwesterly point
// followed by all axes of the more northeasterly point.
type BoundingBox []float64

// South-West corner of Bounding Box
func (bbox BoundingBox) SouthWest() Coord {
	n := len(bbox) / 2
	return Coord(bbox[:n])
}

// North-East corner of Bounding Box
func (bbox BoundingBox) NorthEast() Coord {
	n := len(bbox) / 2
	return Coord(bbox[n:])
}

func (bbox BoundingBox) Join(box BoundingBox) {
	n := len(bbox) / 2
	sw := box.SouthWest()
	ne := box.NorthEast()

	if bbox[0] > sw.Lng() {
		bbox[0] = sw.Lng()
	}
	if bbox[1] > sw.Lat() {
		bbox[1] = sw.Lat()
	}

	if bbox[n] < ne.Lng() {
		bbox[n] = ne.Lng()
	}
	if bbox[n+1] < ne.Lat() {
		bbox[n+1] = ne.Lat()
	}
}

// Helper function to build bounding box
func boundingBox(seed Coord, coords interface{ FMap(f func(Coord)) }) BoundingBox {
	s, w := seed.LatLng()
	n, e := seed.LatLng()

	coords.FMap(func(c Coord) {
		lat, lng := c.LatLng()
		if lng < w {
			w = lng
		}
		if lng > e {
			e = lng
		}

		if lat < s {
			s = lat
		}
		if lat > n {
			n = lat
		}
	})

	return BoundingBox{w, s, e, n}
}
