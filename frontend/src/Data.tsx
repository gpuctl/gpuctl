// These types need to be kept in sync with `internal/webapi/types.go`

export type WorkStationGroup<E = {}> = {
  name: string;
  workstations: (WorkStationData & E)[];
};

export type WorkStationData = {
  name: string;
  cpu: string;
  motherboard: string;
  notes: string;
  owner: string;
  gpus: GPUStats[];
  last_seen: number;
};

// This needs to be kept in sync with `internal/uplink/
export type GPUStats = {
  uuid: string;
  gpu_name: string;
  gpu_brand: string;
  driver_ver: string;
  memory_total: number;

  memory_util: number;
  gpu_util: number;
  memory_used: number;
  fan_speed: number;
  gpu_temp: number;
  memory_temp: number;
  graphics_voltage: number;
  power_draw: number;
  graphics_clock: number;
  max_graphics_clock: number;
  memory_clock: number;
  max_memory_clock: number;
  in_use: boolean;
  user: string;
};

export enum GraphField {
  MEMORY_UTIL = "Memory Utilisation (%)",
  GPU_UTIL = "GPU Utilisation (%)",
  MEMORY_USED = "Memory Used (MB)",
  FAN_SPEED = "Fan Speed (%)",
  GPU_TEMP = "GPU Temperature (°C)",
  MEMORY_TEMP = "Memory Temperature (°C)",
  GRAPHICS_VOLTAGE = "GPU Voltage (mV)",
  POWER_DRAW = "Power Draw (W)",
  GRAPHICS_CLOCK = "GPU Clock (MHz)",
  MEMORY_CLOCK = "Memory Clock (MHz)",
}

export type DurationDeltas = {
  hostname: string;
  seconds_since: number;
};

