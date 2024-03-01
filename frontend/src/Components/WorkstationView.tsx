import { Box, Heading, Text } from "@chakra-ui/react";
import { WorkStationData } from "../Data";

export const WorkstationView = (data: WorkStationData) => {
  return (
    <Box>
      <Heading>{data.name}</Heading>
      {
        // Should it be possible to edit notes while
      }
      <Heading>Notes</Heading>
      <Text>{data.notes}</Text>
    </Box>
  );
};
