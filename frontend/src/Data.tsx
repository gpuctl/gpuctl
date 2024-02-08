// These types need to be kept in sync with `internal/webapi/types.go`

export type WorkStationGroup = {
  name: string;
  // TODO: Consistant cases here
  workStations: WorkStationData[];
};

export type WorkStationData = {
  name: string;
  cpu : string;
  motherboard: string;
  notes: string;
  gpus: GPUStats[];
};

// This needs to be kept in sync with `internal/uplink/
export type GPUStats = {
  gpu_name: string;
  gpu_brand: string;
  driver_ver: string;
  memory_total: number;

  memory_util: number;
  gpu_util: number;
  memory_used: number;
  fan_speed: number;
  gpu_temp: number;
};

export const EXAMPLE_DATA_1: WorkStationGroup[] = [
  {
    name: "Shared",
    workStations: [
      {
        name: "Workstation 1",
        cpu: "i7-7700k",
        motherboard: "ASUS ROG STRIX Z790-A",
        notes: "noisy fan",
        gpus: [
          {
            gpu_name: "NVIDIA GeForce GT 1030",
            gpu_brand: "GeForce",
            driver_ver: "535.146.02",
            memory_total: 2048,
            memory_util: 0,
            gpu_util: 0,
            memory_used: 82,
            fan_speed: 35,
            gpu_temp: 31,
          },
        ],
      },
      {
        name: "Workstation 2",
        cpu: "Ryzen 5800X",
        motherboard: "Z790 AORUS XTREME X",
        notes: "",
        gpus: [
          {
            gpu_name: "NVIDIA TITAN Xp",
            gpu_brand: "Titan",
            driver_ver: "535.146.02",
            memory_total: 12288,
            memory_util: 0,
            gpu_util: 0,
            memory_used: 83,
            fan_speed: 23,
            gpu_temp: 32,
          },
          {
            gpu_name: "NVIDIA TITAN Xp",
            gpu_brand: "Titan",
            driver_ver: "535.146.02",
            memory_total: 12288,
            memory_util: 0,
            gpu_util: 0,
            memory_used: 83,
            fan_speed: 23,
            gpu_temp: 32,
          },
        ],
      },
      {
        name: "Workstation 3",
        cpu: "Intel Pentium 2",
        motherboard: "Acer Veriton M4630G MT",
        notes: "scheduled for replacement 2024",
        gpus: [
          {
            gpu_name: "NVIDIA GeForce GT 730",
            gpu_brand: "GeForce",
            driver_ver: "470.223.02",
            memory_total: 2001,
            memory_util: 0,
            gpu_util: 0,
            memory_used: 220,
            fan_speed: 30,
            gpu_temp: 27,
          },
        ],
      },
      {
        name: "Workstation 5",
        cpu: "Tesla DOJO",
        motherboard: "",
        notes: "We don't particularly like this one, but it always works and we can't really bin it",
        gpus: [
          {
            gpu_name: "NVIDIA TITAN Xp",
            gpu_brand: "Titan",
            driver_ver: "535.146.02",
            memory_total: 12288,
            memory_util: 0,
            gpu_util: 0,
            memory_used: 83,
            fan_speed: 23,
            gpu_temp: 32,
          },
          {
            gpu_name: "NVIDIA TITAN Xp",
            gpu_brand: "Titan",
            driver_ver: "535.146.02",
            memory_total: 12288,
            memory_util: 0,
            gpu_util: 0,
            memory_used: 83,
            fan_speed: 23,
            gpu_temp: 32,
          },
        ],
      },
      {
        name: "Workstation 4",
        cpu: "AMD Z1 Extreme",
        motherboard: "Ticket to Ride: Europe",
        notes: "Please don't use this unless you absolutely have to",
        gpus: [
          {
            gpu_name: "NVIDIA GeForce GT 1030",
            gpu_brand: "GeForce",
            driver_ver: "535.146.02",
            memory_total: 2048,
            memory_util: 0,
            gpu_util: 0,
            memory_used: 82,
            fan_speed: 35,
            gpu_temp: 31,
          },
        ],
      },

      {
        name: "Workstation 6",
        cpu: "AMD 4800S",
        motherboard: "AMD 4800S Desktop Kit",
        notes: "",
        gpus: [
          {
            gpu_name: "NVIDIA GeForce GT 730",
            gpu_brand: "GeForce",
            driver_ver: "470.223.02",
            memory_total: 2001,
            memory_util: 0,
            gpu_util: 0,
            memory_used: 220,
            fan_speed: 30,
            gpu_temp: 27,
          },
        ],
      },
    ],
  },
];
