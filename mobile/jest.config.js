module.exports = {
  preset: "jest-expo",
  setupFilesAfterEnv: ["@testing-library/jest-dom"],
  moduleNameMapper: {
    "^@/(.*)$": "<rootDir>/$1",
  },
  testMatch: ["**/test/**/*.test.{ts,tsx}", "**/test/**/*.spec.{ts,tsx}"],
  collectCoverageFrom: [
    "lib/**/*.{ts,tsx}",
    "components/**/*.{ts,tsx}",
    "hooks/**/*.{ts,tsx}",
    "!lib/api/**/*",
  ],
  coveragePathIgnorePatterns: ["/node_modules/", "/test/"],
};
