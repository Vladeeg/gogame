package main

import (
	"github.com/gen2brain/raylib-go/raylib"
)

type UpdatingSystem struct {
	ComponentKind string
	Action func(c *Component, elapsed float32)
}

type RenderingSystem struct {
	ComponentKind string
	Action func(c *Component)
}

type Component struct {
	Kind string
	EntityID string
	Props interface{}
}

type Entity struct {
	Name string
	Components map[string]Component
}


var (
	components map[string][]*Component
	updatingSystems []UpdatingSystem
	renderingSystems []RenderingSystem
	secondsPerUpdate float32 = 1 / 60
)

type TransformProps struct {
	x, y float32
}

type SpeedProps struct {
	x, y float32
}

type GraphicsRectProps struct {
	offsetX, offsetY, w, h float32
	color rl.Color
}

func Move(c *Component, elapsed float32) {
	props := c.Props.(TransformProps)

	var speedComp *Component
	for _, v := range components["Speed"] {
		if v.EntityID == c.EntityID {
			speedComp = v
		}
	}

	speedProps := speedComp.Props.(SpeedProps)
	props.x += speedProps.x * elapsed
	props.y += speedProps.y * elapsed
	c.Props = props
}

func Bounce(c *Component, elapsed float32) {
	props := c.Props.(SpeedProps)

	var transformComp *Component
	for _, v := range components["Transform"] {
		if v.EntityID == c.EntityID {
			transformComp = v
		}
	}
	var graphicsComp *Component
	for _, v := range components["Graphics:Rect"] {
		if v.EntityID == c.EntityID {
			graphicsComp = v
		}
	}

	transformProps := transformComp.Props.(TransformProps)
	graphicsProps := graphicsComp.Props.(GraphicsRectProps)
	
	w := float32(rl.GetScreenWidth())
	if transformProps.x > w - graphicsProps.w {
		transformProps.x = w - graphicsProps.w
		props.x *= -1
	}
	if transformProps.x < 0 {
		transformProps.x = 0
		props.x *= -1
	}
	h := float32(rl.GetScreenHeight())
	if transformProps.y > h - graphicsProps.h {
		transformProps.y = h - graphicsProps.h
		props.y *= -1
	}
	if transformProps.y < 0 {
		transformProps.y = 0
		props.y *= -1
	}

	transformComp.Props = transformProps
	c.Props = props
}

func RenderGraphicRect(c *Component) {
	props := c.Props.(GraphicsRectProps)
	var transformComp *Component
	for _, v := range components["Transform"] {
		if v.EntityID == c.EntityID {
			transformComp = v
		}
	}

	transformProps := transformComp.Props.(TransformProps)

	rl.DrawRectangle(
		int32(props.offsetX) + int32(transformProps.x),
		int32(props.offsetY) + int32(transformProps.y),
		int32(props.w),
		int32(props.h),
		props.color,
	)
}

func Setup() {
	rl.InitWindow(800, 450, "raylib [core] example - basic window")

	sw := float32(rl.GetScreenWidth())
	sh := float32(rl.GetScreenHeight())

	components = make(map[string][]*Component)
	components["Transform"] = []*Component{
		{
			"Transform",
			"1",
			TransformProps{10, 10},
		},
		{
			"Transform",
			"2",
			TransformProps{sw - 20, 15},
		},
		{
			"Transform",
			"3",
			TransformProps{28, sh - 150},
		},
	}
	components["Speed"] = []*Component{
		{
			"Speed",
			"1",
			SpeedProps{80, 80},
		},
		{
			"Speed",
			"2",
			SpeedProps{-60, 60},
		},
		{
			"Speed",
			"3",
			SpeedProps{-80, -65},
		},
	}
	components["Graphics:Rect"] = []*Component{
		{
			"Graphics:Rect",
			"1",
			GraphicsRectProps{0, 0, 20, 20, rl.Blue},
		},
		{
			"Graphics:Rect",
			"2",
			GraphicsRectProps{0, 0, 25, 10, rl.Red},
		},
		{
			"Graphics:Rect",
			"3",
			GraphicsRectProps{0, 0, 20, 20, rl.Lime},
		},
	}

	updatingSystems = append(updatingSystems, UpdatingSystem{
		"Transform",
		Move,
	})
	updatingSystems = append(updatingSystems, UpdatingSystem{
		"Speed",
		Bounce,
	})
	renderingSystems = append(renderingSystems, RenderingSystem{
		"Graphics:Rect",
		RenderGraphicRect,
	})
}

func ProcessInput() {
//?
}

func Update(elapsed float32) {
	for _, v := range updatingSystems {
		comps := components[v.ComponentKind]
		for _, c := range comps {
			v.Action(c, elapsed)
		}
	}
}

func Render() {
	rl.BeginDrawing()
	rl.ClearBackground(rl.Black)

	for _, v := range renderingSystems {
		comps := components[v.ComponentKind]
		for _, c := range comps {
			v.Action(c)
		}
	}

	rl.EndDrawing()
}

func main() {
	Setup()

	lastTime := rl.GetTime()
	for !rl.WindowShouldClose() {
		current := rl.GetTime()
		elapsed := current - lastTime
		
		ProcessInput()
		Update(elapsed)
		Render()

		lastTime = current
	}

	rl.CloseWindow()
}