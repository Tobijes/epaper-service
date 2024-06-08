package main

import (
	"bytes"
	"image"
	"image/png"
	"log"
	"net/http"
	"strconv"

	"github.com/llgcode/draw2d/draw2dimg"
	"github.com/tobijes/epaper-service/electricity"
)

func main() {
	var img image.Image = electricity.Generate()
	// Save to file initially
	draw2dimg.SaveToPngFile("test.png", img)

	http.HandleFunc("/electricity", handleElectricity)

	log.Println("Listening on 8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

func handleElectricity(w http.ResponseWriter, r *http.Request) {
	var img image.Image = electricity.Generate()
	writeImage(w, &img)
}

// writeImage encodes an image 'img' in jpeg format and writes it into ResponseWriter.
func writeImage(w http.ResponseWriter, img *image.Image) {

	buffer := new(bytes.Buffer)
	if err := png.Encode(buffer, *img); err != nil {
		log.Println("unable to encode image.")
	}

	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Content-Length", strconv.Itoa(len(buffer.Bytes())))
	if _, err := w.Write(buffer.Bytes()); err != nil {
		log.Println("unable to write image.")
	}
}
