import { useState } from "react";
import { View, ScrollView, Pressable, Text } from "react-native";

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

  const weightUnit = userProfile?.unitSystem === "metric" ? "kg" : "lbs";
  const heightUnit = userProfile?.unitSystem === "metric" ? "cm" : "ft";

  return !userProfile ? (
    <LoginForm onLogin={login} />
  ) : (
    <View className="h-full w-full pt-8 bg-stone-50 dark:bg-stone-950">
      <View className="px-4 mb-6 flex-row justify-between items-end">
        <View>
          <Text className="text-2xl font-bold leading-8 text-stone-950 dark:text-stone-50">Profile</Text>
          <Text className="text-stone-500 text-stone-950 dark:text-stone-50">
            @{userProfile.username}
          </Text>
        </View>
        <Pressable
          onPress={logout}
          className="bg-stone-200 dark:bg-stone-800 p-2 rounded-lg"
        >
          <Text className="text-red-500 font-semibold">Logout</Text>
        </Pressable>
      </View>

      <ScrollView className="flex-1 px-4">
        <View className="gap-6 pb-12 bg-stone-50 dark:bg-stone-950">
          <View className="items-center py-4 bg-stone-50 dark:bg-stone-950">
            <View className="w-24 h-24 bg-stone-300 dark:bg-stone-700 rounded-full items-center justify-center">
              <IconSymbol size={48} name="person.fill" color="#9CA3AF" />
            </View>
          </View>

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

          <View className="flex-row gap-4 bg-stone-50 dark:bg-stone-950">
            <View className="flex-1 bg-stone-50 dark:bg-stone-950">
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
            </View>

            <View className="flex-1 bg-stone-50 dark:bg-stone-950">
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
            </View>
          </View>
        </View>
      </ScrollView>
    </View>
  );
}
