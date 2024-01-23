import "./App.css";
import { useJarJar, useOnce } from "./Utils/Hooks";
import { Validated, success, validationElim } from "./Utils/Utils";
import { WorkstationCardNew } from "./Components/WorkstationCardNew";
import {
  Box,
  Button,
  Center,
  ChakraProvider,
  Flex,
  Grid,
  Heading,
  SimpleGrid,
  Spacer,
  Stack,
  useColorModeValue,
  Tabs,
  TabList,
  TabPanels,
  Tab,
  TabPanel,
} from "@chakra-ui/react";
import {} from "@chakra-ui/react";
import { WorkstationCard } from "./Components/WorkstationCard";
import { Navbar } from "./Components/Navbar";

const API_URL = "http://localhost:8000";
export const REFRESH_INTERVAL = 5000;

export type WorkStationGroup = {
  name: string;
  workStations: WorkStationData[];
};

export type WorkStationData = {
  name: string;
  gpus: GPUStats[];
};

export type GPUStats = {
  gpu_name: string;
  gpu_brand: string;
  driver_ver: string;
  memory_total: number;

  memory_util: number;
  gpu_util: number;
  memory_used: number;
  fan_speed: number;
  gpu_temp: number;
};

// Currently does not attempt to do any validation of the returned GPU stats,
// or indeed handle errors that might be thrown by the Promises
const retrieveAllStats: () => Promise<Validated<GPUStats[]>> = async () =>
  success(await (await fetch(API_URL + "/api/stats/all")).json());

function App() {
  const [stats, updateStats] = useJarJar(retrieveAllStats);

  useOnce(() => {
    setInterval(updateStats, REFRESH_INTERVAL);
  });

  return (
    <ChakraProvider>
      <div className="App"></div>
      <Navbar>
        <Center>
          <Box w="100%" m={10} bg={useColorModeValue("gray.100", "gray.800")}>
            <Stack direction={"column"} spacing={5}>
              <Heading size="lg" textAlign="left">
                Group 1: Personal
              </Heading>
              {validationElim(stats, {
                success: (l) => (
                  <Center>
                    <Box w="100%" m={10}>
                      <SimpleGrid minChildWidth={350} spacing={10}>
                        {l.map((row, i) => {
                          return (
                            <WorkstationCard
                              key={i}
                              name={`Workstation ${i}`}
                              gpus={[row]}
                            ></WorkstationCard>
                          );
                        })}
                      </SimpleGrid>
                    </Box>
                  </Center>
                  /* <WorkstationCardNew
                          key={0}
                          name="Workstation 0"
                          gpus={[l[0]]}
                        ></WorkstationCardNew>
                        */
                ),
                loading: () => <p>Retrieving data from API server...</p>,
                failure: (_) => <p>Something has gone wrong!</p>,
              })}
            </Stack>
          </Box>
        </Center>
        <p>Tab le</p>
      </Navbar>
    </ChakraProvider>
  );
}

export default App;
