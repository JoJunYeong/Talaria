package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	//	"github.com/go-vgo/robotgo"

	"github.com/stianeikeland/go-rpio"
)

var micedx float64 = 0
var micedy float64 = 0

var mouse0dx float64 = 0
var mouse0dy float64 = 0
var mouse1dx float64 = 0
var mouse1dy float64 = 0
var mouse2dx float64 = 0
var mouse2dy float64 = 0
var mouse3dx float64 = 0
var mouse3dy float64 = 0
var mouse0x float64 = 0
var mouse0y float64 = 0
var mouse1x float64 = 0
var mouse1y float64 = 0
var mouse2x float64 = 0
var mouse2y float64 = 0
var mouse3x float64 = 0
var mouse3y float64 = 0
var mouse0dx_matrix [5]float64
var mouse0dy_matrix [5]float64
var mouse1dx_matrix [5]float64
var mouse1dy_matrix [5]float64
var mouse2dx_matrix [5]float64
var mouse2dy_matrix [5]float64
var mouse3dx_matrix [5]float64
var mouse3dy_matrix [5]float64
var velocity [8]int64
var i int64
var name string
var mouse_to_sensor string
var left_count float64 = 0
var right_count float64 = 0
var one_sec_DX float64 = 0
var one_sec_DY float64 = 0
var time1 = time.Now().UnixNano()
var time2 = time.Now().UnixNano()
var avg_speed_source_x [1000]float64
var avg_speed_source_y [1000]float64
var avg_speed_source_number_count int64 = 0
var avg_speed_x float64 = 0
var avg_speed_y float64 = 0

//var 0dx_mat_count int64 = 0
//var 0dy_mat_count int64 = 0
//var 1dx_mat_count int64 = 0
//var 1dy_mat_count int64 = 0
//var 2dx_mat_count int64 = 0
//var 2dy_mat_count int64 = 0
//var 3dx_mat_count int64 = 0
//var 3dy_mat_count int64 = 0

var (
	// Use mcu pin 10, correspond to physical pin 19 on the pi
	pin2  = rpio.Pin(2)  // the first pin of direction
	pin3  = rpio.Pin(3)  // the second pin of direction
	pin4  = rpio.Pin(4)  // the first pin of velocity
	pin5  = rpio.Pin(5)  // the second pin of velocity
	pin6  = rpio.Pin(6)  // the third pin of velocity
	pin7  = rpio.Pin(7)  // the next pin of velocity
	pin8  = rpio.Pin(8)  // the next pin of velocity
	pin9  = rpio.Pin(9)  // the next pin of velocity
	pin10 = rpio.Pin(10) // the next pin of velocity
	pin11 = rpio.Pin(11) // the last pin of velocity
	pin12 = rpio.Pin(12) // the signal pin of send all velocity data
	pin13 = rpio.Pin(13) // the signal pin of ready to go
	pin22 = rpio.Pin(22) // the signal pin of received all velocity data
	pin23 = rpio.Pin(23) // the reset signal

)

type binData struct {
	S  uint8
	DX int8
	DY int8
}

type Mouse struct {
	S     uint8
	DX    float64
	DY    float64
	X     float64
	Y     float64
	Left  bool
	Right bool
	Name  string
}

func (m Mouse) String() string {

	time2 = time.Now().UnixNano()
	if strings.Contains(m.Name, "ice") { //

		mouse_to_sensor = "dev/input/sensor"

	} //if
	if strings.Contains(m.Name, "se0") { //

		mouse_to_sensor = "dev/input/sensor0"

	} //if
	if strings.Contains(m.Name, "se1") { //

		mouse_to_sensor = "dev/input/sensor1"

	} //if
	if strings.Contains(m.Name, "se2") { //

		mouse_to_sensor = "dev/input/sensor2"

	} //if
	if strings.Contains(m.Name, "se3") { //

		mouse_to_sensor = "dev/input/sensor3"
	} //if
	return fmt.Sprintf("S:%08b, DX:%4v, DY:%4v, X:%6v, Y:%6v, ONE_DX:%6v, ONE_DY:%6v, avg_X:%6v, avg_Y:%6v, nanosecond:%6v, Name:%20v", m.S, m.DX, m.DY, m.X, m.Y, one_sec_DX, one_sec_DY, avg_speed_x, avg_speed_y, time2, mouse_to_sensor)

} //func

