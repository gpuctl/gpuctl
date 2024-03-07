package gpustats

import (
	"encoding/xml"
	"errors"
	"fmt"
	"os/exec"
	"strconv"
	"strings"

	"github.com/gpuctl/gpuctl/internal/procinfo"
	"github.com/gpuctl/gpuctl/internal/uplink"

	"github.com/google/uuid"
)

const (
	NvidiaNotApplicable = "N/A"
)

var (
	ErrBadField = errors.New("Parsing of data fields failed.")
)

// NvidiaSmiLog was generated 2024-01-16 15:52:50 by https://xml-to-go.github.io/ in Ukraine.
// TODO: this data structure does not support multiple GPUs nor any processes
type NvidiaSmiLog struct {
	XMLName       xml.Name `xml:"nvidia_smi_log"`
	Timestamp     string   `xml:"timestamp"`
	DriverVersion string   `xml:"driver_version"`
	CudaVersion   string   `xml:"cuda_version"`
	AttachedGpus  string   `xml:"attached_gpus"`
	Gpu           []gpu    `xml:"gpu"`
}

type gpu struct {
	ID              string `xml:"id,attr"`
	ProductName     string `xml:"product_name"`
	ProductBrand    string `xml:"product_brand"`
	DisplayMode     string `xml:"display_mode"`
	DisplayActive   string `xml:"display_active"`
	PersistenceMode string `xml:"persistence_mode"`
	MigMode         struct {
		CurrentMig string `xml:"current_mig"`
		PendingMig string `xml:"pending_mig"`
	} `xml:"mig_mode"`
	MigDevices               string `xml:"mig_devices"`
	AccountingMode           string `xml:"accounting_mode"`
	AccountingModeBufferSize string `xml:"accounting_mode_buffer_size"`
	DriverModel              struct {
		CurrentDm string `xml:"current_dm"`
		PendingDm string `xml:"pending_dm"`
	} `xml:"driver_model"`
	Serial         string `xml:"serial"`
	Uuid           string `xml:"uuid"`
	MinorNumber    string `xml:"minor_number"`
	VbiosVersion   string `xml:"vbios_version"`
	MultigpuBoard  string `xml:"multigpu_board"`
	BoardID        string `xml:"board_id"`
	GpuPartNumber  string `xml:"gpu_part_number"`
	GpuModuleID    string `xml:"gpu_module_id"`
	InforomVersion struct {
		ImgVersion string `xml:"img_version"`
		OemObject  string `xml:"oem_object"`
		EccObject  string `xml:"ecc_object"`
		PwrObject  string `xml:"pwr_object"`
	} `xml:"inforom_version"`
	GpuOperationMode struct {
		CurrentGom string `xml:"current_gom"`
		PendingGom string `xml:"pending_gom"`
	} `xml:"gpu_operation_mode"`
	GspFirmwareVersion    string `xml:"gsp_firmware_version"`
	GpuVirtualizationMode struct {
		VirtualizationMode string `xml:"virtualization_mode"`
		HostVgpuMode       string `xml:"host_vgpu_mode"`
	} `xml:"gpu_virtualization_mode"`
	Ibmnpu struct {
		RelaxedOrderingMode string `xml:"relaxed_ordering_mode"`
	} `xml:"ibmnpu"`
	Pci struct {
		PciBus         string `xml:"pci_bus"`
		PciDevice      string `xml:"pci_device"`
		PciDomain      string `xml:"pci_domain"`
		PciDeviceID    string `xml:"pci_device_id"`
		PciBusID       string `xml:"pci_bus_id"`
		PciSubSystemID string `xml:"pci_sub_system_id"`
		PciGpuLinkInfo struct {
			PcieGen struct {
				MaxLinkGen     string `xml:"max_link_gen"`
				CurrentLinkGen string `xml:"current_link_gen"`
			} `xml:"pcie_gen"`
			LinkWidths struct {
				MaxLinkWidth     string `xml:"max_link_width"`
				CurrentLinkWidth string `xml:"current_link_width"`
			} `xml:"link_widths"`
		} `xml:"pci_gpu_link_info"`
		PciBridgeChip struct {
			BridgeChipType string `xml:"bridge_chip_type"`
			BridgeChipFw   string `xml:"bridge_chip_fw"`
		} `xml:"pci_bridge_chip"`
		ReplayCounter         string `xml:"replay_counter"`
		ReplayRolloverCounter string `xml:"replay_rollover_counter"`
		TxUtil                string `xml:"tx_util"`
		RxUtil                string `xml:"rx_util"`
	} `xml:"pci"`
	FanSpeed              string `xml:"fan_speed"`
	PerformanceState      string `xml:"performance_state"`
	ClocksThrottleReasons string `xml:"clocks_throttle_reasons"`
	FbMemoryUsage         struct {
		Total string `xml:"total"`
		Used  string `xml:"used"`
		Free  string `xml:"free"`
	} `xml:"fb_memory_usage"`
	Bar1MemoryUsage struct {
		Total string `xml:"total"`
		Used  string `xml:"used"`
		Free  string `xml:"free"`
	} `xml:"bar1_memory_usage"`
	ComputeMode string `xml:"compute_mode"`
	Utilization struct {
		GpuUtil     string `xml:"gpu_util"`
		MemoryUtil  string `xml:"memory_util"`
		EncoderUtil string `xml:"encoder_util"`
		DecoderUtil string `xml:"decoder_util"`
	} `xml:"utilization"`
	EncoderStats struct {
		SessionCount   string `xml:"session_count"`
		AverageFps     string `xml:"average_fps"`
		AverageLatency string `xml:"average_latency"`
	} `xml:"encoder_stats"`
	FbcStats struct {
		SessionCount   string `xml:"session_count"`
		AverageFps     string `xml:"average_fps"`
		AverageLatency string `xml:"average_latency"`
	} `xml:"fbc_stats"`
	EccMode struct {
		CurrentEcc string `xml:"current_ecc"`
		PendingEcc string `xml:"pending_ecc"`
	} `xml:"ecc_mode"`
	EccErrors struct {
		Volatile struct {
			SingleBit struct {
				DeviceMemory  string `xml:"device_memory"`
				RegisterFile  string `xml:"register_file"`
				L1Cache       string `xml:"l1_cache"`
				L2Cache       string `xml:"l2_cache"`
				TextureMemory string `xml:"texture_memory"`
				TextureShm    string `xml:"texture_shm"`
				Cbu           string `xml:"cbu"`
				Total         string `xml:"total"`
			} `xml:"single_bit"`
			DoubleBit struct {
				DeviceMemory  string `xml:"device_memory"`
				RegisterFile  string `xml:"register_file"`
				L1Cache       string `xml:"l1_cache"`
				L2Cache       string `xml:"l2_cache"`
				TextureMemory string `xml:"texture_memory"`
				TextureShm    string `xml:"texture_shm"`
				Cbu           string `xml:"cbu"`
				Total         string `xml:"total"`
			} `xml:"double_bit"`
		} `xml:"volatile"`
		Aggregate struct {
			SingleBit struct {
				DeviceMemory  string `xml:"device_memory"`
				RegisterFile  string `xml:"register_file"`
				L1Cache       string `xml:"l1_cache"`
				L2Cache       string `xml:"l2_cache"`
				TextureMemory string `xml:"texture_memory"`
				TextureShm    string `xml:"texture_shm"`
				Cbu           string `xml:"cbu"`
				Total         string `xml:"total"`
			} `xml:"single_bit"`
			DoubleBit struct {
				DeviceMemory  string `xml:"device_memory"`
				RegisterFile  string `xml:"register_file"`
				L1Cache       string `xml:"l1_cache"`
				L2Cache       string `xml:"l2_cache"`
				TextureMemory string `xml:"texture_memory"`
				TextureShm    string `xml:"texture_shm"`
				Cbu           string `xml:"cbu"`
				Total         string `xml:"total"`
			} `xml:"double_bit"`
		} `xml:"aggregate"`
	} `xml:"ecc_errors"`
	RetiredPages struct {
		MultipleSingleBitRetirement struct {
			RetiredCount    string `xml:"retired_count"`
			RetiredPagelist string `xml:"retired_pagelist"`
		} `xml:"multiple_single_bit_retirement"`
		DoubleBitRetirement struct {
			RetiredCount    string `xml:"retired_count"`
			RetiredPagelist string `xml:"retired_pagelist"`
		} `xml:"double_bit_retirement"`
		PendingBlacklist  string `xml:"pending_blacklist"`
		PendingRetirement string `xml:"pending_retirement"`
	} `xml:"retired_pages"`
	RemappedRows string `xml:"remapped_rows"`
	Temperature  struct {
		GpuTemp                string `xml:"gpu_temp"`
		GpuTempMaxThreshold    string `xml:"gpu_temp_max_threshold"`
		GpuTempSlowThreshold   string `xml:"gpu_temp_slow_threshold"`
		GpuTempMaxGpuThreshold string `xml:"gpu_temp_max_gpu_threshold"`
		GpuTargetTemperature   string `xml:"gpu_target_temperature"`
		MemoryTemp             string `xml:"memory_temp"`
		GpuTempMaxMemThreshold string `xml:"gpu_temp_max_mem_threshold"`
	} `xml:"temperature"`
	SupportedGpuTargetTemp struct {
		GpuTargetTempMin string `xml:"gpu_target_temp_min"`
		GpuTargetTempMax string `xml:"gpu_target_temp_max"`
	} `xml:"supported_gpu_target_temp"`
	PowerReadings struct {
		PowerState         string `xml:"power_state"`
		PowerManagement    string `xml:"power_management"`
		PowerDraw          string `xml:"power_draw"`
		PowerLimit         string `xml:"power_limit"`
		DefaultPowerLimit  string `xml:"default_power_limit"`
		EnforcedPowerLimit string `xml:"enforced_power_limit"`
		MinPowerLimit      string `xml:"min_power_limit"`
		MaxPowerLimit      string `xml:"max_power_limit"`
	} `xml:"gpu_power_readings"`
	Clocks struct {
		GraphicsClock string `xml:"graphics_clock"`
		SmClock       string `xml:"sm_clock"`
		MemClock      string `xml:"mem_clock"`
		VideoClock    string `xml:"video_clock"`
	} `xml:"clocks"`
	ApplicationsClocks struct {
		GraphicsClock string `xml:"graphics_clock"`
		MemClock      string `xml:"mem_clock"`
	} `xml:"applications_clocks"`
	DefaultApplicationsClocks struct {
		GraphicsClock string `xml:"graphics_clock"`
		MemClock      string `xml:"mem_clock"`
	} `xml:"default_applications_clocks"`
	MaxClocks struct {
		GraphicsClock string `xml:"graphics_clock"`
		SmClock       string `xml:"sm_clock"`
		MemClock      string `xml:"mem_clock"`
		VideoClock    string `xml:"video_clock"`
	} `xml:"max_clocks"`
	MaxCustomerBoostClocks struct {
		GraphicsClock string `xml:"graphics_clock"`
	} `xml:"max_customer_boost_clocks"`
	ClockPolicy struct {
		AutoBoost        string `xml:"auto_boost"`
		AutoBoostDefault string `xml:"auto_boost_default"`
	} `xml:"clock_policy"`
	Voltage struct {
		GraphicsVolt string `xml:"graphics_volt"`
	} `xml:"voltage"`
	SupportedClocks struct {
		SupportedMemClock []struct {
			Value                  string   `xml:"value"`
			SupportedGraphicsClock []string `xml:"supported_graphics_clock"`
		} `xml:"supported_mem_clock"`
	} `xml:"supported_clocks"`
	Processes          processes `xml:"processes"`
	AccountedProcesses string    `xml:"accounted_processes"`
}

