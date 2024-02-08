import { Input } from "@chakra-ui/input";
import { API_URL, AuthToken } from "../App";
import {
  Editable,
  EditableInput,
  EditableTextarea,
  EditablePreview,
  Table,
  TableContainer,
  Tbody,
  Td,
  Th,
  Thead,
  Tr,
} from "@chakra-ui/react";
import { STATS_PATH } from "../Config/Paths";
import { Box, Center, HStack, Heading, VStack } from "@chakra-ui/layout";
import { Button } from "@chakra-ui/button";
import { useState } from "react";
import { WorkStationGroup } from "../Data";

export const ADMIN_PATH = "/admin";
const ADD_MACHINE_URL = "/add_workstation";

const addMachine = async (
  hostname: string,
  group: string,
  token: AuthToken,
) => {
  const resp = await fetch(API_URL + ADMIN_PATH + ADD_MACHINE_URL, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      Authorization: `Bearer ${token.token}`,
    },
    body: JSON.stringify({ hostname, group }),
  });
  if (resp.ok) {
    console.log("Success!");
  } else if (resp.status === 401) {
    console.log("Auth Error!");
  } else {
    console.log("Unknown Error!");
  }
  // We should probably await the response to give feedback on whether adding
  // the machine was successful...
};

type ModifyData = {
  group: string | null;
  motherboard: string | null;
  cpu: string | null;
  notes: string | null;
};

const modifyInfo = (hostname: string, mod: ModifyData) => {
  fetch(API_URL + ADMIN_PATH + STATS_PATH + "/modify", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ hostname, ...mod }),
  });
};

export const AdminPanel = ({
  token,
  groups,
}: {
  token: AuthToken;
  groups: WorkStationGroup[];
}) => {
  const [hostname, setHostname] = useState("");
  const [group, setGroup] = useState("");

  // TODO
  return (
    <VStack>
      <Box w="100%" textAlign="left">
        <Heading size="lg">Add a Machine:</Heading>
      </Box>
      <HStack w="100%">
        <Input
          w="50%"
          placeholder="Hostname (e.g. mira05.doc.ic.ac.uk)"
        ></Input>
        <Input w="45%" placeholder="Group Name (e.g. shared)"></Input>
        <Button
          w="5%"
          onClick={() => {
            addMachine(hostname, group, token);
          }}
        >
          Add
        </Button>
      </HStack>
      <TableContainer>
        <Table variant="striped">
          <Thead>
            <Tr>
              <Th key={0}> Hostname </Th>
              <Th key={1}> Group </Th>
              <Th key={2}> CPU </Th>
              <Th key={3}> Motherboard </Th>
              <Th key={4}> Notes </Th>
            </Tr>
          </Thead>
          <Tbody>
            {groups.map(({ name, workStations }, i) => {
              return workStations.map((workStation, j) => {
                const id = (i * j + j) * 5;
                return (
                  <Tr key={id}>
                    <Th key={id}> {workStation.name} </Th>
                    <Th key={id + 1}> {name} </Th>
                    <Th key={id + 2}> {workStation.cpu} </Th>
                    <Th key={id + 3}> {workStation.motherboard} </Th>
                    <Th key={id + 4}> {workStation.notes} </Th>
                  </Tr>
                );
              });
            })}
          </Tbody>
        </Table>
      </TableContainer>
    </VStack>
  );
};
