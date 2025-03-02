package ui

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"io"
	"io/ioutil"

	_ "golang.org/x/image/bmp"
	_ "golang.org/x/image/tiff"
	_ "golang.org/x/image/webp"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
)

// ----------------------------------------------------------------------
// CONSTANTS FOR DEFAULT TEXTS (USED IF NO TRANSLATION IS PROVIDED)
// ----------------------------------------------------------------------
const DefaultImagePlaceholderText = "Área Central - Visor de Imágenes (Placeholder)"

// ----------------------------------------------------------------------
// ImageArea represents the image viewer with dynamic update support.
// ----------------------------------------------------------------------
type ImageArea struct {
	Container    *fyne.Container
	imageCanvas  *canvas.Rectangle
	imageDisplay *canvas.Image
	textOverlay  *canvas.Text
}

// ----------------------------------------------------------------------
// CreateImageArea generates the central area where images are displayed.
// This function maintains the original name.
// It accepts a translations map for custom text values.
// ----------------------------------------------------------------------
func CreateImageArea(translations map[string]string) *ImageArea {
	// Get translation for placeholder text (or use default)
	placeholderText := DefaultImagePlaceholderText
	if val, ok := translations["image_placeholder"]; ok {
		placeholderText = val
	}

	// Create the background (dark gray)
	imageCanvas := canvas.NewRectangle(color.RGBA{R: 60, G: 60, B: 60, A: 255})

	// Create an empty image display
	imageDisplay := &canvas.Image{}
	imageDisplay.FillMode = canvas.ImageFillContain

	// Create the placeholder text overlay
	textOverlay := canvas.NewText(placeholderText, color.White)
	textOverlay.Alignment = fyne.TextAlignCenter
	textOverlay.TextStyle = fyne.TextStyle{Bold: true}

	// Create the container: background, image display, and centered placeholder text
	containerObj := container.NewMax(imageCanvas, imageDisplay, container.NewCenter(textOverlay))

	return &ImageArea{
		Container:    containerObj,
		imageCanvas:  imageCanvas,
		imageDisplay: imageDisplay,
		textOverlay:  textOverlay,
	}
}

// GetImageDisplay returns the image display widget.
func (ia *ImageArea) GetImageDisplay() *canvas.Image {
	return ia.imageDisplay
}

func decodeImage(r io.Reader) (image.Image, string, error) {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, "", err
	}
	return image.Decode(bytes.NewReader(data))
}

// ----------------------------------------------------------------------
// SetImage updates the displayed image dynamically from an io.Reader.
// ----------------------------------------------------------------------
func (ia *ImageArea) SetImage(imgStream io.Reader) {
	if imgStream == nil {
		ia.textOverlay.Show()
		ia.imageDisplay.Hide()
	} else {

		data, err := ioutil.ReadAll(imgStream)
		if err != nil {
			ia.textOverlay.Text = err.Error()
			ia.textOverlay.Refresh()
			return
		}
		// **DEBUGGING CHECKS (Add these):**
		fmt.Printf("Data length: %d\n", len(data))
		if len(data) > 0 {
			fmt.Printf("First 10 bytes (hex): %x\n", data[:10]) // Show the first 10 bytes in hex
		}

		// Decode the image from the stream
		img, _, err := image.Decode(bytes.NewReader(data))
		if err != nil {
			//ia.textOverlay.Text = "Error loading image"
			ia.textOverlay.Text = err.Error()
			ia.textOverlay.Refresh()
			return
		}
		ia.imageDisplay.Image = img
		ia.imageDisplay.Refresh()
		ia.imageDisplay.Show()
		ia.textOverlay.Hide()
	}
}

// ----------------------------------------------------------------------
// SetImageFromImage updates the displayed image directly from an image.Image.
// ----------------------------------------------------------------------
func (ia *ImageArea) SetImageFromImage(img image.Image) {
	if img == nil {
		ia.textOverlay.Show()
		ia.imageDisplay.Hide()
		return
	}
	ia.imageDisplay.Image = img
	ia.imageDisplay.Refresh()
	ia.imageDisplay.Show()
	ia.textOverlay.Hide()
}
