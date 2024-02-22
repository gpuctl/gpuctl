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
import { cropString, workstationBusy } from "../Utils/Utils";

const LAST_SEEN_WARN_THRESH = 60 * 5;
const GREEN = "#25D36B";

export const WorkstationCardMin = ({
  width,
  data: { name, gpus, last_seen },
}: {
  width: number;
  data: WorkStationData;
}) => {
  const textCol = useColorModeValue("black", "white");
  return (
    <Center>
      <Box
        padding={0}
        w={width}
        rounded={"md"}
        opacity={workstationBusy(gpus) ? 0.4 : 1.0}
        bg={greyed(workstationBusy(gpus))}
      >
        {last_seen !== undefined && last_seen > LAST_SEEN_WARN_THRESH ? (
          <Alert
            roundedTopLeft="md"
            roundedTopRight="md"
            textAlign="center"
            status="warning"
          >
            <AlertIcon></AlertIcon>
            {last_seen === undefined
              ? "Missing provenance data!"
              : `Last seen over ${Math.floor(last_seen / 60)} minutes ago!`}
          </Alert>
        ) : (
          <></>
        )}
        <Box padding={2}>
          <Heading size="lg" color={textCol}>
            {last_seen !== undefined && last_seen <= LAST_SEEN_WARN_THRESH ? (
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

          {gpus.map((s, i) => (
            <Box key={i}>
              <Heading size="md">{`${s.gpu_name} (${(
                s.memory_total / 1000
              ).toFixed(0)} GB)`}</Heading>
              <p>{`${s.in_use ? "🔴 In-use" : "🟢 Free"} (User: ${s.user})`}</p>
              <p>{`${s.gpu_util < 10 ? "🐌" : "🏎️" + "☁️".repeat(Math.ceil(s.gpu_util / 40))} GPU Usage: (${Math.round(s.gpu_util)}% Utilisation)`}</p>
              <p>{`${s.gpu_temp < 80 ? "❄️" : s.gpu_temp < 95 ? "🌡️" : "🔥"} ${Math.round(
                s.gpu_temp,
              )} °C (${Math.round(s.fan_speed)}% Fan Speed)`}</p>
            </Box>
          ))}
        </Box>
      </Box>
    </Center>
  );
};

//Returns color for greyed out components
export const greyed = (b: boolean) => {
  return b ? "white" : "white";
};
