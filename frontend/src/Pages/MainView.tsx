import { Box, Heading, VStack } from "@chakra-ui/react";
import {
  API_URL,
  AUTH_TOKEN_ITEM,
  AuthToken,
  REFRESH_INTERVAL,
  ViewPage,
} from "../App";
import { WorkStationGroup } from "../Data";
import {
  Validated,
  Validation,
  failure,
  success,
  validatedElim,
  validationElim,
} from "../Utils/Utils";
import { ColumnGrid } from "../Components/ColumnGrid";
import { TableTab } from "../Components/DataTable";
import { WorkstationCardMin } from "../Components/WorkstationCardMinimal";
import { Navbar } from "../Components/Navbar";
import { useJarJar, useOnce } from "../Utils/Hooks";
import { ADMIN_PATH, AdminPanel } from "./AdminPanel";

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

const adminView = (stats: WorkStationGroup[], token: AuthToken) => (
  <AdminPanel groups={stats} token={token}></AdminPanel>
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

export const MainView = (props: {
  default: ViewPage;
  authToken: Validated<AuthToken>;
  setAuth: (tok: Validated<AuthToken>) => void;
}) => {
  const [stats, updateStats] = useJarJar(retrieveAllStats);

  useOnce(() => {
    setInterval(updateStats, REFRESH_INTERVAL);
  });

  const signOut = () => {
    props.setAuth(failure(Error("Signed Out")));
    localStorage.removeItem(AUTH_TOKEN_ITEM);
    window.location.reload();
  };

  const signIn = (tok: AuthToken) => {
    localStorage.setItem(AUTH_TOKEN_ITEM, tok.token);
    props.setAuth(success(tok));
    window.location.reload();
  };

  const checkStillAdmin = async (tok: AuthToken) => {
    const resp = await fetch(API_URL + ADMIN_PATH + "/confirm", {
      method: "GET",
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${tok.token}`,
      },
    });
    if (resp.ok) return;
    if (resp.status === 401) signOut();
    console.log("Unknown error occured when checking admin status!");
  };

  validatedElim(props.authToken, {
    success: (tok) => {
      // Just fire the check - it's not super important but doing this
      // every so often is probably good practice
      checkStillAdmin(tok);
    },
    failure: () => {},
  });

  return (
    <Navbar
      initial={props.default}
      authToken={props.authToken}
      signOut={signOut}
      signIn={signIn}
    >
      {displayPartial(stats, cardView)}
      {displayPartial(stats, tableView)}
      {validatedElim(props.authToken, {
        success: (tok) => displayPartial(stats, (s) => adminView(s, tok)),
        failure: () => <></>,
      })}
    </Navbar>
  );
};
