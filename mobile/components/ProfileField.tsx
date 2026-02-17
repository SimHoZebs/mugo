import { TextInputProps } from "react-native";

import { ThemedText } from "@/components/themed-text";
import { ThemedView } from "@/components/themed-view";
import { EditableField } from "@/components/EditableField";

type ProfileFieldProps = TextInputProps & {
  label: string;
  isEditing: boolean;
  suffix?: string;
};

export function ProfileField({ label, isEditing, suffix, multiline, ...rest }: ProfileFieldProps) {
  return (
    <ThemedView>
      <ThemedText className="text-sm font-semibold mb-2">{label}</ThemedText>
      <EditableField
        isEditing={isEditing}
        multiline={multiline}
        numberOfLines={multiline ? 4 : 1}
        textAlignVertical={multiline ? "top" : "center"}
        style={multiline ? { minHeight: 100 } : undefined}
        {...rest}
      />
      {suffix && (
        <ThemedText className="text-sm text-stone-500 dark:text-stone-400 mt-1">
          {suffix}
        </ThemedText>
      )}
    </ThemedView>
  );
}

