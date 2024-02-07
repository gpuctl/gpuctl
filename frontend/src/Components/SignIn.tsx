import { useState } from "react";
import { API_URL } from "../App";
import {
  Validated,
  failure,
  success,
  validatedElim,
  validationElim,
} from "../Utils/Utils";
import { ADMIN_PATH } from "../Pages/AdminPanel";
import {
  Box,
  Button,
  Editable,
  EditableInput,
  EditablePreview,
  Heading,
  Input,
  VStack,
} from "@chakra-ui/react";

const AUTH_PATH = "/auth";
const AUTH_TOKEN_ITEM = "authToken";

type AuthToken = {
  token: string;
};

const DEBUG_AUTH = true;
const DEBUG_USER = "NathanielB";
const DEBUG_PASSWORD = "drowssap";
const DEBUG_TOKEN = "ABCDEFG";

const requestSignIn = async (
  username: string,
  password: string,
): Promise<Validated<AuthToken>> => {
  if (DEBUG_AUTH) {
    if (username === DEBUG_USER && password === DEBUG_PASSWORD) {
      return success({ token: DEBUG_TOKEN });
    }
  }

  const resp = await fetch(API_URL + ADMIN_PATH + AUTH_PATH, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ username, password }),
  });

  if (resp.ok) {
    const tok: AuthToken = await resp.json();
    return success(tok);
  }
  if (resp.status === 401) return failure(Error("Unauthorised!"));
  return failure(Error("Authentication Failed for an Unknown Reason!"));
};

const tryGetAuthToken = (): Validated<AuthToken> => {
  const token = localStorage.getItem(AUTH_TOKEN_ITEM);
  return token == null ? failure(Error("No token :(")) : success({ token });
};

export const SignIn = () => {
  const [authToken, setAuth] =
    useState<Validated<AuthToken>>(tryGetAuthToken());

  const updateAuth = (tok: AuthToken) => {
    localStorage.setItem(AUTH_TOKEN_ITEM, tok.token);
    setAuth(success(tok));
  };

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
        <Button
          bgColor={"white"}
          onClick={async () => {
            const tok = await requestSignIn(username, password);
            validatedElim(tok, {
              success: (t) => {
                updateAuth(t);
                window.location.reload();
              },
              failure: () => {},
            });
          }}
        >
          Sign In
        </Button>
      </VStack>
    </Box>
  );
};
