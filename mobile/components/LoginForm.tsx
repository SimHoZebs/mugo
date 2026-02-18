import { useState } from "react";
import { View, TextInput, Pressable, Text } from "react-native";

import { IconSymbol } from "@/components/ui/icon-symbol";

type LoginFormProps = {
  onLogin: (username: string) => void;
};

export function LoginForm({ onLogin }: LoginFormProps) {
  const [usernameInput, setUsernameInput] = useState("");

  const canSubmit = usernameInput.trim().length > 0;

  return (
    <View className="h-full w-full justify-center px-8 bg-stone-50 dark:bg-stone-950">
      <View className="mb-12 items-center">
        <View className="w-20 h-20 bg-emerald-500 rounded-3xl items-center justify-center mb-6">
          <IconSymbol size={40} name="person.fill" color="white" />
        </View>
        <Text className="text-center mb-2 text-2xl font-bold leading-8 text-stone-950 dark:text-stone-50">
          Welcome to Mugo
        </Text>
        <Text className="text-stone-500 dark:text-stone-400 text-center">
          Enter your username to get started
        </Text>
      </View>

      <View className="gap-4">
        <View>
          <Text className="text-sm font-semibold mb-2 ml-1 text-stone-950 dark:text-stone-50">
            Username
          </Text>
          <TextInput
            className="w-full p-4 rounded-2xl bg-white dark:bg-stone-900 border border-stone-200 dark:border-stone-800 text-stone-950 dark:text-stone-100 text-lg"
            placeholder="e.g. johndoe"
            placeholderTextColor="#9CA3AF"
            value={usernameInput}
            onChangeText={setUsernameInput}
            autoCapitalize="none"
            autoCorrect={false}
            onSubmitEditing={() => canSubmit && onLogin(usernameInput.trim())}
          />
        </View>

        <Pressable
          onPress={() => canSubmit && onLogin(usernameInput.trim())}
          className={`w-full p-5 rounded-2xl items-center justify-center ${
            canSubmit
              ? "bg-emerald-500"
              : "bg-stone-300 dark:bg-stone-800"
          }`}
        >
          <Text
            className={`font-bold text-lg ${
              canSubmit ? "text-white" : "text-stone-500"
            }`}
          >
            Continue
          </Text>
        </Pressable>
      </View>
    </View>
  );
}
