package imaging

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	"sync"
	"time"

	"github.com/davidbyttow/govips/v2/vips"
)

// ------------------------
// Generic Debouncer
// ------------------------

type job[T any, R any] struct {
	payload T
	result  chan jobResult[R]
}

type jobResult[R any] struct {
	result R
	err    error
}

type Debouncer[T any, R any] struct {
	process      func(T) (R, error)
	debounceTime time.Duration
	jobQueue     chan job[T, R]
}

func NewDebouncer[T any, R any](process func(T) (R, error), debounceTime time.Duration, queueSize int) *Debouncer[T, R] {
	d := &Debouncer[T, R]{
		process:      process,
		debounceTime: debounceTime,
		jobQueue:     make(chan job[T, R], queueSize),
	}
	go d.processJobs()
	return d
}

func (d *Debouncer[T, R]) processJobs() {
	var (
		currentJob *job[T, R]
		timer      *time.Timer
	)

	for {
		var timerCh <-chan time.Time
		if timer != nil {
			timerCh = timer.C
		}

		select {
		case newJob := <-d.jobQueue:
			if currentJob != nil {
				currentJob.result <- jobResult[R]{err: errors.New("job cancelled due to debouncing")}
			}
			currentJob = &newJob
			if timer == nil {
				timer = time.NewTimer(d.debounceTime)
			} else {
				if !timer.Stop() {
					select {
					case <-timer.C:
					default:
					}
				}
				timer.Reset(d.debounceTime)
			}

		case <-timerCh:
			if currentJob != nil {
				res, err := d.process(currentJob.payload)
				currentJob.result <- jobResult[R]{result: res, err: err}
				currentJob = nil
			}
		}
	}
}

func (d *Debouncer[T, R]) Do(payload T) (R, error) {
	resultChan := make(chan jobResult[R])
	d.jobQueue <- job[T, R]{payload: payload, result: resultChan}
	result := <-resultChan
	return result.result, result.err
}

// ------------------------
// Brightness Adjustment Using the Generic Debouncer
// ------------------------

type brightnessPayload struct {
	img    image.Image
	factor float64
}

var (
	brightnessDebouncerOnce sync.Once
	brightnessDebouncer     *Debouncer[brightnessPayload, image.Image]
)

func initBrightnessDebouncer() {
	brightnessDebouncerOnce.Do(func() {
		brightnessDebouncer = NewDebouncer[brightnessPayload, image.Image](processBrightnessJob, 10*time.Millisecond, 10)
	})
}

func processBrightnessJob(p brightnessPayload) (image.Image, error) {
	return adjustBrightnessVIPSInternal(p.img, p.factor)
}

func AdjustBrightnessVIPS(img image.Image, factor float64) (image.Image, error) {
	initBrightnessDebouncer()
	return brightnessDebouncer.Do(brightnessPayload{img: img, factor: factor})
}

func adjustBrightnessVIPSInternal(img image.Image, factor float64) (image.Image, error) {
	var buf bytes.Buffer
	if err := jpeg.Encode(&buf, img, nil); err != nil {
		return nil, err
	}
	processedBuffer, err := AdjustBrightnessVIPSFromBuffer(buf.Bytes(), factor)
	if err != nil {
		return nil, err
	}
	processedImg, _, err := image.Decode(bytes.NewReader(processedBuffer))
	if err != nil {
		return nil, err
	}
	return processedImg, nil
}

func AdjustBrightnessVIPSFromBuffer(buf []byte, factor float64) ([]byte, error) {
	vipsImg, err := vips.NewImageFromBuffer(buf)
	if err != nil {
		return nil, fmt.Errorf("failed to create VIPS image from buffer: %w", err)
	}
	defer vipsImg.Close()

	bands := vipsImg.Bands()
	factors := make([]float64, bands)
	for i := 0; i < bands; i++ {
		factors[i] = factor
	}
	if err := vipsImg.Linear(factors, []float64{0, 0, 0}); err != nil {
		return nil, fmt.Errorf("VIPS Linear failed: %w", err)
	}

	processedBuffer, _, err := vipsImg.ExportJpeg(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to export VIPS image: %w", err)
	}
	return processedBuffer, nil
}

