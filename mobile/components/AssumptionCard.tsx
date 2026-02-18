import { View, Text } from "react-native";
import { Assumption } from "@/lib/api/conversationAPI.schemas";

export function AssumptionCard({ assumption }: { assumption: Assumption }) {
  return (
    <View className="p-3 bg-stone-100 dark:bg-stone-800 rounded-lg mb-2">
      <View className="flex-row justify-between items-start mb-1">
        <Text className="text-base leading-6 font-semibold text-stone-950 dark:text-stone-50">
          {assumption.field || "Unknown field"}
        </Text>
        {assumption.confidence && (
          <View className="px-2 py-0.5 bg-amber-100 dark:bg-amber-900 rounded">
            <Text className="text-xs text-amber-700 dark:text-amber-300">
              {assumption.confidence}
            </Text>
          </View>
        )}
      </View>
      <Text className="text-stone-600 dark:text-stone-400">
        Assumed: {assumption.assumed_value} {assumption.unit || ""}
      </Text>
      {assumption.rationale && (
        <Text className="text-sm text-stone-500 mt-1">
          {assumption.rationale}
        </Text>
      )}
    </View>
  );
}