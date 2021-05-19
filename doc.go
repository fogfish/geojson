//
// Copyright (C) 2021 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/fogfish/geojson
//

/*

Package geojson implements a type safe codec for GeoJSON, which is a popular
format for encoding a variety of geographic data structures. It defines a
standard way to express points, curves, and surfaces in coordinate space
together with an application specific metadata about it.

The library allows developers to use Golang struct to define domain models
using a type safe approach of encoding/decoding these models to GeoJSON and
back. The library uses type tagging technique to annotate any structure as
GeoJSON feature:

  type City struct {
    geojson.Feature
    Name      string `json:"name,omitempty"`
  }

*/
package geojson
