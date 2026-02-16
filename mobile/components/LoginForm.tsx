import { useState } from "react";
import { View, TextInput, Pressable } from "react-native";

import { ThemedText } from "@/components/themed-text";
import { ThemedView } from "@/components/themed-view";
import { IconSymbol } from "@/components/ui/icon-symbol";

type LoginFormProps = {
  onLogin: (username: string) => void;
};

export function LoginForm({ onLogin }: LoginFormProps) {
  const [usernameInput, setUsernameInput] = useState("");

  const canSubmit = usernameInput.trim().length > 0;

  return (
    <ThemedView className="h-full w-full justify-center px-8">
      <View className="mb-12 items-center">
        <View className="w-20 h-20 bg-emerald-500 rounded-3xl items-center justify-center mb-6">
          <IconSymbol size={40} name="person.fill" color="white" />
        </View>
        <ThemedText type="title" className="text-center mb-2">
          Welcome to Mugo
        </ThemedText>
        <ThemedText className="text-stone-500 dark:text-stone-400 text-center">
          Enter your username to get started
        </ThemedText>
      </View>

      <View className="gap-4">
        <View>
          <ThemedText className="text-sm font-semibold mb-2 ml-1">
            Username
          </ThemedText>
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
          <ThemedText
            className={`font-bold text-lg ${
              canSubmit ? "text-white" : "text-stone-500"
            }`}
          >
            Continue
          </ThemedText>
        </Pressable>
      </View>
    </ThemedView>
  );
}
