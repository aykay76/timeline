package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/aykay76/timeline/dto"
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
)

const (
	width  = 800.0
	height = 600.0
)

var (
	fps = 20
)

var month *dto.Month
var waypointVao uint32

func init() {
	flag.IntVar(&fps, "fps", fps, "Sets the frames-per-second, used set the speed of the simulation.")
	flag.Parse()

	month = loadMonth("2023_APRIL.json")
}

func loadMonth(path string) *dto.Month {
	var month dto.Month

	fileBytes, err := os.ReadFile(path)
	if err == nil {
		err = json.Unmarshal(fileBytes, &month)
		if err != nil {
			fmt.Println(err)
			return nil
		}
	} else {
		fmt.Println(err)
		return nil
	}

	return &month
}

func main() {
	// GLFW event handling must run on the main OS thread
	runtime.LockOSThread()
	window := initGlfw()
	defer glfw.Terminate()
	prog := initOpenGL()

	waypointVao = makeWaypointPathVao(month.TimelineObjects[0].ActivitySegment.WaypointPath)

	// Make the cells and start the loop
	t := time.Now()
	for !window.ShouldClose() {
		tick()

		if err := draw(prog, window); err != nil {
			panic(err)
		}

		time.Sleep(time.Second/time.Duration(fps) - time.Since(t))
		t = time.Now()
	}
}

func tick() {
}

func draw(prog uint32, window *glfw.Window) error {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.UseProgram(prog)

	gl.BindVertexArray(waypointVao)
	gl.DrawArrays(gl.LINE_STRIP, 0, int32(len(month.TimelineObjects[0].ActivitySegment.WaypointPath.Waypoints)))

	glfw.PollEvents()
	window.SwapBuffers()
	return nil
}
