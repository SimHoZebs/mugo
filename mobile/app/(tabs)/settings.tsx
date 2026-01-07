import { useState } from "react";
import { View, Pressable, ScrollView } from "react-native";

import { ThemedText } from "@/components/themed-text";
import { ThemedView } from "@/components/themed-view";
import useGlobalStore from "@/lib/store";
import { UnitSystem } from "@/lib/types";

export default function SettingsScreen() {
  const userProfile = useGlobalStore((state) => state.userProfile);
  const updateUserProfile = useGlobalStore((state) => state.updateUserProfile);
  const [selectedUnit, setSelectedUnit] = useState<UnitSystem>(
    userProfile.unitSystem,
  );

  const handleUnitChange = (unit: UnitSystem) => {
    setSelectedUnit(unit);
    updateUserProfile({ unitSystem: unit });
  };

  return (
    <ThemedView className="h-full w-full pt-8">
      <View className="px-4 mb-6">
        <ThemedText type="title">Settings</ThemedText>
      </View>

      <ScrollView className="flex-1 px-4">
        <ThemedView className="gap-6">
          <ThemedView>
            <ThemedText className="text-sm font-semibold mb-2">
              Unit System
            </ThemedText>
            <View className="flex-row gap-2">
              {(["metric", "imperial"] as UnitSystem[]).map((system) => (
                <Pressable
                  key={system}
                  onPress={() => handleUnitChange(system)}
                  className={`flex-1 p-4 rounded-xl ${
                    selectedUnit === system
                      ? "bg-emerald-500"
                      : "bg-stone-200 dark:bg-stone-800"
                  }`}
                >
                  <ThemedText
                    className={`text-center font-semibold ${
                      selectedUnit === system
                        ? "text-white"
                        : "text-stone-950 dark:text-stone-100"
                    }`}
                  >
                    {system.charAt(0).toUpperCase() + system.slice(1)}
                  </ThemedText>
                </Pressable>
              ))}
            </View>
            <ThemedText className="text-sm text-stone-500 dark:text-stone-400 mt-2">
              {selectedUnit === "metric"
                ? "Weight in kg, Height in cm"
                : "Weight in lbs, Height in ft"}
            </ThemedText>
          </ThemedView>
        </ThemedView>
      </ScrollView>
    </ThemedView>
  );
}
