import {
  Box,
  Button,
  Menu,
  MenuButton,
  MenuItem,
  MenuList,
  Modal,
  ModalBody,
  ModalCloseButton,
  MenuOptionGroup,
  ModalContent,
  ModalHeader,
  ModalOverlay,
  SimpleGrid,
  Table,
  TableContainer,
  Tbody,
  Td,
  Th,
  Thead,
  Tr,
  MenuItemOption,
  Spacer,
  RangeSlider,
  RangeSliderTrack,
  RangeSliderFilledTrack,
  RangeSliderThumb,
} from "@chakra-ui/react";
import { Navigate, useSearchParams } from "react-router-dom";
import { Graph } from "./Graph";
import { useHistoryStats } from "../Hooks/Hooks";
import { useState } from "react";
import { GraphField, WorkStationData } from "../Data";
import {
  Validation,
  all,
  cropString,
  enumVals,
  mapSuccess,
  success,
  transpose,
  validationElim,
} from "../Utils/Utils";
import { Text } from "@chakra-ui/react";
import { ChevronDownIcon } from "@chakra-ui/icons";
import { useStats } from "../Providers/FetchProvider";
import { GPU_FIELDS, tablify } from "./DataTable";
import { useForceUpdate } from "framer-motion";
import { useAuth } from "../Providers/AuthProvider";
import { EditableField } from "./EditableFields";
import { STATS_PATH } from "../Config/Paths";
import { DEFAULT_VIEW } from "../App";

const USE_FAKE_STATS = false;

const FAKE_STATS = [
  { x: 1, y: 90 },
  { x: 2, y: 12 },
  { x: 3, y: 34 },
  { x: 4, y: 53 },
  { x: 5, y: 98 },
];

const GRAPH_FIELDS = enumVals(GraphField);

export const WorkstationView = ({ hostname }: { hostname: string }) => {
  const [, setPs] = useSearchParams();
  const { allStats } = useStats();

  return validationElim(allStats, {
    success: (stats) => {
      const wstats = stats.flatMap((g) =>
        g.workstations
          .filter((w) => w.name === hostname)
          .map((w) => ({ wstat: w, name: g.name })),
      );

      if (wstats.length !== 1)
        return <Navigate to={STATS_PATH + DEFAULT_VIEW} replace />;

      const { wstat, name } = wstats[0];

      return (
        <Modal
          size="xl"
          isOpen={hostname !== null}
          onClose={() => {
            setPs((ps) => {
              ps.delete("selected");
              return ps;
            });
          }}
          scrollBehavior="inside"
        >
          <ModalOverlay bg="blackAlpha.300" backdropFilter="blur(15px)" />
          <ModalContent minWidth="80%" minHeight="80%">
            <ModalHeader>{`${hostname}`}</ModalHeader>
            <ModalCloseButton />
            <ModalBody>
              <SimpleGrid columns={2} spacing={5}>
                <StatsTable stats={wstat}></StatsTable>
                <Box>
                  <StatsGraphPanel hostname={hostname} />
                  <AdminDetails stats={wstat} group_name={name}></AdminDetails>
                </Box>
              </SimpleGrid>
            </ModalBody>
          </ModalContent>
        </Modal>
      );
    },
    failure: () => <Text>Failure fetching data! Retrying...</Text>,
    loading: () => <Text>Fetching data...</Text>,
  });
};

const StatsTable = ({ stats }: { stats: WorkStationData }) => {
  const [shownColumns, setter] = useState(
    Object.fromEntries(Object.keys(GPU_FIELDS).map((k) => [k, true])),
  );
  const [refresh] = useForceUpdate();

  const tabGpus = transpose(
    stats.gpus.map((g) => tablify(shownColumns, g)),
  ).map((r) => all(r));

  return (
    <Box>
      <Menu closeOnSelect={false}>
        <MenuButton as={Button} colorScheme="blue">
          Columns
        </MenuButton>
        <MenuList overflowY="scroll" maxHeight="200">
          <MenuOptionGroup
            type="checkbox"
            defaultValue={Object.keys(shownColumns).filter(
              (key) => shownColumns[key],
            )}
            onChange={(props) => {
              Object.keys(shownColumns).forEach((col) => {
                shownColumns[col] = props.includes(col);
              });
              setter(shownColumns);
              refresh();
            }}
          >
            {Object.keys(shownColumns).map((col, i) => {
              return (
                <MenuItemOption value={col} key={i}>
                  {` ${col} `}
                </MenuItemOption>
              );
            })}
          </MenuOptionGroup>
        </MenuList>
      </Menu>
      <TableContainer>
        <Table variant="striped">
          <Thead>
            <Tr>
              <Th>Field</Th>
              <Th>GPU 0</Th>
            </Tr>
          </Thead>
          <Tbody>
            {tabGpus.map((fs, i) =>
              fs === null ? (
                <></>
              ) : (
                <Tr key={i}>
                  <Td>{Object.keys(GPU_FIELDS)[i]}</Td>
                  {fs.map((f, j) => (
                    <Td key={j}>{cropString(f, Math.round(35 / fs.length))}</Td>
                  ))}
                </Tr>
              ),
            )}
          </Tbody>
        </Table>
      </TableContainer>
    </Box>
  );
};

