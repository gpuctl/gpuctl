import { Box, Center, Heading, useColorModeValue } from "@chakra-ui/react";
import { WorkStationData } from "../App";

export const WorkstationTab = ({ name, gpus }: WorkStationData) => {
  const textCol = useColorModeValue("black", "white");
  return (
    <Center>
      <Box w="100%" rounded={"md"} bg={useColorModeValue("white", "gray.900")}>
        <Heading color={textCol}>{name}</Heading>
        {gpus.map((s, i) => {
          return (
            <Box key={i}>
              <Heading size="lg" color={textCol}>
                {s.gpu_name}
              </Heading>
            </Box>
          );
        })}
      </Box>
    </Center>
  );
};
