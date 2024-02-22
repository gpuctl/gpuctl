import {
  Alert,
  AlertIcon,
  Box,
  Center,
  HStack,
  Heading,
  Tooltip,
  useColorModeValue,
} from "@chakra-ui/react";
import { WorkStationData } from "../Data";
import { TimeIcon } from "@chakra-ui/icons";
import { cropString } from "../Utils/Utils";

const LAST_SEEN_WARN_THRESH = 60 * 5;
const GREEN = "#25D36B";

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
        padding={0}
        w={width}
        rounded={"md"}
        bg={useColorModeValue("white", "gray.900")}
      >
        {lastSeen !== undefined && lastSeen > LAST_SEEN_WARN_THRESH ? (
          <Alert
            roundedTopLeft="md"
            roundedTopRight="md"
            textAlign="center"
            status="warning"
          >
            <AlertIcon></AlertIcon>
            {lastSeen === undefined
              ? "Missing provenance data!"
              : `Last seen over ${Math.floor(lastSeen / 60)} minutes ago!`}
          </Alert>
        ) : (
          <></>
        )}
        <Box padding={2}>
          <Heading size="lg" color={textCol}>
            {lastSeen !== undefined && lastSeen <= LAST_SEEN_WARN_THRESH ? (
              <HStack>
                <Tooltip placement="right-start" label={`Seen recently!`}>
                  <TimeIcon color={GREEN} />
                </Tooltip>
                <Box>{" " + cropString(name, 15)}</Box>
              </HStack>
            ) : (
              <>{cropString(name, 17)}</>
            )}
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
      </Box>
    </Center>
  );
};
