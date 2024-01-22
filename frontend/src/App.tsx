import "./App.css";
import { WorkstationTab } from "./Components/WorkstationTab";
import { useJarJar, useOnce } from "./Utils/Hooks";
import { Validated, success, validationElim } from "./Utils/Utils";
import { ChakraProvider } from "@chakra-ui/react";

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
        <header className="App-header">
          <p>Welcome to the GPU Control Room!</p>
          {validationElim(stats, {
            success: (l) => (
              <div>
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
              </div>
            ),
            loading: () => <p>Retrieving data from API server...</p>,
            failure: (_) => <p>Something has gone wrong!</p>,
          })}
        </header>
      </div>
    </ChakraProvider>
  );
}

export default App;
