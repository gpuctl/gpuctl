import {
  Alert,
  AlertIcon,
  Box,
  Center,
  HStack,
  Heading,
  Link,
  LinkBox,
  LinkOverlay,
  Tooltip,
  useColorModeValue,
} from "@chakra-ui/react";
import { WorkStationData } from "../Data";
import { TimeIcon } from "@chakra-ui/icons";
import { workstationBusy } from "../Utils/Utils";

import { Link as ReactRouterLink, useSearchParams } from "react-router-dom";

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

  const [params] = useSearchParams();
  params.append("selected", name);

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
                <Link
                  as={ReactRouterLink}
                  to={{ search: params.toString() }}
                  isTruncated
                >
                  {" " + name}
                </Link>
              </HStack>
            ) : (
              <Link
                as={ReactRouterLink}
                to={{ search: params.toString() }}
                isTruncated
              >
                {name}
              </Link>
            )}
          </Heading>
          <LinkBox>
            <LinkOverlay
              as={ReactRouterLink}
              to={{ search: params.toString() }}
            />
            {gpus.map((s, i) => (
              <Box key={i}>
                <Heading size="md">{`${s.gpu_name} (${(
                  s.memory_total / 1000
                ).toFixed(0)} GB)`}</Heading>
                <p>{`${s.in_use ? `🔴 In-use (User: ${s.user})` : "🟢 Available"}`}</p>
                <p>{`${s.gpu_util < 10 ? "🐌" : "🏎️" + "☁️".repeat(Math.ceil(s.gpu_util / 40))} GPU Usage: ${Math.round(s.gpu_util)}%`}</p>
                <p>{`${s.gpu_temp < 75 ? "❄️" : s.gpu_temp < 95 ? "🌡️" : "🔥"} ${Math.round(
                  s.gpu_temp,
                )} °C (${Math.round(s.fan_speed)}% Fan Speed)`}</p>
              </Box>
            ))}
          </LinkBox>
        </Box>
      </Box>
    </Center>
  );
};

//Returns color for greyed out components
export const greyed = (b: boolean) => {
  return b ? "white" : "white";
};
