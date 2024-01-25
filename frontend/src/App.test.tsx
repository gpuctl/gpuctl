import { render, screen, waitFor } from "@testing-library/react";
import { default as App, REFRESH_INTERVAL } from "./App";
import { EXAMPLE_DATA_1, WorkStationGroup } from "./Data";

test("renders welcome message", () => {
  mockFetch(EXAMPLE_DATA_1);
  jest.useFakeTimers();

  render(<App />);
  const welcome = screen.getByText("Welcome to the GPU Control Room!");
  expect(welcome).toBeInTheDocument();
});

test(`before fetch succeeds inform the user that data is being fetched
after fetch succeeds, no longer show that message`, async () => {
  const fetch = mockFetch(EXAMPLE_DATA_1);
  jest.useFakeTimers();

  render(<App />);

  const status = screen.getByText("Retrieving data from API server...");
  expect(status).toBeInTheDocument();

  await waitFor(() => expect(fetch).toHaveBeenCalledTimes(1));

  const new_status = screen.queryByText("Retrieving data from API server...");
  expect(new_status).not.toBeInTheDocument();
});

test("retrieves data from API server and displays correctly", async () => {
  const fetch = mockFetch(EXAMPLE_DATA_1);
  jest.useFakeTimers();

  render(<App />);
  await waitFor(() => expect(fetch).toHaveBeenCalledTimes(1));

  EXAMPLE_DATA_1.forEach((g) => {
    g.workStations.forEach((ws) => {
      ws.gpus.forEach((gpu) => {
        screen
          .getAllByText(gpu.gpu_name, { exact: false })
          .forEach((name) => expect(name).toBeInTheDocument());

        screen
          .getAllByText(gpu.gpu_util + "%", { exact: false })
          .forEach((core_util) => expect(core_util).toBeInTheDocument());

        screen
          .getAllByText(gpu.gpu_temp + " °C", { exact: false })
          .forEach((temp) => expect(temp).toBeInTheDocument());
      });
    });
  });
});

test("data is fetched again after refresh interval", async () => {
  var data = EXAMPLE_DATA_1;
  const fetch = mockFetch(data);
  jest.useFakeTimers();
  jest.spyOn(global, "setInterval");
  const view = render(<App />);

  await waitFor(() => expect(fetch).toHaveBeenCalledTimes(1));

  screen
    .getAllByText("31 °C", { exact: false })
    .forEach((temp) => expect(temp).toBeInTheDocument());
  data[0].workStations[0].gpus[0].gpu_temp = 100;
  jest.advanceTimersByTime(REFRESH_INTERVAL + 1);

  await waitFor(() => expect(fetch).toHaveBeenCalledTimes(1));
  // It appears that it is necessary to request a rerender for every fetch after
  // the first. It is somewhat unclear to me whether it would be better practice
  // to rerender after every fetch (just in case)
  view.rerender(<App />);

  const new_temp = screen.getByText("100 °C", { exact: false });
  expect(new_temp).toBeInTheDocument();
});

const mockFetch = (s: WorkStationGroup[]) => {
  const returnsJSON = jest.fn(() => Promise.resolve(s));
  global.fetch = jest.fn(() => {
    return Promise.resolve({
      json: returnsJSON,
    });
  }) as jest.Mock;
  return returnsJSON;
};
