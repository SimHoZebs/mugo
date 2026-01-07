import { View } from "react-native";
import { ThemedText } from "@/components/themed-text";
import { ThemedView } from "@/components/themed-view";
import { NutritionPayload } from "@/lib/api/conversationAPI.schemas";
import { Meal } from "@/lib/types";

interface MacroItemProps {
  label: string;
  value: number;
  unit: string;
  dotColor: string;
}

function MacroItem(props: MacroItemProps) {
  return (
    <View className="items-center flex-1">
      <View className="flex-row items-center gap-1.5">
        <View className={`w-2 h-2 rounded-full ${props.dotColor}`} />
        <ThemedText className="text-xs text-stone-500 dark:text-stone-400">
          {props.label}
        </ThemedText>
      </View>
      <ThemedText type="defaultSemiBold">
        {Math.round(props.value)}
        {props.unit}
      </ThemedText>
    </View>
  );
}

interface TotalMarcoPanelProps {
  meals: Meal[];
}

export default function TotalMacroPanel(props: TotalMarcoPanelProps) {
  const meals = props.meals;

  const totalCalories = meals.reduce(
    (sum, meal) => sum + meal.nutrition.macros.calories,
    0,
  );
  const totalProtein = meals.reduce(
    (sum, meal) => sum + meal.nutrition.macros.protein,
    0,
  );
  const totalCarbs = meals.reduce(
    (sum, meal) => sum + meal.nutrition.macros.carbs,
    0,
  );
  const totalFat = meals.reduce(
    (sum, meal) => sum + meal.nutrition.macros.fat,
    0,
  );

  return (
    <ThemedView className="p-4 w-full border-b border-stone-200 dark:border-stone-700 rounded-lg">
      <ThemedText className="text-xs text-stone-500 dark:text-stone-400 mb-3">
        Today&apos;s Total
      </ThemedText>

      <View className="flex-row justify-evenly">
        <MacroItem
          label="Calories"
          value={totalCalories}
          unit=""
          dotColor="bg-amber-500"
        />
        <MacroItem
          label="Protein"
          value={totalProtein}
          unit="g"
          dotColor="bg-emerald-500"
        />
        <MacroItem
          label="Carbs"
          value={totalCarbs}
          unit="g"
          dotColor="bg-blue-500"
        />
        <MacroItem
          label="Fat"
          value={totalFat}
          unit="g"
          dotColor="bg-violet-500"
        />
      </View>
    </ThemedView>
  );
}
