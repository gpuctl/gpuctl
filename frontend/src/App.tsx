import "./App.css";
import { TableTab } from "./Components/DataTable";
import { useJarJar, useOnce } from "./Utils/Hooks";
import { Validated, success, validationElim } from "./Utils/Utils";
import { WorkstationCardMin } from "./Components/WorkstationCardMinimal";
import {
  Box,
  ChakraProvider,
  Heading,
  useColorModeValue,
  VStack,
} from "@chakra-ui/react";
import {} from "@chakra-ui/react";
import { Navbar } from "./Components/Navbar";
import { ColumnGrid } from "./Components/ColumnGrid";
import { WorkStationGroup } from "./Data";

const API_URL = "http://localhost:8000";
export const REFRESH_INTERVAL = 5000;

// Currently does not attempt to do any validation of the returned GPU stats,
// or indeed handle errors that might be thrown by the Promises
const retrieveAllStats: () => Promise<
  Validated<WorkStationGroup[]>
> = async () =>
  success(sortData(await (await fetch(API_URL + "/api/stats/all")).json()));
//success(foo);

const sortData = (gs: WorkStationGroup[]) =>
  gs.map(({ name, workStations }) => ({
    name: name,
    workStations: workStations.sort((ws1, ws2) =>
      ws1.name.localeCompare(ws2.name)
    ),
  }));

function App() {
  const [stats, updateStats] = useJarJar(retrieveAllStats);

  useOnce(() => {
    setInterval(updateStats, REFRESH_INTERVAL);
  });

  const bgcol = useColorModeValue("gray.100", "gray.800");

  return (
    <ChakraProvider>
      <div className="App"></div>
      <Navbar>
        {validationElim(stats, {
          success: (g) => (
            <VStack spacing={20}>
              {g.map((l, i) => (
                <Box
                  key={i}
                  w="97%"
                  m={5}
                  bg={bgcol}
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
                      <ColumnGrid
                        minChildWidth={325}
                        hMinSpacing={40}
                        vSpacing={10}
                      >
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
          ),
          loading: () => <p>Retrieving data from API server...</p>,
          failure: (_) => <p>Something has gone wrong!</p>,
        })}
        {validationElim(stats, {
          success: (l: WorkStationGroup[]) => <TableTab groups={l}></TableTab>,
          loading: () => <p>Retrieving data from API server...</p>,
          failure: (_) => <p>Something has gone wrong!</p>,
        })}
      </Navbar>
    </ChakraProvider>
  );
}

export default App;
