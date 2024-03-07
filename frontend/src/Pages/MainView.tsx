import { Box, Heading, VStack } from "@chakra-ui/react";
import { DEFAULT_VIEW, ViewPage } from "../App";
import { WorkStationGroup } from "../Data";
import {
  Validation,
  validateNullable,
  validatedElim,
  validationElim,
} from "../Utils/Utils";
import { ColumnGrid } from "../Components/ColumnGrid";
import { TableTab } from "../Components/DataTable";
import { WorkstationCardMin } from "../Components/WorkstationCardMinimal";
import { Navbar } from "../Components/Navbar";
import { useAuth } from "../Providers/AuthProvider";
import { STATS_PATH } from "../Config/Paths";
import { AdminPanel } from "./AdminPanel";
import { Navigate, useSearchParams } from "react-router-dom";
import { useStats } from "../Providers/FetchProvider";
import { WorkstationView } from "../Components/WorkstationView";
import { AggregateStats } from "../Components/AggregateStats";

const cardView = (stats: WorkStationGroup[]) => (
  <VStack spacing={5}>
    {stats.map((l, i) => (
      <Box key={i} w="97%" m={5} bg="white" paddingTop={5} paddingBottom={10}>
        <VStack spacing={5}>
          <Box w="100%" paddingLeft={5}>
            <Heading size="lg" textAlign="left">
              Group {i + 1}: {l.name}
            </Heading>
          </Box>
          <Box w="100%">
            <ColumnGrid minChildWidth={325} hMinSpacing={40} vSpacing={10}>
              {l.workstations.map((row, i) => {
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
    <Box key={stats.length} m={5} bg="white" paddingTop={5} paddingBottom={10}>
      <AggregateStats
        // TODO: replace with actual aggregate data from backend
        data={{ percent_used: 45.6, total_energy: 4563.2 }}
      />
    </Box>
  </VStack>
);

const tableView = (stats: WorkStationGroup[]) => (
  <TableTab groups={stats}></TableTab>
);

const adminView = (stats: WorkStationGroup[]) => (
  <AdminPanel groups={stats}></AdminPanel>
);

const displayPartial = (
  stats: Validation<WorkStationGroup[]>,
  cont: (gs: WorkStationGroup[]) => JSX.Element,
): JSX.Element => {
  return validationElim(stats, {
    success: cont,
    loading: () => <p>Retrieving data from API server...</p>,
    failure: (_) => <p>Something has gone wrong!</p>,
  });
};

export const MainView = ({ page }: { page: ViewPage }) => {
  const { isSignedIn } = useAuth();
  const [params] = useSearchParams();

  if (page === ViewPage.ADMIN && !isSignedIn()) {
    return <Navigate to={STATS_PATH + DEFAULT_VIEW} replace />;
  }

  const hostname = params.get("selected");

  return (
    <>
      <ConfirmedMainView initial={page} />
      {validatedElim(validateNullable(hostname), {
        success: (h) => <WorkstationView hostname={h} />,
        failure: () => <></>,
      })}
    </>
  );
};

export const ConfirmedMainView = ({ initial }: { initial: ViewPage }) => {
  const { allStats } = useStats();

  return (
    <Navbar initial={initial}>
      {displayPartial(allStats, cardView)}
      {displayPartial(allStats, tableView)}
      {displayPartial(allStats, adminView)}
    </Navbar>
  );
};
