import "./App.css";
import { WorkstationTab } from "./Components/WorkstationCard";
import { useJarJar, useOnce } from "./Utils/Hooks";
import { Validated, success, validationElim } from "./Utils/Utils";
import {
  Box,
  Center,
  ChakraProvider,
  Flex,
  Grid,
  Heading,
  SimpleGrid,
  Stack,
  useColorModeValue,
} from "@chakra-ui/react";

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
      <div className="App">
        <Stack direction={"column"} spacing={30}>
          <Heading size="2xl">Welcome to the GPU Control Room!</Heading>
          <Center>
            <Box w="95%" bg={useColorModeValue("gray.200", "gray.800")}>
              <Stack direction={"column"} spacing={10}>
                <Stack direction={"column"} spacing={5}>
                  <Heading size="xl">Group 1: Personal</Heading>
                  {validationElim(stats, {
                    success: (l) => (
                      <Center>
                        <Box w="95%">
                          <SimpleGrid minChildWidth={300} spacing={20}>
                            {l.map((row, i) => {
                              return (
                                <WorkstationTab
                                  key={i}
                                  name={`Workstation ${i}`}
                                  gpus={[row]}
                                ></WorkstationTab>
                              );
                              /*(
                    <p key={i}>
                      ID: {i}, Name: {row.gpu_name}, Core Utilisation:{" "}
                      {row.gpu_util}
                      %, VRAM Util: {row.memory_util}%, VRAM: {row.memory_total}{" "}
                      GB, Used VRAM: {row.memory_used} GB, Temperature:{" "}
                      {row.gpu_temp} °C
                    </p>
                  );*/
                            })}
                          </SimpleGrid>
                        </Box>
                      </Center>
                    ),
                    loading: () => <p>Retrieving data from API server...</p>,
                    failure: (_) => <p>Something has gone wrong!</p>,
                  })}
                </Stack>
                <Heading>Group 2: Shared</Heading>
                <Heading>Group 3: Remote</Heading>
              </Stack>
            </Box>
          </Center>
        </Stack>
      </div>
    </ChakraProvider>
  );
}

export default App;
