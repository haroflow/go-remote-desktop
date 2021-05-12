package main

import (
	"encoding/gob"
	"image"
	"log"
	"net"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/kbinani/screenshot"
)

var screenImg *canvas.Image

func main() {
	myApp := app.New()
	w := myApp.NewWindow("go-remote-desktop")

	addrInput := widget.NewEntry()
	addrInput.Text = "127.0.0.1:12345"

	screenImg = canvas.NewImageFromImage(nil)

	c := container.NewBorder(
		container.NewVBox(
			widget.NewButton("Start server", func() {
				go startServer()
				w.Resize(fyne.NewSize(600, 400))
			}),
			addrInput,
			widget.NewButton("Connect to server", func() {
				addr := addrInput.Text
				go connectToServer(addr)
			}),
		),

		nil, nil, nil,

		// Center
		screenImg,
	)

	w.SetContent(c)
	w.Resize(fyne.NewSize(200, 150))
	w.ShowAndRun()
}

func startServer() {
	// f, err := os.Create("server.pprof")
	// if err != nil {
	// 	log.Println("Could not create server.pprof:", err)
	// 	return
	// }
	// defer f.Close()
	// pprof.StartCPUProfile(f)
	// defer pprof.StopCPUProfile()

	ln, err := net.Listen("tcp", "0.0.0.0:12345")
	if err != nil {
		log.Println(err)
		return
	}
	defer ln.Close()

	conn, err := ln.Accept()
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	var img *image.RGBA
	d := gob.NewDecoder(conn)
	for {
		err = d.Decode(&img)
		if err != nil {
			log.Println("Could not decode png:", err)
			break
		}

		screenImg.Image = img
		screenImg.Refresh()
	}
}

func connectToServer(addr string) {
	// f, err := os.Create("client.pprof")
	// if err != nil {
	// 	log.Println("Could not create client.pprof:", err)
	// 	return
	// }
	// defer f.Close()
	// pprof.StartCPUProfile(f)
	// defer pprof.StopCPUProfile()

	conn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	var img image.Image
	e := gob.NewEncoder(conn)
	for {
		img, err = screenshot.CaptureDisplay(0)
		if err != nil {
			log.Println(err)
			break
		}

		err = e.Encode(img)
		if err != nil {
			log.Println(err)
			break
		}

		<-time.After(1 / 60 * time.Second)
	}
}
