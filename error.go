//
// Copyright (C) 2021 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/fogfish/geojson
//

package geojson

type Error string

func (err Error) Error() string { return string(err) }

const (
	ErrorUnsupportedType = Error("GeoJSON type is not supported")
)
