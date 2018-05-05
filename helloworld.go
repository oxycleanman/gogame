package main

import (
	"fmt"
	"github.com/veandco/go-sdl2/sdl"
	"time"
	"math"
	"math/rand"
)

const winWidth, winHeight int = 800, 600

type gameState int

const (
	start gameState = iota
	play
)

var state = start

var nums = [][]byte{
	{
		1, 1, 1,
		1, 0, 1,
		1, 0, 1,
		1, 0, 1,
		1, 1, 1,
	},
	{
		1, 1, 0,
		0, 1, 0,
		0, 1, 0,
		0, 1, 0,
		1, 1, 1,
	},
	{
		1, 1, 1,
		0, 0, 1,
		1, 1, 1,
		1, 0, 0,
		1, 1, 1,
	},
	{
		1, 1, 1,
		0, 0, 1,
		1, 1, 1,
		0, 0, 1,
		1, 1, 1,
	},
}

type color struct {
	r, g, b byte
}

type pos struct {
	x, y float32
}

type ball struct {
	pos
	radius float32
	xv     float32
	yv     float32
	speed  float32
	color  color
}

func drawNumber(pos pos, color color, size int, num int, pixels []byte) {
	startX := int(pos.x) - (size*3)/2
	startY := int(pos.y) - (size*5)/2

	for i, v := range nums[num] {
		if v == 1 {
			for y := startY; y < startY+size; y++ {
				for x := startX; x < startX+size; x++ {
					setPixel(x, y, color, pixels)
				}
			}
		}
		startX += size
		if (i+1)%3 == 0 {
			startY += size
			startX -= size * 3
		}
	}
}

func (ball *ball) draw(pixels []byte) {
	for y := -ball.radius; y < ball.radius; y++ {
		for x := -ball.radius; x < ball.radius; x++ {
			if x*x+y*y < ball.radius*ball.radius {
				setPixel(int(ball.x+x), int(ball.y+y), ball.color, pixels)
			}
		}
	}
}

func getCenter() pos {
	return pos{float32(winWidth) / 2, float32(winHeight) / 2}
}

func updateBall(paddle *paddle, keyState []uint8, ball *ball, increase float32) {
	if keyState[sdl.SCANCODE_UP] != 0 {
		ball.yv -= increase
	} else if keyState[sdl.SCANCODE_DOWN] != 0 {
		ball.yv += increase
	}
	ball.xv = -ball.xv
	ball.x = paddle.x + paddle.w/2.0 + ball.radius
}

func resetOnScore(ball *ball, ballColor color, rPaddle *paddle, lPaddle *paddle) {
	ball.pos = getCenter()
	ball.speed = 1.00
	ball.yv = rand.Float32()
	ball.color = color{255,255,255}
	rPaddle.y = float32(winHeight/2)
	lPaddle.y = float32(winHeight/2)
	state = start
}

func (ball *ball) update(leftPaddle *paddle, rightPaddle *paddle, elapsedTime float32, keyState []uint8) {
	ball.x += ball.xv * elapsedTime * ball.speed
	ball.y += ball.yv * elapsedTime * ball.speed
	ball.speed += 0.001
	if math.Mod(float64(ball.speed), float64(0.2)) < 0.001 && ball.color.g > 10 && ball.color.b > 10 {
		ball.color.g -= 10
		ball.color.b -= 10
	}

	if ball.y-ball.radius < 0 || ball.y+ball.radius > float32(winHeight) {
		ball.yv = -ball.yv
	}

	if ball.x < 0 {
		rightPaddle.score++
		resetOnScore(ball, color{255, 255, 255}, rightPaddle, leftPaddle)
	} else if ball.x > float32(winWidth) {
		leftPaddle.score++
		resetOnScore(ball, color{255, 255, 255}, rightPaddle, leftPaddle)
	}

	if ball.x-ball.radius < leftPaddle.x+leftPaddle.w/2 {
		if (ball.y > leftPaddle.y-leftPaddle.h/2 && ball.y < leftPaddle.y-leftPaddle.h/8) ||
			(ball.y >= leftPaddle.y+leftPaddle.h/8 && ball.y < leftPaddle.y+leftPaddle.h/2) {
			updateBall(leftPaddle, keyState, ball, 20)
		}
		if (ball.y >= leftPaddle.y-leftPaddle.h/8 && ball.y < leftPaddle.y-leftPaddle.h/4) ||
			(ball.y >= leftPaddle.y+leftPaddle.h/4 && ball.y < leftPaddle.y+leftPaddle.h/8) {
			updateBall(leftPaddle, keyState, ball, 10)
		}
		if ball.y >= leftPaddle.y-leftPaddle.h/4 && ball.y < leftPaddle.y+leftPaddle.h/4 {
			updateBall(leftPaddle, keyState, ball, 0)
		}
	}

	if ball.x+ball.radius > rightPaddle.x-rightPaddle.w/2 {
		if ball.y > rightPaddle.y-rightPaddle.h/2 && ball.y < rightPaddle.y+rightPaddle.h/2 {
			ball.xv = -ball.xv
			ball.x = rightPaddle.x - rightPaddle.w/2.0 - ball.radius
		}
	}
}

