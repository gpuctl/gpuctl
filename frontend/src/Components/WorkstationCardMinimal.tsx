import {
  Alert,
  AlertIcon,
  AlertDescription,
  Box,
  Card,
  CardHeader,
  CardBody,
  Center,
  HStack,
  Heading,
  LinkBox,
  LinkOverlay,
  Stack,
  StackDivider,
  Text,
  Tooltip,
} from "@chakra-ui/react";
import { GPUStats, WorkStationData } from "../Data";
import { TimeIcon } from "@chakra-ui/icons";

import { Link as ReactRouterLink, useSearchParams } from "react-router-dom";

const LAST_SEEN_WARN_THRESH = 60 * 5;
const GREEN = "#25D36B";

// A mechine is busy if all gpus are in use
const workstationBusy = (gs: GPUStats[]) => {
  return gs.every((g) => {
    return g.in_use;
  });
};

export const WorkstationCardMin = ({
  width,
  data: { name, gpus, last_seen },
}: {
  width: number;
  data: WorkStationData;
}) => {
  const [params] = useSearchParams();

  const newParams = new URLSearchParams(
    Object.fromEntries(Array.from(params.entries())),
  );
  newParams.append("selected", name);

  return (
    <Center>
      <LinkBox>
        <Card
          size="sm"
          w={width}
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
              <AlertIcon />
              <AlertDescription>
                {last_seen === undefined
                  ? "Missing provenance data!"
                  : `Last seen over ${Math.floor(last_seen / 60)} minutes ago!`}
              </AlertDescription>
            </Alert>
          ) : (
            <></>
          )}
          <CardHeader>
            {last_seen !== undefined && last_seen <= LAST_SEEN_WARN_THRESH ? (
              <Tooltip placement="left-start" label={`Seen recently!`}>
                <HStack>
                  <TimeIcon boxSize={10} color={GREEN} />
                  <Heading isTruncated>
                    <LinkOverlay
                      as={ReactRouterLink}
                      to={{ search: newParams.toString() }}
                    >
                      {name}
                    </LinkOverlay>
                  </Heading>
                </HStack>
              </Tooltip>
            ) : (
              <Heading isTruncated>
                <LinkOverlay
                  as={ReactRouterLink}
                  to={{ search: newParams.toString() }}
                >
                  {name}
                </LinkOverlay>
              </Heading>
            )}
          </CardHeader>
          <CardBody>
            <Stack divider={<StackDivider />} spacing="4">
              {gpus.map((s, i) => (
                <Box key={i}>
                  <Heading size="md">{`${s.gpu_name} (${(
                    s.memory_total / 1000
                  ).toFixed(0)} GB)`}</Heading>
                  <Text>{`${s.in_use ? `🔴 In-use (User: ${s.user})` : "🟢 Available"}`}</Text>
                  <Text>{`${s.gpu_util < 10 ? "🐌" : "🏎️" + "☁️".repeat(Math.ceil(s.gpu_util / 40))} GPU Usage: ${Math.round(s.gpu_util)}%`}</Text>
                  <Text>{`${s.gpu_temp < 75 ? "❄️" : s.gpu_temp < 95 ? "🌡️" : "🔥"} ${Math.round(
                    s.gpu_temp,
                  )} °C (${Math.round(s.fan_speed)}% Fan Speed)`}</Text>
                </Box>
              ))}
            </Stack>
          </CardBody>
        </Card>
      </LinkBox>
    </Center>
  );
};

//Returns color for greyed out components
export const greyed = (b: boolean) => {
  return b ? "white" : "white";
};
