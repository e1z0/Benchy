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
	"fmt"
	"runtime"
	"sync"
	"time"

	"github.com/klauspost/compress/zstd"
)

func RunCPUZstd(ctx context.Context, dur time.Duration, threads, level int) Result {
	if threads <= 0 {
		threads = runtime.NumCPU()
	}
	if level < 1 || level > 19 {
		level = 3
	}
	ctx, cancel := context.WithTimeout(ctx, dur)
	defer cancel()

	var wg sync.WaitGroup
	var bytes uint64
	block := make([]byte, 4*1024*1024)

	for i := 0; i < threads; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			enc, _ := zstd.NewWriter(nil, zstd.WithEncoderLevel(zstd.EncoderLevelFromZstd(level)))
			var local uint64
			for {
				select {
				case <-ctx.Done():
					bytes += local
					return
				default:
					_ = enc.EncodeAll(block, nil)
					local += uint64(len(block))
				}
			}
		}()
	}
	wg.Wait()
	return Result{Name: "Zstd Compress", Threads: threads, Duration: dur, Bytes: bytes, Unit: "B/s", Notes: "level=" + fmt.Sprint(level)}
}
