import { Box, Heading, VStack } from "@chakra-ui/react";
import { API_URL, DEFAULT_VIEW, REFRESH_INTERVAL, ViewPage } from "../App";
import { DurationDeltas, WorkStationGroup, WorkStationData } from "../Data";
import { Validated, Validation, success, validationElim } from "../Utils/Utils";
import { ColumnGrid } from "../Components/ColumnGrid";
import { TableTab } from "../Components/DataTable";
import { WorkstationCardMin } from "../Components/WorkstationCardMinimal";
import { Navbar } from "../Components/Navbar";
import { useJarJar, useOnce } from "../Utils/Hooks";
import { useAuth } from "../Providers/AuthProvider";
import { STATS_PATH } from "../Config/Paths";
import { AdminPanel } from "./AdminPanel";
import { Navigate } from "react-router-dom";
import { partition } from "lodash";
const API_ALL_STATS_PATH = "/stats/all";
const API_LAST_SEEN_PATH = "/stats/since_last_seen";

// Currently does not attempt to do any validation of the returned GPU stats,
// or indeed handle errors that might be thrown by the Promises
const retrieveAllStats: () => Promise<
  Validated<WorkStationGroup[]>
> = async () =>
  success(preProcess(await (await fetch(API_URL + API_ALL_STATS_PATH)).json()));

// We will consider a machine to be in use if any of it's GPUs are
const inUse = (machine: WorkStationData) =>
  machine.gpus.some(({ in_use }) => in_use);

const sortData = <T,>(ws: (WorkStationData & T)[]) =>
  ws
    .map(({ name, gpus, ...rest }) => ({
      name,
      gpus: gpus.sort((g1, g2) => g1.uuid.localeCompare(g2.uuid)),
      ...rest,
    }))
    .sort((ws1, ws2) => ws1.name.localeCompare(ws2.name));

const tagFree = (ws: WorkStationData[], free: boolean) =>
  ws.map((data) => ({
    free,
    ...data,
  }));

const retrieveLastSeen: () => Promise<Validated<DurationDeltas[]>> = async () =>
  success(await (await fetch(API_URL + API_LAST_SEEN_PATH)).json());

const preProcess = (
  gs: WorkStationGroup[],
): WorkStationGroup<{ free: boolean }>[] =>
  gs.map(({ name, workstations }) => {
    const [used, free] = partition(sortData(workstations), inUse);

    return {
      name,
      workstations: tagFree(free, true).concat(tagFree(used, false)),
    };
  });

const cardView = (
  stats: WorkStationGroup[],
  lastSeen: Record<string, number>,
) => (
  <VStack spacing={20}>
    {stats.map((l, i) => (
      <Box
        key={i}
        w="97%"
        m={5}
        bg={"gray.100"}
        paddingTop={5}
        paddingBottom={10}
      >
        <VStack spacing={5}>
          <Box w="100%" paddingLeft={5}>
            <Heading size="lg" textAlign="left">
              Group {i + 1}: {l.name}
            </Heading>
          </Box>
          <Box w="100%">
            <ColumnGrid minChildWidth={325} hMinSpacing={40} vSpacing={10}>
              {l.workstations.map((row, i) => {
                return (
                  <WorkstationCardMin
                    key={i}
                    width={325}
                    data={{ lastSeen: lastSeen[row.name], ...row }}
                  ></WorkstationCardMin>
                );
              })}
            </ColumnGrid>
          </Box>
        </VStack>
      </Box>
    ))}
  </VStack>
);

const tableView = (stats: WorkStationGroup[]) => (
  <TableTab groups={stats}></TableTab>
);

const adminView = (stats: WorkStationGroup[]) => (
  <AdminPanel groups={stats}></AdminPanel>
);

const displayPartial = (
  stats: Validation<WorkStationGroup[]>,
  lastSeen: Validation<DurationDeltas[]>,
  cont: (gs: WorkStationGroup[], ls: Record<string, number>) => JSX.Element,
): JSX.Element => {
  const ls = Object.fromEntries(
    validationElim(lastSeen, {
      success: (ls) => ls,
      failure: () => [],
      loading: () => [],
    }).map(({ hostname, seconds_since }) => [hostname, seconds_since]),
  );
  return validationElim(stats, {
    success: (gs: WorkStationGroup[]) => cont(gs, ls),
    loading: () => <p>Retrieving data from API server...</p>,
    failure: (_) => <p>Something has gone wrong!</p>,
  });
};

export const MainView = ({ page }: { page: ViewPage }) => {
  const { isSignedIn } = useAuth();
  if (page === ViewPage.ADMIN && !isSignedIn()) {
    return <Navigate to={STATS_PATH + DEFAULT_VIEW} replace />;
  }
  return <ConfirmedMainView initial={page} />;
};

export const ConfirmedMainView = ({ initial }: { initial: ViewPage }) => {
  const [stats, updateStats] = useJarJar(retrieveAllStats);
  const [lastSeen, updateLastSeen] = useJarJar(retrieveLastSeen);

  useOnce(() => {
    setInterval(updateStats, REFRESH_INTERVAL);
  });

  useOnce(() => {
    setInterval(updateLastSeen, REFRESH_INTERVAL);
  });

  return (
    <Navbar initial={initial}>
      {displayPartial(stats, lastSeen, cardView)}
      {displayPartial(stats, lastSeen, tableView)}
      {displayPartial(stats, lastSeen, adminView)}
    </Navbar>
  );
};
