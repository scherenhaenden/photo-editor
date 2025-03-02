package main

import (
	_ "bytes"
	"fmt"
	"fyne.io/fyne/v2/storage"
	"image"
	"image/color"
	_ "image/color"
	_ "image/jpeg"
	"photo-editor/imaging"
	_ "photo-editor/imaging"
	"photo-editor/ui"
	"runtime"
	_ "runtime"
	"sync"
	_ "sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/davidbyttow/govips/v2/vips"
)

/*
-------------------------------------------------------------
🔹 Image Processing: OpenImageIO & VIPS Integration
-------------------------------------------------------------
📌 This program integrates OpenImageIO (OIIO) and VIPS
    for **high-performance image processing** in Go.

✅ **Why OpenImageIO?**
   - Supports **RAW image formats** (CR2, NEF, ARW, DNG, etc.).
   - Optimized for **film & VFX workflows**.
   - Multi-threaded processing with **GPU acceleration** (via OpenCL/CUDA).
   - Provides **color management (ICC profiles)**.

✅ **Why VIPS?**
   - Extremely **fast & memory-efficient** for large images.
   - Supports **deep color processing & HDR images**.
   - Built-in **image transformations (resize, rotate, crop, etc.)**.
   - Uses **multi-threading & SIMD (AVX, NEON, SSE)** for speed.

-------------------------------------------------------------
📌 Installation Requirements:
-------------------------------------------------------------
🔹 **For OpenImageIO (OIIO):**
   - Install OIIO using Homebrew (macOS) or your package manager:
     ```sh
     brew install openimageio
     ```
   - If using Linux:
     ```sh
     sudo apt install libopenimageio-dev
     ```

🔹 **For VIPS:**
   - Install VIPS with:
     ```sh
     brew install vips
     ```
   - If using Linux:
     ```sh
     sudo apt install libvips-dev
     ```

-------------------------------------------------------------
📌 Using OpenImageIO in Go:
-------------------------------------------------------------
🔹 OpenImageIO provides a C++ API, so we need Go bindings.
   - Use the `github.com/owulveryck/go-openimageio` package:
     ```sh
     go get github.com/owulveryck/go-openimageio
     ```
   - Example Usage:
     ```go
     import "github.com/owulveryck/go-openimageio"

     func processRAWImage(filename string) {
         img, err := oiio.OpenImage(filename)
         if err != nil {
             fmt.Println("Error loading image:", err)
             return
         }
         defer img.Close()

         // Apply operations...
         img.Save("output.jpg")
     }
     ```

-------------------------------------------------------------
📌 Using VIPS in Go:
-------------------------------------------------------------
🔹 VIPS has official Go bindings:
   ```sh
   go get github.com/davidbyttow/govips/v2


/*
export CGO_CFLAGS=$(pkg-config --cflags MagickWand | sed 's/-Xpreprocessor //')
export CGO_LDFLAGS=$(pkg-config --libs MagickWand)
go get -u gopkg.in/gographics/imagick.v3/imagick

*/

// AdjustBrightnessVIPSFromBuffer applies a brightness adjustment using VIPS
// and returns the new image data in the same format.
/*func AdjustBrightnessVIPSFromBuffer(buffer []byte, factor float64) ([]byte, error) {
	// Create a VIPS image from the input buffer (automatically detects format).
	vipsImg, err := vips.NewImageFromBuffer(buffer)
	if err != nil {
		return nil, fmt.Errorf("failed to create VIPS image from buffer: %w", err)
	}

	// Apply a linear brightness transformation: newPixel = factor * oldPixel + 0
	if err := vipsImg.Linear1(factor, 0); err != nil {
		return nil, fmt.Errorf("VIPS Linear1 failed: %w", err)
	}

	// Export the modified image back in the same format (JPEG, PNG, etc.).
	newBuffer, _, err := vipsImg.Export(nil) // No hardcoded format
	if err != nil {
		return nil, fmt.Errorf("failed to export VIPS image: %w", err)
	}

	return newBuffer, nil
}

// AdjustBrightnessVIPS processes an image.Image with VIPS (No PNG conversion).
func AdjustBrightnessVIPS(img image.Image, factor float64) (image.Image, error) {
	// Convert the input image to a byte buffer (use original format, not PNG).
	var buf bytes.Buffer
	err := jpeg.Encode(&buf, img, nil) // Use JPEG instead of PNG
	if err != nil {
		return nil, fmt.Errorf("failed to encode image to buffer: %w", err)
	}

	// Process the image buffer with VIPS.
	processedBuffer, err := AdjustBrightnessVIPSFromBuffer(buf.Bytes(), factor)
	if err != nil {
		return nil, fmt.Errorf("failed to adjust brightness via VIPS: %w", err)
	}

	// Decode the processed image back into an image.Image.
	processedImg, _, err := image.Decode(bytes.NewReader(processedBuffer))
	if err != nil {
		return nil, fmt.Errorf("failed to decode processed image: %w", err)
	}

	return processedImg, nil
}*/

