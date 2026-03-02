import { useState } from "react";
import { View, TextInput } from "react-native";

import { IconSymbol } from "@/components/ui/icon-symbol";
import { Text } from "@/components/ui/Text";
import { Button } from "@/components/ui/Button";

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
        <Text variant="h1" className="text-center mb-2">
          Welcome to Mugo
        </Text>
        <Text variant="caption" className="text-center">
          Enter your username to get started
        </Text>
      </View>

      <View className="gap-4">
        <View>
          <Text variant="caption" className="font-semibold mb-2 ml-1">
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

        <Button
          variant="primary"
          label="Continue"
          disabled={!canSubmit}
          onPress={() => canSubmit && onLogin(usernameInput.trim())}
        />
      </View>
    </View>
  );
}
