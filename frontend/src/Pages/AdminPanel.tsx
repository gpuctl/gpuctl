import React, { PropsWithChildren } from "react";
import { useNavigate } from "react-router-dom";
import { Button, Box, StyleProps, VStack } from "@chakra-ui/react";
import { AddMachineForm } from "../Components/AddMachineForm";
import { GroupInfoManagement } from "../Components/GroupInfoManagement";
import { WorkStationGroup } from "../Data";
import {
  AutoComplete,
  AutoCompleteItem,
  AutoCompleteList,
} from "@choc-ui/chakra-autocomplete";
import { useAuth } from "../Providers/AuthProvider";

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
  const { isSignedIn, logout } = useAuth();
  const nav = useNavigate();

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
      <Button
        colorScheme="red"
        onClick={() => {
          logout();
          nav("/");
        }}
      >
        Sign Out
      </Button>
    </VStack>
  );
};
