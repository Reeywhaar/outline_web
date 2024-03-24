import React, {
  FunctionComponent,
  useCallback,
  useEffect,
  useState,
} from "react";
import { sleep } from "@app/utils";
import { ServersResponse } from "@app/services/Api";
import { useAppContext } from "@app/providers/AppContext";
import { Server } from "../Server/Server";
import { AuthForm } from "../AuthForm/AuthForm";

import classes from "./App.module.scss";

export const App: FunctionComponent = () => {
  const [servers, setServers] = useState<ServersResponse>([]);
  const { api } = useAppContext();
  const [authChecked, setAuthChecked] = useState(false);
  const [loggedIn, setLoggedIn] = useState(false);
  const [error, setError] = useState<Error | null>(null);

  const handleAuth = useCallback(async ({ password }: { password: string }) => {
    setError(null);
    const abortController = new AbortController();
    try {
      const resp = await api.auth(password, abortController.signal);
      if (resp.success) {
        setLoggedIn(true);
      }
    } catch (e) {
      setError(e);
    }
  }, []);

  useEffect(() => {
    (async () => {
      const abortController = new AbortController();
      const check = await api.authCheck(abortController.signal);
      setLoggedIn(check);
      setAuthChecked(true);
    })();
  }, []);

  useEffect(() => {
    if (!loggedIn) return;
    const abortController = new AbortController();
    const get = async () => {
      const servers = await api.fetchServers(abortController.signal);
      setServers((oldServers) =>
        oldServers.length === servers.length ? oldServers : servers
      );
    };

    (async () => {
      while (!abortController.signal.aborted) {
        await get();
        await sleep(60000);
      }
    })();

    return () => {
      abortController.abort();
    };
  }, [loggedIn]);

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
