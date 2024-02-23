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
      <AddMachineForm GroupSelect={GroupSelect} groups={groups} />
      <GroupInfoManagement GroupSelect={GroupSelect} groups={groups} />
    </VStack>
  );
};
