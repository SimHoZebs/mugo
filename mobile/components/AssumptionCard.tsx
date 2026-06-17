import { View } from "react-native";
import { Assumption } from "@/lib/api/mugoAPI.schemas";
import { Text } from "@/components/ui/Text";
import { Badge } from "@/components/ui/Badge";

export function AssumptionCard({ assumption }: { assumption: Assumption }) {
  return (
    <View className="p-3 bg-stone-100 dark:bg-stone-800 rounded-lg mb-2">
      <View className="flex-row justify-between items-start mb-1">
        <Text variant="h3">{assumption.field || "Unknown field"}</Text>
        {assumption.confidence && (
          <Badge variant="warning" label={assumption.confidence} />
        )}
      </View>
      <Text variant="caption">
        Assumed: {assumption.assumed_value} {assumption.unit || ""}
      </Text>
      {assumption.rationale && (
        <Text variant="caption" className="mt-1">
          {assumption.rationale}
        </Text>
      )}
    </View>
  );
}
