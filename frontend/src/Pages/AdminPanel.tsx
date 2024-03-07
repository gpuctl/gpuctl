import React, { PropsWithChildren } from "react";
import { Box, StyleProps, VStack } from "@chakra-ui/react";
import { AddMachineForm } from "../Components/AddMachineForm";
import { GroupInfoManagement } from "../Components/GroupInfoManagement";
import { WorkStationGroup } from "../Data";
import {
  AutoComplete,
  AutoCompleteItem,
  AutoCompleteList,
} from "@choc-ui/chakra-autocomplete";

export const ADMIN_PATH = "/admin";

type AdminPanelProps = {
  groups: WorkStationGroup[];
};

export type GS = ({
  onChange,
  ...props
}: {
  children?: React.ReactNode;
} & StyleProps & {
    onChange: (v: string) => void;
  }) => JSX.Element;

export const AdminPanel: React.FC<AdminPanelProps> = ({ groups }) => {
  // In a perfect world, we would use Choc auto-complete for group selection.
  // For now though, we just use plain text entry, because the auto-complete
  // component is being a huge pain.
  const GroupSelect = ({
    onChange,
    ...props
  }: PropsWithChildren & StyleProps & { onChange: (v: string) => void }) => (
    <AutoComplete
      openOnFocus
      creatable
      onChange={onChange}
      values={groups.map((g) => g.name)}
    >
      <Box {...props} />
      <AutoCompleteList maxHeight="200px" overflow="scroll">
        {groups.map((g, i) => (
          <AutoCompleteItem key={i} value={g.name}></AutoCompleteItem>
        ))}
      </AutoCompleteList>
    </AutoComplete>
  );

  return (
    <VStack padding={10} spacing={10}>
<<<<<<< Updated upstream
      <AddMachineForm GroupSelect={GroupSelect} groups={groups} />
      <GroupInfoManagement GroupSelect={GroupSelect} groups={groups} />
<<<<<<< Updated upstream
=======
      <Button
        colorScheme="red"
        onClick={() => {
          logout();
          nav("/");
        }}
      >
        Sign Out
      </Button>
=======
      <Box w="100%" textAlign="left">
        <Heading size="lg">Add a Machine:</Heading>
      </Box>
      <HStack w="100%">
        <Input
          w="50%"
          placeholder="Hostname (e.g. mira05.doc.ic.ac.uk)"
          onChange={(e) => setHostname(e.target.value)}
        ></Input>
        <AutoComplete
          openOnFocus
          creatable
          onChange={(e) => setGroup(e.target.value)}
          values={groups.map((g) => g.name)}
        >
          <AutoCompleteInput
            w="45%"
            placeholder="Group Name (e.g. shared)"
          ></AutoCompleteInput>
          <AutoCompleteList maxHeight="200px" overflow="scroll">
            {groups.map((g, i) => (
              <AutoCompleteItem key={i} value={g.name}></AutoCompleteItem>
            ))}
          </AutoCompleteList>
        </AutoComplete>
        <Button
          w="5%"
          onClick={() => {
            addMachine(hostname, group);
          }}
        >
          Add
        </Button>
      </HStack>
      <Box w="100%" textAlign="left">
        <Heading size="lg">Group & Info Management:</Heading>
      </Box>
      <TableContainer>
        <Table variant="striped">
          <Thead>
            <Tr>
              <Th width="5rem"> Hostname </Th>
              <Th width="5rem"> Group </Th>
              <Th> CPU </Th>
              <Th width="5rem"> Motherboard </Th>
              <Th width="5rem"> Notes </Th>
            </Tr>
          </Thead>
          <Tbody>
            {instKeys(
              groups.flatMap(({ name: group, workstations }, i) => {
                return workstations.map((workstation, j) => {
                  return (k: number) => (
                    <Tr key={k}>
                      <Td> {workstation.name} </Td>
                      <Td>
                        <Editable
                          isPreviewFocusable={true}
                          submitOnBlur={true}
                          placeholder="Unknown"
                          defaultValue={group}
                          textColor={pickCol(group)}
                          onSubmit={(s) =>
                            modifyInfo(workstation.name, {
                              group: s,
                              cpu: null,
                              motherboard: null,
                              notes: null,
                            })
                          }
                        >
                          {(props) => (
                            <>
                              <EditablePreview />
                              <div>
                              <EditableInput textColor={"gray.600"} />
                              <EditableButton {...props} />
                              </div>
                            </>
                          )}
                        </Editable>
                      </Td>
                      <Td>
                        <Editable
                          isPreviewFocusable={true}
                          submitOnBlur={true}
                          placeholder="Unknown"
                          textColor={pickCol(workstation.cpu)}
                          defaultValue={workstation.cpu}
                          onSubmit={(s) =>
                            modifyInfo(workstation.name, {
                              group: null,
                              cpu: s,
                              motherboard: null,
                              notes: null,
                            })
                          }
                        >
                          {(props) => (
                            <>
                              <EditablePreview />
                              <EditableInput textColor={"gray.600"} />
                              <EditableButton {...props} />
                            </>
                          )}
                        </Editable>
                      </Td>
                      <Td>
                        <Editable
                          isPreviewFocusable={true}
                          submitOnBlur={true}
                          placeholder="Unknown"
                          defaultValue={workstation.motherboard}
                          textColor={pickCol(workstation.motherboard)}
                          onSubmit={(s) =>
                            modifyInfo(workstation.name, {
                              group: null,
                              cpu: null,
                              motherboard: s,
                              notes: null,
                            })
                          }
                        >
                          {(props) => (
                            <>
                              <EditablePreview />
                              <EditableInput textColor={"gray.600"} />
                              <EditableButton {...props} />
                            </>
                          )}
                        </Editable>
                      </Td>
                      <Td>
                        <Editable
                          isPreviewFocusable={true}
                          submitOnBlur={true}
                          placeholder="None"
                          defaultValue={workstation.notes}
                          textColor={pickCol(workstation.notes)}
                          onSubmit={(s) =>
                            modifyInfo(workstation.name, {
                              group: null,
                              cpu: null,
                              motherboard: null,
                              notes: s,
                            })
                          }
                        >
                          {(props) => (
                            <>
                              <EditablePreview />
                              <EditableInput textColor={"gray.600"} />
                              <EditableButton {...props} />
                            </>
                          )}
                        </Editable>
                      </Td>
                      <Td>
                        <Button
                          bgColor={"red.300"}
                          onClick={() => {
                            removeMachine(hostname);
                          }}
                        >
                          Kill Satellite
                        </Button>
                      </Td>
                    </Tr>
                  );
                });
              }),
            )}
          </Tbody>
        </Table>
      </TableContainer>
>>>>>>> Stashed changes
>>>>>>> Stashed changes
    </VStack>
  );
};
