package ui

import (
	"image"
	"image/color"
	"io"

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

// ----------------------------------------------------------------------
// SetImage updates the displayed image dynamically.
// If imgStream is nil, the placeholder text is shown.
// ----------------------------------------------------------------------
func (ia *ImageArea) SetImage(imgStream io.Reader) {
	if imgStream == nil {
		ia.textOverlay.Show()
		ia.imageDisplay.Hide()
	} else {
		// Decode the image from the stream
		img, _, err := image.Decode(imgStream)
		if err != nil {
			ia.textOverlay.Text = "Error loading image"
			ia.textOverlay.Refresh()
			return
		}
		// Update the image display
		ia.imageDisplay.Image = img
		ia.imageDisplay.Refresh()
		ia.imageDisplay.Show()
		ia.textOverlay.Hide()
	}
}
