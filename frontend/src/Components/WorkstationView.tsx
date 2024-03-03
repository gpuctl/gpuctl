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
} from "@chakra-ui/react";
import { useSearchParams } from "react-router-dom";
import { Graph } from "./Graph";
import { useHistoryStats } from "../Hooks/Hooks";
import { useState } from "react";
import { GraphField } from "../Data";
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
import { COLS, tablify } from "./DataTable";
import { useForceUpdate } from "framer-motion";

const USE_FAKE_STATS = true;

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
            <StatsTable hostname={hostname}></StatsTable>
            <StatsGraph hostname={hostname} />
          </SimpleGrid>
        </ModalBody>
      </ModalContent>
    </Modal>
  );
};

const StatsTable = ({ hostname }: { hostname: string }) => {
  const [shownColumns, setter] = useState(
    Object.fromEntries(Object.keys(COLS).map((k) => [k, true])),
  );

  const { allStats } = useStats();
  const [refresh] = useForceUpdate();

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

      {validationElim(allStats, {
        success: (s) => {
          const gpus = s.flatMap((g) =>
            g.workstations
              .filter((w) => w.name === hostname)
              .flatMap((w) => w.gpus),
          );

          const tabGpus = transpose(
            gpus.map((g) => tablify(shownColumns, g)),
          ).map((r) => all(r));

          return (
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
                        <Td>{Object.keys(COLS)[i]}</Td>
                        {fs.map((f, j) => (
                          <Td key={j}>
                            {cropString(f, Math.round(35 / fs.length))}
                          </Td>
                        ))}
                      </Tr>
                    ),
                  )}
                </Tbody>
              </Table>
            </TableContainer>
          );
        },
        failure: () => <Text>Failure fetching data! Retrying...</Text>,
        loading: () => <Text>Fetching data...</Text>,
      })}
    </Box>
  );
};

const StatsGraph = ({ hostname }: { hostname: string }) => {
  const [field, setField] = useState<GraphField>(GraphField.GPU_UTIL);

  const historyStats = useHistoryStats(hostname);

  const statsToDisplay: Validation<
    {
      x: number;
      y: number;
    }[][]
  > = USE_FAKE_STATS
    ? success([FAKE_STATS])
    : mapSuccess(historyStats, (s) =>
        s.map(({ timestamp, sample }) =>
          sample.map((s) => ({ x: timestamp, y: s[field] })),
        ),
      );

  return (
    <Box>
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
      {validationElim(statsToDisplay, {
        success: (s) => <Graph data={s}></Graph>,
        failure: () => (
          <Text>Failed to fetch historical data for graph! Retrying...</Text>
        ),
        loading: () => <Text>Fetching data...</Text>,
      })}
    </Box>
  );
};