type processes struct {
	ProcessInfo []struct {
		GpuInstanceID     string `xml:"gpu_instance_id"`
		ComputeInstanceID string `xml:"compute_instance_id"`
		Pid               string `xml:"pid"`
		Type              string `xml:"type"`
		ProcessName       string `xml:"process_name"`
		UsedMemory        string `xml:"used_memory"`
	} `xml:"process_info"`
}

func (xml NvidiaSmiLog) ExtractGPUInfo() ([]uplink.GPUInfo, error) {
	var res []uplink.GPUInfo

	for _, gpu := range xml.Gpu {
		mem, err := parseUIntWithUnit(gpu.FbMemoryUsage.Total, "mem usage")
		if err != nil {
			return nil, err
		}
		uuid, err := parseUuid(gpu.Uuid)
		if err != nil {
			return nil, err
		}

		res = append(res,
			uplink.GPUInfo{
				Uuid:          uuid,
				Name:          gpu.ProductName,
				Brand:         gpu.ProductBrand,
				DriverVersion: xml.DriverVersion,
				MemoryTotal:   mem,
			})
	}

	return res, nil
}

// Filter down the relevant information from our nvidia-smi dump
func (xml NvidiaSmiLog) ExtractGPUStatSample() ([]uplink.GPUStatSample, error) {
	var res []uplink.GPUStatSample

	// Go over gpus in xml
	var errs error
	for _, gpu := range xml.Gpu {
		// TODO: add the new information as well
		fanSpeed, err := parseFloatWithUnit(gpu.FanSpeed, "fan speed")
		memoryUsed, err_ := parseFloatWithUnit(gpu.FbMemoryUsage.Used, "mem used")
		err = errors.Join(err, err_)
		temp, err_ := parseFloatWithUnit(gpu.Temperature.GpuTemp, "gpu temp")
		err = errors.Join(err, err_)
		gpuUtil, err_ := parseFloatWithUnit(gpu.Utilization.GpuUtil, "gpu util")
		err = errors.Join(err, err_)
		memUtil, err_ := parseFloatWithUnit(gpu.Utilization.MemoryUtil, "mem util")
		err = errors.Join(err, err_)
		memTemp, err_ := parseFloatWithUnit(gpu.Temperature.MemoryTemp, "mem temp")
		err = errors.Join(err, err_)
		graphicsVolt, err_ := parseFloatWithUnit(gpu.Voltage.GraphicsVolt, "gpu volt")
		err = errors.Join(err, err_)
		powerDraw, err_ := parseFloatWithUnit(gpu.PowerReadings.PowerDraw, "power draw") // SUS
		err = errors.Join(err, err_)
		gFreq, err_ := parseFloatWithUnit(gpu.Clocks.GraphicsClock, "gpu clock")
		err = errors.Join(err, err_)
		maxGFreq, err_ := parseFloatWithUnit(gpu.MaxClocks.GraphicsClock, "gpu max clock")
		err = errors.Join(err, err_)
		mFreq, err_ := parseFloatWithUnit(gpu.Clocks.MemClock, "mem clock")
		err = errors.Join(err, err_)
		maxMFreq, err_ := parseFloatWithUnit(gpu.MaxClocks.MemClock, "max mem clock")
		err = errors.Join(err, err_)
		uuid, err_ := parseUuid(gpu.Uuid)
		err = errors.Join(err, err_)


		// Filter processes
		procs := []uplink.GPUProcInfo{}

		for _, proc := range gpu.Processes.ProcessInfo {
			procmem, err := parseFloatWithUnit(proc.UsedMemory, "proc used mem")
			pid, err2 := strconv.ParseUint(proc.Pid, 10, 0)

			err = errors.Join(err, err2)
			if err != nil {
				return nil, err
			}

			procs = append(procs,
				uplink.GPUProcInfo{
					Pid:     pid,
					Name:    proc.ProcessName,
					MemUsed: procmem,
				})
		}

		curr := uplink.GPUStatSample{
			Uuid:              uuid,
			MemoryUtilisation: memUtil,
			GPUUtilisation:    gpuUtil,
			MemoryUsed:        memoryUsed,
			FanSpeed:          fanSpeed,
			Temp:              temp,
			MemoryTemp:        memTemp,
			GraphicsVoltage:   graphicsVolt,
			PowerDraw:         powerDraw,
			GraphicsClock:     gFreq,
			MaxGraphicsClock:  maxGFreq,
			MemoryClock:       mFreq,
			MaxMemoryClock:    maxMFreq,
			RunningProcesses:  procs,
		}
		errs = errors.Join(errs, err)
		res = append(res, curr)
	}
	return res, errs
}

