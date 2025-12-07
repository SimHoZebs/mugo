import { useState, createContext, useContext, ReactNode } from "react";
import { View, Pressable, TextInput } from "react-native";
import { ThemedView } from "./themed-view";
import { ThemedText } from "./themed-text";

// Context for InputBar state
interface InputBarContextValue {
  text: string;
  setText: (text: string) => void;
  onSubmit?: (text: string) => void;
  isLoading: boolean;
  disabled: boolean;
}

const InputBarContext = createContext<InputBarContextValue | null>(null);

function useInputBarContext() {
  const context = useContext(InputBarContext);
  if (!context) {
    throw new Error("InputBar components must be used within an InputBar");
  }
  return context;
}

// Main InputBar container
interface InputBarProps {
  children: ReactNode;
  onSubmit?: (text: string) => void;
  isLoading?: boolean;
  disabled?: boolean;
}

function InputBar({
  children,
  onSubmit,
  isLoading = false,
  disabled = false,
}: InputBarProps) {
  const [text, setText] = useState("");

  return (
    <InputBarContext.Provider
      value={{ text, setText, onSubmit, isLoading, disabled }}
    >
      <ThemedView className="p-4 border-t border-stone-200 dark:border-stone-800">
        <View className="flex-row items-center gap-3">{children}</View>
      </ThemedView>
    </InputBarContext.Provider>
  );
}

// Text input with submit button
interface InputProps {
  placeholder?: string;
}

function Input({ placeholder = "Describe your meals..." }: InputProps) {
  const { text, setText, onSubmit, isLoading, disabled } = useInputBarContext();

  const handleSubmit = () => {
    if (!text.trim() || isLoading || disabled) return;
    onSubmit?.(text.trim());
    setText("");
  };

  const hasText = text.trim().length > 0;
  const canSubmit = hasText && !isLoading && !disabled;

  return (
    <View className="flex-1 flex-row justify-center items-center bg-stone-200 dark:bg-stone-800 rounded-2xl px-4 py-2">
      <TextInput
        className="flex-1 w-full text-stone-950 dark:text-stone-100"
        placeholder={placeholder}
        value={text}
        onChangeText={setText}
        onSubmitEditing={handleSubmit}
        editable={!disabled && !isLoading}
        multiline
      />

      <Pressable
        onPress={handleSubmit}
        disabled={!canSubmit}
        className={`w-8 h-8 rounded-full items-center justify-center ${
          canSubmit
            ? "bg-emerald-500 active:bg-emerald-600"
            : "bg-stone-300 dark:bg-stone-700"
        }`}
      >
        <ThemedText className="text-white text-sm font-bold">
          {isLoading ? "•••" : "↑"}
        </ThemedText>
      </Pressable>
    </View>
  );
}

// Action button (for mic, camera, etc.)
interface ActionProps {
  onPress?: () => void;
  children: ReactNode;
}

function Action({ onPress, children }: ActionProps) {
  const { isLoading, disabled } = useInputBarContext();

  return (
    <Pressable
      onPress={onPress}
      disabled={disabled || isLoading}
      className="w-11 h-11 bg-stone-200 dark:bg-stone-800 justify-center items-center rounded-full active:bg-stone-300 dark:active:bg-stone-700"
    >
      {children}
    </Pressable>
  );
}

// Attach subcomponents
InputBar.Input = Input;
InputBar.Action = Action;

export default InputBar;
