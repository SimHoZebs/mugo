import { Pressable, View } from "react-native";
import { useRouter } from "expo-router";
import { Meal } from "@/lib/types";
import { Text } from "@/components/ui/Text";
import { Badge } from "@/components/ui/Badge";

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
      {isLoading ? (
        <View className="h-6 w-full rounded-full bg-stone-200 dark:bg-stone-700" />
      ) : (
        <Text variant="micro" className="font-medium text-white">
          {Math.round(value)}
          {unit}
        </Text>
      )}
    </View>
  );
}

export default function MealCard({ meal, loading }: MealCardProps) {
  const router = useRouter();

  const handlePress = () => {
    if (!meal) return;
    router.push({
      pathname: "/meal-detail",
      params: { id: meal.id },
    });
  };

  const isLoading = Boolean(loading);
  const isReady = Boolean(meal && !loading);
  const isVisible = isLoading || isReady;
  const macros = meal?.nutrition.macros ?? null;
  const assumptions = meal?.nutrition.assumptions ?? null;
  const hasAssumptions = Boolean(assumptions && assumptions.length > 0);

  return (
    <Pressable
      onPress={isReady ? handlePress : undefined}
      disabled={!isReady}
      className={isVisible ? undefined : "hidden"}
    >
      <View className="p-4 border border-stone-300 dark:border-stone-700 rounded-xl bg-stone-50 dark:bg-stone-950">
        <View
          className={`flex-row justify-between items-start ${
            isLoading ? "mb-3" : "mb-1"
          }`}
        >
          {!isReady ? (
            <View className="h-5 w-2/3 rounded bg-stone-200 dark:bg-stone-700" />
          ) : (
            <Text variant="h3" className="flex-1">
              {meal?.nutrition.name}
            </Text>
          )}
          {!isReady && !hasAssumptions ? (
            <View className="h-4 w-16 rounded bg-stone-200 dark:bg-stone-700" />
          ) : hasAssumptions ? (
            <Badge
              variant="warning"
              label={`${assumptions?.length} assumption${assumptions && assumptions.length > 1 ? "s" : ""}`}
            />
          ) : null}
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
      </View>
    </Pressable>
  );
}
