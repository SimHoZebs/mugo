import { View, TextInput, Pressable, TextInputProps } from "react-native";

import { IconSymbol } from "@/components/ui/icon-symbol";

type EditableFieldProps = TextInputProps & {
  isEditing: boolean;
  onRemove?: () => void;
};

export function EditableField(props: EditableFieldProps) {
  const { isEditing, onRemove, ...textInputProps } = props;
  return (
    <View className="relative">
      <TextInput
        className={`w-full p-4 pr-10 rounded-xl text-stone-950 dark:text-stone-100 text-base ${
          isEditing
            ? "bg-stone-100 dark:bg-stone-900 border-2 border-emerald-500"
            : "bg-stone-200 dark:bg-stone-800"
        }`}
        placeholderTextColor="#9CA3AF"
        {...textInputProps}
      />
      {isEditing && (
        <View className="absolute right-3 top-3">
          <IconSymbol size={20} name="checkmark" color="#10B981" />
        </View>
      )}
      {!isEditing && onRemove && (
        <Pressable onPress={onRemove} className="absolute right-3 top-3">
          <IconSymbol size={20} name="xmark" color="#9CA3AF" />
        </Pressable>
      )}
    </View>
  );
}

