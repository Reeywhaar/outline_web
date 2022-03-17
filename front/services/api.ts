export type ApiData = {
  name: string;
  server_id: string;
  users: {
    name: string;
    usage: number;
  }[];
};

export class Api {
  async fetchData(abortSignal?: AbortSignal) {
    const resp = await fetch("/api/data", { signal: abortSignal });
    if (resp.status >= 400) throw new Error(await resp.text());
    return (await resp.json()) as ApiData;
  }
}
