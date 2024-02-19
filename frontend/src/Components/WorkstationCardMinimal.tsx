import { Box, Center, Heading, useColorModeValue } from "@chakra-ui/react";
import { WorkStationData } from "../Data";

export const WorkstationCardMin = ({
  width,
  data: { name, gpus, lastSeen },
}: {
  width: number;
  data: WorkStationData & { lastSeen: number | undefined };
}) => {
  const textCol = useColorModeValue("black", "white");
  return (
    <Center>
      <Box
        padding={2}
        w={width}
        rounded={"md"}
        bg={useColorModeValue("white", "gray.900")}
      >
        <Heading size="lg" color={textCol}>
          {name}
        </Heading>
        {gpus.map((s, i) => {
          return (
            <Box key={i}>
              <Heading size="md">{`${s.gpu_name} (${(
                s.memory_total / 1000
              ).toFixed(0)} GB)`}</Heading>
              <p>{`${s.gpu_util < 25 ? "🟢 Free" : "🔴 In-use"} (${Math.round(
                s.gpu_util,
              )}% Utilisation)`}</p>
              <p>{`${s.gpu_temp < 80 ? "❄️" : s.gpu_temp < 95 ? "🌡️" : "🔥"} ${Math.round(
                s.gpu_temp,
              )} °C (${Math.round(s.fan_speed)}% Fan Speed)`}</p>
            </Box>
          );
        })}
      </Box>
    </Center>
  );
};
