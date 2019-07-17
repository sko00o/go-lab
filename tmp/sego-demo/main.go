package main

import (
	"fmt"

	"github.com/huichen/sego"
)

func main() {
	var segmenter sego.Segmenter
	segmenter.LoadDictionary("/mnt/d/GoDev/src/github.com/huichen/sego/data/dictionary.txt")

	text := []byte("中华人民共和国是中央人民政府")
	segments := segmenter.Segment(text)

	fmt.Println(sego.SegmentsToString(segments, false))
}
