package main

import (
	"fyne.io/fyne/v2/storage"
	"photo-editor/ui"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func main() {
	// 1. Crear la aplicación con tema oscuro.
	myApp := app.New()
	myApp.Settings().SetTheme(theme.DarkTheme())

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

	//----------------------------------------------------------------------
	// PANEL LATERAL DERECHO: AJUSTES, FILTROS, CANALES, NAVEGADOR
	//----------------------------------------------------------------------
	// Ajustes de Color
	colorAdjustContent := container.NewVBox(
		widget.NewLabelWithStyle("Ajustes de Color", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewLabel("Brillo"),
		widget.NewSlider(0, 100),
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

	imageView := ui.CreateImageArea(translations)

	// Button to load a new image
	loadButton := widget.NewButton("Load Image", func() {
		openDialog := dialog.NewFileOpen(func(uri fyne.URIReadCloser, err error) {
			if err != nil {
				dialog.ShowError(err, myWindow)
				return
			}
			if uri == nil {
				return
			}
			// Actualizar la imagen usando el stream del archivo seleccionado
			imageView.SetImage(uri)
			uri.Close()
		}, myWindow)
		// Filtrar para que solo se muestren imágenes
		openDialog.SetFilter(storage.NewExtensionFileFilter([]string{".jpg", ".jpeg", ".png"}))
		openDialog.Show()
	})

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
		imageView.Container, // center
	)

	myWindow.SetContent(mainContent)
	myWindow.ShowAndRun()
}
