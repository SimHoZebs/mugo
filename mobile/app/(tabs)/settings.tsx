import { useState, useEffect } from "react";
import { View, Pressable, ScrollView, Text } from "react-native";

import useGlobalStore from "@/lib/store";
import { UnitSystem } from "@/lib/types";

export default function SettingsScreen() {
  const userProfile = useGlobalStore((state) => state.userProfile);
  const updateUserProfile = useGlobalStore((state) => state.updateUserProfile);
  const [selectedUnit, setSelectedUnit] = useState<UnitSystem>(
    userProfile?.unitSystem || "metric",
  );

  useEffect(() => {
    if (userProfile?.unitSystem) {
      setSelectedUnit(userProfile.unitSystem);
    }
  }, [userProfile?.unitSystem]);

  const handleUnitChange = (unit: UnitSystem) => {
    setSelectedUnit(unit);
    updateUserProfile({ unitSystem: unit });
  };

  return !userProfile ? (
    <View className="h-full w-full justify-center items-center px-4 bg-stone-50 dark:bg-stone-950">
      <Text className="text-stone-500 text-center">
        Please sign in on the Profile tab to access settings.
      </Text>
    </View>
  ) : (
    <View className="h-full w-full pt-8 bg-stone-50 dark:bg-stone-950">
      <View className="px-4 mb-6">
        <Text className="text-2xl font-bold leading-8 text-stone-950 dark:text-stone-50">Settings</Text>
      </View>

      <ScrollView className="flex-1 px-4">
        <View className="gap-6 bg-stone-50 dark:bg-stone-950">
          <View className="bg-stone-50 dark:bg-stone-950">
            <Text className="text-sm font-semibold mb-2 text-stone-950 dark:text-stone-50">
              Unit System
            </Text>
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
                  <Text
                    className={`text-center font-semibold ${
                      selectedUnit === system
                        ? "text-white"
                        : "text-stone-950 dark:text-stone-100"
                    }`}
                  >
                    {system.charAt(0).toUpperCase() + system.slice(1)}
                  </Text>
                </Pressable>
              ))}
            </View>
            <Text className="text-sm text-stone-500 dark:text-stone-400 mt-2">
              {selectedUnit === "metric"
                ? "Weight in kg, Height in cm"
                : "Weight in lbs, Height in ft"}
            </Text>
          </View>
        </View>
      </ScrollView>
    </View>
  );
}
