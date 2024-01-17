import React from "react";
import "./App.css";
import { useAsync } from "./Utils/Hooks";

const API_URL = "http://localhost:8000";

type GPUStats = {
  gpu_name: string;
  brand: string;
  driver_version: string;
  memory_total: number;

  memory_util: number;
  gpu_util: number;
  memory_used: number;
  fan_speed: number;
  gpu_temp: number;
};

const retrieveAllStats: () => Promise<GPUStats[]> = async () =>
  (await fetch(API_URL + "/api/stats/all")).json();

function App() {
  const stats = useAsync(retrieveAllStats());

  return (
    <div className="App">
      <header className="App-header">
        <p>Welcome to the GPU Control Room!</p>
        {stats?.map((row, i) => {
          return (
            <p key={i}>
              ID: {i}, Name: {row.gpu_name}, Core Utilisation: {row.gpu_util}%,
              VRAM: {row.memory_total} GB, Used VRAM: {row.memory_used},
              Temperature: {row.gpu_temp}
            </p>
          );
        })}
      </header>
    </div>
  );
}

export default App;
