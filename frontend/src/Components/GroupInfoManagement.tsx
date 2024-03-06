import { useEffect, useRef, useState } from "react";
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
  List,
  ListItem,
  Flex,
  Link,
} from "@chakra-ui/react";
import { EditableField } from "./EditableFields";
import { WorkStationGroup } from "../Data";
import { instKeys, validateNullable, validatedElim } from "../Utils/Utils";
import {
  useGetAllFiles,
  useGetSpecificFile,
  useRemoveFile,
  useRemoveMachine,
  useUploadFile,
} from "../Hooks/Hooks";
import { GS } from "../Pages/AdminPanel";

import { Link as ReactRouterLink, useSearchParams } from "react-router-dom";
import React from "react";

export const GroupInfoManagement = ({
  GroupSelect,
  groups,
}: {
  GroupSelect: GS;
  groups: WorkStationGroup[];
}) => {
  const {
    isOpen: isFilesModalOpen,
    onOpen: onFilesModalOpen,
    onClose: onFilesModalClose,
  } = useDisclosure();
  const { isOpen, onOpen, onClose } = useDisclosure();
  const [copied, setCopied] = useState(false);
  const [currentMachine, setCurrentMachine] = useState("");
  const [currentFile, setCurrentFile] = useState("");
  const [params] = useSearchParams();

  const [files, setFiles] = useState<any[]>([]); // TODO: Change this any
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

  const useDownload = () => {
    const downloader = useGetSpecificFile(
      currentMachine,
      currentFile,
      (response) => {
        validatedElim(response, {
          success: async (r) => {
            // here the response should have given an attached file through the attachment header stuff
            // let's download

            if (r.status == 200) {
              const downloadLink = document.createElement("a");
              downloadLink.href = URL.createObjectURL(
                new Blob([await r.arrayBuffer()]),
              );
              downloadLink.download = currentFile;
              downloadLink.click();
            }
          },
          failure: () => {
            // alert("Failure");
          },
        });
      },
    );

    return () => downloader();
  };

  const download = useDownload();

  const handleSpecificFileDownload = (filename: string) => {
    return () => {
      setCurrentFile(filename);
    };
  };

  useEffect(() => {
    download();
  }, [currentFile]);

  const useViewFiles = () => {
    const getFiles = useGetAllFiles(currentMachine, (response) => {
      validatedElim(response, {
        success: async (r) => {
          setFiles(await r.json());
        },
        failure: () => {
          // alert("Failure");
        },
      });
    });
    return () => getFiles();
  };

  const view = useViewFiles();

  const handleViewFiles = (name: string) => {
    return () => {
      setCurrentMachine(name);
      onFilesModalOpen();
    };
  };

  useEffect(() => {
    if (currentMachine) {
      view();
    }
  }, [currentMachine]);

  // ...

  const uploadFile = useUploadFile((response) => {
    validatedElim(response, {
      success: () => {
        view();
      },
      failure: () => {
        // alert("Failure");
      },
    });
  });

  const removeFile = useRemoveFile();

  const handleFileUpload = (name: string) => {
    return async (event: React.ChangeEvent<HTMLInputElement>) => {
      if (!event.target.files) return;
      const file = event.target.files[0];

      const arrayBuffer = await file.arrayBuffer();
      const uint8Array = new Uint8Array(arrayBuffer);

      uploadFile(name, file.type, file.name, uint8Array);
    };
  };

  const handleFileRemoval = (name: string, filename: string) => {
    return (event: React.MouseEvent<HTMLButtonElement, MouseEvent>) => {
      removeFile(name, filename);
      setFiles(files.filter((i) => i != filename));
    };
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
              <Th>Files</Th>
            </Tr>
          </Thead>
          <Tbody>
            {instKeys(
              groups.flatMap((group) =>
                group.workstations.map((workstation) => (k) => {
                  const newParams = new URLSearchParams(
                    Object.fromEntries(Array.from(params.entries())),
                  );
                  newParams.append("selected", workstation.name);

                  return (
                    <Tr key={workstation.name}>
                      <Td>
                        {" "}
                        <Link
                          as={ReactRouterLink}
                          to={{ search: newParams.toString() }}
                        >
                          {workstation.name}{" "}
                        </Link>
                      </Td>
                      <EditableField
                        group={group.name}
                        workstation={workstation}
                        fieldKey="group"
                        placeholder="unknown"
                        isEven={k % 2 === 0}
                      />
                      <EditableField
                        group={group.name}
                        workstation={workstation}
                        fieldKey="cpu"
                        placeholder="unknown"
                        isEven={k % 2 === 0}
                      />
                      <EditableField
                        group={group.name}
                        workstation={workstation}
                        fieldKey="motherboard"
                        placeholder="unknown"
                        isEven={k % 2 === 0}
                      />
                      <EditableField
                        group={group.name}
                        workstation={workstation}
                        fieldKey="notes"
                        placeholder="none"
                        isEven={k % 2 === 0}
                      />
                      <EditableField
                        group={group.name}
                        workstation={workstation}
                        fieldKey="owner"
                        placeholder="none"
                        isEven={k % 2 === 0}
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
                      <Td>
                        <Button
                          colorScheme="green"
                          onClick={handleViewFiles(workstation.name)}
                        >
                          Files
                        </Button>
                      </Td>
                    </Tr>
                  );
                }),
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
      <Modal
        isOpen={isFilesModalOpen}
        onClose={onFilesModalClose}
        size="xl"
        isCentered
      >
        <ModalOverlay />
        <ModalContent>
          <ModalHeader>Files</ModalHeader>
          <ModalCloseButton />
          <ModalBody>
            <List spacing={3}>
              {files.map((file, index) => (
                <ListItem
                  key={index}
                  display="flex"
                  justifyContent="space-between"
                  alignItems="center"
                >
                  <Box flex="1">{file}</Box>
                  <Flex alignItems="center" gap="2">
                    <Button
                      colorScheme="blue"
                      onClick={handleSpecificFileDownload(file)}
                    >
                      Download
                    </Button>
                    <Button
                      colorScheme="red"
                      size="sm"
                      onClick={handleFileRemoval(currentMachine, file)}
                    >
                      X
                    </Button>
                  </Flex>
                </ListItem>
              ))}
            </List>
          </ModalBody>
          <ModalFooter>
            <Input type="file" onChange={handleFileUpload(currentMachine)} />
          </ModalFooter>
        </ModalContent>
      </Modal>
    </Box>
  );
};
