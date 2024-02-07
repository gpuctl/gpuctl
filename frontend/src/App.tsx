import "./App.css";
import { TableTab } from "./Components/DataTable";
import { useJarJar, useOnce } from "./Utils/Hooks";
import {
  Validated,
  Validation,
  failure,
  success,
  validationElim,
  enumVals,
} from "./Utils/Utils";
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
  CARD = "/card_view",
  TABLE = "/table_view",
  ADMIN = "/admin_view",
}

enum Tricky {
  CARD = "/card_view",
  TABLE = "/table_viewadsdas",
  ADMIN = "/admin_view",
}

enum Num {
  C = "D'",
  A = 0,
  B = 1,
  D = "C'",
}

type Wahh = ViewPage;

const DEFAULT_VIEW = ViewPage.CARD;

export const VIEW_PAGE_INDEX = {
  [ViewPage.CARD]: 0,
  [ViewPage.TABLE]: 1,
  [ViewPage.ADMIN]: 2,
};

export type AuthToken = {
  token: string;
};

export const AUTH_TOKEN_ITEM = "authToken";

const tryGetAuthToken = (): Validated<AuthToken> => {
  const token = localStorage.getItem(AUTH_TOKEN_ITEM);
  return token == null ? failure(Error("No token :(")) : success({ token });
};

const App = () => {
  const [authToken, setAuth] =
    useState<Validated<AuthToken>>(tryGetAuthToken());

  const ooh = enumVals(ViewPage);
  const ohBaby = enumVals(Num);

  console.log(ooh);
  console.log(ohBaby);

  return (
    <ChakraProvider>
      <div className="App"></div>
      <BrowserRouter>
        <Routes>
          <Route index element={<Navigate to={STATS_PATH} replace />} />
          <Route
            path={STATS_PATH}
            element={<Navigate to={STATS_PATH + DEFAULT_VIEW} replace />}
          />

          <Route
            path={STATS_PATH + ViewPage.CARD}
            element={<MainView default={ViewPage.CARD} />}
          />
          <Route
            path={STATS_PATH + ViewPage.TABLE}
            element={<MainView default={ViewPage.TABLE} />}
          />
          <Route
            path={STATS_PATH + ViewPage.ADMIN}
            element={<MainView default={ViewPage.ADMIN} />}
          />
        </Routes>
      </BrowserRouter>
    </ChakraProvider>
  );
};

export default App;
