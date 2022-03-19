import "./style.css";

import React from "react";
import { render } from "react-dom";
import { App } from "@app/components/App/App";
import { ContextProvider } from "@app/providers/AppContext";

render(
  <ContextProvider>
    <App />
  </ContextProvider>,
  document.getElementsByTagName("main")[0]
);
