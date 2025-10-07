/* SPDX-License-Identifier: GPL-3.0-or-later
 *
 * Benchy
 * Copyright (C) 2025 e1z0 <e1z0@icloud.com>
 *
 * This file is part of Benchy.
 *
 * Benchy is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * Benchy is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with Benchy. If not, see <https://www.gnu.org/licenses/>.
 */
package benchmarks

import (
	"context"
	"math"
	"runtime"
	"sync"
	"time"
)

type Image struct {
	W, H int
	Pix  []float32
}

func makeNoise(w, h int) *Image {
	pix := make([]float32, w*h)
	for i := range pix {
		pix[i] = float32((i*1664525+1013904223)&0xffff) / 65535.0
	}
	return &Image{W: w, H: h, Pix: pix}
}

func gaussianKernel(r int, sigma float64) []float32 {
	k := make([]float32, 2*r+1)
	var sum float64
	for i := -r; i <= r; i++ {
		v := math.Exp(-float64(i*i) / (2 * sigma * sigma))
		k[i+r] = float32(v)
		sum += v
	}
	for i := range k {
		k[i] = float32(float64(k[i]) / sum)
	}
	return k
}

func blur1D(dst, src []float32, w, h, r int, k []float32, horizontal bool) {
	if horizontal {
		for y := 0; y < h; y++ {
			o := y * w
			for x := 0; x < w; x++ {
				var acc float32
				for i := -r; i <= r; i++ {
					xi := x + i
					if xi < 0 {
						xi = 0
					} else if xi >= w {
						xi = w - 1
					}
					acc += src[o+xi] * k[i+r]
				}
				dst[o+x] = acc
			}
		}
	} else {
		for y := 0; y < h; y++ {
			for x := 0; x < w; x++ {
				var acc float32
				for i := -r; i <= r; i++ {
					yi := y + i
					if yi < 0 {
						yi = 0
					} else if yi >= h {
						yi = h - 1
					}
					acc += src[yi*w+x] * k[i+r]
				}
				dst[y*w+x] = acc
			}
		}
	}
}

func RunImageBlur(ctx context.Context, dur time.Duration, threads int) Result {
	if threads <= 0 {
		threads = runtime.NumCPU()
	}
	ctx, cancel := context.WithTimeout(ctx, dur)
	defer cancel()

	img := makeNoise(1920, 1080)
	r := 5
	k := gaussianKernel(r, 2.0)
	workPerThread := (img.H + threads - 1) / threads

	var wg sync.WaitGroup
	var px uint64
	start := time.Now()
	for y0 := 0; y0 < img.H; y0 += workPerThread {
		y1 := y0 + workPerThread
		if y1 > img.H {
			y1 = img.H
		}
		wg.Add(1)
		go func(y0, y1 int) {
			defer wg.Done()
			src := img.Pix[y0*img.W : y1*img.W]
			tmp := make([]float32, len(src))
			dst := make([]float32, len(src))
			for {
				select {
				case <-ctx.Done():
					return
				default:
					blur1D(tmp, src, img.W, y1-y0, r, k, true)
					blur1D(dst, tmp, img.W, y1-y0, r, k, false)
					px += uint64((y1 - y0) * img.W)
				}
			}
		}(y0, y1)
	}
	wg.Wait()
	elapsed := time.Since(start)
	if elapsed == 0 {
		elapsed = dur
	}
	return Result{Name: "Gaussian Blur 1080p", Threads: threads, Duration: dur, Ops: px / uint64(max64(1, int64(elapsed.Seconds()))), Unit: "px/s"}
}