type paddle struct {
	pos
	w     float32
	h     float32
	speed float32
	score int
	color color
}

func lerp(a float32, b float32, pct float32) float32 {
	return a + pct*(b-a)
}

func (paddle *paddle) draw(pixels []byte) {
	startX := int(paddle.x - paddle.w/2)
	startY := int(paddle.y - paddle.h/2)

	for y := 0; y < int(paddle.h); y++ {
		for x := 0; x < int(paddle.w); x++ {
			setPixel(startX+x, startY+y, paddle.color, pixels)
		}
	}

	numX := lerp(paddle.x, getCenter().x, 0.2)
	drawNumber(pos{numX, 35}, paddle.color, 10, paddle.score, pixels)
}

func (paddle *paddle) update(keyState []uint8, elapsedTime float32) {
	if keyState[sdl.SCANCODE_UP] != 0 && !(paddle.y - paddle.h/2 <= 0) {
		paddle.y -= paddle.speed * elapsedTime
	}
	if keyState[sdl.SCANCODE_DOWN] != 0 && !(paddle.y + paddle.h/2 >= float32(winHeight)){
		paddle.y += paddle.speed * elapsedTime
	}
}

func (paddle *paddle) aiUpdate(ball *ball, elapsedTime float32) {
	if paddle.y + paddle.h/2 * 0.7 < ball.y + ball.radius {
		paddle.y += paddle.speed * elapsedTime
	} else if paddle.y - paddle.h/2 * 0.7 > ball.y + ball.radius {
		paddle.y -= paddle.speed * elapsedTime
	}
}

func clear(pixels []byte) {
	for i := range pixels {
		pixels[i] = 0
	}
}

func setPixel(x, y int, c color, pixels []byte) {
	index := (y*winWidth + x) * 4

	if index < len(pixels)-4 && index >= 0 {
		pixels[index] = c.r
		pixels[index+1] = c.g
		pixels[index+2] = c.b
	}
}

func main() {

	err := sdl.Init(sdl.INIT_EVERYTHING)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer sdl.Quit()

	window, err := sdl.CreateWindow("Testing", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		int32(winWidth), int32(winHeight), sdl.WINDOW_SHOWN)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer window.Destroy()
	window.SetResizable(true)

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer renderer.Destroy()

	tex, err := renderer.CreateTexture(sdl.PIXELFORMAT_ABGR8888, sdl.TEXTUREACCESS_STREAMING, int32(winWidth), int32(winHeight))
	if err != nil {
		fmt.Println(err)
		return
	}
	defer tex.Destroy()

	pixels := make([]byte, winWidth*winHeight*4)

	player1 := paddle{pos{50, float32(winHeight/2)}, float32(winWidth)*0.02, float32(winHeight)*0.20, 500, 0, color{255, 255, 255}}
	player2 := paddle{pos{float32(winWidth - 50), float32(winHeight/2)}, float32(winWidth)*0.02, float32(winHeight)*0.20, 500, 0, color{255, 255, 255}}
	ball := ball{getCenter(), float32(winHeight)*0.02, 100, 0, 1.00, color{255, 255, 255}}

	keyState := sdl.GetKeyboardState()

	var frameStart time.Time
	var elapsedTime float32

	for {
		frameStart = time.Now()
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				return
			}
		}

		if state == play {
			drawNumber(getCenter(), color{255, 255, 255}, 20, 2, pixels)
			ball.update(&player1, &player2, elapsedTime, keyState)
			player1.update(keyState, elapsedTime)
			player2.aiUpdate(&ball, elapsedTime)
		} else if state == start {
			if keyState[sdl.SCANCODE_SPACE] != 0 {
				if player1.score == 3 || player2.score == 3 {
					player1.score = 0
					player2.score = 0
				}
				state = play
			}
		}

		clear(pixels)
		player1.draw(pixels)
		player2.draw(pixels)
		ball.draw(pixels)

		tex.Update(nil, pixels, winWidth*4)
		renderer.Copy(tex, nil, nil)
		renderer.Present()

		elapsedTime = float32(time.Since(frameStart).Seconds())
		if elapsedTime < .005 {
			sdl.Delay(5 - uint32(elapsedTime/1000.0))
			elapsedTime = float32(time.Since(frameStart).Seconds())
		}
	}
}
