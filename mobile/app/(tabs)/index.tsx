import {
  KeyboardAwareScrollView,
  KeyboardStickyView,
} from "react-native-keyboard-controller";
import "react-native-get-random-values";
import { v7 as uuid7 } from "uuid";

import { useRef, useState } from "react";
import { ScrollView, View } from "react-native";

import TotalMacroPanel from "@/components/TotalMacroPanel";
import MealCard from "@/components/MealCard";
import InputBar from "@/components/InputBar";
import { Text } from "@/components/ui/Text";
import useGlobalStore from "@/lib/store";
import { createMealsBatch } from "@/lib/api/logs/logs";
import type { CreateMealResponseBody } from "@/lib/api/mugoAPI.schemas";
import { toLocalMealLog } from "@/lib/types";

export default function HomeScreen() {
  const meals = useGlobalStore((state) => state.meals);
  const setMeals = useGlobalStore((state) => state.setMeals);
  const [pendingMealId, setPendingMealId] = useState<string | null>(null);
  const scrollViewRef = useRef<ScrollView | null>(null);

  const handleSubmitNutrition = async (text: string) => {
    const newSessionId = uuid7();
    const loadingMealId = `loading-${newSessionId}`;

    const payload = {
      description: text,
      session_id: newSessionId,
      user_id: "user-1",
    };

    try {
      setPendingMealId(loadingMealId);
      requestAnimationFrame(() =>
        scrollViewRef.current?.scrollToEnd({ animated: true }),
      );
      const response = await createMealsBatch(payload);

      if (response.status < 200 || response.status >= 300) {
        throw response.data;
      }

      const data = response.data as CreateMealResponseBody;
      const sessionId = data.session_id || newSessionId;
      const newMealLogs = (data.meals ?? []).map((meal) =>
        toLocalMealLog(meal, sessionId),
      );

      setMeals([...useGlobalStore.getState().meals, ...newMealLogs]);
      setPendingMealId(null);
    } catch (error) {
      console.error("Error submitting request to:", getCreateMealLogUrl());
      console.error("Error:", error);
      setPendingMealId(null);
    }
  };

  return (
    <View className="h-full w-full gap-4 pt-8 bg-stone-50 dark:bg-stone-950">
      <Text variant="h1" className="px-4">
        Tuesday
      </Text>

      <TotalMacroPanel meals={meals} />

      <KeyboardAwareScrollView ref={scrollViewRef}>
        <View className="gap-4 px-4 bg-stone-50 dark:bg-stone-950">
          {meals.map((meal) => (
            <MealCard key={meal.id} meal={meal} />
          ))}
          {pendingMealId && <MealCard key={pendingMealId} loading />}
        </View>
      </KeyboardAwareScrollView>

      <KeyboardStickyView>
        <InputBar onSubmit={handleSubmitNutrition}>
          <InputBar.Action onPress={() => console.log("Mic pressed")}>
            <Text className="text-lg">🎤</Text>
          </InputBar.Action>
          <InputBar.Input placeholder="Log what you ate..." />
          <InputBar.Action onPress={() => console.log("Camera pressed")}>
            <Text className="text-lg">📷</Text>
          </InputBar.Action>
        </InputBar>
      </KeyboardStickyView>
    </View>
  );
}
