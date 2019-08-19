package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"html/template"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"io/ioutil"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

func main() {
	mux := http.NewServeMux()
	files := http.FileServer(http.Dir("public"))
	mux.Handle("/static/", http.StripPrefix("/static/", files))
	mux.HandleFunc("/", upload)
	mux.HandleFunc("/mosaic", mosaic)
	server := &http.Server{
		Addr:    "127.0.0.1:8080",
		Handler: mux,
	}

	TILESDB = tilesDB(500)
	fmt.Println("Mosaic server started.")
	server.ListenAndServe()
}

func upload(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("public/upload.html")
	t.Execute(w, nil)
}

func mosaic(w http.ResponseWriter, r *http.Request) {
	t0 := time.Now()

	r.ParseMultipartForm(4096)
	file, _, err := r.FormFile("image")
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	defer file.Close()
	tileSize, _ := strconv.Atoi(r.FormValue("tile_size"))

	original, _, err := image.Decode(file)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	bounds := original.Bounds()
	newImage := image.NewNRGBA(image.Rect(bounds.Min.X, bounds.Min.Y, bounds.Max.X, bounds.Max.Y))

	db := cloneTilesDB()
	sp := image.Point{0, 0}
	for y := bounds.Min.Y; y < bounds.Max.Y; y += tileSize {
		for x := bounds.Min.X; x < bounds.Max.X; x += tileSize {
			r, g, b, _ := original.At(x, y).RGBA()
			color := [3]float64{float64(r), float64(g), float64(b)}
			nearest := nearest(color, &db)
			file, err := os.Open(nearest)
			if err != nil {
				fmt.Println("open file error:", err, nearest)
				continue
			}

			img, _, err := image.Decode(file)
			if err != nil {
				fmt.Println("decode image error:", err)
				continue
			}

			t := resize(img, tileSize)
			tile := t.SubImage(t.Bounds())
			tileBounds := image.Rect(x, y, x+tileSize, y+tileSize)
			draw.Draw(newImage, tileBounds, tile, sp, draw.Src)
			file.Close()
		}
	}

	buf1 := new(bytes.Buffer)
	jpeg.Encode(buf1, original, nil)
	originalStr := base64.StdEncoding.EncodeToString(buf1.Bytes())

	buf2 := new(bytes.Buffer)
	jpeg.Encode(buf2, newImage, nil)
	mosaic := base64.StdEncoding.EncodeToString(buf2.Bytes())

	images := map[string]string{
		"original": originalStr,
		"mosaic":   mosaic,
		"duration": time.Since(t0).String(),
	}
	t, _ := template.ParseFiles("public/result.html")
	t.Execute(w, images)
}

func averageColor(img image.Image) [3]float64 {
	bounds := img.Bounds()
	var r, g, b float64
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r1, g1, b1, _ := img.At(x, y).RGBA()
			r, g, b = r+float64(r1), g+float64(g1), b+float64(b1)
		}
	}

	totalPixels := float64(bounds.Max.X * bounds.Max.Y)
	return [3]float64{r / totalPixels, g / totalPixels, b / totalPixels}
}

func resize(in image.Image, newWidth int) image.NRGBA {
	bounds := in.Bounds()
	ratio := bounds.Dx() / newWidth
	out := image.NewNRGBA(image.Rect(bounds.Min.X/ratio, bounds.Min.Y/ratio, bounds.Max.X/ratio, bounds.Max.Y/ratio))
	for y, j := bounds.Min.Y, bounds.Min.Y; y < bounds.Max.Y; y, j = y+ratio, j+1 {
		for x, i := bounds.Min.X, bounds.Min.X; x < bounds.Max.X; x, i = x+ratio, i+1 {
			r, g, b, a := in.At(x, y).RGBA()
			out.SetNRGBA(i, j, color.NRGBA{R: uint8(r >> 8), G: uint8(g >> 8), B: uint8(b >> 8), A: uint8(a >> 8)})
		}
	}
	return *out
}

func tilesDB(ct int) map[string][3]float64 {
	fmt.Println("Start populating tiles db ...")
	db := make(map[string][3]float64)
	files, _ := ioutil.ReadDir(downloadDir)
	if len(files) < ct {
		var wg sync.WaitGroup
		for ct > 0 {
			wg.Add(1)
			go func() {
				download(100)
				wg.Done()
			}()
			ct -= 100
		}
		wg.Wait()
		files, _ = ioutil.ReadDir(downloadDir)
	}
	for _, f := range files {
		name := filepath.Join(downloadDir, f.Name())
		file, err := os.Open(name)
		if err != nil {
			fmt.Println("cannot open file", name, err)
			continue
		}

		img, _, err := image.Decode(file)
		if err != nil {
			fmt.Println("error in populating TILESB:", err, name)
			continue
		}

		db[name] = averageColor(img)
		file.Close()
	}

	fmt.Println("Finished populating tiles db.")
	return db
}

func nearest(target [3]float64, db *map[string][3]float64) string {
	var filename string
	smallest := 1000000.0

	for k, v := range *db {
		dist := distance(target, v)
		if dist < smallest {
			smallest, filename = dist, k
		}
	}
	delete(*db, filename)
	return filename
}

func distance(a, b [3]float64) float64 {
	sq := func(x float64) float64 { return x * x }
	return math.Sqrt(sq(a[0]-b[0]) + sq(a[1]-b[1]) + sq(a[2]-b[2]))
}

var TILESDB map[string][3]float64

func cloneTilesDB() map[string][3]float64 {
	db := make(map[string][3]float64)
	for k, v := range TILESDB {
		db[k] = v
	}
	return db
}

var API = "https://api.thecatapi.com/v1/images/search?limit=%d"

var downloadDir = "tiles"

func download(ct int) {
	resp, err := http.Get(fmt.Sprintf(API, ct))
	if err != nil {
		panic(err)
	}

	type result struct {
		URL string `json:"url"`
	}
	var results []result
	b, _ := ioutil.ReadAll(resp.Body)
	if err := json.Unmarshal(b, &results); err != nil {
		panic(err)
	}

	var wg sync.WaitGroup
	wg.Add(len(results))
	for i := range results {
		url := results[i].URL
		filename := strings.ToLower(url[strings.LastIndex(url, "/")+1:])

		go func(filename string) {
			defer wg.Done()

			if !strings.HasSuffix(filename, ".jpg") {
				return
			}

			resp, err := http.Get(url)
			if err != nil {
				return
			}
			b, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return
			}

			f, err := os.Create(filepath.Join(downloadDir, filename))
			if err != nil {
				return
			}
			defer f.Close()
			f.Write(b)

		}(filename)
	}
	wg.Wait()
}
