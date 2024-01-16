import React from "react";
import "./App.css";
import { useAsync } from "./Utils/Hooks";

type Stats = {
  clock_speed: number;
  util: number;
  gpu_mem: number;
  gpu_mem_used: number;
};

const retrieveAllStats = async () => {
  return [{ clock_speed: 0, util: 0, gpu_mem: 0, gpu_mem_used: 0 }];
};

function App() {
  const stats = useAsync(retrieveAllStats);

  return (
    <div className="App">
      <header className="App-header">
        <p>Welcome to the GPU Control Room!</p>
        {stats?.map((row, i) => {
          return (
            <p>
              ID: {i}, Core Utilisation: {row.util}% Clock Speed:{" "}
              {row.clock_speed} MHz, VRAM: {row.gpu_mem} GB, Used VRAM:{" "}
              {row.gpu_mem_used}
            </p>
          );
        })}
      </header>
    </div>
  );
}

export default App;
