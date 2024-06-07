import { expect, test } from "vitest";
import { formatSecondsToHMS } from "./utils";

test("formats seconds to HH:MM:SS format", () => {
  expect(formatSecondsToHMS(60)).toBe("00:01:00");
  expect(formatSecondsToHMS(80)).toBe("00:01:20");
  expect(formatSecondsToHMS(8)).toBe("00:00:08");
  expect(formatSecondsToHMS(3640)).toBe("01:00:40");
});
