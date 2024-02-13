import React, { ReactNode, createContext, useContext, useState } from "react";
import {
  Validated,
  Validation,
  discard,
  failure,
  loading,
  success,
} from "../Utils/Utils";
import { API_URL } from "../App";
import { ADMIN_PATH } from "../Pages/AdminPanel";

const AUTH_PATH = "/auth";

/** Debug authentication is NOT a vulnerability: We don't get a token and so
 * the WebAPI server will still reject any admin requests. The purpose is to
 * just test sign-in functionality without the back-end running.
 */
const DEBUG_AUTH = true;
const DEBUG_USER = "NathanielB";
const DEBUG_PASSWORD = "drowssap";

type AuthCtx = {
  user: Validated<string>;
  login: (username: string, password: string) => void;
  logout: () => void;
  useAuthFetch: (
    path: string,
    init?: RequestInit | undefined,
  ) => Validation<Response>;
};

const AuthContext = createContext<AuthCtx>({
  user: failure(Error("No auth context provided")),
  login: () => {},
  logout: () => {},
  useAuthFetch: () => failure(Error("No auth context provided")),
});

export const useAuth = () => useContext(AuthContext);

export const AuthProvider = ({ children }: { children: ReactNode[] }) => {
  const [user, setUser] = useState<Validated<string>>(
    failure(Error("Not logged in!")),
  );

  const authFetch = (path: string, init?: RequestInit | undefined) =>
    fetch(API_URL + ADMIN_PATH + path, init);

  /** Feedback about if the login was successful should be retrieved by reading
   *  'user'
   */
  const login = (username: string, password: string) => {
    discard(async () => {
      if (
        DEBUG_AUTH &&
        username === DEBUG_USER &&
        password === DEBUG_PASSWORD
      ) {
        setUser(success(username));
        return;
      }

      const r = await authFetch(AUTH_PATH, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ username, password }),
      });

      if (r.ok) {
        setUser(success(username));
      } else if (r.status === 401) {
        setUser(failure(Error("Username or password was incorrect!")));
      } else {
        setUser(failure(Error("Auth failed for an unknown reason")));
      }
    });
  };

  const logout = () => {
    setUser(failure(Error("Logged out!")));
  };

  const useAuthFetch = (path: string, init?: RequestInit | undefined) => {
    const [resp, setResp] = useState<Validation<Response>>(loading());

    discard(async () => {
      const r = await authFetch(path, init);
      if (!r.ok && r.status === 403) {
        logout();
      }
      setResp(success(r));
    });

    return resp;
  };

  const value = { user, login, logout, useAuthFetch };

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
};
