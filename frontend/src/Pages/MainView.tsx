import { Box, Heading, VStack } from "@chakra-ui/react";
import { API_URL, DEFAULT_VIEW, REFRESH_INTERVAL, ViewPage } from "../App";
import { WorkStationGroup } from "../Data";
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

const API_ALL_STATS_PATH = "/stats/all";

// Currently does not attempt to do any validation of the returned GPU stats,
// or indeed handle errors that might be thrown by the Promises
const retrieveAllStats: () => Promise<
  Validated<WorkStationGroup[]>
> = async () =>
  success(sortData(await (await fetch(API_URL + API_ALL_STATS_PATH)).json()));
//success(foo);

const sortData = (gs: WorkStationGroup[]) =>
  gs.map(({ name, workStations }) => ({
    name: name,
    workStations: workStations.sort((ws1, ws2) =>
      ws1.name.localeCompare(ws2.name),
    ),
  }));

const cardView = (stats: WorkStationGroup[]) => (
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
              {l.workStations.map((row, i) => {
                return (
                  <WorkstationCardMin
                    key={i}
                    width={325}
                    data={row}
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
  cont: (gs: WorkStationGroup[]) => JSX.Element,
): JSX.Element =>
  validationElim(stats, {
    success: (l: WorkStationGroup[]) => cont(l),
    loading: () => <p>Retrieving data from API server...</p>,
    failure: (_) => <p>Something has gone wrong!</p>,
  });

export const MainView = ({ page }: { page: ViewPage }) => {
  const { isSignedIn } = useAuth();
  if (page === ViewPage.ADMIN && !isSignedIn()) {
    return <Navigate to={STATS_PATH + DEFAULT_VIEW} replace />;
  }
  return <ConfirmedMainView initial={page} />;
};

export const ConfirmedMainView = ({ initial }: { initial: ViewPage }) => {
  const [stats, updateStats] = useJarJar(retrieveAllStats);

  useOnce(() => {
    setInterval(updateStats, REFRESH_INTERVAL);
  });

  return (
    <Navbar initial={initial}>
      {displayPartial(stats, cardView)}
      {displayPartial(stats, tableView)}
      {displayPartial(stats, (s) => adminView(s))}
    </Navbar>
  );
};
