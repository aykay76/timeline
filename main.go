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
var minlat float32
var minlng float32
var maxlat float32
var maxlng float32
var activitySegments []dto.ActivitySegment

func init() {
	flag.IntVar(&fps, "fps", fps, "Sets the frames-per-second, used set the speed of the simulation.")
	flag.Parse()

	files := []string{"2024_MARCH.json"} //, "2024_FEBRUARY.json", "2024_JANUARY.json", "2023_DECEMBER.json", "2023_NOVEMBER.json", "2023_OCTOBER.json", "2023_SEPTEMBER.json", "2023_AUGUST.json", "2023_JULY.json", "2023_JUNE.json", "2023_MAY.json"}
	for _, file := range files {
		month = loadMonth(file)
		for i := 0; i < len(month.TimelineObjects); i++ {
			o := &month.TimelineObjects[i]

			if o.PlaceVisit.CentreLatE7 == 0 {
				// y, m, d := o.ActivitySegment.Duration.StartTimestamp.Date()

				// if y == 2023 && m == 12 && d == 6 && o.ActivitySegment.Duration.StartTimestamp.Hour() == 15 {
				fmt.Println(o.ActivitySegment.Duration.StartTimestamp)
				activitySegments = append(activitySegments, o.ActivitySegment)
				// }
			}
		}
	}
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

	minlat = float32(999.0)
	minlng = float32(999.0)
	maxlat = float32(-999.0)
	maxlng = float32(-999.0)

	for i := 0; i < len(activitySegments); i++ {
		activitySegment := &activitySegments[i]
		lat := float32(activitySegment.StartLocation.LatitudeE7) / 1e7
		lng := float32(activitySegment.StartLocation.LongitudeE7) / 1e7
		if lng < minlng {
			minlng = lng
		}
		if lng > maxlng {
			maxlng = lng
		}
		if lat < minlat {
			minlat = lat
		}
		if lat > maxlat {
			maxlat = lat
		}

		for _, p := range activitySegment.WaypointPath.Waypoints {
			if p.LatE7 >= 550000000 {
				lat := float32(p.LatE7) / 1e7
				lng := float32(p.LngE7) / 1e7

				if lng < minlng {
					minlng = lng
				}
				if lng > maxlng {
					maxlng = lng
				}
				if lat < minlat {
					minlat = lat
				}
				if lat > maxlat {
					maxlat = lat
				}
			}
		}

		lat = float32(activitySegment.EndLocation.LatitudeE7) / 1e7
		lng = float32(activitySegment.EndLocation.LongitudeE7) / 1e7
		if lng < minlng {
			minlng = lng
		}
		if lng > maxlng {
			maxlng = lng
		}
		if lat < minlat {
			minlat = lat
		}
		if lat > maxlat {
			maxlat = lat
		}

		makeWaypointPathVao(activitySegment)
	}
	fmt.Println(minlat, minlng, maxlat, maxlng)

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

	for i := 0; i < len(activitySegments); i++ {
		waypointPath := &activitySegments[i]
		if waypointPath.Vao != 0xffff {
			gl.BindVertexArray(waypointPath.Vao)
			gl.DrawArrays(gl.LINE_STRIP, 0, int32(len(waypointPath.WaypointPath.Waypoints)+2))
		}
	}

	glfw.PollEvents()
	window.SwapBuffers()
	return nil
}
