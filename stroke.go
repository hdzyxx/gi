// Copyright (c) 2018, Randall C. O'Reilly. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gi

import (
	// "fmt"
	"image/color"
	"log"
)

// end-cap of a line: stroke-linecap property in SVG
type LineCap int

const (
	LineCapButt LineCap = iota
	LineCapRound
	LineCapSquare
)

// contrary to some docs, apparently need to run go generate manually
//go:generate stringer -type=LineCap

// the way in which lines are joined together: stroke-linejoin property in SVG
type LineJoin int

const (
	LineJoinMiter     LineJoin = iota
	LineJoinMiterClip          // SVG2 -- not yet supported
	LineJoinRound
	LineJoinBevel
	LineJoinArcs // SVG2 -- not yet supported
)

// contrary to some docs, apparently need to run go generate manually
//go:generate stringer -type=LineJoin

// PaintStroke contains all the properties specific to painting a line -- the svg elements define the corresponding SVG style attributes, which are processed in StrokeStyle
type PaintStroke struct {
	On         bool        `desc:"is stroke active -- if property is none then false"`
	Color      color.Color `desc:"default stroke color when such a color is needed -- Server could be anything"`
	Server     PaintServer `svg:"stroke",desc:"paint server for the stroke -- if solid color, defines the stroke color"`
	Width      float64     `svg:"stroke-width",desc:"line width"`
	Dashes     []float64   `svg:"stroke-dasharray",desc:"dash pattern"`
	Cap        LineCap     `svg:"stroke-linecap",desc:"how to draw the end cap of lines"`
	Join       LineJoin    `svg:"stroke-linejoin",desc:"how to join line segments"`
	MiterLimit float64     `svg:"stroke-miterlimit,min:"1",desc:"limit of how far to miter -- must be 1 or larger"`
}

// initialize default values for paint stroke
func (ps *PaintStroke) Defaults() {
	ps.On = false // svg says default is off
	ps.Server = NewSolidcolorPaintServer(color.Black)
	ps.Width = 1.0
	ps.Cap = LineCapButt
	ps.Join = LineJoinMiter // Miter not yet supported, but that is the default -- falls back on bevel
	ps.MiterLimit = 1.0
}

// todo: figure out more elemental, generic de-stringer kind of thing

// update the stroke settings from the style info on the node
func (ps *PaintStroke) SetFromNode(g *GiNode2D) {
	// always check if property has been set before setting -- otherwise defaults to empty -- true = inherit props
	if c, got := g.PropColor("stroke"); got { // todo: support url's to paint server elements!
		if c == nil {
			ps.On = false
		} else {
			ps.On = true
			ps.Color = c // todo: only if color
			ps.Server = NewSolidcolorPaintServer(c)
		}
	}
	if w, got := g.PropLength("stroke-width"); got {
		ps.Width = w
	}
	if _, got := g.PropNumber("stroke-opacity"); got {
		// todo: need to set the color alpha according to value
	}
	if es, got := g.PropEnum("stroke-linecap"); got {
		var lc LineCap = -1
		switch es { // first go through short-hand codes
		case "round":
			lc = LineCapRound
		case "butt":
			lc = LineCapButt
		case "square":
			lc = LineCapSquare
		}
		if lc == -1 {
			i, err := StringToLineCap(es) // stringer gen
			if err != nil {
				ps.Cap = i
			} else {
				log.Print(err)
			}
		} else {
			ps.Cap = lc
		}
	}
	if es, got := g.PropEnum("stroke-linejoin"); got {
		var lc LineJoin = -1
		switch es { // first go through short-hand codes
		case "miter":
			lc = LineJoinMiter
		case "miter-clip":
			lc = LineJoinMiterClip
		case "round":
			lc = LineJoinRound
		case "bevel":
			lc = LineJoinBevel
		case "arcs":
			lc = LineJoinArcs
		}
		if lc == -1 {
			i, err := StringToLineJoin(es) // stringer gen
			if err != nil {
				ps.Join = i
			} else {
				log.Print(err)
			}
		} else {
			ps.Join = lc
		}
	}
	if l, got := g.PropNumber("miter-limit"); got {
		ps.MiterLimit = l
	}
}