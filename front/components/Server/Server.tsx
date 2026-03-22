import React, { useEffect, useState, FunctionComponent } from "react";
import { ServerResponse } from "@app/services/Api";
import { useAppContext } from "@app/providers/AppContext";
import { sleep, spawn } from "@app/utils";
import { useVisible } from "@app/hooks/useVisible";

import classes from "./Server.module.scss";
import { useDisposableEffect } from "@app/hooks/useDisposableEffect";

export const Server: FunctionComponent<{ id: number }> = ({ id }) => {
  const [data, setData] = useState<ServerResponse | null>(null);
  const [err, setError] = useState(null);
  const { api } = useAppContext();
  const visible = useVisible();

  const total = data?.users.reduce((c, x) => c + x.usage, 0) ?? 0;

  useDisposableEffect(
    (stack) => {
      if (!visible) return;

      stack.use(
        spawn(async (stack) => {
          const abortController = stack.adopt(new AbortController(), (ab) =>
            ab.abort(),
          );

          while (!abortController.signal.aborted) {
            try {
              setData(await api.fetchServer(id, abortController.signal));
              setError(null);
            } catch (e) {
              setError(e.message);
            }
            await sleep(5000, abortController.signal);
          }
        }),
      );
    },
    [visible],
  );

  return (
    <div className={classes.root}>
      <h1>{data ? data.name : "Loading..."}</h1>

      {data?.users.map((item) => (
        <div
          className={classes.item}
          key={item.name}
          title={`Total: ${item.usage}`}
        >
          <span className={item.name ? undefined : classes.dim}>
            {item.name || "Unknown"}
          </span>
          : {humanSize(item.usage)}
        </div>
      ))}
      {!!data?.users.length && (
        <>
          <div className={classes.spacer}></div>
          <div title={`Total: ${total}`}>
            <b>Total:</b> {humanSize(total)}
          </div>
        </>
      )}
      {!!err && (
        <>
          <div className={classes.spacer}></div>
          <div className={classes.error}>Error: {err}</div>
        </>
      )}
    </div>
  );
};

const KB = 1024;
const MB = KB * 1024;
const GB = MB * 1024;

const humanSize = (value: number) => {
  if (value >= GB) return `${(value / GB).toFixed(2)}GB`;
  if (value >= MB) return `${(value / MB).toFixed(2)}MB`;
  if (value >= KB) return `${(value / KB).toFixed(2)}KB`;
  return `${value} bytes`;
};
