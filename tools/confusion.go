package tools

import (
	"image/jpeg"
	"image/png"
	"math"
	"os"

	"github.com/nfnt/resize"
)

// ThumbImage 图片剪切
func ThumbImage(path, newpath, ext string) error {
	if ext == "png" || ext == "PNG" {
		err := Imag_thumbpng(path, newpath)
		return err
	}
	if ext == "jpg" || ext == "jpeg" || ext == "JPG" || ext == "JPEG" {
		err := Imag_thumbjpg(path, newpath)
		return err
	}
	return nil
}

//jpge小图
func Imag_thumbjpg(file string, to string) error {
	// 打开图片并解码
	file_origin, err := os.Open(file)
	defer file_origin.Close()
	if err != nil {
		return err
	}
	origin, err := jpeg.Decode(file_origin)
	if err != nil {
		return err
	}

	xWight := math.Ceil(float64(origin.Bounds().Size().X) / float64(2))
	yHeight := math.Ceil(float64(origin.Bounds().Size().Y) / float64(2))

	canvas := resize.Resize(uint(xWight), uint(yHeight), origin, resize.Lanczos3)
	file_out, err := os.Create(to)
	defer file_out.Close()
	if err != nil {
		return err
	}
	err = jpeg.Encode(file_out, canvas, &jpeg.Options{80})
	return err
}

// png小图
func Imag_thumbpng(file string, to string) error {
	// 打开图片并解码
	file_origin, err := os.Open(file)
	defer file_origin.Close()
	if err != nil {
		return err
	}
	origin, err := png.Decode(file_origin)
	if err != nil {
		return err
	}
	xWight := math.Ceil(float64(origin.Bounds().Size().X) / float64(2))
	yHeight := math.Ceil(float64(origin.Bounds().Size().Y) / float64(2))
	canvas := resize.Resize(uint(xWight), uint(yHeight), origin, resize.Lanczos3)
	file_out, err := os.Create(to)
	defer file_out.Close()
	if err != nil {
		return err
	}
	err = png.Encode(file_out, canvas)
	return err
}

// 检查文件或目录是否存在
// 如果由 filename 指定的文件或目录存在则返回 true，否则返回 false
func Exist(filename string) (bool, error) {
	_, err := os.Stat(filename)
	boo := (err == nil || os.IsExist(err))
	return boo, err
}
