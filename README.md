<p align="center">
  <h3 align="center">{ ㆞ }</h3>
  <h3 align="center">GeoJSON</h3>
  <p align="center"><strong>GeoJSON codec for Go structs</strong></p>

  <p align="center">
    <!-- Version -->
    <a href="https://github.com/fogfish/geojson/releases">
      <img src="https://img.shields.io/github/v/tag/fogfish/geojson?label=version" />
    </a>
    <!-- Documentation -->
    <a href="https://pkg.go.dev/github.com/fogfish/geojson">
      <img src="https://pkg.go.dev/badge/github.com/fogfish/geojson" />
    </a>
    <!-- Build Status  -->
    <a href="https://github.com/fogfish/geojson/actions/">
      <img src="https://github.com/fogfish/geojson/workflows/build/badge.svg" />
    </a>
    <!-- GitHub -->
    <a href="http://github.com/fogfish/geojson">
      <img src="https://img.shields.io/github/last-commit/fogfish/geojson.svg" />
    </a>
    <!-- Coverage -->
    <a href="https://coveralls.io/github/fogfish/geojson?branch=main">
      <img src="https://coveralls.io/repos/github/fogfish/geojson/badge.svg?branch=main" />
    </a>
    <!-- Go Card -->
    <a href="https://goreportcard.com/report/github.com/fogfish/geojson">
      <img src="https://goreportcard.com/badge/github.com/fogfish/geojson" />
    </a>
  </p>
</p>

---

The library implements a type safe codec for [GeoJSON](https://geojson.org) with the focus on encoding application specific data.  

## Inspiration

GeoJSON is a popular format for encoding a variety of geographic data structures. It defines a standard way to express points, curves, and surfaces in coordinate space together with an application specific metadata about it.

```json
{
  "type": "Feature",
  "id": "[wikipedia:Helsinki]",
  "geometry": {
    "type": "Point",
    "coordinates": [24.9384, 60.1699]
  },
  "properties": {
    "name": "Helsinki"
  }
}
```

Unfortunately, efficient and type-safe implementation of GeoJSON codec can be challenging:

(i) Pure structs are verbose. The `properties` is an application specific and it's type is controlled outside of the codec library. Usage of duck type (`interface{}`) is a common trait used by other GeoJSON Golang libraries. As the results, developers are misses an ability to caught errors at compile time, any mistake becomes visible at run time as a panic. [interface{} says nothing.](https://youtu.be/PAAkCSZUG1c?t=7m40s).

(ii) Implementing GeoJSON types using generics is requires overly complex type definitions. It can lead to complex type hierarchies, especially in nested GeoJSON structures like FeatureCollection. It suffers from usability for client application. Specifying types for properties at every level (e.g., Feature or FeatureCollection) adds boilerplate and increases the learning curve.



## Key features

The library allows developers to use pure Golang struct to define domain models using a type safe approach of encoding/decoding these models to GeoJSON and back. The library uses type tagging technique to annotate any structure as GeoJSON feature:  

```go
type City struct {
  geojson.Feature
  Name      string `json:"name,omitempty"`
}
```


## Getting started

The latest version of the library is available at `main` branch of this repository. All development, including new features and bug fixes, take place on the `main` branch using forking and pull requests as described in contribution guidelines. The stable version is available via Golang modules.

The following code snippet demonstrates a typical usage scenario.

```go
import "github.com/fogfish/geojson"

//
// declare any domain type and annotate as a geojson.Feature
type City struct {
  geojson.Feature
  Name      string `json:"name,omitempty"`
}

//
// Each GeoJSON type declares JSON codes using helper functions.
func (x City) MarshalJSON() ([]byte, error) {
	type tStruct City
	return x.Feature.EncodeGeoJSON(tStruct(x))
}

func (x *City) UnmarshalJSON(b []byte) error {
	type tStruct *City
	return x.Feature.DecodeGeoJSON(b, tStruct(x))
}

//
// Create new instance of the type
city := City{
  Feature: geojson.NewPoint(
    "[wikipedia:Helsinki]",
    geojson.Coord{24.9384, 60.1699},
  ),
  Name: "Helsinki",
}

//
// Use type checks to validate the type of the Geometry 
city.Feature.Geometry.Coords.(*geojson.Point)
```

### Feature Collection

The library support feature collection through the collection type. It represents a collection of spatially bounded elements, as defined by the GeoJSON FeatureCollection standard. This construct is designed to support ["foreign members"](https://www.rfc-editor.org/rfc/rfc7946#section-6.1) for improved exchange of geospatial data. The value of a "foreign member" is determined by the application.

**This library reuses the "properties" attribute, which acts as a foreign member within the context of collections.**

```go
type Cities struct {
  geojson.Collection[City]
  Country string `json:"country,omitempty"`
}

//
// Each GeoJSON type declares JSON codes using helper functions.
func (x Cities) MarshalJSON() ([]byte, error) {
	type tStruct Cities
	return x.Collection.EncodeGeoJSON(tStruct(x))
}

func (x *Cities) UnmarshalJSON(b []byte) error {
	type tStruct *Cities
	return x.Collection.DecodeGeoJSON(b, tStruct(x))
}
```


## How To Contribute

The library is [MIT](LICENSE) licensed and accepts contributions via GitHub pull requests:

1. Fork it
2. Create your feature branch (`git checkout -b my-new-feature`)
3. Commit your changes (`git commit -am 'Added some feature'`)
4. Push to the branch (`git push origin my-new-feature`)
5. Create new Pull Request

The build and testing process requires [Go](https://golang.org) version 1.16 or later.

**build** and **test** library.

```bash
git clone https://github.com/fogfish/geojson
cd geojson
go test
```

### commit message

The commit message helps us to write a good release note, speed-up review process. The message should address two question what changed and why. The project follows the template defined by chapter [Contributing to a Project](http://git-scm.com/book/ch5-2.html) of Git book.

### bugs

If you experience any issues with the library, please let us know via [GitHub issues](https://github.com/fogfish/geojson/issue). We appreciate detailed and accurate reports that help us to identity and replicate the issue. 


## License

[![See LICENSE](https://img.shields.io/github/license/fogfish/geojson.svg?style=for-the-badge)](LICENSE)
