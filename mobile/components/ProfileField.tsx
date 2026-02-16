import { View, TextInput } from "react-native";

import { ThemedText } from "@/components/themed-text";
import { ThemedView } from "@/components/themed-view";
import { IconSymbol } from "@/components/ui/icon-symbol";
import { TextInputProps } from "react-native";

type ProfileFieldProps = TextInputProps & {
  label: string
  isEditing:boolean
  suffix?: string
};

export function ProfileField( props: ProfileFieldProps) {
  return (
    <ThemedView>
      <ThemedText className="text-sm font-semibold mb-2">{props.label}</ThemedText>
      <View className="relative">
        <TextInput
          className={`w-full p-4 pr-10 rounded-xl text-stone-950 dark:text-stone-100 text-base ${
            props.isEditing
              ? "bg-stone-100 dark:bg-stone-900 border-2 border-emerald-500"
              : "bg-stone-200 dark:bg-stone-800"
          }`}
          value={props.value}
          onChangeText={props.onChangeText}
          onFocus={props.onFocus}
          onBlur={props.onBlur}
          placeholder={props.placeholder}
          placeholderTextColor="#9CA3AF"
          keyboardType={props.keyboardType}
          multiline={props.multiline}
          numberOfLines={props.multiline ? 4 : 1}
          textAlignVertical={props.multiline ? "top" : "center"}
          style={props.multiline ? { minHeight: 100 } : undefined}
        />
        {props.isEditing && (
          <View className="absolute right-3 top-3">
            <IconSymbol size={20} name="checkmark" color="#10B981" />
          </View>
        )}
      </View>
      {props.suffix && (
        <ThemedText className="text-sm text-stone-500 dark:text-stone-400 mt-1">
          {props.suffix}
        </ThemedText>
      )}
    </ThemedView>
  );
}
