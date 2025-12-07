import { View } from "react-native";
import { ThemedText } from "@/components/themed-text";
import { Assumption } from "@/lib/api/conversationAPI.schemas";

export function AssumptionCard({ assumption }: { assumption: Assumption }) {
  return (
    <View className="p-3 bg-stone-100 dark:bg-stone-800 rounded-lg mb-2">
      <View className="flex-row justify-between items-start mb-1">
        <ThemedText type="defaultSemiBold">
          {assumption.field || "Unknown field"}
        </ThemedText>
        {assumption.confidence && (
          <View className="px-2 py-0.5 bg-amber-100 dark:bg-amber-900 rounded">
            <ThemedText className="text-xs text-amber-700 dark:text-amber-300">
              {assumption.confidence}
            </ThemedText>
          </View>
        )}
      </View>
      <ThemedText className="text-stone-600 dark:text-stone-400">
        Assumed: {assumption.assumed_value} {assumption.unit || ""}
      </ThemedText>
      {assumption.rationale && (
        <ThemedText className="text-sm text-stone-500 dark:text-stone-500 mt-1">
          {assumption.rationale}
        </ThemedText>
      )}
    </View>
  );
}