import { Box, Heading, useColorModeValue } from "@chakra-ui/react";
import { WorkStationData } from "../App";

export const WorkstationTab = ({ name, gpus }: WorkStationData) => {
  const textCol = useColorModeValue("black", "white");
  return (
    <Box rounded={"md"} bg={useColorModeValue("white", "gray.900")}>
      <Heading color={textCol}>{name}</Heading>
      {gpus.map((s, i) => {
        return (
          <Box key={i}>
            <Heading color={textCol}>{s.gpu_name}</Heading>
          </Box>
        );
      })}
    </Box>
  );
};
