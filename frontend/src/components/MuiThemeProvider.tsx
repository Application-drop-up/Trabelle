"use client";

import { ThemeProvider, createTheme } from "@mui/material/styles";

const theme = createTheme({ palette: { mode: "dark" } });

export function MuiThemeProvider({ children }: { children: React.ReactNode }) {
  return <ThemeProvider theme={theme}>{children}</ThemeProvider>;
}
