import { TextInputProps } from "react-native";

import { ThemedText } from "@/components/themed-text";
import { ThemedView } from "@/components/themed-view";
import { EditableField } from "@/components/EditableField";

type ProfileFieldProps = TextInputProps & {
  label: string;
  isEditing: boolean;
  suffix?: string;
};

export function ProfileField(props: ProfileFieldProps) {
  return (
    <ThemedView>
      <ThemedText className="text-sm font-semibold mb-2">{props.label}</ThemedText>
      <EditableField
        isEditing={props.isEditing}
        value={props.value}
        onChangeText={props.onChangeText}
        onFocus={props.onFocus}
        onBlur={props.onBlur}
        placeholder={props.placeholder}
        keyboardType={props.keyboardType}
        multiline={props.multiline}
        numberOfLines={props.multiline ? 4 : 1}
        textAlignVertical={props.multiline ? "top" : "center"}
        style={props.multiline ? { minHeight: 100 } : undefined}
      />
      {props.suffix && (
        <ThemedText className="text-sm text-stone-500 dark:text-stone-400 mt-1">
          {props.suffix}
        </ThemedText>
      )}
    </ThemedView>
  );
}

