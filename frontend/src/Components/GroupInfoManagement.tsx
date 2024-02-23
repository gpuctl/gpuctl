import React, { useState } from "react";
import {
  Box,
  Button,
  Heading,
  Table,
  TableContainer,
  Tbody,
  Td,
  Th,
  Thead,
  Tr,
  useDisclosure,
  Modal,
  ModalOverlay,
  ModalContent,
  ModalHeader,
  ModalFooter,
  ModalBody,
  ModalCloseButton,
  Input,
} from "@chakra-ui/react";
import { EditableField } from "./EditableFields";
import { WorkStationGroup } from "../Data";
import { instKeys } from "../Utils/Utils";
import { useRemoveMachine } from "../Hooks/Hooks";
import { GS } from "../Pages/AdminPanel";

export const GroupInfoManagement = ({
  GroupSelect,
  groups,
}: {
  GroupSelect: GS;
  groups: WorkStationGroup[];
}) => {
  const { isOpen, onOpen, onClose } = useDisclosure();
  const [copied, setCopied] = useState(false);
  const [currentMachine, setCurrentMachine] = useState("");

  const removeMachine = useRemoveMachine();

  const copyToClipboard = (username: string) => {
    const command = `ssh ${username}@${currentMachine} shutdown now`;
    navigator.clipboard.writeText(command);
    setCopied(true);
    setTimeout(() => {
      setCopied(false);
    }, 3000);
    onClose();
  };

  const handleShutdownClick = (machineName: string) => {
    setCurrentMachine(machineName);
    onOpen();
  };

  return (
    <Box w="100%">
      <Heading size="lg">Group & Info Management:</Heading>
      <TableContainer mt={4}>
        <Table variant="striped">
          <Thead>
            <Tr>
              <Th>Hostname</Th>
              <Th>Group</Th>
              <Th>CPU</Th>
              <Th>Motherboard</Th>
              <Th>Notes</Th>
              <Th>Action</Th>
              <Th>Shutdown</Th>
            </Tr>
          </Thead>
          <Tbody>
            {instKeys(
              groups.flatMap((group) =>
                group.workstations.map((workstation) => (k) => (
                  <Tr key={k}>
                    <Td>{workstation.name}</Td>
                    <EditableField
                      GroupSelect={GroupSelect}
                      workstation={workstation}
                      fieldKey="group"
                      placeholder="unknown"
                    />
                    <EditableField
                      GroupSelect={GroupSelect}
                      workstation={workstation}
                      fieldKey="cpu"
                      placeholder="unknown"
                    />
                    <EditableField
                      GroupSelect={GroupSelect}
                      workstation={workstation}
                      fieldKey="motherboard"
                      placeholder="unknown"
                    />
                    <EditableField
                      GroupSelect={GroupSelect}
                      workstation={workstation}
                      fieldKey="notes"
                      placeholder="none"
                    />
                    <Td>
                      <Button
                        colorScheme="red"
                        onClick={() => removeMachine(workstation.name)}
                      >
                        Remove
                      </Button>
                    </Td>
                    <Td>
                      <Button
                        colorScheme="blue"
                        onClick={() => handleShutdownClick(workstation.name)}
                        disabled={copied}
                      >
                        {copied ? "Copied" : "Copy Shutdown Command"}
                      </Button>
                    </Td>
                  </Tr>
                )),
              ),
            )}
          </Tbody>
        </Table>
      </TableContainer>
      <Modal isOpen={isOpen} onClose={onClose}>
        <ModalOverlay />
        <ModalContent>
          <ModalHeader>Enter Username</ModalHeader>
          <ModalCloseButton />
          <ModalBody>
            <Input placeholder="Username" id="usernameInput" />
          </ModalBody>
          <ModalFooter>
            <Button
              colorScheme="blue"
              mr={3}
              onClick={() =>
                copyToClipboard(
                  document.getElementById("usernameInput")?.innerText || "",
                )
              }
            >
              Copy Command
            </Button>
            <Button variant="ghost" onClick={onClose}>
              Cancel
            </Button>
          </ModalFooter>
        </ModalContent>
      </Modal>
    </Box>
  );
};
