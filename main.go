package main

import (
	"bytes"
	"image"
	"image/color"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/disintegration/imaging"
	"github.com/llgcode/draw2d/draw2dimg"
	"github.com/tobijes/epaper-service/electricity"
	"golang.org/x/image/bmp"
)

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func main() {
	var img image.Image = electricity.Generate()
	// Save to file initially
	draw2dimg.SaveToPngFile("test.png", img)

	http.HandleFunc("/electricity", handleElectricity)

	port := getEnv("PORT", "8080")
	log.Printf("Listening on %s \n", port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

func handleElectricity(w http.ResponseWriter, r *http.Request) {
	var img image.Image = electricity.Generate()
	rotate := r.URL.Query().Get("rotate")
	if rotate == "true" {
		img = imaging.Rotate(img, 90, color.Black)
	}
	writeImage(w, &img)
}

// writeImage encodes an image 'img' in jpeg format and writes it into ResponseWriter.
func writeImage(w http.ResponseWriter, img *image.Image) {

	buffer := new(bytes.Buffer)
	if err := bmp.Encode(buffer, *img); err != nil {
		log.Println("unable to encode image.")
	}

	w.Header().Set("Content-Type", "image/bmp")
	w.Header().Set("Content-Length", strconv.Itoa(len(buffer.Bytes())))
	if _, err := w.Write(buffer.Bytes()); err != nil {
		log.Println("unable to write image.")
	}
}
