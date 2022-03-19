import React, { createContext, FunctionComponent, useContext } from "react";
import { Api } from "../services/Api";

export type AppContext = {
  api: Api;
};

const defaultContext: AppContext = {
  api: new Api(),
};

const Context = createContext<AppContext>(defaultContext);

export const ContextProvider: FunctionComponent = ({ children }) => (
  <Context.Provider value={defaultContext}>{children}</Context.Provider>
);

export const useAppContext = () => {
  return useContext(Context);
};
