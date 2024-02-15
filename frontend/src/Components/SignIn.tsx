import { useState } from "react";
import { Box, Button, Heading, Input, VStack } from "@chakra-ui/react";
import { useAuth } from "../Providers/AuthProvider";
import { Failure, VTag } from "../Utils/Utils";
 
export const SignIn = () => {
  const { login, user : authName } = useAuth();

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
        { authName.tag == VTag.Failure && (authName as Failure).error.message == "Username or password was incorrect!" ? (<p> Username or passowrd incorrect</p>) : (<p> test</p>)}
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
