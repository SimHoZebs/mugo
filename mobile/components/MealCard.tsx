import { Pressable, View } from "react-native";
import { useRouter } from "expo-router";
import { ThemedText } from "./themed-text";
import { ThemedView } from "./themed-view";
import { Meal } from "@/lib/types";

interface MealCardProps {
  meal?: Meal;
  loading?: boolean;
}

interface MacroPillProps {
  value: number;
  unit: string;
  colorClass: string;
  loading?: boolean;
  className?: string;
}

function MacroPill({
  value,
  unit,
  colorClass,
  loading,
  className,
}: MacroPillProps) {
  const isLoading = Boolean(loading);
  const containerClassName = [
    "rounded-full",
    isLoading ? "p-0" : "px-2 py-1",
    isLoading ? "" : colorClass,
    className ?? "",
  ]
    .join(" ")
    .trim();

  return (
    <View className={containerClassName}>
      <ThemedText
        loading={isLoading}
        className={
          isLoading
            ? "h-6 w-full rounded-full"
            : "text-xs font-medium text-white"
        }
      >
        {Math.round(value)}
        {unit}
      </ThemedText>
    </View>
  );
}

export default function MealCard({ meal, loading }: MealCardProps) {
  const router = useRouter();

  const handlePress = () => {
    router.push({
      pathname: "/meal-detail",
      params: { id: meal.id },
    });
  };

  const isLoading = Boolean(loading);
  const isReady = Boolean(meal && !loading);
  const isVisible = isLoading || isReady;
  const macros = isReady ? meal.nutrition.macros : null;
  const assumptions = isReady ? meal.nutrition.assumptions : null;
  const hasAssumptions = Boolean(assumptions && assumptions.length > 0);

  return (
    <Pressable
      onPress={isReady ? handlePress : undefined}
      disabled={!isReady}
      className={isVisible ? undefined : "hidden"}
    >
      <ThemedView className="p-4 border border-stone-300 dark:border-stone-700 rounded-xl">
        <View
          className={`flex-row justify-between items-start ${
            isLoading ? "mb-3" : "mb-1"
          }`}
        >
          {isLoading && !isReady ? (
            <ThemedText loading className="h-5 w-2/3" />
          ) : (
            <ThemedText type="defaultSemiBold" className="flex-1">
              {meal.nutrition.name}
            </ThemedText>
          )}
          {isLoading && !isReady && !hasAssumptions ? (
            <ThemedText loading className="h-4 w-16" />
          ) : (
            <View className="px-2 py-0.5 bg-amber-100 dark:bg-amber-900/50 rounded">
              <ThemedText className="text-xs text-amber-700 dark:text-amber-300">
                {assumptions.length} assumption
                {assumptions.length > 1 ? "s" : ""}
              </ThemedText>
            </View>
          )}
        </View>

        {(isLoading || macros) && (
          <View className={`flex-row flex-wrap gap-2 ${macros ? "mt-2" : ""}`}>
            {isLoading && (
              <>
                <MacroPill
                  value={0}
                  unit=""
                  colorClass=""
                  loading
                  className="w-16"
                />
                <MacroPill
                  value={0}
                  unit=""
                  colorClass=""
                  loading
                  className="w-14"
                />
                <MacroPill
                  value={0}
                  unit=""
                  colorClass=""
                  loading
                  className="w-14"
                />
                <MacroPill
                  value={0}
                  unit=""
                  colorClass=""
                  loading
                  className="w-14"
                />
              </>
            )}
            {isReady && macros && (
              <>
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
              </>
            )}
          </View>
        )}
      </ThemedView>
    </Pressable>
  );
}
