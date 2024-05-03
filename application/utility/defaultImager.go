package utility

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"os"
)

func GetDefaultImageBytes(path string, format string) ([]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}

	var imgBytes []byte
	imgBytes, err = encodeToBytes(img, format)
	if err != nil {
		return nil, err
	}

	return imgBytes, nil
}

func encodeToBytes(img image.Image, format string) ([]byte, error) {
	buf := new(bytes.Buffer)
	var err error

	switch format {
	case "jpeg":
		err = jpeg.Encode(buf, img, nil)
	case "png":
		err = png.Encode(buf, img)
	default:
		return nil, fmt.Errorf("unsupported image format: %s", format)
	}

	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
