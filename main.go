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
	width  = 700.0
	height = 700.0
	scaleX = 0.03
	scaleY = 0.018
	shiftX = 1.9
	shiftY = 3.6
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
var placeVisits []dto.PlaceVisit
var placeVisitsVao uint32

func init() {
	flag.IntVar(&fps, "fps", fps, "Sets the frames-per-second, used set the speed of the simulation.")
	flag.Parse()

	files := []string{"2024_MARCH.json", "2024_FEBRUARY.json", "2024_JANUARY.json", "2023_DECEMBER.json", "2023_NOVEMBER.json", "2023_OCTOBER.json", "2023_SEPTEMBER.json", "2023_AUGUST.json", "2023_JULY.json", "2023_JUNE.json", "2023_MAY.json"}
	for _, file := range files {
		month = loadMonth(file)
		for i := 0; i < len(month.TimelineObjects); i++ {
			o := &month.TimelineObjects[i]

			if o.PlaceVisit.CentreLatE7 == 0 {
				// y, m, _ := o.ActivitySegment.Duration.StartTimestamp.Date()

				// if y == 2023 && m == 12 {
				// fmt.Println(o.ActivitySegment.Duration.StartTimestamp)
				if o.ActivitySegment.StartLocation.LongitudeE7 > 1e7 && o.ActivitySegment.EndLocation.LongitudeE7 > 1e7 {
					activitySegments = append(activitySegments, o.ActivitySegment)
				}
				// }
			} else {
				placeVisits = append(placeVisits, o.PlaceVisit)
				for i := 0; i < len(o.PlaceVisit.ChildVisits); i++ {
					placeVisits = append(placeVisits, o.PlaceVisit.ChildVisits[i])
				}
			}
		}
	}

	fmt.Println(len(activitySegments), " activity segments; ", len(placeVisits), " place visits.")
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

	var lat float32
	var lng float32

	for i := 0; i < len(activitySegments); i++ {
		activitySegment := &activitySegments[i]
		if activitySegment.StartLocation.LongitudeE7 > 123000000 && activitySegment.StartLocation.LongitudeE7 < 127000000 {
			lat = float32(activitySegment.StartLocation.LatitudeE7) / 1e7
			lng = float32(activitySegment.StartLocation.LongitudeE7) / 1e7
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

		for _, p := range activitySegment.WaypointPath.Waypoints {
			if p.LngE7 > 123000000 && p.LngE7 < 127000000 {
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

		if activitySegment.EndLocation.LongitudeE7 > 123000000 && activitySegment.EndLocation.LongitudeE7 < 127000000 {
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
		}
	}

	for i := 0; i < len(activitySegments); i++ {
		activitySegment := &activitySegments[i]
		err := makeWaypointPathVao(activitySegment)
		if err != nil {
			fmt.Println(err)
		}
	}

	points := make([]float32, 2*len(placeVisits))
	pidx := 0
	for i := 0; i < len(placeVisits); i++ {
		points[pidx] = float32(placeVisits[i].Location.LongitudeE7) / 1e7
		pidx++
		points[pidx] = float32(placeVisits[i].Location.LatitudeE7) / 1e7
		pidx++
	}
	// fmt.Println(points)
	// largest := maxlat - minlat
	// if maxlng-minlng > largest {
	// 	largest = maxlng - minlng
	// }
	for i, p := range points {
		if i%2 == 0 {
			points[i] = ((p - minlng) / scaleX) - shiftX
		} else {
			points[i] = ((p - minlat) / scaleY) - shiftY
		}
	}
	// fmt.Println(points)

	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(points), gl.Ptr(points), gl.STATIC_DRAW)

	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)
	gl.EnableVertexAttribArray(0)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.VertexAttribPointer(0, 2, gl.FLOAT, false, 0, nil)
	placeVisitsVao = vao

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
			// fmt.Println(waypointPath)
			gl.BindVertexArray(waypointPath.Vao)
			gl.DrawArrays(gl.LINE_STRIP, 0, int32(len(waypointPath.WaypointPath.Waypoints)+2))
		}
	}

	gl.PointSize(5.0)
	gl.BindVertexArray(placeVisitsVao)
	gl.DrawArrays(gl.POINTS, 0, int32(len(placeVisits)))

	glfw.PollEvents()
	window.SwapBuffers()
	return nil
}
