package helpers

import (
	"bytes"
	"errors"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"

	_ "golang.org/x/image/webp"

	"github.com/chai2010/webp"
)

func CheckSize(r io.Reader, maxKB int) (bool, []byte, error) {
	var buf bytes.Buffer
	size, err := io.Copy(&buf, r)
	if err != nil {
		return false, nil, err
	}
	return size <= int64(maxKB*1024), buf.Bytes(), nil
}

func CompressToWebP(r io.Reader, maxKB int) ([]byte, error) {
	img, _, err := image.Decode(r)
	if err != nil {
		return nil, err
	}

	var (
		buf     bytes.Buffer
		quality = 90
	)

	for quality >= 40 {
		buf.Reset()
		err := webp.Encode(&buf, img, &webp.Options{Quality: float32(quality)})
		if err != nil {
			return nil, err
		}
		if buf.Len() <= maxKB*1024 {
			return buf.Bytes(), nil
		}
		quality -= 5
	}

	return nil, errors.New("cannot compress image under target size")
}
