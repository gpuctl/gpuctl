import React from 'react';
import { VStack } from '@chakra-ui/react';
import { AddMachineForm } from '../Components/AddMachineForm';
import { GroupInfoManagement } from '../Components/GroupInfoManagement';
import { WorkStationGroup } from '../Data';

export const ADMIN_PATH = "/admin";


type AdminPanelProps = {
  groups: WorkStationGroup[];
};

export const AdminPanel: React.FC<AdminPanelProps> = ({ groups }) => {
  return (
    <VStack padding={10} spacing={10}>
      <AddMachineForm groups={groups} />
      <GroupInfoManagement groups={groups} />
    </VStack>
  );
};