import "./App.css";
import { enumVals, enumIndex } from "./Utils/Utils";
import { ChakraProvider } from "@chakra-ui/react";
import {} from "@chakra-ui/react";
import { BrowserRouter, Navigate, Route, Routes } from "react-router-dom";
import { MainView } from "./Pages/MainView";
import { STATS_PATH } from "./Config/Paths";
import { AuthProvider } from "./Providers/AuthProvider";
import { FetchStatsProvider } from "./Providers/FetchProvider";

export const API_URL =
  process.env.NODE_ENV === "production" ? "/api" : "http://localhost:8000/api";
export const REFRESH_INTERVAL = 5000;

export enum ViewPage {
  CARD = "/card_view",
  TABLE = "/table_view",
  ADMIN = "/admin_view",
}

export const DEFAULT_VIEW = ViewPage.CARD;

export const VIEW_PAGE_INDEX = enumIndex(ViewPage);

const App = () => (
  <ChakraProvider>
    <AuthProvider>
      <FetchStatsProvider>
        <div className="App"></div>
        <BrowserRouter>
          <Routes>
            <Route index element={<Navigate to={STATS_PATH} replace />} />
            <Route
              path={STATS_PATH}
              element={<Navigate to={STATS_PATH + DEFAULT_VIEW} replace />}
            />
            {enumVals(ViewPage).map((page, i) => (
              <Route
                key={i}
                path={STATS_PATH + page}
                element={<MainView page={page} />}
              />
            ))}
          </Routes>
        </BrowserRouter>
      </FetchStatsProvider>
    </AuthProvider>
  </ChakraProvider>
);

export default App;
