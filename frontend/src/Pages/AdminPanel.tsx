import { Input } from "@chakra-ui/input";
import { API_URL, AuthToken } from "../App";
import { STATS_PATH } from "../Config/Paths";
import { Box, Center, HStack, Heading, VStack } from "@chakra-ui/layout";
import { Button } from "@chakra-ui/button";
import { useState } from "react";

export const ADMIN_PATH = "/admin";
const ADD_MACHINE_URL = "/add_workstation";

const addMachine = (hostname: string, group: string, token: AuthToken) => {
  fetch(API_URL + ADMIN_PATH + ADD_MACHINE_URL, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      Authorization: `Bearer ${token.token}`,
    },
    body: JSON.stringify({ hostname, group }),
  });
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

export const AdminPanel = ({ token }: { token: AuthToken }) => {
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
    </VStack>
  );
};
