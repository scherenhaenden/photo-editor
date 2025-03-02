package ui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
)

// ----------------------------------------------------------------------
// ÁREA CENTRAL: VISOR DE IMÁGENES
// ----------------------------------------------------------------------
// CreateImageArea genera el área central donde se mostrará la imagen.
func CreateImageArea() fyne.CanvasObject {
	// Fondo gris oscuro simulando el visor de imágenes.
	imageCanvas := canvas.NewRectangle(color.RGBA{R: 60, G: 60, B: 60, A: 255})

	// Texto placeholder para indicar el visor de imágenes.
	textOverlay := canvas.NewText("Área Central - Visor de Imágenes (Placeholder)", color.White)
	textOverlay.Alignment = fyne.TextAlignCenter
	textOverlay.TextStyle = fyne.TextStyle{Bold: true}

	// Crear un contenedor que llena todo el espacio y coloca el texto en el centro.
	imageViewContainer := container.NewStack(imageCanvas, container.NewCenter(textOverlay))

	return imageViewContainer
}
