import { useState } from "react";
import { Box, Button, Heading, Input, Text, VStack } from "@chakra-ui/react";
import { useAuth } from "../Providers/AuthProvider";
import { Failure, VTag, validatedElim } from "../Utils/Utils";

export const SignIn = () => {
  const { login, user: authName } = useAuth();

  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");

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
        <Button
          bgColor={"white"}
          onClick={async () => {
            login(username, password);
          }}
        >
          Sign In
        </Button>
      </VStack>
    </Box>
  );
};
