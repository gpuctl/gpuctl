
// These types need to be kept in sync with `internal/webapi/types.go`

export type WorkStationGroup = {
  name: string;
  // TODO: Consistant cases here
  workStations: WorkStationData[];
};

export type WorkStationData = {
  name: string;
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