func Follow(rc chan Mouse, dev string) { //
	fp, err := os.Open(dev)
	m := Mouse{Name: dev}

	if err != nil {
		panic(err)
	} //if

	defer fp.Close()

	for { //

		thing := binData{}
		err := binary.Read(fp, binary.LittleEndian, &thing)

		if err == io.EOF { //
			break
		} //if
		m.DX = float64(thing.DX)
		m.DY = float64(thing.DY)
		m.X += m.DX
		m.Y += m.DY
		m.S = thing.S
		rc <- m

		if strings.Contains(m.Name, "mice") { //

			micedx = m.DX
			micedy = m.DY

		} //if
		if strings.Contains(m.Name, "se0") { //

			mouse0dx = m.DX
			mouse0dy = m.DY
			mouse0x = m.X
			mouse0y = m.Y

			for i := 0; i <= 3; i++ { //
				mouse0dx_matrix[i+1] = mouse0dx_matrix[i]
				mouse0dy_matrix[i+1] = mouse0dy_matrix[i]
			} //for
			mouse0dx_matrix[0] = mouse0dx
			mouse0dy_matrix[0] = mouse0dy

		} //if
		if strings.Contains(m.Name, "se1") { //

			mouse1dx = m.DX
			mouse1dy = m.DY
			mouse1x = m.X
			mouse1y = m.Y

			for i := 0; i <= 3; i++ { //
				mouse1dx_matrix[i+1] = mouse1dx_matrix[i]
				mouse1dy_matrix[i+1] = mouse1dy_matrix[i]
			} //for
			mouse1dx_matrix[0] = mouse1dx
			mouse1dy_matrix[0] = mouse1dy

		} //if
		if strings.Contains(m.Name, "se2") { //

			mouse2dx = m.DX
			mouse2dy = m.DY
			mouse2x = m.X
			mouse2y = m.Y
			for i := 0; i <= 3; i++ { //
				mouse2dx_matrix[i+1] = mouse2dx_matrix[i]
				mouse2dy_matrix[i+1] = mouse2dy_matrix[i]
			} //for
			mouse2dx_matrix[0] = mouse2dx
			mouse2dy_matrix[0] = mouse2dy

		} //if
		if strings.Contains(m.Name, "se3") { //

			mouse3dx = m.DX
			mouse3dy = m.DY
			mouse3x = m.X
			mouse3y = m.Y
			for i := 0; i <= 3; i++ { //
				mouse3dx_matrix[i+1] = mouse3dx_matrix[i]
				mouse3dy_matrix[i+1] = mouse3dy_matrix[i]
			} //for
			mouse3dx_matrix[0] = mouse3dx
			mouse3dy_matrix[0] = mouse3dy

		} //if

		///////////////////////////////// time check & appending to slice for calculate average
		/*
			time2 = time.Now().UnixNano()
			if (time2 - time1) > 1000000 {

				time1 = time2
				if strings.Contains(m.Name, "se0") { //
					one_sec_DY = mouse0y
					one_sec_DX = mouse0x
				} else if strings.Contains(m.Name, "se1") { //
					one_sec_DY = mouse1y
					one_sec_DX = mouse1x
				} else if strings.Contains(m.Name, "se2") { //
					one_sec_DY = mouse2y
					one_sec_DX = mouse2x
				} else if strings.Contains(m.Name, "se3") { //
					one_sec_DY = mouse3y
					one_sec_DX = mouse3x
				} //if

				m.Y = 0
				m.X = 0 //
				avg_speed_source_number_count++
				avg_speed_source_x[avg_speed_source_number_count] = mouse2dx
				avg_speed_source_y[avg_speed_source_number_count] = mouse2dy
				if avg_speed_source_number_count >= 100 {

					sum_x := 0.0
					sum_y := 0.0
					for i := 0; i < int(avg_speed_source_number_count); i++ {
						sum_x += avg_speed_source_x[avg_speed_source_number_count]
						sum_y += avg_speed_source_y[avg_speed_source_number_count]
					}
					avg_speed_x = sum_x / float64(avg_speed_source_number_count)
					avg_speed_y = sum_y / float64(avg_speed_source_number_count)
					avg_speed_source_number_count = 0
				}
			}
		*/
		///////////////////////////////////////////

	} //for

} //func

