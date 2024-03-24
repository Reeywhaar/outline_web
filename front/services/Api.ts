export type AuthResponse = { success: boolean };
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
  async auth(password: string, abortSignal?: AbortSignal) {
    const password_hash = await makeHash(password);
    const resp = await fetch("/api/auth", {
      method: "POST",
      body: JSON.stringify({ password_hash }),
      signal: abortSignal,
    });
    if (resp.status >= 400) throw await respToError(resp);
    return (await resp.json()) as AuthResponse;
  }
  async authCheck(abortSignal?: AbortSignal) {
    const resp = await fetch("/api/servers", {
      signal: abortSignal,
    });
    if (resp.status === 401 || resp.status === 403) return false;
    if (resp.status >= 400) throw await respToError(resp);
    return true;
  }
  async fetchServers(abortSignal?: AbortSignal) {
    const resp = await fetch("/api/servers", { signal: abortSignal });
    if (resp.status >= 400) throw await respToError(resp);
    return (await resp.json()) as ServersResponse;
  }
  async fetchServer(id: number, abortSignal?: AbortSignal) {
    const resp = await fetch(`/api/servers/${encodeURIComponent(id)}`, {
      signal: abortSignal,
    });
    if (resp.status >= 400) throw await respToError(resp);
    return (await resp.json()) as ServerResponse;
  }
}

async function makeHash(message: string) {
  const msgUint8 = new TextEncoder().encode(message); // encode as (utf-8) Uint8Array
  const hashBuffer = await crypto.subtle.digest("SHA-256", msgUint8); // hash the message
  const hashArray = Array.from(new Uint8Array(hashBuffer)); // convert buffer to byte array
  const hashHex = hashArray
    .map((b) => b.toString(16).padStart(2, "0"))
    .join(""); // convert bytes to hex string
  return hashHex;
}

const respToError = async (resp: Response) => {
  try {
    const text = await resp.text();
    const json = JSON.parse(text);
    if (json.error && json.message) {
      const err = new ResponseError(json.message);
      err.status = resp.status;
      err.response = resp;
      err.raw = text;
      return err;
    }
  } catch (e) {
    return new Error(resp.statusText);
  }
};

class ResponseError extends Error {
  public status: number;
  public response: Response;
  public raw: string;
}
