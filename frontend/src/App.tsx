import React from "react";
import "./App.css";
import { useAsync } from "./Utils/Hooks";

const API_URL = "http://localhost:8000";

type Stats = {
  clock_speed: number;
  util: number;
  gpu_mem: number;
  gpu_mem_used: number;
};

const retrieveAllStats: () => Promise<Stats[]> = async () => {
  const stats = await fetch(API_URL + "/api/stats/all", { mode: "cors" });
  const jason = await stats.json();
  return jason;
};

function App() {
  const stats = useAsync(retrieveAllStats());

  return (
    <div className="App">
      <header className="App-header">
        <p>Welcome to the GPU Control Room!</p>
        {stats?.map((row, i) => {
          return (
            <p key={i}>
              ID: {i}, Core Utilisation: {row.util}%, Clock Speed:{" "}
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
