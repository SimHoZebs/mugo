import { View } from "react-native";
import { MacroDisplay } from "@/components/MacroDisplay";
import { MealLog } from "@/lib/types";
import { Text } from "@/components/ui/Text";

interface TotalMarcoPanelProps {
  meals: MealLog[];
}

export default function TotalMacroPanel(props: TotalMarcoPanelProps) {
  const meals = props.meals;

  const { totalCalories, totalProtein, totalCarbs, totalFat } = meals.reduce(
    (acc, meal) => {
      acc.totalCalories += meal.nutrition.macros.calories;
      acc.totalProtein += meal.nutrition.macros.protein;
      acc.totalCarbs += meal.nutrition.macros.carbs;
      acc.totalFat += meal.nutrition.macros.fat;
      return acc;
    },
    { totalCalories: 0, totalProtein: 0, totalCarbs: 0, totalFat: 0 },
  );

  return (
    <View className="p-4 w-full border-b border-stone-200 dark:border-stone-700 rounded-lg bg-stone-50 dark:bg-stone-950">
      <Text variant="micro" className="mb-3">
        Today&apos;s Total
      </Text>

      <View className="flex-row justify-evenly">
        <MacroDisplay
          variant="column"
          label="Calories"
          value={totalCalories}
          unit=""
          colorClass="bg-amber-500"
        />
        <MacroDisplay
          variant="column"
          label="Protein"
          value={totalProtein}
          unit="g"
          colorClass="bg-emerald-500"
        />
        <MacroDisplay
          variant="column"
          label="Carbs"
          value={totalCarbs}
          unit="g"
          colorClass="bg-blue-500"
        />
        <MacroDisplay
          variant="column"
          label="Fat"
          value={totalFat}
          unit="g"
          colorClass="bg-violet-500"
        />
      </View>
    </View>
  );
}
