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
	"crypto/rand"
	"io"
	"os"
	"time"
)

func RunDiskSeq(ctx context.Context, dur time.Duration, path string) Result {
	ctx, cancel := context.WithTimeout(ctx, dur)
	defer cancel()

	chunk := make([]byte, 64*1024*1024)
	_, _ = io.ReadFull(rand.Reader, chunk)

	start := time.Now()
	f, err := os.Create(path)
	if err != nil {
		return Result{Name: "Disk seq R/W", Err: err.Error()}
	}
	var wBytes uint64
	for time.Since(start) < dur {
		select {
		case <-ctx.Done():
			break
		default:
		}
		n, err := f.Write(chunk)
		if err != nil {
			_ = f.Close()
			return Result{Name: "Disk seq R/W", Err: err.Error()}
		}
		wBytes += uint64(n)
	}
	_ = f.Sync()
	_ = f.Close()

	rf, err := os.Open(path)
	if err != nil {
		return Result{Name: "Disk seq R/W", Err: err.Error()}
	}
	defer rf.Close()
	var rBytes uint64
	start = time.Now()
	buf := make([]byte, len(chunk))
	for time.Since(start) < dur {
		select {
		case <-ctx.Done():
			break
		default:
		}
		n, err := rf.Read(buf)
		if n > 0 {
			rBytes += uint64(n)
		}
		if err == io.EOF {
			_, _ = rf.Seek(0, 0)
			continue
		}
		if err != nil {
			return Result{Name: "Disk seq R/W", Err: err.Error()}
		}
	}
	_ = os.Remove(path)

	notes := "write " + humanBytes(wBytes) + "/s, read " + humanBytes(rBytes) + "/s"
	return Result{Name: "Disk seq R/W", Duration: dur, Bytes: rBytes, Unit: "B/s", Notes: notes}
}