func main() { //
	// Open and map memory to access gpio, check for errors
	if err := rpio.Open(); err != nil { //
		fmt.Println(err)
		os.Exit(1)
	} //if

	// Unmap gpio memory when done
	defer rpio.Close()

	// Set pin to output mode
	pin2.Output()
	pin3.Output()
	pin4.Output()
	pin5.Output()
	pin6.Output()
	pin7.Output()
	pin8.Output()
	pin9.Output()
	pin10.Output()
	pin11.Output()
	pin12.Output()
	pin13.Input()
	pin22.Input()
	pin23.Output()

	pin2.Low()
	pin3.Low()
	pin4.Low()
	pin5.Low()
	pin6.Low()
	pin7.Low()
	pin8.Low()
	pin9.Low()
	pin10.Low()
	pin11.Low()
	pin12.High()
	pin23.Low()
	//	pin13.Low()
	lastStatePin13 := pin13.Read()
	lastStatePin22 := pin22.Read()
	//	currentStatePin13 := pin13.Read()
	//	currentStatePin13 := pin13.High()
	pin12.Low()
	ch := make(chan Mouse)

	ms, _ := filepath.Glob("/dev/input/mouse*")
	go Follow(ch, "/dev/input/mice")

	for _, path := range ms { //
		go Follow(ch, path)
	} //for

	all := make(map[string]Mouse)
	var keys []string
	var ok bool

	for m := range ch { //

		fmt.Printf("\r\033[%vA", len(all))

		if _, ok = all[m.Name]; !ok { //
			keys = append(keys, m.Name)
			sort.Strings(keys)
		} //if

		all[m.Name] = m

		for _, key := range keys {
			val, _ := all[key]
			fmt.Printf("\n-%25v  %v    ", key, val)
			pin23.Low()
			i = 0
			for i < 8 {
				velocity[i] = 0
				i++
			}
			var count2 int64

			//if ((mouse0dx <= -10 && mouse0dy <= 10 && mouse0dy >= -10 ) && (mouse1dx <= -10 && mouse1dy <= 10 && mouse1dy >= -10)) || ((mouse2dx <= -10 && mouse2dy <= 10 && mouse2dy >= -10) && (mouse3dx <= -10 && mouse3dy <= 10 && mouse3dy >= -10))  {		// Right conditional statement
			////	robotgo.TypeStr("d")	// Right keyboard macro
			//pin5.High()
			//time.Sleep(time.Second / 120)
			//pin5.Low()
			//} else if ((mouse0dx >= 10 && mouse0dy <= 10 && mouse0dy >= -10) && (mouse1dx >= 10 && mouse1dy <=10 && mouse1dy >= -10)) || ((mouse2dx >= 10 && mouse2dy <= 10 && mouse2dy >= -10) && (mouse3dx >= 10 && mouse3dy <= 10 && mouse3dy >= -10)) {		// Left conditional statement
			////	robotgo.TypeStr("a")	//Left keyboard macro
			//pin3.High()
			//time.Sleep(time.Second / 120)
			//pin3.Low()
			//} else

			//if ((mouse0dy >= 10 ) || (mouse1dy >= 10 ) || (mouse2dy >= 10 ) || (mouse3dy >= 10 )) {//		// Down conditional statement
			//	robotgo.TypeStr("s")	// Down keyboard macro
			//pin2.Low()
			//pin4.High()
			//time.Sleep(time.Second / 240)
			//pin3.Low()

			//} else

			if micedy <= -50 { //	// Up conditional statement
				name = m.Name

				lastStatePin13 = pin13.Read()
				lastStatePin22 = pin22.Read()
				if lastStatePin13 == 0 {
					//	fmt.Printf("\n\nlastStatePin13 is 0  ")
				}
				if lastStatePin13 == 1 {
					//	fmt.Printf("\n\nlastStatePin13 is 1  ")
					pin12.High()

					lastStatePin13 = pin13.Read()
					for lastStatePin13 == 1 {
						count2++
						lastStatePin13 = pin13.Read()
						if count2%10000000 == 0 && count2 != 0 {
							fmt.Printf("\nlastStatePin13 is still 1  ")
						}
					}
					pin12.Low()

					if strings.Contains(name, "ice") { //
						if micedy <= -255 { // if dy<-255, then rearrange dy to -255
							micedy = -255
						}
						micedy = math.Abs(micedy)
						//fmt.Printf("mouse0dy : %4v",mouse0dy)

						i = 0

						var dy_int int64

						dy_int = int64(micedy)

						for i < 8 {
							velocity[7-i] = dy_int % 2
							dy_int = dy_int / 2
							i++
						}

						i = 0
						for i < 8 {
							fmt.Print(velocity[i])
							i++
						}

						//	fmt.Printf("\n0dy Data: %4v", mouse0dy)

					} //if

					if strings.Contains(name, "se0") { //

						name = "mouse0"

						/*
							if mouse0dy <= -255 { // if dy<-255, then rearrange dy to -255
								mouse0dy = -255
							}
							mouse0dy = math.Abs(mouse0dy)
							//fmt.Printf("mouse0dy : %4v",mouse0dy)

							i = 0

							var dy_int int64

							dy_int = int64(mouse0dy)

							for i < 8 {
								velocity[7-i] = dy_int % 2
								dy_int = dy_int / 2
								i++
							}
						*/

						/*
							i = 0
							for i < 8 {
								fmt.Print(velocity[i])
								i++
							}
						*/

						//	fmt.Printf("\n0dy Data: %4v", mouse0dy)

					} //if
					if strings.Contains(name, "se1") { //
						right_count = 0
						name = "mouse1"
						left_count++
						//	fmt.Printf("left_count=   %4v", left_count)
						/*
							if mouse1dy <= -255 { // if dy<-255, then rearrange dy to 255
								mouse1dy = -255
							}
							mouse1dy = math.Abs(mouse1dy)
							//fmt.Printf("mouse1dy : %4v",mouse1dy)

							i = 0

							var dy_int int64

							dy_int = int64(mouse1dy)

							for i < 8 {

								velocity[7-i] = dy_int % 2

								dy_int = dy_int / 2
								i++
							}
						*/
						/*
							i = 0
							for i < 8 {
								fmt.Print(velocity[i])
								i++
							}
						*/
						//	fmt.Printf("\n1dy Data: %4v", mouse1dy)

					} //if
					if strings.Contains(name, "se2") { //
						left_count = 0
						name = "mouse2"
						right_count++
						//	fmt.Printf("right_count=   %4v", right_count)
						/*
							if mouse2dy <= -255 { // if dy<-255, then rearrange dy to 255
								mouse2dy = -255
							}
							mouse2dy = math.Abs(mouse2dy)
							//fmt.Printf("mouse2dy : %4v",mouse2dy)

							i = 0

							var dy_int int64
							mouse2dy = math.Abs(mouse2dy)
							dy_int = int64(mouse2dy)

							for i < 8 {

								velocity[7-i] = dy_int % 2

								dy_int = dy_int / 2
								i++
							}
						*/
						/*
							i = 0
							for i < 8 {
								fmt.Print(velocity[i])
								i++
							}
						*/
						//	fmt.Printf("\n2dy Data: %4v", mouse2dy)

					} //if
					if strings.Contains(name, "se3") { //
						name = "mouse3"

						/*
							if mouse3dy <= -255 { // if dy<-255, then rearrange dy to 255
								mouse3dy = -255
							}
							mouse3dy = math.Abs(mouse3dy)
							//fmt.Printf("mouse3dy : %4v",mouse3dy)

							i = 0

							var dy_int int64

							dy_int = int64(mouse3dy)

							for i < 8 {

								velocity[7-i] = dy_int % 2

								dy_int = dy_int / 2
								i++
							}

							i = 0
							for i < 8 {
								fmt.Print(velocity[i])
								i++
							}
						*/
						//	fmt.Printf("\n3dy Data: %4v", mouse3dy)

					} //if

					//	fmt.Printf("another enter")

					//	robotgo.TypeStr("w")	// Up keyboard macro
					//pin4.Low()
					//pin2.High()

					lastStatePin13 = pin13.Read()
					for lastStatePin13 == 1 {
						lastStatePin13 = pin13.Read()
					}

					pin2.Low()
					pin3.Low()
					pin4.Low()
					pin5.Low()
					pin6.Low()
					pin7.Low()
					pin8.Low()
					pin9.Low()
					pin10.Low()
					pin11.Low()

					//fmt.Printf("\nData send start")
					//pin4.High()
					if velocity[7] == 1 {
						pin4.High()
					}
					if velocity[6] == 1 {
						pin5.High()
					}
					if velocity[5] == 1 {
						pin6.High()
					}
					if velocity[4] == 1 {
						pin7.High()
					}
					if velocity[3] == 1 {
						pin8.High()
					}
					if velocity[2] == 1 {
						pin9.High()
					}
					if velocity[1] == 1 {
						pin10.High()
					}
					if velocity[0] == 1 {
						pin11.High()
					}
					pin12.High() // signal pin of send all data.
					//fmt.Printf("\nData send complete.")

					//pin1.Low()
					var count int64 = 0
					var count3 int64 = 0
					lastStatePin22 = pin22.Read()
					for lastStatePin22 == 0 {
						count++
						lastStatePin22 = pin22.Read()
						if count%1000 == 0 && count != 0 {

							//	fmt.Printf("\nlastStatePin22==0 count=%v", count)
							count3++
							if count3%1 == 0 {
								pin12.Low()
								pin23.High()

							}
						}

					}
					pin12.Low()
					//////////// signal reset ////////////

					pin2.Low()
					pin3.Low()
					pin4.Low()
					pin5.Low()
					pin6.Low()
					pin7.Low()
					pin8.Low()
					pin9.Low()
					pin10.Low()
					pin11.Low()
					pin12.Low()

					//////////////////////////////////////

				} //if
			} else { //if (lastStatePin13==1)

				//	fmt.Printf("else lastStatePin13==0")

			}
			if ((mouse0dx >= -45 && mouse0dx <= 45 && mouse0dy <= 45 && mouse0dy >= -45) && (mouse1dx >= -45 && mouse1dx <= 45 && mouse1dy <= 45 && mouse1dy >= -45)) && ((mouse2dx >= -45 && mouse2dx <= 45 && mouse2dy <= 45 && mouse2dy >= -45) && (mouse3dx >= -45 && mouse3dx <= 45 && mouse3dy <= 45 && mouse3dy >= -45)) { //
				pin2.Low()
				//	fmt.Printf("enter")
				//pin4.Low()
			} //if

		} // for
	} //for
} //main
