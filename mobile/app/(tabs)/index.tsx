import {
  KeyboardAwareScrollView,
  KeyboardStickyView,
} from "react-native-keyboard-controller";
import { useState } from "react";
import "react-native-get-random-values";
import { v7 as uuid7 } from "uuid";

import { ThemedText } from "@/components/themed-text";
import { ThemedView } from "@/components/themed-view";
import TotalMacroPanel from "@/components/TotalMacroPanel";
import MealCard from "@/components/MealCard";
import InputBar from "@/components/InputBar";
import { postAgentsNutrition } from "@/lib/api/default/default";
import { NutritionResponseBody } from "@/lib/api/conversationAPI.schemas";
import { ScrollView } from "react-native";
import { defaultMeals } from "@/constants/defaultMeals";

export default function HomeScreen() {
  const [sessionId] = useState(() => uuid7());
  const [meals, setMeals] =
    useState<Array<NutritionResponseBody["analysis"]>>(defaultMeals);

  const handleSubmitNutrition = async (text: string) => {
    try {
      const response = await postAgentsNutrition({
        text,
        session_id: sessionId,
        user_id: "user-1",
      });
      if (response.status !== 200) {
        throw new Error(`Server error: ${response.data}`);
      }
      const nutritionData = response.data.analysis;
      setMeals((prevMeals) => [nutritionData, ...prevMeals]);
    } catch (error) {
      console.error("Error submitting nutrition:", error);
    }
  };

  return (
    <ThemedView className="h-full w-full gap-4 pt-8">
      <ThemedText className="px-4" type="title">
        Tuesday
      </ThemedText>

      <TotalMacroPanel />
      <KeyboardAwareScrollView>
        <ThemedView className="gap-4">
          {meals.map((meal, index) => (
            <MealCard key={index} mealData={meal} />
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
