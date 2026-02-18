import { TextInputProps, View, Text } from "react-native";

import { EditableField } from "@/components/EditableField";

type ProfileFieldProps = TextInputProps & {
  label: string;
  isEditing: boolean;
  suffix?: string;
};

export function ProfileField(props: ProfileFieldProps) {
  return (
    <View className="bg-stone-50 dark:bg-stone-950">
      <Text className="text-sm font-semibold mb-2 text-stone-950 dark:text-stone-50">{props.label}</Text>
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
        <Text className="text-sm text-stone-500 dark:text-stone-400 mt-1">
          {props.suffix}
        </Text>
      )}
    </View>
  );
}