// adjustBrightnessParallel applies a brightness factor to the image in parallel.
// A factor of 1.0 leaves the image unchanged; <1.0 darkens it; >1.0 brightens it.
func adjustBrightnessParallel(img image.Image, factor float64) image.Image {
	bounds := img.Bounds()
	newImg := image.NewRGBA(bounds)
	numWorkers := runtime.NumCPU()
	rows := bounds.Dy()
	rowsPerWorker := (rows + numWorkers - 1) / numWorkers // ceiling division

	var wg sync.WaitGroup
	for w := 0; w < numWorkers; w++ {
		startY := bounds.Min.Y + w*rowsPerWorker
		endY := startY + rowsPerWorker
		if endY > bounds.Max.Y {
			endY = bounds.Max.Y
		}
		wg.Add(1)
		go func(startY, endY int) {
			defer wg.Done()
			for y := startY; y < endY; y++ {
				for x := bounds.Min.X; x < bounds.Max.X; x++ {
					origColor := color.NRGBAModel.Convert(img.At(x, y)).(color.NRGBA)
					r := float64(origColor.R) * factor
					g := float64(origColor.G) * factor
					b := float64(origColor.B) * factor

					// Clamp values to [0, 255]
					if r > 255 {
						r = 255
					}
					if g > 255 {
						g = 255
					}
					if b > 255 {
						b = 255
					}
					if r < 0 {
						r = 0
					}
					if g < 0 {
						g = 0
					}
					if b < 0 {
						b = 0
					}

					newImg.Set(x, y, color.NRGBA{
						R: uint8(r),
						G: uint8(g),
						B: uint8(b),
						A: origColor.A,
					})
				}
			}
		}(startY, endY)
	}
	wg.Wait()
	return newImg
}

// goImageToBytes converts a Go image.Image to a byte slice in a format suitable for VIPS.
// This function assumes an RGBA format, which works for most common image types.
func goImageToBytes(img image.Image) []byte {
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	// Pre-allocate the byte slice for efficiency
	pixels := make([]byte, 0, width*height*4)

	// Iterate through the image and extract the RGBA values
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			color := img.At(x, y)
			r, g, b, a := color.RGBA()
			pixels = append(pixels, uint8(r>>8), uint8(g>>8), uint8(b>>8), uint8(a>>8))
		}
	}
	return pixels
}

var uiQueue = make(chan func(), 100)

