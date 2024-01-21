import { render, screen, waitFor } from "@testing-library/react";
import { default as App, GPUStats, REFRESH_INTERVAL } from "./App";

test("renders welcome message", () => {
  mockFetch(EXAMPLE_GPU_DATA_1);
  jest.useFakeTimers();

  render(<App />);
  const welcome = screen.getByText("Welcome to the GPU Control Room!");
  expect(welcome).toBeInTheDocument();
});

test(`before fetch succeeds inform the user that data is being fetched
after fetch succeeds, no longer show that message`, async () => {
  const fetch = mockFetch(EXAMPLE_GPU_DATA_1);
  jest.useFakeTimers();

  render(<App />);

  const status = screen.getByText("Retrieving data from API server...");
  expect(status).toBeInTheDocument();

  await waitFor(() => expect(fetch).toHaveBeenCalledTimes(1));

  const new_status = screen.queryByText("Retrieving data from API server...");
  expect(new_status).not.toBeInTheDocument();
});

test("retrieves data from API server and displays correctly", async () => {
  const fetch = mockFetch(EXAMPLE_GPU_DATA_1);
  jest.useFakeTimers();

  render(<App />);
  await waitFor(() => expect(fetch).toHaveBeenCalledTimes(1));

  EXAMPLE_GPU_DATA_1.forEach((row) => {
    screen
      .getAllByText(row.gpu_name, { exact: false })
      .forEach((name) => expect(name).toBeInTheDocument());

    screen
      .getAllByText(row.gpu_util + "%", { exact: false })
      .forEach((core_util) => expect(core_util).toBeInTheDocument());

    screen
      .getAllByText(row.memory_util + "%", {
        exact: false,
      })
      .forEach((mem_util) => expect(mem_util).toBeInTheDocument());

    screen
      .getAllByText(row.gpu_temp + " °C", { exact: false })
      .forEach((temp) => expect(temp).toBeInTheDocument());
  });
});

test("data is fetched again after refresh interval", async () => {
  var data = EXAMPLE_GPU_DATA_1;
  const fetch = mockFetch(data);
  jest.useFakeTimers();
  jest.spyOn(global, "setInterval");
  const view = render(<App />);

  await waitFor(() => expect(fetch).toHaveBeenCalledTimes(1));

  const old_temp = screen.getByText("31 °C", { exact: false });
  expect(old_temp).toBeInTheDocument();
  data[0].gpu_temp = 100;
  jest.advanceTimersByTime(REFRESH_INTERVAL + 1);

  await waitFor(() => expect(fetch).toHaveBeenCalledTimes(1));
  // It appears that it is necessary to request a rerender for every fetch after
  // the first. It is somewhat unclear to me whether it would be better practice
  // to rerender after every fetch (just in case)
  view.rerender(<App />);

  const new_temp = screen.getByText("100 °C", { exact: false });
  expect(new_temp).toBeInTheDocument();
});

const mockFetch = (s: GPUStats[]) => {
  const returnsJSON = jest.fn(() => Promise.resolve(s));
  global.fetch = jest.fn(() => {
    return Promise.resolve({
      json: returnsJSON,
    });
  }) as jest.Mock;
  return returnsJSON;
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
