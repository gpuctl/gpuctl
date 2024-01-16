import React from "react";
import { render, screen } from "@testing-library/react";
import App from "./App";

test("renders welcome message", () => {
  render(<App />);
  const welcome = screen.getByText(/Welcome to the GPU Control Room!/i);
  expect(welcome).toBeInTheDocument();
});
