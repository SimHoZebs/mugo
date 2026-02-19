import { TextInputProps, View } from "react-native";

import { EditableField } from "@/components/EditableField";
import { Text } from "@/components/ui/Text";

type ProfileFieldProps = TextInputProps & {
  label: string;
  isEditing: boolean;
  suffix?: string;
};

export function ProfileField(props: ProfileFieldProps) {
  return (
    <View className="bg-stone-50 dark:bg-stone-950">
      <Text variant="caption" className="font-semibold mb-2">{props.label}</Text>
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
        <Text variant="caption" className="mt-1">
          {props.suffix}
        </Text>
      )}
    </View>
  );
}

