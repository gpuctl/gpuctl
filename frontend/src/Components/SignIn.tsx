import { useState } from "react";
import { Box, Button, Heading, Input, Text, VStack } from "@chakra-ui/react";
import { useAuth } from "../Providers/AuthProvider";
import { validatedElim, isSuccess } from "../Utils/Utils";

export const SignIn = ({ panelCallback }: { panelCallback: () => void }) => {
  const { login, user: authName } = useAuth();

  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");

  const attemptLogin = async () => {
    const success = await login(username, password);
    if (success) {
      panelCallback();
    }
  };

  return (
    <Box padding={4} bgColor={"gray.100"}>
      <VStack spacing={2}>
        <Box w="100%">
          <Heading textAlign={"left"} size="l">
            Username
          </Heading>
        </Box>
        <Input
          value={username}
          onChange={(e) => setUsername(e.target.value)}
          bgColor={"white"}
          onKeyDown={(e) => {
            if (e.key === "Enter") {
              attemptLogin();
            }
          }}
        ></Input>
        <Box w="100%">
          <Heading textAlign={"left"} size="l">
            Password
          </Heading>
        </Box>
        <Input
          type="password"
          value={password}
          onChange={(e) => setPassword(e.target.value)}
          bgColor={"white"}
          onKeyDown={(e) => {
            if (e.key === "Enter") {
              attemptLogin();
            }
          }}
        ></Input>
        {validatedElim(authName, {
          success: () => <></>,
          failure: (e) =>
            e.message === "Username or password was incorrect!" ? (
              <Text color="tomato"> Username or password incorrect</Text>
            ) : (
              <></>
            ),
        })}
        ;
        <Button bgColor={"white"} onClick={attemptLogin}>
          Sign In
        </Button>
      </VStack>
    </Box>
  );
};
