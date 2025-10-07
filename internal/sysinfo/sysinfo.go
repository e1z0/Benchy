package sysinfo

import (
	"fmt"
	"runtime"
)

type Info struct {
	GoVersion string `json:"go_version"`
	OS        string `json:"os"`
	Arch      string `json:"arch"`
	// CPU
	CPUModel      string `json:"cpu_model"`       // e.g., "Apple M2 Pro" / "Intel(R) Core(TM) i7-12700"
	CPUVendor     string `json:"cpu_vendor"`      // Apple / Intel / AMD / etc
	PhysicalCores int    `json:"physical_cores"`  // where available (fallback 0 if unknown)
	LogicalCPUs   int    `json:"logical_cpus"`    // NumCPU()
	NominalFreqHz uint64 `json:"nominal_freq_hz"` // nominal/base frequency if known (0 if unknown)

	// Machine / Board
	MachineModel    string `json:"machine_model"`    // e.g., "MacBookPro16,1" or "Precision 3460"
	SystemVendor    string `json:"system_vendor"`    // e.g., "Apple" / "Dell Inc."
	ProductName     string `json:"product_name"`     // e.g., "Precision 3460"
	ProductVersion  string `json:"product_version"`  // e.g., "Not Specified" / device revision
	BoardName       string `json:"board_name"`       // e.g., "X570 AORUS PRO"
	FirmwareVersion string `json:"firmware_version"` // BIOS/EFI version if available

	// Memory (physical)
	TotalRAMBytes uint64 `json:"total_ram_bytes"`
}

func Collect() Info {
	inf := Info{
		GoVersion:   runtime.Version(),
		OS:          runtime.GOOS,
		Arch:        runtime.GOARCH,
		LogicalCPUs: runtime.NumCPU(),
	}
	populateExtra(&inf) // implemented in per-OS files
	return inf
}

// String is a short friendly summary for UI.
func (i Info) String() string {
	freq := ""
	if i.NominalFreqHz > 0 {
		freq = fmt.Sprintf(" @ %.2f GHz", float64(i.NominalFreqHz)/1e9)
	}
	mem := ""
	if i.TotalRAMBytes > 0 {
		mem = fmt.Sprintf("\nMemory: %.1f GB", float64(i.TotalRAMBytes)/1073741824.0)
	}
	model := i.MachineModel
	if model == "" {
		model = i.ProductName
	}
	if model != "" && i.SystemVendor != "" {
		model = i.SystemVendor + " " + model
	}

	return fmt.Sprintf(
		"Go: %s\nOS/Arch: %s/%s\nCPU: %s%s\nCores: %d physical / %d logical\nMachine: %s%s",
		i.GoVersion, i.OS, i.Arch,
		coalesce(i.CPUVendor+" ", "")+coalesce(i.CPUModel, "Unknown CPU"), freq,
		i.PhysicalCores, i.LogicalCPUs,
		coalesce(model, "Unknown"), mem,
	)
}

func coalesce(a, b string) string {
	if a != "" {
		return a
	}
	return b
}
