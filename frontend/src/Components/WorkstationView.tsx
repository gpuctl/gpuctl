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
  ModalContent,
  ModalHeader,
  ModalOverlay,
  SimpleGrid,
  Table,
  TableContainer,
  Tbody,
  Th,
  Thead,
  Tr,
} from "@chakra-ui/react";
import { useSearchParams } from "react-router-dom";
import { Graph } from "./Graph";
import { useHistoryStats } from "../Hooks/Hooks";
import { useState } from "react";
import { GraphField } from "../Data";
import { enumVals, mapSuccess, success, validationElim } from "../Utils/Utils";
import { Text } from "@chakra-ui/react";
import { ChevronDownIcon } from "@chakra-ui/icons";
// import { COLS } from "./DataTable";

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

  // const [shownColumns, setter] = useState(COLS);

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
    >
      <ModalOverlay bg="blackAlpha.300" backdropFilter="blur(15px)" />
      <ModalContent>
        <ModalHeader>{`${hostname}`}</ModalHeader>
        <ModalCloseButton />
        <ModalBody>
          <SimpleGrid columns={2} spacing={5}>
            <Box>
              <TableContainer>
                <Table variant="striped">
                  <Thead>
                    <Tr>
                      <Th>Field</Th>
                      <Th>GPU 0</Th>
                    </Tr>
                  </Thead>
                  <Tbody></Tbody>
                </Table>
              </TableContainer>
            </Box>
            <StatsGraph hostname={hostname} />
          </SimpleGrid>
        </ModalBody>
      </ModalContent>
    </Modal>
  );
};

const StatsGraph = ({ hostname }: { hostname: string }) => {
  const [field, setField] = useState<GraphField>(GraphField.GPU_UTIL);

  const historyStats = useHistoryStats(hostname);

  const statsToDisplay = USE_FAKE_STATS
    ? success(FAKE_STATS)
    : mapSuccess(historyStats, (s) =>
        s.map(({ timestamp, sample }) => ({
          x: timestamp,
          // TODO: Support machines with multiple GPUs/no GPUs
          y: sample[0][field],
        })),
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
        success: (s) => <Graph data={[s]}></Graph>,
        failure: () => <Text>Failed to fetch historical data for graph!</Text>,
        loading: () => <Text>Fetching data...</Text>,
      })}
    </Box>
  );
};
