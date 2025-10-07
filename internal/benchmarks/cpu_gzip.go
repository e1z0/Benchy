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
	"bytes"
	"compress/gzip"
	"context"
	"fmt"
	"math/rand"
	"runtime"
	"sync"
	"time"
)

func RunCPUGzip(ctx context.Context, dur time.Duration, threads, level int) Result {
	if threads <= 0 {
		threads = runtime.NumCPU()
	}
	if level < gzip.HuffmanOnly || level > gzip.BestCompression {
		level = gzip.DefaultCompression
	}
	ctx, cancel := context.WithTimeout(ctx, dur)
	defer cancel()

	var wg sync.WaitGroup
	var bytesTotal uint64
	bufSize := 4 * 1024 * 1024

	seed := rand.New(rand.NewSource(42))
	srcTemplate := make([]byte, bufSize)
	for i := range srcTemplate {
		srcTemplate[i] = byte((i*31 + 7) ^ (i >> 3))
		if seed.Intn(5) == 0 {
			srcTemplate[i] ^= 0xFF
		}
	}

	for i := 0; i < threads; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			src := make([]byte, bufSize)
			copy(src, srcTemplate)
			var local uint64
			for {
				select {
				case <-ctx.Done():
					bytesTotal += local
					return
				default:
					var out bytes.Buffer
					zw, _ := gzip.NewWriterLevel(&out, level)
					_, _ = zw.Write(src)
					_ = zw.Close()
					local += uint64(len(src))
				}
			}
		}()
	}
	wg.Wait()
	return Result{Name: "Gzip Compress", Threads: threads, Duration: dur, Bytes: bytesTotal, Unit: "B/s", Notes: "level=" + fmt.Sprint(level)}
}
