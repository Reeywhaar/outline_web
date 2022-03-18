export type ServersResponse = number[];

export type ServerResponse = {
  name: string;
  server_id: string;
  users: {
    name: string;
    usage: number;
  }[];
};

export class Api {
  async fetchServers(abortSignal?: AbortSignal) {
    const resp = await fetch("/api/servers", { signal: abortSignal });
    if (resp.status >= 400) throw new Error(await resp.text());
    return (await resp.json()) as ServersResponse;
  }
  async fetchServer(id: number, abortSignal?: AbortSignal) {
    const resp = await fetch(`/api/servers/${encodeURIComponent(id)}`, {
      signal: abortSignal,
    });
    if (resp.status >= 400) throw new Error(await resp.text());
    return (await resp.json()) as ServerResponse;
  }
}
