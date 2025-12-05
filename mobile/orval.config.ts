import { defineConfig } from "orval";

export default defineConfig({
  api: {
    input: {
      // Huma usually serves the spec here by default
      target: "http://localhost:8888/openapi.json",
    },
    output: {
      mode: "tags-split",
      target: "./lib/api",
      client: "fetch", // Uses native fetch (perfect for React Native)
      baseUrl: process.env.EXPO_PUBLIC_API_URL || "http://192.168.1.216:8888",
    },
  },
});