// convert nvidia uuid to standard uuid
// nvidia include an unneeded GPU- on the front of all their uuids
func parseUuid(input string) (uuid.UUID, error) {
	after, _ := strings.CutPrefix(input, "GPU-")
	return uuid.Parse(after)
}

// Helper function for extracting GPU data such as available memory and fan speed
func parseUIntWithUnit(input string, name string) (uint64, error) {
	part, _, _ := strings.Cut(input, " ")
	// Interpret N/A as 0
	if part == NvidiaNotApplicable {
		return 0, nil
	}

	// wrap error up in name of component on failure
	val, err := strconv.ParseUint(part, 10, 0)
	if err != nil {
		err = fmt.Errorf("parsing %s: %w", name, err)
	}

	return val, err
}

func parseFloatWithUnit(input string, name string) (float64, error) {
	part, _, _ := strings.Cut(input, " ")
	// Interpret N/A as 0
	if part == NvidiaNotApplicable {
		return 0, nil
	}

	val, err := strconv.ParseFloat(part, 10)
	if err != nil {
		err = fmt.Errorf("parsing %s: %w", name, err)
	}

	return val, err
}

// Helper function to unmarshal Nvidia XML dump
func ParseNvidiaSmi(input []byte) (NvidiaSmiLog, error) {
	var result NvidiaSmiLog
	if e := xml.Unmarshal(input, &result); e != nil {
		return NvidiaSmiLog{}, e
	}
	return result, nil
}

