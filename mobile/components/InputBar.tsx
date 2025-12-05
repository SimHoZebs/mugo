import { useState } from "react";
import {
  Button,
  Pressable,
  TextInput,
  TextInputSubmitEditingEvent,
} from "react-native";
import { ThemedView } from "./themed-view";
import { ThemedText } from "./themed-text";

interface Props {
  onSubmitEditing: (event: TextInputSubmitEditingEvent) => void;
  disabled?: boolean;
}

const InputBar = (props: Props) => {
  const [text, onChangeText] = useState("");

  return (
    <ThemedView className="flex-row w-full justify-evenly gap-4">
      <Pressable className="w-16 h-16 dark:bg-stone-900 justify-center items-center rounded-full">
        <ThemedText>Mic</ThemedText>
      </Pressable>

      <TextInput
        className="dark:bg-stone-900 dark:text-stone-100 bg-stone-200 flex-1"
        onChangeText={onChangeText}
        value={text}
        onSubmitEditing={props.onSubmitEditing}
        placeholder="Describe your meals..."
        editable={!props.disabled}
      />

      <Pressable className="w-16 h-16 dark:bg-stone-900 justify-center items-center rounded-full">
        <ThemedText>Cam</ThemedText>
      </Pressable>
    </ThemedView>
  );
};

export default InputBar;
