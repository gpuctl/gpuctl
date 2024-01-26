import {
  Box,
  Center,
  Heading,
  Table,
  useColorModeValue,
  Tbody,
  Td,
  Tr,
} from "@chakra-ui/react";
import { WorkStationData } from "../Data";
import { isFree } from "../Utils/Utils";

export const WorkstationCard = ({ name, gpus }: WorkStationData) => {
  const textCol = useColorModeValue("black", "white");
  return (
    <Center>
      <Box w="100%" rounded={"md"} bg={useColorModeValue("white", "gray.900")}>
        <Heading size="md" color={textCol}>
          {name}
        </Heading>
        {gpus.map((s, i) => {
          return (
            <Box key={i}>
              <Table size="sm" className="table-tiny" variant="striped">
                <Tbody>
                  <Tr>
                    <Td>
                      <Heading size="sm">GPU Model:</Heading>
                    </Td>
                    <Td>
                      <Heading size="sm">{s.gpu_name}</Heading>
                    </Td>
                  </Tr>
                  <Tr>
                    <Td>
                      <Heading size="sm">Free?</Heading>
                    </Td>
                    <Td>
                      <Heading size="sm">
                        {isFree(s) ? `Yes ✅ (${s.gpu_util}% used)` : "No ❌"}
                      </Heading>
                    </Td>
                  </Tr>
                  <Tr>
                    <Td>
                      <Heading size="sm">GPU Memory:</Heading>
                    </Td>
                    <Td>
                      <Heading size="sm">
                        {(s.memory_total / 1000).toFixed(0)} GB (
                        {s.memory_total - s.memory_used} MB free)
                      </Heading>
                    </Td>
                  </Tr>
                  <Tr>
                    <Td>
                      <Heading size="sm">GPU Temperature:</Heading>
                    </Td>
                    <Td>
                      <Heading size="sm">{s.gpu_temp} °C</Heading>
                    </Td>
                  </Tr>
                </Tbody>
              </Table>
            </Box>
          );
        })}
      </Box>
    </Center>
  );
};
