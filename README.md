# gi

GoGi is part of the GoKi Go language (golang) full strength tree structure system (ki = tree in Japanese)

`package gi` -- scenegraph based 2D and 3D GUI / graphics interface (Gi) in Go

GoDoc documentation: https://godoc.org/github.com/rcoreilly/goki/gi


# Code Map

* `ginode.go` -- basic `GiNode`, `GiNode2D`, `3D` structs and interfaces -- all Gi nodes are of this type
* `geom2d.go` -- All 2D geometry: Point2D, Size2D, etc
* `paint.go` -- `Paint` struct that does all the direct rendering, based on `gg`
	+ `stroke.go`, `fill.go` -- `PaintStroke` and `PaintFill` structs for stroke, fill settings
* `viewport2d.go` -- `Viewport2D` that has an `Image.RGBA` that `Paint` renders onto
* `window.go` -- `Window` is the top-level window that uses `go.wde` to open a gui window
* `shapes2d.go` -- All the basic 2D shapes: `GiRect`, `GiCircle` etc
* `text2d.go` -- Font, text rendering nodes
* `path.go` -- path-rendering nodes

# Design notes

* Excellent rendering code on top of freetype rasterizer, does almost everything we need: https://github.com/fogleman/gg -- borrowed heavily from that!

The 2D Gi is based entirely on the SVG2 spec: https://www.w3.org/TR/SVG2/Overview.html, and renders directly to an Image struct (`Viewport2D`)

The 3D Gi is based on TBD (will be impl later) and renders directly into a `Viewport3D` offscreen image buffer (OpenGL for now, but with generalization to Vulkan etc).

Any number of such (nested or otherwise) Viewport nodes can be created and they are all then composted into the final underlying bitmap of the Window.

Within a given rendering parent (Viewport2D or Viewport3D), only nodes of the appropriate type (`GiNode2D` or `GiNode3D`) may be used -- each has a pointer to their immediate parent viewport (indirectly through a ViewBox in 2D)

There are nodes to embed a Viewport2D bitmap within a Viewport3D scene, and vice-versa.  For a 2D viewport in a 3D scene, it acts like a texture and can be mapped onto a plane or anything else.  For a 3D viewport in a 2D scene, it just composts into the bitmap directly.

The overall parent Window can either provide a 2D or 3D viewport, which map directly into the underlying pixels of the window, and provide the "native" mode of the window, for efficiency.

## 2D Design

* Using the basic 64bit geom from fogleman/gg -- the `github.com/go-gl/mathgl/mgl32/` math elements (vectors, matricies) which build on the basic `golang.org/x/image/math/f32` did not have appropriate 2D rendering transforms etc.

* Using 64bit floats for coordinates etc because the spec says you need those for the "high quality" transforms, and Go defaults to them, and it just makes life easier -- won't have so crazy many coords in 2D space as we might have in 3D, where 32bit makes more sense and optimizes GPU hardware.

* SVG names use the lower-case starting camelCase convention, and often contain -'s, so we use a struct tag of svg: to map the UpperCase Go field names to their underlying SVG names for parsing etc.  E.g., use `"{min-x,min-y}"` to label the vector elements of the `Min Vec2` field.

* The SVG default coordinate system has 0,0 at the upper-left.  The default 3D coordinate system flips the Y axis so 0,0 is at the lower left effectively (actually it uses center-based coordinates so 0,0 is in the center of the image, effectively -- everything is defined by the camera anyway)

## 3D Design

* keep all the elements separate: geometry, material, transform, etc.  Including shader programs.  Maximum combinatorial flexibility.  not clear if Qt3D really obeys this principle, but Inventor does, and probably other systems do to.


# Links

* SVG *text* generator in go: https://github.com/ajstarks/svgo
* cairo wrapper in go: https://github.com/ungerik/go-cairo -- maybe needed for PDF generation from SVG?
* https://github.com/jyotiska/go-webcolors -- need this for parsing colors
* https://golang.org/pkg/image/
* https://godoc.org/golang.org/x/image/vector
* https://godoc.org/github.com/golang/freetype/raster
* https://github.com/fogleman/gg -- key lib using above -- 2D rendering!

* Shiny (not much progress recently, only works on android?):  https://github.com/golang/go/issues/11818 https://github.com/golang/exp/tree/master/shiny

* Current plans for GUI based on OpenGL: https://docs.google.com/document/d/1mXev7TyEnvM4t33lnqoji-x7EqGByzh4RpE4OqEZck4

* Window events: https://github.com/skelterjohn/go.wde, Material gui https://github.com/skelterjohn/go.wde

* scenegraphs / 3D game engines: 
	+ https://github.com/g3n/engine
	+ https://github.com/oakmound/oak
	+ https://github.com/walesey/go-engine