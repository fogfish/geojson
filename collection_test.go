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
	"github.com/fogfish/it/v2"
)

func TestCollection(t *testing.T) {
	spb := GeoJsonCity{
		Feature: geojson.NewPoint("city:spb", geojson.Coord{100.0, 0.0}),
		City:    City{Name: "Saint-Petersburg"},
	}

	hel := GeoJsonCity{
		Feature: geojson.NewPoint("city:hel", geojson.Coord{101.0, 1.0}),
		City:    City{Name: "Helsinki"},
	}

	sto := GeoJsonCity{
		Feature: geojson.NewPoint("city:sto", geojson.Coord{102.0, 2.0}),
		City:    City{Name: "Stockholm"},
	}

	seq := geojson.Collection[GeoJsonCity]{
		Features: []GeoJsonCity{spb, hel, sto},
	}

	bin, err := json.Marshal(seq)
	it.Then(t).Should(it.Nil(err))

	var c geojson.Collection[GeoJsonCity]
	err = json.Unmarshal([]byte(bin), &c)

	it.Then(t).Should(
		it.Nil(err),
		it.Equiv(c.Features[0], spb),
		it.Equiv(c.Features[1], hel),
		it.Equiv(c.Features[2], sto),
	)
}
