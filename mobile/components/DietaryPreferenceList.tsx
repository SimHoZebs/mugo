import { View, TextInput, Pressable } from "react-native";

import { ThemedText } from "@/components/themed-text";
import { ThemedView } from "@/components/themed-view";
import { IconSymbol } from "@/components/ui/icon-symbol";
import { DietaryPreference } from "@/lib/types";

type DietaryPreferenceListProps = {
  preferences: DietaryPreference[];
  editingId: string | null;
  onEditStart: (id: string) => void;
  onEditEnd: () => void;
  onUpdate: (id: string, text: string) => void;
  onRemove: (id: string) => void;
  onAdd: () => void;
};

export function DietaryPreferenceList({
  preferences,
  editingId,
  onEditStart,
  onEditEnd,
  onUpdate,
  onRemove,
  onAdd,
}: DietaryPreferenceListProps) {
  return (
    <ThemedView>
      <ThemedText className="text-sm font-semibold mb-2">
        Dietary Preferences
      </ThemedText>
      <ThemedView className="gap-2">
        {preferences.map((pref) => {
          const isEditing = editingId === pref.id;

          return (
            <View key={pref.id} className="relative">
              <TextInput
                className={`w-full p-4 pr-10 rounded-xl text-stone-950 dark:text-stone-100 text-base ${
                  isEditing
                    ? "bg-stone-100 dark:bg-stone-900 border-2 border-emerald-500"
                    : "bg-stone-200 dark:bg-stone-800"
                }`}
                value={pref.text}
                onChangeText={(text) => onUpdate(pref.id, text)}
                onFocus={() => onEditStart(pref.id)}
                onBlur={onEditEnd}
                placeholder="Add a preference..."
                placeholderTextColor="#9CA3AF"
              />
              {isEditing && (
                <View className="absolute right-3 top-3">
                  <IconSymbol size={20} name="checkmark" color="#10B981" />
                </View>
              )}
              {!isEditing && (
                <Pressable
                  onPress={() => onRemove(pref.id)}
                  className="absolute right-3 top-3"
                >
                  <IconSymbol size={20} name="xmark" color="#9CA3AF" />
                </Pressable>
              )}
            </View>
          );
        })}
        <Pressable
          onPress={onAdd}
          className="w-full p-4 bg-stone-200 dark:bg-stone-800 rounded-xl items-center justify-center border-2 border-dashed border-stone-400 dark:border-stone-600"
        >
          <ThemedText className="text-stone-500 dark:text-stone-400 font-semibold">
            + Add Preference
          </ThemedText>
        </Pressable>
      </ThemedView>
    </ThemedView>
  );
}
