import "./App.css";
import { TableTab } from "./Components/DataTable";
import { useJarJar, useOnce } from "./Utils/Hooks";
import { Validated, Validation, success, validationElim } from "./Utils/Utils";
import { WorkstationCardMin } from "./Components/WorkstationCardMinimal";
import { Box, ChakraProvider, Heading, VStack } from "@chakra-ui/react";
import {} from "@chakra-ui/react";
import { Navbar } from "./Components/Navbar";
import { ColumnGrid } from "./Components/ColumnGrid";
import { WorkStationGroup } from "./Data";
import React, { useState } from "react";
import { BrowserRouter, Navigate, Route, Routes } from "react-router-dom";
import { MainView } from "./Pages/MainView";
import { STATS_PATH } from "./Config/Paths";

export const API_URL = "http://localhost:8000";
export const REFRESH_INTERVAL = 5000;

export enum ViewPage {
  CARD = "card_view",
  TABLE = "table_view",
}

export const VIEW_PAGE_INDEX = { [ViewPage.CARD]: 0, [ViewPage.TABLE]: 1 };

const App = () => {
  const [viewPage, setViewPage] = useState(ViewPage.CARD);

  return (
    <ChakraProvider>
      <div className="App"></div>
      <BrowserRouter>
        <Routes>
          <Route index element={<Navigate to={STATS_PATH} replace />} />
          <Route path={STATS_PATH + "/:viewPage"} element={<MainView />} />
          <Route path={"/admin_sign_in"} element={<></>} />
          <Route path={"/admin_panel"} element={<></>} />
        </Routes>
      </BrowserRouter>
    </ChakraProvider>
  );
};

export default App;
