import React from "react";
import { cleanup, render, screen } from "@testing-library/react";
import { default as App, GPUStats, REFRESH_INTERVAL } from "./App";

test("renders welcome message", () => {
  mockFetch(EXAMPLE_GPU_DATA_1);
  jest.useFakeTimers();

  render(<App />);
  const welcome = screen.getByText("Welcome to the GPU Control Room!");
  expect(welcome).toBeInTheDocument();
});

test(`before fetch succeeds inform the user that data is being fetched
after fetch succeeds, no longer show that message`, () => {
  mockFetch(EXAMPLE_GPU_DATA_1);
  jest.useFakeTimers();

  render(<App />);

  const status = screen.getByText("Retrieving data from API server...");
  expect(status).toBeInTheDocument();

  cleanup();
  render(<App />);

  const new_status = screen.getByText("Retrieving data from API server...");
  expect(new_status).not.toBeInTheDocument();
});

test("retrieves data from API server and displays correctly", () => {
  mockFetch(EXAMPLE_GPU_DATA_1);
  jest.useFakeTimers();

  render(<App />);
  cleanup();
  render(<App />);

  const gpu = screen.getByText("1030");
  expect(gpu).toBeInTheDocument();

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
  mockFetch(EXAMPLE_GPU_DATA_1);
  jest.useFakeTimers();
  jest.spyOn(global, "setInterval");

  const old_temp = screen.getByText("31 °C");
  expect(old_temp).toBeInTheDocument();

  data[0].gpu_temp = 100;

  jest.advanceTimersByTime(REFRESH_INTERVAL + 1);

  const new_temp = screen.getByText("100 °C");
  expect(new_temp).toBeInTheDocument();
});

const mockFetch = (s: GPUStats[]) => {
  global.fetch = jest.fn(() => {
    return {
      json: () => s,
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