// Get the Nvidia GPU status directly from the computer using `nvidia-smi`
func getNvidiaGPUStatus() (NvidiaSmiLog, error) {
	output, err := exec.Command("nvidia-smi", "-q", "-x").Output()
	if err != nil {
		return NvidiaSmiLog{}, err
	}

	return ParseNvidiaSmi(output)
}

// Struct to act as our adapter
type NvidiaGPUHandler struct {
	Lookup     procinfo.UidLookup            // Which Uid maps to which user
	ProcFilter func(uplink.GPUProcInfo) bool // How we determine if a process is worth considering
}

// Run the whole pipeline of getting GPU information
func (h NvidiaGPUHandler) GetGPUStatus() ([]uplink.GPUStatSample, error) {
	smi, err := getNvidiaGPUStatus()
	if err != nil {
		return nil, err
	}

	samples, err := smi.ExtractGPUStatSample()
	if err != nil {
		return nil, err
	}

	uplink.FilterProcesses(samples, h.ProcFilter)
	uplink.PopulateNames(samples, h.Lookup)
	return samples, nil
}

// Run the whole pipeline of getting GPU information
func (h NvidiaGPUHandler) GetGPUInformation() ([]uplink.GPUInfo, error) {
	smi, err := getNvidiaGPUStatus()
	if err != nil {
		return nil, err
	}
	return smi.ExtractGPUInfo()
}
