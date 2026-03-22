import "./style.css";

import React from "react";
import { createRoot } from "react-dom/client";
import { App } from "@app/components/App/App";
import { ContextProvider } from "@app/providers/AppContext";

const root = createRoot(document.getElementsByTagName("main")[0]);
root.render(
  <ContextProvider>
    <App />
  </ContextProvider>,
);
