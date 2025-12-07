import { Pressable, View } from "react-native";
import { useRouter } from "expo-router";
import { NutritionPayload } from "@/lib/api/conversationAPI.schemas";
import { ThemedText } from "./themed-text";
import { ThemedView } from "./themed-view";

interface Props {
  title?: string;
  description?: string;
  mealData?: NutritionPayload;
}

interface MacroPillProps {
  value: number;
  unit: string;
  colorClass: string;
}

function MacroPill({ value, unit, colorClass }: MacroPillProps) {
  return (
    <View className={`px-2 py-1 rounded-full ${colorClass}`}>
      <ThemedText className="text-xs font-medium text-white">
        {Math.round(value)}
        {unit}
      </ThemedText>
    </View>
  );
}

export default function MealCard({ mealData }: Props) {
  const router = useRouter();

  const handlePress = () => {
    router.push({
      pathname: "/meal-detail",
      params: {
        title: mealData?.name || "Meal Detail",
        calories: mealData?.macros?.calories?.toString() || "0",
        protein: mealData?.macros?.protein?.toString() || "0",
        carbs: mealData?.macros?.carbs?.toString() || "0",
        fat: mealData?.macros?.fat?.toString() || "0",
        assumptions: JSON.stringify(mealData?.assumptions || []),
      },
    });
  };

  const macros = mealData?.macros;
  const assumptions = mealData?.assumptions;
  const hasAssumptions = assumptions && assumptions.length > 0;

  return (
    <Pressable onPress={handlePress} className="mb-6">
      <ThemedView className="p-4 border border-stone-300 dark:border-stone-700 rounded-xl">
        <View className="flex-row justify-between items-start mb-1">
          <ThemedText type="defaultSemiBold" className="flex-1">
            {mealData?.name || "Unnamed Meal"}
          </ThemedText>
          {hasAssumptions && (
            <View className="px-2 py-0.5 bg-amber-100 dark:bg-amber-900/50 rounded">
              <ThemedText className="text-xs text-amber-700 dark:text-amber-300">
                {assumptions.length} assumption
                {assumptions.length > 1 ? "s" : ""}
              </ThemedText>
            </View>
          )}
        </View>

        {macros && (
          <View className="flex-row flex-wrap gap-2 mt-2">
            <MacroPill
              value={macros.calories}
              unit=" kcal"
              colorClass="bg-amber-500"
            />
            <MacroPill
              value={macros.protein}
              unit="g P"
              colorClass="bg-emerald-500"
            />
            <MacroPill
              value={macros.carbs}
              unit="g C"
              colorClass="bg-blue-500"
            />
            <MacroPill
              value={macros.fat}
              unit="g F"
              colorClass="bg-violet-500"
            />
          </View>
        )}
      </ThemedView>
    </Pressable>
  );
}
