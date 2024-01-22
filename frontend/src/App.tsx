import "./App.css";
import { WorkstationTab } from "./Components/WorkstationCard";
import { useJarJar, useOnce } from "./Utils/Hooks";
import { Validated, success, validationElim } from "./Utils/Utils";
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
        <Tabs variant='soft-rounded'>
          <TabList>
            <Tab>Simple View</Tab>
            <Tab>Table View</Tab>
            <Spacer />
            <Button mr={5}> Sign in </Button>
          </TabList>
          <Heading size="2xl">Welcome to the GPU Control Room!</Heading>
          <TabPanels>
            <TabPanel>
              <Center>
                <Box
                  w="100%"
                  m={10}
                  bg={useColorModeValue("gray.100", "gray.800")}
                >
                  <Stack direction={"column"} spacing={10}>
                    <Center>
                      <Box
                        w="100%"
                        m={10}
                        bg={useColorModeValue("gray.200", "gray.800")}
                      >
                        <Stack direction={"column"} spacing={5}>
                          {/* <Box></Box> */}
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
                            loading: () => (
                              <p>Retrieving data from API server...</p>
                            ),
                            failure: (_) => <p>Something has gone wrong!</p>,
                          })}
                        </Stack>
                      </Box>
                    </Center>
                    <Heading size="lg" textAlign="left">
                      Group 2: Shared
                    </Heading>
                    <Heading size="lg" textAlign="left">
                      Group 3: Remote
                    </Heading>
                  </Stack>
                </Box>
              </Center>
            </TabPanel>
            <TabPanel>
              <p>Tab le</p>
            </TabPanel>
          </TabPanels>
        </Tabs>
      </div>
    </ChakraProvider>
  );
}

export default App;
