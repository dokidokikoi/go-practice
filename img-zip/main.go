package main

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"os"

	"github.com/h2non/filetype"
)

func main() {
	file, err := os.Open("wallhaven-28kdom.png")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		panic(err)
	}
	kind, _ := filetype.Match(data)

	target, err := os.Create("wallhaven-28kdom(1).png")
	if err != nil {
		panic(err)
	}
	defer target.Close()

	switch kind.MIME.Value {
	case "image/png":
		err = PngZip(40, bytes.NewBuffer(data), target)
		if err != nil {
			panic(err)
		}
	case "image/jpeg":
		err = JpgZip(40, bytes.NewBuffer(data), target)
		if err != nil {
			panic(err)
		}
	case "image/gif":
		err = GifZip(bytes.NewBuffer(data), target)
		if err != nil {
			panic(err)
		}
	}
}

func PngZip(quality int, file io.Reader, target io.Writer) error {
	// 解码图片
	img, err := png.Decode(file)
	if err != nil {
		return err
	}

	// 压缩图片
	return jpeg.Encode(target, img, &jpeg.Options{Quality: quality})
}

func JpgZip(quality int, file io.Reader, target io.Writer) error {
	// 解码图片
	img, _, err := image.Decode(file)
	if err != nil {
		return err
	}

	// 压缩图片
	return jpeg.Encode(target, img, &jpeg.Options{Quality: quality})
}

func GifZip(f io.Reader, target io.Writer) error {
	// 解码GIF文件
	gifImg, err := gif.DecodeAll(f)
	if err != nil {
		return err
	}

	// 设置压缩质量参数
	newPalette := []color.Color{color.White, color.Black}
	newGIF := &gif.GIF{}
	newGIF.Delay = gifImg.Delay
	newGIF.Image = make([]*image.Paletted, len(gifImg.Image))
	for i, img := range gifImg.Image {
		rect := img.Bounds()
		newGIF.Image[i] = image.NewPaletted(rect, newPalette)
		draw.Draw(newGIF.Image[i], rect, img, rect.Min, draw.Src)
	}

	// 将压缩后的GIF写入文件
	return gif.EncodeAll(target, newGIF)
}