const StatsGraphPanel = ({ hostname }: { hostname: string }) => {
  const [field, setField] = useState<GraphField>(GraphField.GPU_UTIL);

  const historyStats = useHistoryStats(hostname);

  const statsToDisplay: Validation<
    {
      x: number;
      y: number;
    }[][]
  > = USE_FAKE_STATS
    ? success([FAKE_STATS])
    : mapSuccess(historyStats, (hist) => {
        const minTS = Math.min(
          ...hist.flatMap((h) => h.map(({ timestamp }) => timestamp)),
        );
        return hist.map((h) =>
          h.map(({ timestamp, sample }) => ({
            x: timestamp - minTS,
            y: sample[GPU_FIELDS[field]],
          })),
        );
      });

  return (
    <Box w="100%">
      <Menu>
        <MenuButton
          as={Button}
          rightIcon={<ChevronDownIcon />}
        >{`${field}`}</MenuButton>
        <MenuList>
          {GRAPH_FIELDS.map((f, i) => (
            <MenuItem
              key={i}
              onClick={() => {
                setField(f);
              }}
            >
              {f}
            </MenuItem>
          ))}
        </MenuList>
      </Menu>
      <Spacer height={5} />
      {validationElim(statsToDisplay, {
        success: (s) => <StatsGraph stats={s}></StatsGraph>,
        failure: () => (
          <Text>Failed to fetch historical data for graph! Retrying...</Text>
        ),
        loading: () => <Text>Fetching data...</Text>,
      })}
    </Box>
  );
};

const AdminDetails = ({
  stats,
  group_name,
}: {
  stats: WorkStationData;
  group_name: string;
}) => {
  const { isSignedIn } = useAuth();

  if (!isSignedIn()) return <></>;

  return (
    <Table variant="striped">
      <Tbody>
        <Tr key={0}>
          <Td> Group </Td>
          <Td>
            <EditableField
              group={group_name}
              workstation={stats}
              fieldKey="group"
              placeholder="unknown"
              isEven={true}
            />
          </Td>
        </Tr>
        <Tr key={1}>
          <Td>CPU </Td>
          <Td>
            <EditableField
              group={group_name}
              workstation={stats}
              fieldKey="cpu"
              placeholder="unknown"
              isEven={false}
            />
          </Td>
        </Tr>
        <Tr key={2}>
          <Td>Motherboard</Td>
          <Td>
            <EditableField
              group={group_name}
              workstation={stats}
              fieldKey="motherboard"
              placeholder="unknown"
              isEven={true}
            />
          </Td>
        </Tr>
        <Tr key={3}>
          <Td>Notes</Td>
          <Td>
            <EditableField
              group={group_name}
              workstation={stats}
              fieldKey="notes"
              placeholder="none"
              isEven={false}
            />
          </Td>
        </Tr>
        <Tr key={4}>
          <Td>Owner</Td>
          <Td>
            <EditableField
              group={group_name}
              workstation={stats}
              fieldKey="owner"
              placeholder="none"
              isEven={true}
            />
          </Td>
        </Tr>
      </Tbody>
    </Table>
  );
};

const StatsGraph = ({ stats }: { stats: { x: number; y: number }[][] }) => {
  const maxTS =
    stats.length === 0 ? 0 : Math.max(...stats[0].map(({ x }) => x));

  const [startTS, setStartTS] = useState(0);
  const [endTS, setEndTS] = useState(maxTS);

  const filtered = stats.map((s) =>
    s.filter(({ x }) => startTS <= x && x <= endTS),
  );

  return (
    <>
      <Graph data={filtered} xlabel="Seconds Since Added"></Graph>
      <RangeSlider
        defaultValue={[0, maxTS]}
        min={0}
        max={maxTS}
        onChange={([min, max]) => {
          setStartTS(min);
          setEndTS(max);
        }}
      >
        <RangeSliderTrack>
          <RangeSliderFilledTrack />
        </RangeSliderTrack>
        <RangeSliderThumb index={0} />
        <RangeSliderThumb index={1} />
      </RangeSlider>
    </>
  );
};