// ------------------------
// Sharpening Using the Generic Debouncer (Single Factor API)
// ------------------------

type sharpenPayload struct {
	img           image.Image
	sigma, x1, m2 float64
}

var (
	sharpenDebouncerOnce sync.Once
	sharpenDebouncer     *Debouncer[sharpenPayload, image.Image]
)

func initSharpenDebouncer() {
	sharpenDebouncerOnce.Do(func() {
		sharpenDebouncer = NewDebouncer[sharpenPayload, image.Image](processSharpenJob, 10*time.Millisecond, 10)
	})
}

// processSharpenJob maps the single sharpening factor to VIPS's three parameters.
func processSharpenJob(p sharpenPayload) (image.Image, error) {
	return adjustSharpenVIPSInternal(p.img, p.sigma, p.x1, p.m2)
}

// AdjustSharpenVIPS applies a sharpening effect using a single factor parameter.
// Internally, this factor is mapped to the three parameters required by VIPS.
func AdjustSharpenVIPS(img image.Image, factor float64) (image.Image, error) {
	initSharpenDebouncer()
	sigma, x1, m2 := mapSharpenFactor(factor)
	return sharpenDebouncer.Do(sharpenPayload{img: img, sigma: sigma, x1: x1, m2: m2})
}

// AdjustSharpenVIPS applies a sharpening effect using a single factor parameter.
// Internally, this factor is mapped to the three parameters required by VIPS.
func AdjustSharpenVIPSFull(img image.Image, sigma float64, x1 float64, m2 float64) (image.Image, error) {
	initSharpenDebouncer()
	return sharpenDebouncer.Do(sharpenPayload{img: img, sigma: sigma, x1: x1, m2: m2})
}

// adjustSharpenVIPSInternal converts the single factor into the three required parameters.
func adjustSharpenVIPSInternal(img image.Image, sigma float64, x1 float64, m2 float64) (image.Image, error) {
	// Map the single factor to sigma, x1, and m2.

	var buf bytes.Buffer
	if err := jpeg.Encode(&buf, img, nil); err != nil {
		return nil, err
	}
	processedBuffer, err := AdjustSharpenVIPSFromBuffer(buf.Bytes(), sigma, x1, m2)
	if err != nil {
		return nil, err
	}
	processedImg, _, err := image.Decode(bytes.NewReader(processedBuffer))
	if err != nil {
		return nil, err
	}
	return processedImg, nil
}

// mapSharpenFactor maps a single sharpening factor to the three parameters for VIPS.
// You can adjust this mapping as needed to fine-tune the sharpening effect.
func mapSharpenFactor(factor float64) (sigma, x1, m2 float64) {
	// These are example mappings; tweak them to suit your desired effect.
	sigma = factor    // Controls the radius of the blur used in edge detection.
	x1 = factor * 0.8 // Acts as a threshold for sharpening.
	m2 = factor * 0.5 // Controls the intensity of edge enhancement.
	return
}

// AdjustSharpenVIPSFromBuffer applies VIPS's sharpen using the given parameters.
func AdjustSharpenVIPSFromBuffer(buf []byte, sigma, x1, m2 float64) ([]byte, error) {
	vipsImg, err := vips.NewImageFromBuffer(buf)
	if err != nil {
		return nil, fmt.Errorf("failed to create VIPS image from buffer: %w", err)
	}
	defer vipsImg.Close()

	if err := vipsImg.Sharpen(sigma, x1, m2); err != nil {
		return nil, fmt.Errorf("failed to sharpen VIPS image: %w", err)
	}

	processedBuffer, _, err := vipsImg.ExportJpeg(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to export VIPS image: %w", err)
	}
	return processedBuffer, nil
}
