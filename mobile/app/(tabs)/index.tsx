import {
  KeyboardAwareScrollView,
  KeyboardStickyView,
} from "react-native-keyboard-controller";
import "react-native-get-random-values";
import { v7 as uuid7 } from "uuid";

import { ThemedText } from "@/components/themed-text";
import { ThemedView } from "@/components/themed-view";
import TotalMacroPanel from "@/components/TotalMacroPanel";
import MealCard from "@/components/MealCard";
import InputBar from "@/components/InputBar";
import { postAgentsNutrition } from "@/lib/api/default/default";
import useGlobalStore from "@/lib/store";

export default function HomeScreen() {
  const meals = useGlobalStore((state) => state.meals);
  const setMeals = useGlobalStore((state) => state.setMeals);

  const handleSubmitNutrition = async (text: string) => {
    const newMealId = uuid7();
    const newSessionId = uuid7();

    try {
      const response = await postAgentsNutrition({
        text,
        session_id: newSessionId,
        user_id: "user-1",
      });
      if (response.status !== 200) {
        throw new Error(`Server error: ${response.data}`);
      }
      const newMeal = {
        id: newMealId,
        sessionId: newSessionId,
        nutrition: response.data.analysis,
      };
      setMeals([...meals, newMeal]);
    } catch (error) {
      console.error("Error submitting nutrition:", error);
    }
  };

  return (
    <ThemedView className="h-full w-full gap-4 pt-8">
      <ThemedText className="px-4" type="title">
        Tuesday
      </ThemedText>

      <TotalMacroPanel meals={meals} />

      <KeyboardAwareScrollView>
        <ThemedView className="gap-4 px-4">
          {meals.map((meal) => (
            <MealCard key={meal.id} meal={meal} />
          ))}
        </ThemedView>
      </KeyboardAwareScrollView>

      <KeyboardStickyView offset={{ closed: 0, opened: 80 }}>
        <InputBar onSubmit={handleSubmitNutrition}>
          <InputBar.Action onPress={() => console.log("Mic pressed")}>
            <ThemedText className="text-lg">🎤</ThemedText>
          </InputBar.Action>
          <InputBar.Input placeholder="Describe your meals..." />
          <InputBar.Action onPress={() => console.log("Camera pressed")}>
            <ThemedText className="text-lg">📷</ThemedText>
          </InputBar.Action>
        </InputBar>
      </KeyboardStickyView>
    </ThemedView>
  );
}
