# Giraffe [![Build Status](https://travis-ci.org/benweidig/giraffe.svg?branch=master)](https://travis-ci.org/benweidig/giraffe)

A Go template renderer I've created for usage with [Gin Gonic](https://gin-gonic.github.io/gin/), supporting layouts, partials and different datasources like [packr](https://github.com/gobuffalo/packr). It should be usable without Gin, at least I'm going to try in another project.


## Get it

```
go get -u github.com/benweidig/giraffe
```

## Use it

### The simple way

You can simply use `giraffe.Default()` to get a usable `giraffe`:

```
import (
    "github.com/benweidig/giraffe"
)

func yourSetupMethodMaybe() {
    r := gin.New()

    r.HTMLRender = giraffe.Default()

    <or>

    r.HTMLRender = giraffe.Debug()
}
```

Both versions have some sensible defaults:

| Setting      | Default()      | Debug()        | Description                                     |
| ------------ | -------------- | -------------- | ----------------------------------------------- |
| Datasource   | `fs.Default()` | `fs.Default()` | A filesystem-based datasourse, see below        |
| Layout       | `layout`       | `layout`       | Filename/path of layout file, without extension |
| Funcs        | `[]`           | `[]`           | User-supplied template functions                |
| DisableCache | `false`        | `true`         | Caching                                         |

### The harder way

You can also create a `giraffe` with your own config:

```
import (
    "github.com/benweidig/giraffe"
    // You should rename the import or it will collide with packr
    gPacker "github.com/benweidig/giraffe/datasources/packr"
)

func yourSetupMethodMaybe() {
    r := gin.New()

    config := &giraffe.Config{
        Datasource:   &gPacker.Datasource{
            Box:        myPackrBox,
            Extensions: ".tpl",
        },
		Layout:       "master",
		Funcs:        make(template.FuncMap),
		DisableCache: false,
    }

    r.HTMLRender = giraffe.New(config)
}
```

## Datasources

The `giraffe` uses a `giraffe.Datasource` to load the actual template content. Two datasources are included in the box:

### fs.Datasource

```
import (
    "github.com/benweidig/giraffe/datasources/fs"
)
```

A simple filesystem datasource with 2 settings, but you could also use fs.Default() to get one with sensible defaults.

| Setting   | Default() | Description                                                         |
| --------- | --------- | ------------------------------------------------------------------- |
| Root      | `views`   | The root folder for the views                                       |
| Extension | `.html`   | The template extension (incl. the dot), so we can use shortes names |

### packr.Datasource

```
import (
    "github.com/benweidig/giraffe/datasources/packr"
)
```

A [packr.Box](https://github.com/gobuffalo/packr)-based datasource, no default is available.


| Setting   | Description                                                         |
| --------- | ------------------------------------------------------------------- |
| Box       | The `packr.Box`                                                     |
| Extension | The template extension (incl. the dot), so we can use shortes names |

## Convenience names

Templates are loaded by a more "convenient name" instead of the full/relative path. Partials are even shorter, you can ignore the partial folder name:

```
{{ partial "mypartial" . }}
```

To load the partial `mypartial` in the folder `partials`. Partials must reside in `partials` right now, so there's no need to write it all the time.

## License

MIT. See [LICENSE](LICENSE).
