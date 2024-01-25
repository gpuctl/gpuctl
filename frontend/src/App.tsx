import "./App.css";
import { useJarJar, useOnce } from "./Utils/Hooks";
import { Validated, success, validationElim } from "./Utils/Utils";
import { WorkstationCardMin } from "./Components/WorkstationCardMinimal";
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
  VStack,
} from "@chakra-ui/react";
import {} from "@chakra-ui/react";
import { WorkstationCard } from "./Components/WorkstationCard";
import { Navbar } from "./Components/Navbar";
import { ColumnGrid } from "./Components/ColumnGrid";

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

const foo: WorkStationGroup[] = [
  {
    name: "Shared",
    workStations: [
      {
        name: "Workstation 1",
        gpus: [
          {
            gpu_name: "NVIDIA GeForce GT 1030",
            gpu_brand: "GeForce",
            driver_ver: "535.146.02",
            memory_total: 2048,
            memory_util: 0,
            gpu_util: 0,
            memory_used: 82,
            fan_speed: 35,
            gpu_temp: 31,
          },
        ],
      },
      {
        name: "Workstation 2",
        gpus: [
          {
            gpu_name: "NVIDIA TITAN Xp",
            gpu_brand: "Titan",
            driver_ver: "535.146.02",
            memory_total: 12288,
            memory_util: 0,
            gpu_util: 0,
            memory_used: 83,
            fan_speed: 23,
            gpu_temp: 32,
          },
          {
            gpu_name: "NVIDIA TITAN Xp",
            gpu_brand: "Titan",
            driver_ver: "535.146.02",
            memory_total: 12288,
            memory_util: 0,
            gpu_util: 0,
            memory_used: 83,
            fan_speed: 23,
            gpu_temp: 32,
          },
        ],
      },
      {
        name: "Workstation 3",
        gpus: [
          {
            gpu_name: "NVIDIA GeForce GT 730",
            gpu_brand: "GeForce",
            driver_ver: "470.223.02",
            memory_total: 2001,
            memory_util: 0,
            gpu_util: 0,
            memory_used: 220,
            fan_speed: 30,
            gpu_temp: 27,
          },
        ],
      },
      {
        name: "Workstation 5",
        gpus: [
          {
            gpu_name: "NVIDIA TITAN Xp",
            gpu_brand: "Titan",
            driver_ver: "535.146.02",
            memory_total: 12288,
            memory_util: 0,
            gpu_util: 0,
            memory_used: 83,
            fan_speed: 23,
            gpu_temp: 32,
          },
          {
            gpu_name: "NVIDIA TITAN Xp",
            gpu_brand: "Titan",
            driver_ver: "535.146.02",
            memory_total: 12288,
            memory_util: 0,
            gpu_util: 0,
            memory_used: 83,
            fan_speed: 23,
            gpu_temp: 32,
          },
        ],
      },
      {
        name: "Workstation 4",
        gpus: [
          {
            gpu_name: "NVIDIA GeForce GT 1030",
            gpu_brand: "GeForce",
            driver_ver: "535.146.02",
            memory_total: 2048,
            memory_util: 0,
            gpu_util: 0,
            memory_used: 82,
            fan_speed: 35,
            gpu_temp: 31,
          },
        ],
      },

      {
        name: "Workstation 6",
        gpus: [
          {
            gpu_name: "NVIDIA GeForce GT 730",
            gpu_brand: "GeForce",
            driver_ver: "470.223.02",
            memory_total: 2001,
            memory_util: 0,
            gpu_util: 0,
            memory_used: 220,
            fan_speed: 30,
            gpu_temp: 27,
          },
        ],
      },
    ],
  },
];

// Currently does not attempt to do any validation of the returned GPU stats,
// or indeed handle errors that might be thrown by the Promises
const retrieveAllStats: () => Promise<
  Validated<WorkStationGroup[]>
> = async () => success(await (await fetch(API_URL + "/api/stats/all")).json());
//success(foo);

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
              {g.map((l) => (
                <Box
                  w="97%"
                  m={5}
                  bg={bgcol}
                  paddingTop={10}
                  paddingBottom={10}
                >
                  <VStack spacing={5}>
                    <Heading size="lg" textAlign="left">
                      Group 1: Personal
                    </Heading>
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
        <p>Tab le</p>
      </Navbar>
    </ChakraProvider>
  );
}

export default App;
