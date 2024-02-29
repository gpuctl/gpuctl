import React, { ReactNode, createContext, useContext, useState } from "react";
import {
  Validated,
  Validation,
  discard,
  failure,
  fire,
  isSuccess,
  loading,
  success,
} from "../Utils/Utils";
import { API_URL } from "../App";
import { ADMIN_PATH } from "../Pages/AdminPanel";
import { useOnce } from "../Utils/Hooks";

const AUTH_PATH = "/auth";

/**
 * Debug authentication is NOT a vulnerability: We don't get a token and so
 * the WebAPI server will still reject any admin requests. The purpose is to
 * just test sign-in functionality without the back-end running.
 */
const DEBUG_AUTH = true;
const DEBUG_USER = "NathanielB";
const DEBUG_PASSWORD = "drowssap";

/**
 * Ideally no part of our site should *rely* on the page being reloaded to
 * update. The reload is primarily to make it obvious to the user than signing
 * in/out actually did something.
 */
const RELOAD_ON_LOG_CHANGE = false;

type AuthCtx = {
  user: Validated<string>;
  isSignedIn: () => boolean;
  login: (username: string, password: string) => void;
  logout: () => void;
  useAuthFetch: (
    path: string,
    f?: (r: Validated<Response>) => void,
  ) => [Validation<Response>, (init?: RequestInit) => void];
};

type UsernameReminder = { username: string };

const AuthContext = createContext<AuthCtx>({
  user: failure(Error("No auth context provided")),
  isSignedIn: () => false,
  login: () => {},
  logout: () => {},
  useAuthFetch: () => [failure(Error("No auth context provided")), () => {}],
});

const authFetch = (path: string, init?: RequestInit) => {
  if (init === undefined) {
    init = { credentials: "include" };
  } else {
    init["credentials"] = "include";
  }
  return fetch(API_URL + ADMIN_PATH + path, init);
};

export const useAuth = () => useContext(AuthContext);

export const AuthProvider = ({
  children,
}: {
  children: ReactNode[] | ReactNode;
}) => {
  const [user, setUserDirect] = useState<Validated<string>>(
    success("Checking log-in status"),
  );

  // On first page load, we would like to check if we are currently signed in
  useOnce(
    discard(async () => {
      const r = await authFetch("/confirm", {
        method: "GET",
      });
      if (r.ok) {
        const remind: UsernameReminder = await r.json();
        setUserDirect(success(remind.username));
      } else {
        setUserDirect(failure(Error("Not logged in!")));
      }
    }),
  );

  const setUser = (u: Validated<string>) => {
    const changed = isSuccess(u) !== isSuccess(user);
    setUserDirect(u);

    if (changed && RELOAD_ON_LOG_CHANGE) {
      window.location.reload();
    }
  };

  /**
   * Feedback about if the login was successful should be retrieved by reading
   * 'user'
   */
  const login = (username: string, password: string) => {
    fire(async () => {
      console.log("Tried to log in!");

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

      if (!r.ok) {
        setUser(
          failure(
            r.status === 401
              ? Error("Username or password was incorrect!")
              : Error("Auth failed for an unknown reason"),
          ),
        );
        return;
      }

      setUser(success(username));
      return;
    });
  };

  const logout = () => {
    setUser(failure(Error("Logged out!")));
  };

  // This maybe should return the firing function...
  const useAuthFetch = (
    path: string,
    callback?: (r: Validated<Response>) => void,
  ): [Validation<Response>, () => void] => {
    const [resp, setResp] = useState<Validation<Response>>(loading());

    const onResp = (r: Validated<Response>) => {
      setResp(r);
      if (callback) callback(r);
    };

    const f = (init?: RequestInit | undefined) => {
      setResp(loading());

      fire(async () => {
        const r = await authFetch(path, init);
        if (r.status === 401) {
          logout();
          onResp(failure(Error(`Request Failed: ${r.statusText}`)));
        } else onResp(success(r));
      });
    };

    return [resp, f];
  };

  const isSignedIn = () => isSuccess(user);

  return (
    <AuthContext.Provider
      value={{ user, isSignedIn, login, logout, useAuthFetch }}
    >
      {children}
    </AuthContext.Provider>
  );
};
