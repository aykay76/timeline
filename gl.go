package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/aykay76/timeline/dto"
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
)

const (
	vertexShaderSource = `
		#version 400

		in vec3 vp;
		void main() {
			gl_Position = vec4(vp, 1.0);
		}
	` + "\x00"

	fragmentShaderSource = `
		#version 400

		out vec4 frag_colour;
		void main() {
  			frag_colour = vec4(1, 1, 1, 1.0);
		}
	` + "\x00"
)

// initGlfw initializes glfw and returns a Window to use.
func initGlfw() *glfw.Window {
	if err := glfw.Init(); err != nil {
		panic(err)
	}
	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	window, err := glfw.CreateWindow(width, height, "Go-OpenGL", nil, nil)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()

	return window
}

// initOpenGL initializes OpenGL and returns an intiialized program.
func initOpenGL() uint32 {
	if err := gl.Init(); err != nil {
		panic(err)
	}
	version := gl.GoStr(gl.GetString(gl.VERSION))
	log.Println("OpenGL version", version)

	vertexShader, err := compileShader(vertexShaderSource, gl.VERTEX_SHADER)
	if err != nil {
		panic(err)
	}

	fragmentShader, err := compileShader(fragmentShaderSource, gl.FRAGMENT_SHADER)
	if err != nil {
		panic(err)
	}

	prog := gl.CreateProgram()
	gl.AttachShader(prog, vertexShader)
	gl.AttachShader(prog, fragmentShader)
	gl.LinkProgram(prog)
	return prog
}

// https://github.com/go-gl/examples/blob/master/gl41core-cube/cube.go
func compileShader(source string, shaderType uint32) (uint32, error) {
	shader := gl.CreateShader(shaderType)

	csources, free := gl.Strs(source)
	gl.ShaderSource(shader, 1, csources, nil)
	free()
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to compile %v: %v", source, log)
	}

	return shader, nil
}

func makeWaypointPathVao(path *dto.ActivitySegment) error {
	if path.StartLocation.LongitudeE7 < 123000000 || path.StartLocation.LongitudeE7 > 127000000 {
		return fmt.Errorf("start location is west of 1 degree")
	}
	if path.EndLocation.LongitudeE7 < 123000000 || path.EndLocation.LongitudeE7 > 127000000 {
		return fmt.Errorf("end location is west of 1 degree")
	}
	for _, point := range path.WaypointPath.Waypoints {
		if point.LngE7 < 123000000 || point.LngE7 > 127000000 {
			return fmt.Errorf("waypoint is west of 1 degree")
		}
	}

	pidx := 0
	points := make([]float32, 2*len(path.WaypointPath.Waypoints)+4)
	points[pidx] = float32(path.StartLocation.LongitudeE7) / 1e7
	pidx++
	points[pidx] = float32(path.StartLocation.LatitudeE7) / 1e7
	pidx++
	for _, point := range path.WaypointPath.Waypoints {
		points[pidx] = float32(point.LngE7) / 1e7
		pidx++
		points[pidx] = float32(point.LatE7) / 1e7
		pidx++
	}
	points[pidx] = float32(path.EndLocation.LongitudeE7) / 1e7
	pidx++
	points[pidx] = float32(path.EndLocation.LatitudeE7) / 1e7
	pidx++
	// fmt.Println(points)
	// largest := maxlat - minlat
	// if maxlat-minlat > largest {
	// 	largest = maxlat - minlat
	// }
	// fmt.Println(largest)
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
	path.Vao = vao

	return nil
}