func main() {

	// Initialize VIPS (Enable multi-threading, set memory cache)
	vips.Startup(&vips.Config{
		ConcurrencyLevel: 0,                 // 0 = Use all available CPU cores
		MaxCacheMem:      100 * 1024 * 1024, // 100MB cache (adjust as needed)
	})
	defer vips.Shutdown() // Ensure VIPS is properly shut down when the program exits

	// 1. Crear la aplicación con tema oscuro.
	myApp := app.New()
	myApp.Settings().SetTheme(theme.DarkTheme())

	// Start the UI update queue processor.
	go func() {
		for updateFn := range uiQueue {
			updateFn() // Execute the update function.
		}
	}()

	// 2. Crear la ventana principal.
	myWindow := myApp.NewWindow("Editor de Fotos - Vista Increíble")
	myWindow.Resize(fyne.NewSize(1600, 1000))

	//----------------------------------------------------------------------
	// MENÚ PRINCIPAL
	//----------------------------------------------------------------------
	menuArchivo := fyne.NewMenu("Archivo",
		fyne.NewMenuItem("Nuevo", func() {
			dialog.ShowInformation("Nuevo", "Crear un nuevo documento", myWindow)
		}),
		fyne.NewMenuItem("Abrir...", func() {
			dialog.ShowInformation("Abrir", "Seleccionar archivo a abrir", myWindow)
		}),
		fyne.NewMenuItem("Abrir Recientes", func() {
			dialog.ShowInformation("Abrir Recientes", "Lista de archivos recientes", myWindow)
		}),
		fyne.NewMenuItemSeparator(),
		fyne.NewMenuItem("Guardar", func() {}),
		fyne.NewMenuItem("Guardar Como...", func() {}),
		fyne.NewMenuItem("Exportar...", func() {}),
		fyne.NewMenuItemSeparator(),
		fyne.NewMenuItem("Salir", func() {
			myApp.Quit()
		}),
	)

	menuEditar := fyne.NewMenu("Editar",
		fyne.NewMenuItem("Deshacer (Undo)", nil),
		fyne.NewMenuItem("Rehacer (Redo)", nil),
		fyne.NewMenuItemSeparator(),
		fyne.NewMenuItem("Copiar", nil),
		fyne.NewMenuItem("Cortar", nil),
		fyne.NewMenuItem("Pegar", nil),
		fyne.NewMenuItemSeparator(),
		fyne.NewMenuItem("Preferencias", func() {
			dialog.ShowInformation("Preferencias", "Atajos de teclado, idioma, tema, etc.", myWindow)
		}),
	)

	menuVer := fyne.NewMenu("Ver",
		fyne.NewMenuItem("Zoom In", nil),
		fyne.NewMenuItem("Zoom Out", nil),
		fyne.NewMenuItem("Ajustar a Pantalla", nil),
		fyne.NewMenuItemSeparator(),
		fyne.NewMenuItem("Mostrar/ocultar Histograma", nil),
		fyne.NewMenuItem("Mostrar/ocultar Guías", nil),
	)

	menuImagen := fyne.NewMenu("Imagen",
		fyne.NewMenuItem("Ajustes Rápidos", nil),
		fyne.NewMenuItem("Filtros", nil),
		fyne.NewMenuItemSeparator(),
		fyne.NewMenuItem("Transformar", nil),
		fyne.NewMenuItem("Recortar", nil),
		fyne.NewMenuItem("Cambiar Tamaño Lienzo", nil),
	)

	menuHerramientas := fyne.NewMenu("Herramientas",
		fyne.NewMenuItem("Pincel", nil),
		fyne.NewMenuItem("Borrador", nil),
		fyne.NewMenuItem("Varita Mágica", nil),
		fyne.NewMenuItem("Texto", nil),
		fyne.NewMenuItem("Selección Libre", nil),
	)

	menuAyuda := fyne.NewMenu("Ayuda",
		fyne.NewMenuItem("Documentación", nil),
		fyne.NewMenuItem("Acerca de", func() {
			dialog.ShowInformation("Acerca de", "Editor de Fotos Pro v3.0", myWindow)
		}),
	)

	// Asignar menús a la ventana.
	myWindow.SetMainMenu(
		fyne.NewMainMenu(
			menuArchivo,
			menuEditar,
			menuVer,
			menuImagen,
			menuHerramientas,
			menuAyuda,
		),
	)

	//----------------------------------------------------------------------
	// TOOLBAR SUPERIOR (PRINCIPAL)
	//----------------------------------------------------------------------
	topToolbar := widget.NewToolbar(
		widget.NewToolbarAction(theme.DocumentCreateIcon(), func() {
			dialog.ShowInformation("Nuevo", "Crear un nuevo documento", myWindow)
		}),
		widget.NewToolbarAction(theme.FileIcon(), func() {
			dialog.ShowInformation("Abrir", "Seleccionar archivo a abrir", myWindow)
		}),
		widget.NewToolbarSeparator(),
		widget.NewToolbarAction(theme.ContentCutIcon(), func() {}),
		widget.NewToolbarAction(theme.ContentCopyIcon(), func() {}),
		widget.NewToolbarAction(theme.ContentPasteIcon(), func() {}),
		widget.NewToolbarSeparator(),
		widget.NewToolbarAction(theme.ZoomFitIcon(), func() {
			dialog.ShowInformation("Zoom", "Ajustar a pantalla", myWindow)
		}),
		widget.NewToolbarSpacer(),
		widget.NewToolbarAction(theme.HelpIcon(), func() {
			dialog.ShowInformation("Ayuda", "Sección de ayuda o tutorial", myWindow)
		}),
	)

	//----------------------------------------------------------------------
	// SEGUNDA BARRA (BUSCADOR / ATALAJES RÁPIDOS)
	//----------------------------------------------------------------------
	searchEntry := widget.NewEntry()
	searchEntry.SetPlaceHolder("Buscar herramientas, filtros, etc...")

	quickActions := widget.NewToolbar(
		widget.NewToolbarAction(theme.SearchIcon(), func() {
			// Simular búsqueda (placeholder)
			dialog.ShowInformation("Búsqueda", "Buscando: "+searchEntry.Text, myWindow)
		}),
		widget.NewToolbarAction(theme.ViewRefreshIcon(), func() {
			dialog.ShowInformation("Actualizar", "Refrescar vista o recargar imagen", myWindow)
		}),
		widget.NewToolbarAction(theme.InfoIcon(), func() {
			dialog.ShowInformation("Información", "Información de la imagen o del proyecto", myWindow)
		}),
	)

	secondBar := container.NewBorder(
		nil,          // top
		nil,          // bottom
		nil,          // left
		quickActions, // right
		searchEntry,  // center
	)

	// Combinar ambas barras superiores en un VBox
	topBars := container.NewVBox(topToolbar, secondBar)

	//----------------------------------------------------------------------
	// PANEL LATERAL IZQUIERDO: HERRAMIENTAS / CAPAS / HISTORIAL
	//----------------------------------------------------------------------
	// "Herramientas"
	toolsContent := container.NewVBox(
		widget.NewLabelWithStyle("Herramientas", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewButton("Pincel", func() {}),
		widget.NewButton("Borrador", func() {}),
		widget.NewButton("Varita Mágica", func() {}),
		widget.NewButton("Texto", func() {}),
		widget.NewButton("Selección Libre", func() {}),
	)

	// "Capas"
	layersContent := container.NewVBox(
		widget.NewLabelWithStyle("Capas", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		// Lista de capas (placeholder)
		widget.NewLabel("Capa 1: Fondo"),
		widget.NewLabel("Capa 2: Ajustes"),
		widget.NewLabel("Capa 3: Texto"),
		// Botones para gestionar capas
		widget.NewButton("Nueva Capa", func() {}),
		widget.NewButton("Eliminar Capa", func() {}),
	)

	// "Historial" (Undo/Redo)
	historyContent := container.NewVBox(
		widget.NewLabelWithStyle("Historial de Cambios", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewLabel("1. Añadir Capa de Ajuste"),
		widget.NewLabel("2. Borrar con Pincel"),
		widget.NewLabel("3. Ajustar Brillo"),
		widget.NewLabel("4. Aplicar Filtro Sepia"),
	)

	leftTabs := container.NewAppTabs(
		container.NewTabItem("Herramientas", toolsContent),
		container.NewTabItem("Capas", layersContent),
		container.NewTabItem("Historial", historyContent),
	)
	leftTabs.SetTabLocation(container.TabLocationLeading)
	sliderBrightness := widget.NewSlider(-100, 100)
	sliderSharpness := widget.NewSlider(-1, 1)

	//----------------------------------------------------------------------
	// PANEL LATERAL DERECHO: AJUSTES, FILTROS, CANALES, NAVEGADOR
	//----------------------------------------------------------------------
	// Ajustes de Color
	colorAdjustContent := container.NewVBox(
		widget.NewLabelWithStyle("Ajustes de Color", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewLabel("Brillo"),
		sliderBrightness,
		widget.NewLabel("Sharpness"),
		sliderSharpness,
		widget.NewLabel("Contraste"),
		widget.NewSlider(0, 100),
		widget.NewLabel("Saturación"),
		widget.NewSlider(0, 100),
		widget.NewLabel("Temperatura"),
		widget.NewSlider(-50, 50),
	)

	// Filtros
	filtersContent := container.NewVBox(
		widget.NewLabelWithStyle("Filtros", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewButton("Blanco y Negro", func() {}),
		widget.NewButton("Sepia", func() {}),
		widget.NewButton("Desenfoque Gaussiano", func() {}),
		widget.NewButton("Viñeta", func() {}),
		widget.NewButton("Pop Art", func() {}),
	)

	// Canales (RGB, Alfa, etc.)
	channelsContent := container.NewVBox(
		widget.NewLabelWithStyle("Canales", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewCheck("Mostrar Canal R", func(bool) {}),
		widget.NewCheck("Mostrar Canal G", func(bool) {}),
		widget.NewCheck("Mostrar Canal B", func(bool) {}),
		widget.NewCheck("Mostrar Canal Alfa", func(bool) {}),
	)

	// Navegador (Minimapa)
	navigatorContent := container.NewVBox(
		widget.NewLabelWithStyle("Navegador", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewLabel("Vista en miniatura de la imagen"),
		// Aquí podrías implementar un canvas.Image con una vista reducida
		widget.NewButton("Mover a la Sección Superior Izquierda", func() {}),
		widget.NewButton("Mover a la Sección Central", func() {}),
	)

	rightTabs := container.NewAppTabs(
		container.NewTabItem("Color", colorAdjustContent),
		container.NewTabItem("Filtros", filtersContent),
		container.NewTabItem("Canales", channelsContent),
		container.NewTabItem("Navegador", navigatorContent),
	)
	rightTabs.SetTabLocation(container.TabLocationTrailing)

	// Example translations
	translations := map[string]string{
		"image_placeholder": "Cargando imagen...", // Spanish translation
	}

	// 4. Variable to store the original loaded image.
	var originalImage image.Image

	imageArea := ui.CreateImageArea(translations)

	// 5. Button to load a new image via a file dialog.
	// Button to load a new image via a file dialog.
	loadButton := widget.NewButton("Load Image", func() {
		openDialog := dialog.NewFileOpen(func(uri fyne.URIReadCloser, err error) {
			if err != nil {
				dialog.ShowError(err, myWindow)
				return
			}
			if uri == nil {
				return
			}
			// Load the image from the selected file.
			img, _, err := image.Decode(uri)
			if err != nil {
				dialog.ShowError(err, myWindow)
				return
			}
			originalImage = img // Store the original image.
			// Update the UI with the loaded image.
			fmt.Printf("URI scheme: %s\n", uri.URI().Scheme())
			fmt.Printf("URI path: %s\n", uri.URI().Path())
			//imageArea.SetImage(uri)
			imageArea.GetImageDisplay().Image = img
			imageArea.GetImageDisplay().Refresh()
			uri.Close()
		}, myWindow)
		openDialog.SetFilter(storage.NewExtensionFileFilter([]string{".jpg", ".jpeg", ".png"}))
		openDialog.Show()
	})

	// 6. Slider for brightness.
	// Here, we do the brightness processing externally (in main).

	sliderBrightness.SetValue(0) // 50 is the "normal" brightness (factor 1.0)
	sliderBrightness.OnChanged = func(v float64) {
		if originalImage == nil {
			return // No image loaded yet.
		}
		factor := v / 50

		go func() {
			adjusted, err := imaging.AdjustBrightnessVIPS(originalImage, factor)
			if err != nil {
				fmt.Println("Error adjusting brightness:", err)
				return
			}
			// Instead of calling RunOnMain, send a function to the uiQueue.
			uiQueue <- func() {
				imageArea.GetImageDisplay().Image = adjusted
				imageArea.GetImageDisplay().Refresh()
			}
		}()
	}

	sliderSharpness.SetValue(0) // 50 is the "normal" brightness (factor 1.0)
	sliderSharpness.OnChanged = func(v float64) {
		if originalImage == nil {
			return // No image loaded yet.
		}
		factor := v / 50

		go func() {
			adjusted, err := imaging.AdjustSharpenVIPS(originalImage, factor)
			if err != nil {
				fmt.Println("Error adjusting brightness:", err)
				return
			}
			// Instead of calling RunOnMain, send a function to the uiQueue.
			uiQueue <- func() {
				imageArea.GetImageDisplay().Image = adjusted
				imageArea.GetImageDisplay().Refresh()
			}
		}()
	}

	//----------------------------------------------------------------------
	// BARRA INFERIOR: FILMSTRIP + STATUS
	//----------------------------------------------------------------------
	filmstripLabel := widget.NewLabel("Carrete de imágenes (Placeholder)")
	/*filmstripThumbnails := container.NewHBox(
		widget.NewButton("Foto 1", func() {}),
		widget.NewButton("Foto 2", func() {}),
		widget.NewButton("Foto 3", func() {}),
		widget.NewButton("Foto 4", func() {}),
	)*/

	filmstripContainer := container.NewVBox(
		loadButton,
		widget.NewLabel("Carrete de imágenes (Placeholder)"),
		container.NewHBox(
			widget.NewButton("Foto 1", func() {}),
			widget.NewButton("Foto 2", func() {}),
			widget.NewButton("Foto 3", func() {}),
			widget.NewButton("Foto 4", func() {}),
		),
	)

	// Status bar
	statusLabel := widget.NewLabel("Listo. No se ha cargado ninguna imagen.")
	statusBar := container.NewHBox(statusLabel)

	// Combinar filmstrip y status en un VBox
	bottomBar := container.NewVBox(
		container.NewVBox(
			filmstripLabel,
			filmstripContainer,
		),
		statusBar,
	)

	//----------------------------------------------------------------------
	// CONTENEDOR PRINCIPAL (BORDER)
	//----------------------------------------------------------------------
	mainContent := container.NewBorder(
		topBars,             // top
		bottomBar,           // bottom
		leftTabs,            // left
		rightTabs,           // right
		imageArea.Container, // center
	)

	myWindow.SetContent(mainContent)
	myWindow.ShowAndRun()
}
