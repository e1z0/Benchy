//go:build windows

// internal/sysinfo/sysinfo_windows.go
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
	"github.com/yusufpapurcu/wmi"
)

type win32_Processor struct {
	Name                      string
	Manufacturer              string
	NumberOfCores             uint32
	NumberOfLogicalProcessors uint32
	MaxClockSpeed             uint32 // MHz (approx base clock)
}
type win32_ComputerSystem struct {
	Manufacturer        string
	Model               string
	TotalPhysicalMemory uint64
}
type win32_BaseBoard struct {
	Product string
}
type win32_BIOS struct {
	SMBIOSBIOSVersion string
}

func populateExtra(i *Info) {
	// CPU
	var cpus []win32_Processor
	_ = wmi.Query("SELECT Name, Manufacturer, NumberOfCores, NumberOfLogicalProcessors, MaxClockSpeed FROM Win32_Processor", &cpus)
	if len(cpus) > 0 {
		c := cpus[0]
		i.CPUModel = c.Name
		i.CPUVendor = c.Manufacturer
		if c.NumberOfCores > 0 {
			i.PhysicalCores = int(c.NumberOfCores)
		}
		if c.NumberOfLogicalProcessors > 0 {
			i.LogicalCPUs = int(c.NumberOfLogicalProcessors)
		}
		if c.MaxClockSpeed > 0 {
			i.NominalFreqHz = uint64(c.MaxClockSpeed) * 1_000_000 // MHz -> Hz
		}
	}

	// System
	var sys []win32_ComputerSystem
	_ = wmi.Query("SELECT Manufacturer, Model, TotalPhysicalMemory FROM Win32_ComputerSystem", &sys)
	if len(sys) > 0 {
		s := sys[0]
		i.SystemVendor = s.Manufacturer
		i.ProductName = s.Model
		i.MachineModel = s.Model
		i.TotalRAMBytes = s.TotalPhysicalMemory
	}

	// Board (optional)
	var boards []win32_BaseBoard
	_ = wmi.Query("SELECT Product FROM Win32_BaseBoard", &boards)
	if len(boards) > 0 {
		i.BoardName = boards[0].Product
	}

	// BIOS/Firmware
	var bios []win32_BIOS
	_ = wmi.Query("SELECT SMBIOSBIOSVersion FROM Win32_BIOS", &bios)
	if len(bios) > 0 {
		i.FirmwareVersion = bios[0].SMBIOSBIOSVersion
	}
}
