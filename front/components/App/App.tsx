import React, {
  FunctionComponent,
  useCallback,
  useEffect,
  useState,
} from "react";
import { sleep, spawn } from "@app/utils";
import { ServersResponse } from "@app/services/Api";
import { useAppContext } from "@app/providers/AppContext";
import { Server } from "../Server/Server";
import { AuthForm } from "../AuthForm/AuthForm";

import classes from "./App.module.scss";
import { useDisposableEffect } from "@app/hooks/useDisposableEffect";

export const App: FunctionComponent = () => {
  const [servers, setServers] = useState<ServersResponse>([]);
  const { api } = useAppContext();
  const [authChecked, setAuthChecked] = useState(false);
  const [loggedIn, setLoggedIn] = useState(false);
  const [error, setError] = useState<Error | null>(null);

  const handleAuth = useCallback(async ({ password }: { password: string }) => {
    spawn(
      async (stack) => {
        setError(null);
        const abortController = stack.adopt(new AbortController(), (ac) =>
          ac.abort(),
        );
        const resp = await api.auth(password, abortController.signal);
        if (resp.success) {
          setLoggedIn(true);
        }
      },
      (e) => {
        setError(e);
      },
    );
  }, []);

  useDisposableEffect((stack) => {
    stack.use(
      spawn(async (stack) => {
        const abortController = stack.adopt(new AbortController(), (ac) =>
          ac.abort(),
        );
        const check = await api.authCheck(abortController.signal);
        setLoggedIn(check);
        setAuthChecked(true);
      }),
    );
  }, []);

  useDisposableEffect(
    (stack) => {
      if (!loggedIn) return;
      stack.use(
        spawn(async (stack) => {
          const abortController = stack.adopt(new AbortController(), (ac) =>
            ac.abort(),
          );

          while (!abortController.signal.aborted) {
            const servers = await api.fetchServers(abortController.signal);
            setServers((oldServers) =>
              oldServers.length === servers.length ? oldServers : servers,
            );
            await sleep(60000, abortController.signal);
          }
        }),
      );
    },
    [loggedIn],
  );

  return (
    <div>
      {!authChecked ? (
        "loading..."
      ) : !loggedIn ? (
        <AuthForm onSubmit={handleAuth} />
      ) : (
        servers.map((id) => <Server key={id} id={id} />)
      )}
      {error ? <div className={classes.error}>{error.message}</div> : null}
    </div>
  );
};
