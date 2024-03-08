import { API_URL, REFRESH_INTERVAL } from "../App";
import { EXAMPLE_DATA_1, WorkStationData, WorkStationGroup } from "../Data";
import { partition } from "lodash";
import { Validated, Validation, failure, success } from "../Utils/Utils";
import { ReactNode, createContext, useContext, useState } from "react";
import { useJarJar } from "../Utils/Hooks";
import { useInterval } from "@chakra-ui/react";
import { workstationBusy } from "../Components/WorkstationCardMinimal";

const API_ALL_STATS_PATH = "/stats/all";

type StatsCB = (s: Validated<WorkStationGroup[]>) => boolean;

type FetchStatsCtx = {
  allStats: Validation<WorkStationGroup[]>;
  setShouldFetch: (b: boolean) => void;
  onNextFetch: (f: StatsCB) => void;
};

const FetchStatsContext = createContext<FetchStatsCtx>({
  allStats: failure(Error("No fetch stats provider!")),
  setShouldFetch: () => {},
  onNextFetch: () => {},
});

export const useStats = () => useContext(FetchStatsContext);

export const FetchStatsProvider = ({
  children,
}: {
  children: ReactNode[] | ReactNode;
}) => {
  const [shouldFetch, setShouldFetch] = useState<boolean>(true);
  const [fetchCallback, setFetchCallback] = useState<StatsCB>(() => () => true);
  const [stats, updateStats] = useJarJar(() => {
    const cb = fetchCallback;
    setFetchCallback(() => () => true);
    return retrieveAllStats(cb);
  });

  useInterval(() => {
    if (shouldFetch) {
      updateStats();
    }
  }, REFRESH_INTERVAL);

  const mergeCallback = (f: StatsCB) => {
    setFetchCallback(() => (stats: Validated<WorkStationGroup[]>) => {
      fetchCallback(stats);
      if (!f(stats)) {
        mergeCallback(f);
      }
      return true;
    });
  };

  return (
    <FetchStatsContext.Provider
      value={{
        allStats: stats,
        setShouldFetch,
        onNextFetch: mergeCallback,
      }}
    >
      {children}
    </FetchStatsContext.Provider>
  );
};

// Currently does not attempt to do any validation of the returned GPU stats,
// or indeed handle errors that might be thrown by the Promises
const retrieveAllStats: (
  cb: StatsCB,
) => Promise<Validated<WorkStationGroup[]>> = async (cb) => {
  const stats = success(preProcess(EXAMPLE_DATA_1));
  // const stats = success(
  //   preProcess(await (await fetch(API_URL + API_ALL_STATS_PATH)).json()),
  // );
  cb(stats);
  return stats;
};
// USEFUL FOR TESTING, DON'T DELETE PLS
// success(preProcess(EXAMPLE_DATA_1));

const inUse = (machine: WorkStationData) => workstationBusy(machine.gpus);

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
