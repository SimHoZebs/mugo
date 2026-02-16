import { useState } from "react";
import { View, ScrollView, Pressable } from "react-native";

import { ThemedText } from "@/components/themed-text";
import { ThemedView } from "@/components/themed-view";
import { IconSymbol } from "@/components/ui/icon-symbol";
import { LoginForm } from "@/components/LoginForm";
import { ProfileField } from "@/components/ProfileField";
import { DietaryPreferenceList } from "@/components/DietaryPreferenceList";
import useGlobalStore from "@/lib/store";
import { v7 as uuid7 } from "uuid";

export default function ProfileScreen() {
  const userProfile = useGlobalStore((state) => state.userProfile);
  const updateUserProfile = useGlobalStore((state) => state.updateUserProfile);
  const addDietaryPreference = useGlobalStore(
    (state) => state.addDietaryPreference,
  );
  const updateDietaryPreference = useGlobalStore(
    (state) => state.updateDietaryPreference,
  );
  const removeDietaryPreference = useGlobalStore(
    (state) => state.removeDietaryPreference,
  );
  const login = useGlobalStore((state) => state.login);
  const logout = useGlobalStore((state) => state.logout);

  const [editingField, setEditingField] = useState<string | null>(null);

  if (!userProfile) {
    return <LoginForm onLogin={login} />;
  }

  const weightUnit = userProfile.unitSystem === "metric" ? "kg" : "lbs";
  const heightUnit = userProfile.unitSystem === "metric" ? "cm" : "ft";

  return (
    <ThemedView className="h-full w-full pt-8">
      <View className="px-4 mb-6 flex-row justify-between items-end">
        <View>
          <ThemedText type="title">Profile</ThemedText>
          <ThemedText className="text-stone-500">
            @{userProfile.username}
          </ThemedText>
        </View>
        <Pressable
          onPress={logout}
          className="bg-stone-200 dark:bg-stone-800 p-2 rounded-lg"
        >
          <ThemedText className="text-red-500 font-semibold">Logout</ThemedText>
        </Pressable>
      </View>

      <ScrollView className="flex-1 px-4">
        <ThemedView className="gap-6 pb-12">
          <ThemedView className="items-center py-4">
            <View className="w-24 h-24 bg-stone-300 dark:bg-stone-700 rounded-full items-center justify-center">
              <IconSymbol size={48} name="person.fill" color="#9CA3AF" />
            </View>
          </ThemedView>

          <ProfileField
            label="Name"
            value={userProfile.name}
            placeholder="Your Name"
            isEditing={editingField === "name"}
            onChangeText={(text) => updateUserProfile({ name: text })}
            onFocus={() => setEditingField("name")}
            onBlur={() => setEditingField(null)}
          />

          <DietaryPreferenceList
            preferences={userProfile.dietaryPreferences}
            editingId={
              editingField?.startsWith("pref-")
                ? editingField.replace("pref-", "")
                : null
            }
            onEditStart={(id) => setEditingField(`pref-${id}`)}
            onEditEnd={() => setEditingField(null)}
            onUpdate={updateDietaryPreference}
            onRemove={removeDietaryPreference}
            onAdd={() => addDietaryPreference({ id: uuid7(), text: "" })}
          />

          <ThemedView className="flex-row gap-4">
            <ThemedView className="flex-1">
              <ProfileField
                label="Weight"
                value={userProfile.weight === 0 ? "" : String(userProfile.weight)}
                placeholder="0"
                keyboardType="numeric"
                isEditing={editingField === "weight"}
                onChangeText={(text) => {
                  const numValue = text === "" ? 0 : Number(text);
                  updateUserProfile({ weight: numValue });
                }}
                onFocus={() => setEditingField("weight")}
                onBlur={() => setEditingField(null)}
                suffix={weightUnit}
              />
            </ThemedView>

            <ThemedView className="flex-1">
              <ProfileField
                label="Height"
                value={userProfile.height === 0 ? "" : String(userProfile.height)}
                placeholder="0"
                keyboardType="numeric"
                isEditing={editingField === "height"}
                onChangeText={(text) => {
                  const numValue = text === "" ? 0 : Number(text);
                  updateUserProfile({ height: numValue });
                }}
                onFocus={() => setEditingField("height")}
                onBlur={() => setEditingField(null)}
                suffix={heightUnit}
              />
            </ThemedView>
          </ThemedView>
        </ThemedView>
      </ScrollView>
    </ThemedView>
  );
}
