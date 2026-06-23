import { defineConfig } from "orval";

const specUrl = "http://localhost:8888/openapi.json";
const baseUrl = process.env.EXPO_PUBLIC_API_URL || "http://localhost:8888";

export default defineConfig({
  api: {
    input: {
      target: specUrl,
    },
    output: {
      mode: "tags-split",
      target: "./lib/api",
      client: "fetch",
      baseUrl: baseUrl,
    },
  },
});
