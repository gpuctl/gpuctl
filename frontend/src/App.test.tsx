import React from "react";
import { render, screen } from "@testing-library/react";
import { default as App, GPUStats, REFRESH_INTERVAL } from "./App";

test("renders welcome message", () => {
  mockFetch(EXAMPLE_GPU_DATA_1);
  jest.useFakeTimers();

  render(<App />);
  const welcome = screen.getByText(/Welcome to the GPU Control Room!/i);
  expect(welcome).toBeInTheDocument();
});

test("should fail", () => {
  mockFetch(EXAMPLE_GPU_DATA_1);
  const name = screen.getByText(/Retrieving data.../);
  expect(name).toBeInTheDocument();
});

test("retrieves data from API server and displays correctly", () => {
  mockFetch(EXAMPLE_GPU_DATA_1);
  jest.useFakeTimers();
  EXAMPLE_GPU_DATA_1.forEach((row) => {
    const name = screen.getByText(row.gpu_name);
    expect(name).toBeInTheDocument();

    const core_util = screen.getByText(row.gpu_util + "%");
    expect(core_util).toBeInTheDocument();

    const memory_util = screen.getByText(row.memory_util + "%");
    expect(memory_util).toBeInTheDocument();

    const temp = screen.getByText(row.gpu_temp + " °C");
    expect(temp).toBeInTheDocument();
  });
});

test("data is fetched again after refresh interval", () => {
  var data = EXAMPLE_GPU_DATA_1;
  mockFetch(data);
  jest.useFakeTimers();
  jest.spyOn(global, "setInterval");

  const old_temp = screen.getByText("31 °C");
  expect(old_temp).toBeInTheDocument();

  data[0].gpu_temp = 100;

  jest.advanceTimersByTime(6000);

  const new_temp = screen.getByText("100 °C");
  expect(new_temp).toBeInTheDocument();
});

const mockFetch = (s: GPUStats[]) => {
  global.fetch = jest.fn(() => {
    return {
      json: () => Promise.resolve(s),
    };
  }) as jest.Mock;
};

const EXAMPLE_GPU_DATA_1 = [
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
];
