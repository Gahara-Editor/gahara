import { describe, it } from "vitest";
import { render } from "@testing-library/svelte";
import App from "./App.svelte";

describe("App.svelte tests", () => {
  it("displays main view on load", () => {
    const { getByText } = render(App);
    getByText("Gahara");
  });
});
