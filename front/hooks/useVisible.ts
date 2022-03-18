import { useEffect, useState } from "react";

export function useVisible() {
  const [visible, setVisible] = useState(!document.hidden);

  useEffect(() => {
    const handleChange = () => {
      setVisible(!document.hidden);
    };

    document.addEventListener("visibilitychange", handleChange);

    return () => {
      document.removeEventListener("visibilitychange", handleChange);
    };
  }, []);

  return visible;
}