export const EXAMPLE_DATA_1: WorkStationGroup[] = [
  {
    name: "Shared",
    workstations: [
      {
        name: "Workstation 1",
        cpu: "i7-7700k",
        motherboard: "ASUS ROG STRIX Z790-A",
        notes: "noisy fan",
        owner: "",
        last_seen: 500,
        gpus: [
          {
            uuid: "AAAAA",
            gpu_name: "NVIDIA GeForce GT 1030",
            gpu_brand: "GeForce",
            driver_ver: "535.146.02",
            memory_total: 2048,
            memory_util: 0,
            gpu_util: 0,
            memory_used: 82,
            fan_speed: 35,
            gpu_temp: 31,
            memory_temp: 1,
            graphics_voltage: 2,
            power_draw: 3,
            graphics_clock: 4,
            max_graphics_clock: 5,
            memory_clock: 6,
            max_memory_clock: 7,
            in_use: false,
            user: "",
          },
        ],
      },
      {
        name: "Workstation 2",
        cpu: "Ryzen 5800X",
        motherboard: "Z790 AORUS XTREME X",
        notes: "",
        owner: "Bob Barker",
        last_seen: 260,
        gpus: [
          {
            uuid: "BBBBB",
            gpu_name: "NVIDIA TITAN Xp",
            gpu_brand: "Titan",
            driver_ver: "535.146.02",
            memory_total: 12288,
            memory_util: 0,
            gpu_util: 0,
            memory_used: 83,
            fan_speed: 23,
            gpu_temp: 32,
            memory_temp: 8,
            graphics_voltage: 9,
            power_draw: 10,
            graphics_clock: 11,
            max_graphics_clock: 12,
            memory_clock: 13,
            max_memory_clock: 14,
            in_use: false,
            user: "",
          },
          {
            uuid: "CCCCC",
            gpu_name: "NVIDIA TITAN Xp",
            gpu_brand: "Titan",
            driver_ver: "535.146.02",
            memory_total: 12288,
            memory_util: 0,
            gpu_util: 0,
            memory_used: 83,
            fan_speed: 23,
            gpu_temp: 32,
            memory_temp: 15,
            graphics_voltage: 16,
            power_draw: 17,
            graphics_clock: 18,
            max_graphics_clock: 19,
            memory_clock: 20,
            max_memory_clock: 21,
            in_use: true,
            user: "",
          },
        ],
      },
      {
        name: "Workstation 3",
        cpu: "Intel Pentium 2",
        motherboard: "Acer Veriton M4630G MT",
        notes: "scheduled for replacement 2024",
        owner: "Kermit",
        last_seen: 30000000000,
        gpus: [
          {
            uuid: "DDDDD",
            gpu_name: "NVIDIA GeForce GT 730",
            gpu_brand: "GeForce",
            driver_ver: "470.223.02",
            memory_total: 2001,
            memory_util: 0,
            gpu_util: 0,
            memory_used: 220,
            fan_speed: 30,
            gpu_temp: 27,
            memory_temp: 22,
            graphics_voltage: 23,
            power_draw: 24,
            graphics_clock: 25,
            max_graphics_clock: 26,
            memory_clock: 27,
            max_memory_clock: 28,
            in_use: true,
            user: "",
          },
        ],
      },
      {
        name: "Workstation 5",
        cpu: "Tesla DOJO",
        motherboard: "",
        notes:
          "We don't particularly like this one, but it always works and we can't really bin it",
        owner: "",
        last_seen: 90000000000,
        gpus: [
          {
            uuid: "EEEEE",
            gpu_name: "NVIDIA TITAN Xp",
            gpu_brand: "Titan",
            driver_ver: "535.146.02",
            memory_total: 12288,
            memory_util: 0,
            gpu_util: 0,
            memory_used: 83,
            fan_speed: 23,
            gpu_temp: 32,
            memory_temp: 29,
            graphics_voltage: 30,
            power_draw: 31,
            graphics_clock: 32,
            max_graphics_clock: 33,
            memory_clock: 34,
            max_memory_clock: 35,
            in_use: false,
            user: "",
          },
          {
            uuid: "FFFFF",
            gpu_name: "NVIDIA TITAN Xp",
            gpu_brand: "Titan",
            driver_ver: "535.146.02",
            memory_total: 12288,
            memory_util: 0,
            gpu_util: 0,
            memory_used: 83,
            fan_speed: 23,
            gpu_temp: 32,
            memory_temp: 36,
            graphics_voltage: 37,
            power_draw: 38,
            graphics_clock: 39,
            max_graphics_clock: 40,
            memory_clock: 41,
            max_memory_clock: 42,
            in_use: false,
            user: "",
          },
        ],
      },
      {
        name: "Workstation 4",
        cpu: "AMD Z1 Extreme",
        motherboard: "Ticket to Ride: Europe",
        notes: "Please don't use this unless you absolutely have to",
        owner: "Mort",
        last_seen: 40000000000,
        gpus: [
          {
            uuid: "GGGGG",
            gpu_name: "NVIDIA GeForce GT 1030",
            gpu_brand: "GeForce",
            driver_ver: "535.146.02",
            memory_total: 2048,
            memory_util: 0,
            gpu_util: 0,
            memory_used: 82,
            fan_speed: 35,
            gpu_temp: 31,
            memory_temp: 43,
            graphics_voltage: 44,
            power_draw: 45,
            graphics_clock: 46,
            max_graphics_clock: 47,
            memory_clock: 48,
            max_memory_clock: 49,
            in_use: false,
            user: "",
          },
        ],
      },
    ],
  },
  {
    name: "Other",
    workstations: [
      {
        name: "Workstation 6",
        cpu: "AMD 4800S",
        motherboard: "AMD 4800S Desktop Kit",
        notes: "",
        owner: "",
        last_seen: 180000000123,
        gpus: [
          {
            uuid: "HHHHH",
            gpu_name: "NVIDIA GeForce GT 730",
            gpu_brand: "GeForce",
            driver_ver: "470.223.02",
            memory_total: 2001,
            memory_util: 0,
            gpu_util: 0,
            memory_used: 220,
            fan_speed: 30,
            gpu_temp: 27,
            memory_temp: 50,
            graphics_voltage: 51,
            power_draw: 52,
            graphics_clock: 53,
            max_graphics_clock: 54,
            memory_clock: 55,
            max_memory_clock: 56,
            in_use: true,
            user: "",
          },
        ],
      },
    ],
  },
];
