import { useState, useEffect } from "react";
import { View, ScrollView } from "react-native";

import { Text } from "@/components/ui/Text";
import { Button } from "@/components/ui/Button";
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
      <Text variant="caption" className="text-center">
        Please sign in on the Profile tab to access settings.
      </Text>
    </View>
  ) : (
    <View className="h-full w-full pt-8 bg-stone-50 dark:bg-stone-950">
      <View className="px-4 mb-6">
        <Text variant="h1">Settings</Text>
      </View>

      <ScrollView className="flex-1 px-4">
        <View className="gap-6 bg-stone-50 dark:bg-stone-950">
          <View className="bg-stone-50 dark:bg-stone-950">
            <Text variant="caption" className="font-semibold mb-2">
              Unit System
            </Text>
            <View className="flex-row gap-2">
              {(["metric", "imperial"] as UnitSystem[]).map((system) => (
                <Button
                  key={system}
                  variant={selectedUnit === system ? "primary" : "ghost"}
                  label={system.charAt(0).toUpperCase() + system.slice(1)}
                  onPress={() => handleUnitChange(system)}
                  className="flex-1"
                />
              ))}
            </View>
            <Text variant="caption" className="mt-2">
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
