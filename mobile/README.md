# Welcome to your Expo app 👋

This is an [Expo](https://expo.dev) project created with [`create-expo-app`](https://www.npmjs.com/package/create-expo-app).

## Get started

1. Install dependencies

   ```bash
   npm install
   ```

2. Start the app

   ```bash
   npx expo start
   ```

In the output, you'll find options to open the app in a

- [development build](https://docs.expo.dev/develop/development-builds/introduction/)
- [Android emulator](https://docs.expo.dev/workflow/android-studio-emulator/)
- [iOS simulator](https://docs.expo.dev/workflow/ios-simulator/)
- [Expo Go](https://expo.dev/go), a limited sandbox for trying out app development with Expo

You can start developing by editing the files inside the **app** directory. This project uses [file-based routing](https://docs.expo.dev/router/introduction).

## Get a fresh project

When you're ready, run:

```bash
npm run reset-project
```

This command will move the starter code to the **app-example** directory and create a blank **app** directory where you can start developing.

## Learn more

To learn more about developing your project with Expo, look at the following resources:

- [Expo documentation](https://docs.expo.dev/): Learn fundamentals, or go into advanced topics with our [guides](https://docs.expo.dev/guides).
- [Learn Expo tutorial](https://docs.expo.dev/tutorial/introduction/): Follow a step-by-step tutorial where you'll create a project that runs on Android, iOS, and the web.

## Join the community

Join our community of developers creating universal apps.

- [Expo on GitHub](https://github.com/expo/expo): View our open source platform and contribute.
- [Discord community](https://chat.expo.dev): Chat with Expo users and ask questions.

## End-to-End Testing with Detox

This app uses [Detox](https://wix.github.io/Detox/) for End-to-End (E2E) testing. Due to the high cost and slowness of running mobile emulators in CI/CD pipelines, E2E tests are configured to be run **manually and locally** before major releases, rather than automatically on every pull request.

### Prerequisites

1.  **Detox CLI:** Install the Detox CLI globally on your machine:
    ```bash
    npm install -g detox-cli
    ```
2.  **Environment Setup:** Follow the [Detox Environment Setup Guide](https://wix.github.io/Detox/docs/introduction/environment-setup) to ensure you have the necessary tools for iOS (Xcode, applesimutils) or Android (Android Studio, Java, Emulators).
3.  **Running Backend:** The app needs to communicate with the backend. Ensure the Go server and database are running locally:
    ```bash
    make db
    make server
    ```

### Running Tests

**For Android:**
1.  Ensure you have an Android emulator created (e.g., `Pixel_3a_API_30_x86`). If your emulator has a different name, update `mobile/.detoxrc.js` to match it.
2.  Start your Android emulator.
3.  Build the app for testing:
    ```bash
    pnpm e2e:build-android
    ```
4.  Run the tests:
    ```bash
    pnpm e2e:test-android
    ```

**For iOS:**
1.  Ensure you have an iOS Simulator installed (e.g., `iPhone 15`). Update `mobile/.detoxrc.js` if you are using a different simulator name or iOS version.
2.  Build the app for testing:
    ```bash
    pnpm e2e:build-ios
    ```
3.  Run the tests:
    ```bash
    pnpm e2e:test-ios
    ```
