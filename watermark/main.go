package main

import (
	"errors"
	"fmt"
	"github.com/nfnt/resize"
	"image"
	"image/draw"
	_ "image/gif"
	"image/jpeg"
	_ "image/png"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	_, err := AddWatermark("./img.jpg", "./watermark.png")
	if err != nil {
		panic(err)
	}
}

// AddWatermark 给图片添加水印 参数和返回值为图片路径
func AddWatermark(imgPath, waterPath string) (dst string, err error) {
	// 1.读取原始图片
	var file *os.File
	file, err = os.Open(imgPath)
	if err != nil {
		panic(err)
	}
	filename := file.Name()
	filename = strings.TrimSuffix(filepath.Base(filename), filepath.Ext(filename))

	if err != nil {
		return "", err
	}
	defer file.Close()

	// 2.将原始图片解码为 image.Image 格式
	var img image.Image
	img, err = file2Image(file)
	if err != nil {
		return "", err
	}

	// 3.读取水印图片
	watermarkFile, err := os.Open(waterPath)
	if err != nil {
		return "", err
	}
	defer watermarkFile.Close()

	// 4.将水印图片解码为 image.Image 格式
	var watermarkImg image.Image
	watermarkImg, err = file2Image(watermarkFile)
	if err != nil {
		return "", err
	}

	// 5.计算水印图像的缩放比例，使其适合原始图像的大小
	scale := float64(img.Bounds().Dx()) / float64(watermarkImg.Bounds().Dx())
	watermarkWidth := int(float64(watermarkImg.Bounds().Dx()) * scale / 8)
	watermarkHeight := int(float64(watermarkImg.Bounds().Dy()) * scale / 8)

	// 6.缩放水印图像
	watermarkImg = resize.Resize(uint(watermarkWidth), uint(watermarkHeight), watermarkImg, resize.Lanczos3)

	// 7.在原始图像上添加水印
	// 水印距离边界距离
	distance := img.Bounds().Dx() / 100
	offset := image.Pt(img.Bounds().Dx()-watermarkWidth-distance, distance)
	b := img.Bounds()
	m := image.NewRGBA(b)
	draw.Draw(m, b, img, image.Point{}, draw.Src)
	draw.Draw(m, watermarkImg.Bounds().Add(offset), watermarkImg, image.Point{}, draw.Over)

	// 将添加了水印的图片保存到文件
	imgWithWatermark := fmt.Sprintf("./%s-%s.jpg", filename, "watermark")
	var output *os.File
	output, err = os.Create(imgWithWatermark)
	if err != nil {
		panic(err)
	}
	defer output.Close()

	jpeg.Encode(output, m, &jpeg.Options{Quality: jpeg.DefaultQuality})
	return imgWithWatermark, nil
}

// 将 文件转换成 image.Image 格式
func file2Image(f *os.File) (img image.Image, err error) {
	img, _, err = image.Decode(f)
	if err != nil {
		return nil, errors.New("图片解码出错")
	}
	return
}

type ImgType int

const (
	JPG ImgType = iota + 1
	PNG
	BMP
	GIF
)

// typeOfImage 通过文件头部信息判断图片类型
func typeOfImage(f *os.File) (t ImgType, err error) {
	var data = make([]byte, 5)
	var n int
	n, err = f.Read(data)
	if n > 2 && data[0] == 0xff && data[1] == 0xd8 {
		return JPG, nil
	}
	if n > 3 && data[0] == 0x89 && data[1] == 0x50 && data[2] == 0x4e && data[3] == 0x47 {
		return PNG, nil
	}
	if n > 2 && data[0] == 0x42 && data[1] == 0x4d {
		return BMP, nil
	}
	if n > 3 && data[0] == 0x47 && data[1] == 0x49 && data[2] == 0x46 && data[3] == 0x38 {
		return GIF, nil
	}
	return 0, errors.New("无法判断的类型")
}
