//go:build linux

// internal/sysinfo/sysinfo_linux.go
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
package sysinfo

import (
	"os"
	"strconv"
	"strings"
)

func populateExtra(i *Info) {
	// CPU vendor/model from /proc/cpuinfo (first processor)
	if b, _ := os.ReadFile("/proc/cpuinfo"); len(b) > 0 {
		text := string(b)
		i.CPUVendor = firstKV(text, "vendor_id\t:")
		if i.CPUVendor == "" {
			i.CPUVendor = firstKV(text, "Hardware\t:") // ARM
		}
		i.CPUModel = firstKV(text, "model name\t:")
		if i.CPUModel == "" {
			i.CPUModel = firstKV(text, "Processor\t:")
		}
		// Physical cores is tricky on Linux; leave 0 if unknown.
	}

	// Nominal frequency: try /sys/devices/system/cpu/cpu0/cpufreq/cpuinfo_max_freq (kHz)
	if b, err := os.ReadFile("/sys/devices/system/cpu/cpu0/cpufreq/cpuinfo_max_freq"); err == nil {
		if khz, _ := strconv.ParseUint(strings.TrimSpace(string(b)), 10, 64); khz > 0 {
			i.NominalFreqHz = khz * 1000
		}
	}

	// DMI / machine model & vendor
	i.SystemVendor = readFirst("/sys/class/dmi/id/sys_vendor")
	i.ProductName = readFirst("/sys/class/dmi/id/product_name")
	i.ProductVersion = readFirst("/sys/class/dmi/id/product_version")
	i.MachineModel = i.ProductName // common mapping
	i.BoardName = readFirst("/sys/class/dmi/id/board_name")
	i.FirmwareVersion = readFirst("/sys/class/dmi/id/bios_version")

	// Memory
	if b, err := os.ReadFile("/proc/meminfo"); err == nil {
		// MemTotal: kB
		if kb, _ := parseMemKB(string(b)); kb > 0 {
			i.TotalRAMBytes = uint64(kb) * 1024
		}
	}
}

func firstKV(text, key string) string {
	for _, line := range strings.Split(text, "\n") {
		if strings.HasPrefix(line, key) {
			return strings.TrimSpace(strings.TrimPrefix(line, key))
		}
	}
	return ""
}

func readFirst(path string) string {
	if b, err := os.ReadFile(path); err == nil {
		return strings.TrimSpace(string(b))
	}
	return ""
}

func parseMemKB(meminfo string) (uint64, error) {
	for _, line := range strings.Split(meminfo, "\n") {
		if strings.HasPrefix(line, "MemTotal:") {
			f := strings.Fields(line) // ["MemTotal:", "16337756", "kB"]
			if len(f) >= 2 {
				return strconv.ParseUint(f[1], 10, 64)
			}
		}
	}
	return 0, nil
}
