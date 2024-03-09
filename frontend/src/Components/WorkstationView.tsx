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
  Tr,
  MenuItemOption,
  Spacer,
  RangeSlider,
  RangeSliderTrack,
  RangeSliderFilledTrack,
  RangeSliderThumb,
  HStack,
  VStack,
  Thead,
  Icon,
  IconProps,
  Stack,
} from "@chakra-ui/react";
import { Navigate, useSearchParams } from "react-router-dom";
import { Graph } from "./Graph";
import { useHistoryStats } from "../Hooks/Hooks";
import { useState } from "react";
import { GraphField, WorkStationData } from "../Data";
import {
  Validation,
  all,
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
import { range } from "d3";
import { NotesPopout } from "./NotesPopout";

const USE_FAKE_STATS = false;

const FAKE_STATS = [
  { x: 1, y: 90 },
  { x: 2, y: 12 },
  { x: 3, y: 34 },
  { x: 4, y: 53 },
  { x: 5, y: 98 },
];

const GRAPH_FIELDS = enumVals(GraphField);

export const GRAPH_COLS = [
  "#5DA5DA",
  "#FAA43A",
  "#60BD68",
  "#F17CB0",
  "#B2912F",
  "#B276B2",
  "#DECF3F",
  "#F15854",
];

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
                <VStack spacing={10}>
                  <StatsGraphPanel
                    hostname={hostname}
                    numGPUs={wstat.gpus.length}
                  />
                  <AdminDetails stats={wstat} group_name={name}></AdminDetails>
                </VStack>
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
      <TableContainer overflowX="scroll" css="transform: scaleY(-1)">
        <Table variant="striped" css="transform: scaleY(-1)">
          <Thead>
            <Tr>
              <Th>Field</Th>
              {range(0, stats.gpus.length).map((i) => (
                <Th key={i}>GPU {i + 1}</Th>
              ))}
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
                    <Td key={j}>{f}</Td>
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

const StatsGraphPanel = ({
  hostname,
  numGPUs,
}: {
  hostname: string;
  numGPUs: number;
}) => {
  const [field, setField] = useState<GraphField>(GraphField.GPU_TEMP);

  const historyStats = useHistoryStats(hostname);

  const statsToDisplay: Validation<
    {
      x: number;
      y: number;
    }[][]
  > = USE_FAKE_STATS
    ? success([FAKE_STATS])
    : mapSuccess(historyStats, (hist) => {
        return hist.map((h) =>
          h.map(({ timestamp, sample }) => ({
            x: timestamp,
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
        success: (s) =>
          s.length > 0 || s[0].length > 0 ? (
            <StatsGraph stats={s} numGPUs={numGPUs}></StatsGraph>
          ) : (
            <Text>No historical data found!</Text>
          ),
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
    <Table variant="simple">
      <Tbody>
        <Tr>
          <Td> Group </Td>
          <EditableField
            group={group_name}
            workstation={stats}
            fieldKey="group"
            placeholder="unknown"
            isEven={true}
          />
        </Tr>
        <Tr>
          <Td>CPU </Td>
          <EditableField
            group={group_name}
            workstation={stats}
            fieldKey="cpu"
            placeholder="unknown"
            isEven={true}
          />
        </Tr>
        <Tr>
          <Td>Motherboard</Td>
          <EditableField
            group={group_name}
            workstation={stats}
            fieldKey="motherboard"
            placeholder="unknown"
            isEven={true}
          />
        </Tr>
        <Tr>
          <Td>Notes</Td>
          <Td>
            <HStack>
              <Text isTruncated maxWidth="15rem">
                {stats.notes}
              </Text>
              <NotesPopout
                wname={stats.name}
                notes={stats.notes}
                isEven={true}
              />
            </HStack>
          </Td>
        </Tr>
        <Tr>
          <Td>Owner</Td>
          <EditableField
            group={group_name}
            workstation={stats}
            fieldKey="owner"
            placeholder="none"
            isEven={true}
          />
        </Tr>
      </Tbody>
    </Table>
  );
};

const StatsGraph = ({
  stats,
  numGPUs,
}: {
  stats: { x: number; y: number }[][];
  numGPUs: number;
}) => {
  const minTS =
    stats.length === 0 ? 0 : Math.min(...stats[0].map(({ x }) => x));
  const maxTS =
    stats.length === 0 ? 0 : Math.max(...stats[0].map(({ x }) => x));

  const [startTS, setStartTS] = useState(minTS);
  const [endTS, setEndTS] = useState(maxTS);
  const [atEnd, setAtEnd] = useState(true);

  if (atEnd && endTS !== maxTS) {
    setEndTS(maxTS);
  }

  const filtered = stats
    .map((s) => s.filter(({ x }) => startTS <= x && x <= endTS))
    .map((line, i) =>
      line.length > 0 ? { off: stats[i].indexOf(line[0]), line } : null,
    );

  return (
    <>
      <Graph data={filtered} xlabel="Time" maxPoints={50}></Graph>
      <RangeSlider
        defaultValue={[minTS, maxTS]}
        value={[startTS, endTS]}
        min={minTS}
        max={maxTS}
        onChange={([min, max]) => {
          if (min === max) return;
          setStartTS(min);
          setEndTS(max);
          setAtEnd(max === maxTS);
        }}
      >
        <RangeSliderTrack>
          <RangeSliderFilledTrack />
        </RangeSliderTrack>
        <RangeSliderThumb index={0} />
        <RangeSliderThumb index={1} />
      </RangeSlider>
      {numGPUs > 1 ? (
        <>
          <Stack direction={["column", "row"]} spacing={5}>
            {range(0, numGPUs).map((i) => (
              <HStack key={i}>
                <CircleIcon color={GRAPH_COLS[i]}></CircleIcon>
                <Text fontSize="xl">{`GPU ${i + 1}`}</Text>
              </HStack>
            ))}
          </Stack>
        </>
      ) : (
        <></>
      )}
    </>
  );
};

const CircleIcon = (props: IconProps) => (
  <Icon viewBox="0 0 200 200" {...props}>
    <path
      fill="currentColor"
      d="M 100, 100 m -75, 0 a 75,75 0 1,0 150,0 a 75,75 0 1,0 -150,0"
    />
  </Icon>
);
