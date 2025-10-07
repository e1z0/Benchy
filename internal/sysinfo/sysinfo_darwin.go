//go:build darwin

// internal/sysinfo/sysinfo_darwin.go
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
	"os/exec"
	"strconv"
	"strings"
)

func populateExtra(i *Info) {
	// CPU brand & vendor (brand string includes vendor typically)
	i.CPUModel = sysctlStr("machdep.cpu.brand_string")
	if i.CPUModel == "" {
		// Apple Silicon doesn't expose machdep.* — fall back to chip name
		i.CPUModel = sysctlStr("hw.model") // often shows SoC family on Apple Silicon
	}
	i.CPUVendor = sysctlStr("machdep.cpu.vendor") // empty on Apple Silicon; that's fine

	// Cores
	i.PhysicalCores = sysctlInt("hw.physicalcpu")
	if i.PhysicalCores == 0 {
		i.PhysicalCores = sysctlInt("hw.physicalcpu_max")
	}

	// Nominal CPU frequency (Hz)
	i.NominalFreqHz = uint64(sysctlInt64("hw.cpufrequency"))

	// Memory
	i.TotalRAMBytes = uint64(sysctlInt64("hw.memsize"))

	// Machine model (MacBookPro16,1 etc)
	i.MachineModel = sysctlStr("hw.model")

	// Firmware/ioreg (optional): try to get board name/system vendor via ioreg
	// These may be empty on newer Macs but harmless to try.
	if out := run("ioreg", "-rd1", "-c", "IOPlatformExpertDevice"); out != "" {
		i.ProductName = parseAfter(out, `"product-name" = "`, `"`)
		i.BoardName = parseAfter(out, `"board-id" = "`, `"`)
		i.FirmwareVersion = parseAfter(out, `"boot-rom-version" = "`, `"`)
		i.SystemVendor = "Apple"
	}
}

func sysctlStr(key string) string { return strings.TrimSpace(run("sysctl", "-n", key)) }
func sysctlInt(key string) int {
	v, _ := strconv.Atoi(strings.TrimSpace(run("sysctl", "-n", key)))
	return v
}
func sysctlInt64(key string) int64 {
	v, _ := strconv.ParseInt(strings.TrimSpace(run("sysctl", "-n", key)), 10, 64)
	return v
}

func run(name string, args ...string) string {
	out, _ := exec.Command(name, args...).Output()
	return string(out)
}
func parseAfter(s, a, z string) string {
	if p := strings.Index(s, a); p >= 0 {
		p += len(a)
		if q := strings.Index(s[p:], z); q >= 0 {
			return s[p : p+q]
		}
	}
	return ""
}
