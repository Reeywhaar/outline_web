import { useEffect, useState } from "react";
import { useDisposableEffect } from "./useDisposableEffect";

export function useVisible() {
  const [visible, setVisible] = useState(!document.hidden);

  useDisposableEffect((stack) => {
    const handleChange = () => {
      setVisible(!document.hidden);
    };

    document.addEventListener("visibilitychange", handleChange);
    stack.defer(() => {
      document.removeEventListener("visibilitychange", handleChange);
    });
  }, []);

  return visible;
}
