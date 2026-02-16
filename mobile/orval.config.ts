import { defineConfig } from "orval";

export default defineConfig({
  api: {
    input: {
      // Huma usually serves the spec here by default
      target: process.env.API_SERVER_URL + "/openapi.json",
    },
    output: {
      mode: "tags-split",
      target: "./lib/api",
      client: "fetch", // Uses native fetch (perfect for React Native)
      baseUrl: process.env.API_SERVER_URL
    },
  },
});
