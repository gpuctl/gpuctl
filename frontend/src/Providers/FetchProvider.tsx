import { API_URL, REFRESH_INTERVAL } from "../App";
import { WorkStationData, WorkStationGroup } from "../Data";
import { partition } from "lodash";
import { Validated, Validation, failure, success } from "../Utils/Utils";
import { ReactNode, createContext, useContext, useState } from "react";
import { useJarJar } from "../Utils/Hooks";
import { useInterval } from "@chakra-ui/react";

const API_ALL_STATS_PATH = "/stats/all";

type FetchStatsCtx = {
  allStats: Validation<WorkStationGroup[]>;
  setShouldFetch: (b: boolean) => void;
};

const FetchStatsContext = createContext<FetchStatsCtx>({
  allStats: failure(Error("No fetch stats provider!")),
  setShouldFetch: () => {},
});

export const useStats = () => useContext(FetchStatsContext);

export const FetchStatsProvider = ({
  children,
}: {
  children: ReactNode[] | ReactNode;
}) => {
  const [shouldFetch, setShouldFetch] = useState<boolean>(true);
  const [stats, updateStats] = useJarJar(retrieveAllStats);
  useInterval(() => {
    if (shouldFetch) {
      updateStats();
    }
  }, REFRESH_INTERVAL);

  return (
    <FetchStatsContext.Provider
      value={{
        allStats: stats,
        setShouldFetch,
      }}
    >
      {children}
    </FetchStatsContext.Provider>
  );
};

// Currently does not attempt to do any validation of the returned GPU stats,
// or indeed handle errors that might be thrown by the Promises
const retrieveAllStats: () => Promise<
  Validated<WorkStationGroup[]>
> = async () =>
  success(preProcess(await (await fetch(API_URL + API_ALL_STATS_PATH)).json()));
// USEFUL FOR TESTING, DON'T DELETE PLS
// success(preProcess(EXAMPLE_DATA_1));

// We will consider a machine to be in use if any of it's GPUs are
const inUse = (machine: WorkStationData) =>
  machine.gpus.some(({ in_use }) => in_use);

const sortData = <T,>(ws: (WorkStationData & T)[]) =>
  ws
    .map(({ name, gpus, last_seen, ...rest }) => ({
      name,
      last_seen: last_seen / 1_000_000_000,
      gpus: gpus.sort((g1, g2) => g1.uuid.localeCompare(g2.uuid)),
      ...rest,
    }))
    .sort((ws1, ws2) => ws1.name.localeCompare(ws2.name));

const tagFree = (ws: WorkStationData[], free: boolean) =>
  ws.map((data) => ({
    free,
    ...data,
  }));

const preProcess = (
  gs: WorkStationGroup[],
): WorkStationGroup<{ free: boolean }>[] =>
  gs
    .map(({ name, workstations }) => {
      const [used, free] = partition(sortData(workstations), inUse);

      return {
        name,
        workstations: tagFree(free, true).concat(tagFree(used, false)),
      };
    })
    .sort((g1, g2) =>
      g1.name === "Shared"
        ? -1
        : g2.name === "Shared"
          ? 1
          : g1.name.localeCompare(g2.name),
    );
