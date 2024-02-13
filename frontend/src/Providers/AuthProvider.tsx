import React, { ReactNode, createContext, useContext, useState } from "react";
import {
  Validated,
  Validation,
  discard,
  failure,
  loading,
  success,
} from "../Utils/Utils";
import { useOnce } from "../Utils/Hooks";
import { API_URL } from "../App";
import { ADMIN_PATH } from "../Pages/AdminPanel";

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

  const login = (username: string, password: string) => {
    // Hit /api/auth endpoint
    setUser(success(username));
  };

  const logout = () => {
    setUser(failure(Error("Logged out!")));
  };

  const useAuthFetch = (path: string, init?: RequestInit | undefined) => {
    const [resp, setResp] = useState<Validation<Response>>(loading());

    discard(async () => {
      const r = await fetch(API_URL + ADMIN_PATH + path, init);
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
