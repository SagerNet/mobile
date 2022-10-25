// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build darwin || linux
// +build darwin linux

// Flappy Gopher is a simple one-button game that uses the
// mobile framework and the experimental sprite engine.
package main

import (
	"math/rand"
	"time"

	"github.com/sagernet/gomobile/app"
	"github.com/sagernet/gomobile/event/key"
	"github.com/sagernet/gomobile/event/lifecycle"
	"github.com/sagernet/gomobile/event/paint"
	"github.com/sagernet/gomobile/event/size"
	"github.com/sagernet/gomobile/event/touch"
	"github.com/sagernet/gomobile/exp/gl/glutil"
	"github.com/sagernet/gomobile/exp/sprite"
	"github.com/sagernet/gomobile/exp/sprite/clock"
	"github.com/sagernet/gomobile/exp/sprite/glsprite"
	"github.com/sagernet/gomobile/gl"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	app.Main(func(a app.App) {
		var glctx gl.Context
		var sz size.Event
		for e := range a.Events() {
			switch e := a.Filter(e).(type) {
			case lifecycle.Event:
				switch e.Crosses(lifecycle.StageVisible) {
				case lifecycle.CrossOn:
					glctx, _ = e.DrawContext.(gl.Context)
					onStart(glctx)
					a.Send(paint.Event{})
				case lifecycle.CrossOff:
					onStop()
					glctx = nil
				}
			case size.Event:
				sz = e
			case paint.Event:
				if glctx == nil || e.External {
					continue
				}
				onPaint(glctx, sz)
				a.Publish()
				a.Send(paint.Event{}) // keep animating
			case touch.Event:
				if down := e.Type == touch.TypeBegin; down || e.Type == touch.TypeEnd {
					game.Press(down)
				}
			case key.Event:
				if e.Code != key.CodeSpacebar {
					break
				}
				if down := e.Direction == key.DirPress; down || e.Direction == key.DirRelease {
					game.Press(down)
				}
			}
		}
	})
}

var (
	startTime = time.Now()
	images    *glutil.Images
	eng       sprite.Engine
	scene     *sprite.Node
	game      *Game
)

func onStart(glctx gl.Context) {
	images = glutil.NewImages(glctx)
	eng = glsprite.Engine(images)
	game = NewGame()
	scene = game.Scene(eng)
}

func onStop() {
	eng.Release()
	images.Release()
	game = nil
}

func onPaint(glctx gl.Context, sz size.Event) {
	glctx.ClearColor(1, 1, 1, 1)
	glctx.Clear(gl.COLOR_BUFFER_BIT)
	now := clock.Time(time.Since(startTime) * 60 / time.Second)
	game.Update(now)
	eng.Render(scene, now, sz)
}
