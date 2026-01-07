import { useState } from "react";
import {
  View,
  ScrollView,
  TextInput,
  Pressable,
} from "react-native";

import { ThemedText } from "@/components/themed-text";
import { ThemedView } from "@/components/themed-view";
import { IconSymbol } from "@/components/ui/icon-symbol";
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

  const [editingField, setEditingField] = useState<string | null>(null);

  const getWeightUnit = () => (userProfile.unitSystem === "metric" ? "kg" : "lbs");
  const getHeightUnit = () => (userProfile.unitSystem === "metric" ? "cm" : "ft");

  const handleSave = (field: string, value: any) => {
    updateUserProfile({ [field]: value });
    setEditingField(null);
  };

  const InputField = ({
    label,
    field,
    value,
    placeholder,
    keyboardType = "default",
    multiline = false,
    onSubmit,
  }: {
    label: string;
    field: string;
    value: any;
    placeholder: string;
    keyboardType?: "default" | "numeric";
    multiline?: boolean;
    onSubmit?: () => void;
  }) => {
    const isEditing = editingField === field;

    return (
      <ThemedView>
        <ThemedText className="text-sm font-semibold mb-2">{label}</ThemedText>
        <View className="relative">
          <TextInput
            className={`w-full p-4 pr-10 rounded-xl text-stone-950 dark:text-stone-100 text-base ${
              isEditing
                ? "bg-stone-100 dark:bg-stone-900 border-2 border-emerald-500"
                : "bg-stone-200 dark:bg-stone-800"
            }`}
            value={value}
            onChangeText={(text) => {
              if (multiline) {
                handleSave(field, text);
              } else if (keyboardType === "numeric") {
                const numValue = text === "" ? 0 : Number(text);
                handleSave(field, numValue);
              } else {
                handleSave(field, text);
              }
            }}
            onFocus={() => setEditingField(field)}
            onBlur={() => setEditingField(null)}
            placeholder={placeholder}
            placeholderTextColor="#9CA3AF"
            keyboardType={keyboardType}
            multiline={multiline}
            numberOfLines={multiline ? 4 : 1}
            textAlignVertical={multiline ? "top" : "center"}
            style={multiline ? { minHeight: 100 } : undefined}
          />
          {isEditing && (
            <View className="absolute right-3 top-3">
              <IconSymbol size={20} name="checkmark" color="#10B981" />
            </View>
          )}
        </View>
      </ThemedView>
    );
  };

  return (
    <ThemedView className="h-full w-full pt-8">
      <View className="px-4 mb-6">
        <ThemedText type="title">Profile</ThemedText>
      </View>

      <ScrollView className="flex-1 px-4">
        <ThemedView className="gap-6">
          <ThemedView className="items-center py-8">
            <View className="w-24 h-24 bg-stone-300 dark:bg-stone-700 rounded-full items-center justify-center">
              <IconSymbol size={48} name="person.fill" color="#9CA3AF" />
            </View>
          </ThemedView>

          <InputField
            label="Name"
            field="name"
            value={userProfile.name}
            placeholder="Your Name"
          />

          <ThemedView>
            <ThemedText className="text-sm font-semibold mb-2">
              Dietary Preferences
            </ThemedText>
            <ThemedView className="gap-2">
              {userProfile.dietaryPreferences.map((pref, index) => (
                <View key={pref.id} className="relative">
                  <TextInput
                    className={`w-full p-4 pr-10 rounded-xl text-stone-950 dark:text-stone-100 text-base ${
                      editingField === `pref-${pref.id}`
                        ? "bg-stone-100 dark:bg-stone-900 border-2 border-emerald-500"
                        : "bg-stone-200 dark:bg-stone-800"
                    }`}
                    value={pref.text}
                    onChangeText={(text) =>
                      updateDietaryPreference(pref.id, text)
                    }
                    onFocus={() => setEditingField(`pref-${pref.id}`)}
                    onBlur={() => setEditingField(null)}
                    placeholder="Add a preference..."
                    placeholderTextColor="#9CA3AF"
                  />
                  {editingField === `pref-${pref.id}` && (
                    <View className="absolute right-3 top-3">
                      <IconSymbol size={20} name="checkmark" color="#10B981" />
                    </View>
                  )}
                  {editingField !== `pref-${pref.id}` && (
                    <Pressable
                      onPress={() => removeDietaryPreference(pref.id)}
                      className="absolute right-3 top-3"
                    >
                      <IconSymbol size={20} name="xmark" color="#9CA3AF" />
                    </Pressable>
                  )}
                </View>
              ))}
              <Pressable
                onPress={() =>
                  addDietaryPreference({
                    id: uuid7(),
                    text: "",
                  })
                }
                className="w-full p-4 bg-stone-200 dark:bg-stone-800 rounded-xl items-center justify-center border-2 border-dashed border-stone-400 dark:border-stone-600"
              >
                <ThemedText className="text-stone-500 dark:text-stone-400 font-semibold">
                  + Add Preference
                </ThemedText>
              </Pressable>
            </ThemedView>
          </ThemedView>

          <ThemedView className="flex-row gap-4">
            <ThemedView className="flex-1">
              <ThemedText className="text-sm font-semibold mb-2">
                Weight
              </ThemedText>
              <View className="relative">
                <TextInput
                  className={`w-full p-4 rounded-xl text-stone-950 dark:text-stone-100 text-base ${
                    editingField === "weight"
                      ? "bg-stone-100 dark:bg-stone-900 border-2 border-emerald-500"
                      : "bg-stone-200 dark:bg-stone-800"
                  }`}
                  value={userProfile.weight === 0 ? "" : String(userProfile.weight)}
                  onChangeText={(text) => {
                    const numValue = text === "" ? 0 : Number(text);
                    handleSave("weight", numValue);
                  }}
                  onFocus={() => setEditingField("weight")}
                  onBlur={() => setEditingField(null)}
                  placeholder="0"
                  keyboardType="numeric"
                  placeholderTextColor="#9CA3AF"
                />
                {editingField === "weight" && (
                  <View className="absolute right-3 top-3">
                    <IconSymbol size={20} name="checkmark" color="#10B981" />
                  </View>
                )}
              </View>
              <ThemedText className="text-sm text-stone-500 dark:text-stone-400 mt-1">
                {getWeightUnit()}
              </ThemedText>
            </ThemedView>

            <ThemedView className="flex-1">
              <ThemedText className="text-sm font-semibold mb-2">
                Height
              </ThemedText>
              <View className="relative">
                <TextInput
                  className={`w-full p-4 rounded-xl text-stone-950 dark:text-stone-100 text-base ${
                    editingField === "height"
                      ? "bg-stone-100 dark:bg-stone-900 border-2 border-emerald-500"
                      : "bg-stone-200 dark:bg-stone-800"
                  }`}
                  value={userProfile.height === 0 ? "" : String(userProfile.height)}
                  onChangeText={(text) => {
                    const numValue = text === "" ? 0 : Number(text);
                    handleSave("height", numValue);
                  }}
                  onFocus={() => setEditingField("height")}
                  onBlur={() => setEditingField(null)}
                  placeholder="0"
                  keyboardType="numeric"
                  placeholderTextColor="#9CA3AF"
                />
                {editingField === "height" && (
                  <View className="absolute right-3 top-3 top-[50%] -translate-y-1/2">
                    <IconSymbol size={20} name="checkmark" color="#10B981" />
                  </View>
                )}
              </View>
              <ThemedText className="text-sm text-stone-500 dark:text-stone-400 mt-1">
                {getHeightUnit()}
              </ThemedText>
            </ThemedView>
          </ThemedView>
        </ThemedView>
      </ScrollView>
    </ThemedView>
  );
}
