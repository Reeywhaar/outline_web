import React, { FunctionComponent, useEffect, useState } from "react";
import { sleep } from "@app/utils";
import { ServersResponse } from "@app/services/Api";
import { useAppContext } from "@app/providers/AppContext";
import { Server } from "../Server/Server";

export const App: FunctionComponent = () => {
  const [servers, setServers] = useState<ServersResponse>([]);
  const { api } = useAppContext();

  useEffect(() => {
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
  }, []);

  return (
    <div>
      {servers.map((id) => (
        <Server key={id} id={id} />
      ))}
    </div>
  );
};
