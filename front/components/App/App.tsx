import React, { FunctionComponent, useEffect, useState } from "react";
import { Api, ServersResponse } from "../../services/api";
import { sleep } from "../../utils";
import { Server } from "../Server/Server";

export const App: FunctionComponent = () => {
  const [servers, setServers] = useState<ServersResponse>([]);

  useEffect(() => {
    const abortController = new AbortController();
    const api = new Api();
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
