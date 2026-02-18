import { Pressable } from "react-native";

import { ThemedText } from "@/components/themed-text";
import { ThemedView } from "@/components/themed-view";
import { EditableField } from "@/components/EditableField";
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

export function DietaryPreferenceList(props: DietaryPreferenceListProps) {
  return (
    <ThemedView>
      <ThemedText className="text-sm font-semibold mb-2">
        Dietary Preferences
      </ThemedText>
      <ThemedView className="gap-2">
        {props.preferences.map((pref) => (
          <EditableField
            key={pref.id}
            isEditing={props.editingId === pref.id}
            value={pref.text}
            onChangeText={(text) => props.onUpdate(pref.id, text)}
            onFocus={() => props.onEditStart(pref.id)}
            onBlur={props.onEditEnd}
            placeholder="Add a preference..."
            onRemove={() => props.onRemove(pref.id)}
          />
        ))}
        <Pressable
          onPress={props.onAdd}
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

