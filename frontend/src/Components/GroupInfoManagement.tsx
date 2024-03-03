import { useRef, useState } from "react";
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

  const inputRef = useRef<HTMLInputElement>(null);
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
              <Th>Owner</Th>
              <Th>Action</Th>
              <Th>Shutdown</Th>
            </Tr>
          </Thead>
          <Tbody>
            {instKeys(
              groups.flatMap((group) =>
                group.workstations.map((workstation) => (k) => (
                  <Tr key={workstation.name}>
                    <Td>{workstation.name}</Td>
                    <EditableField
                      GroupSelect={GroupSelect}
                      group={group.name}
                      workstation={workstation}
                      fieldKey="group"
                      placeholder="unknown"
                      isEven={k % 2 == 0}
                    />
                    <EditableField
                      GroupSelect={GroupSelect}
                      group={group.name}
                      workstation={workstation}
                      fieldKey="cpu"
                      placeholder="unknown"
                      isEven={k % 2 == 0}
                    />
                    <EditableField
                      GroupSelect={GroupSelect}
                      group={group.name}
                      workstation={workstation}
                      fieldKey="motherboard"
                      placeholder="unknown"
                      isEven={k % 2 == 0}
                    />
                    <EditableField
                      GroupSelect={GroupSelect}
                      group={group.name}
                      workstation={workstation}
                      fieldKey="notes"
                      placeholder="none"
                      isEven={k % 2 == 0}
                    />
                    <EditableField
                      GroupSelect={GroupSelect}
                      group={group.name}
                      workstation={workstation}
                      fieldKey="owner"
                      placeholder="none"
                      isEven={k % 2 == 0}
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
            <Input placeholder="Username" ref={inputRef} />
          </ModalBody>
          <ModalFooter>
            <Button
              colorScheme="blue"
              mr={3}
              onClick={() => copyToClipboard(inputRef.current?.value || "")}
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
