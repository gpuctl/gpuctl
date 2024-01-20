import "./App.css";
import { useJarJar, useOnce } from "./Utils/Hooks";
import { inlineLog } from "./Utils/Utils";

const API_URL = "http://localhost:8000";
export const REFRESH_INTERVAL = 5000;

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

const retrieveAllStats: () => Promise<GPUStats[]> = async () => {
  const yeow = fetch(API_URL + "/api/stats/all");
  console.log(yeow);
  const yeow2 = await yeow;
  console.log(yeow2);
  return yeow2.json();
};
// (await fetch(API_URL + "/api/stats/all")).json();

function App() {
  const [stats, updateStats] = useJarJar(retrieveAllStats);

  useOnce(() => {
    setInterval(updateStats, REFRESH_INTERVAL);
  });

  return (
    <div className="App">
      <header className="App-header">
        <p>Welcome to the GPU Control Room!</p>
        {stats?.map((row, i) => {
          return (
            <p key={i}>
              ID: {i}, Name: {inlineLog(row.gpu_name)}, Core Utilisation:{" "}
              {row.gpu_util}%, VRAM Util: {row.memory_util}%, VRAM:{" "}
              {row.memory_total} GB, Used VRAM: {row.memory_used} GB,
              Temperature: {row.gpu_temp} °C
            </p>
          );
        }) ?? (
          <div>
            <p>Retrieving data from API server...</p>
            <p>NVIDIA yeow</p>
          </div>
        )}
      </header>
    </div>
  );
}

export default App;
