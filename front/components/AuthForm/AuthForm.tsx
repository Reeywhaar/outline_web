import React, { FunctionComponent, useCallback, useState } from "react";

import classes from "./AuthForm.module.scss";

export const AuthForm: FunctionComponent<{
  onSubmit: (data: { password: string }) => unknown;
}> = ({ onSubmit }) => {
  const [password, setPassword] = useState("");

  const handlePasswordChange = useCallback(
    (e: { target: { value: string } }) => {
      setPassword(e.target.value);
    },
    []
  );

  const handleSubmit = useCallback(
    (e: { preventDefault?: () => unknown }) => {
      e.preventDefault?.();
      onSubmit({ password });
    },
    [password]
  );

  return (
    <form className={classes.root} onSubmit={handleSubmit}>
      <input type="password" value={password} onChange={handlePasswordChange} />
      <button type="submit">Login</button>
    </form>
  );
};
