import { View } from "react-native";
import { Text } from "@/components/ui/Text";

interface MacroDisplayProps {
  label: string;
  value: number;
  unit: string;
  colorClass: string;
  variant?: "row" | "column";
}

export function MacroDisplay(props: MacroDisplayProps) {
  const isColumn = props.variant === "column";

  return (
    <View
      className={
        isColumn
          ? "items-center flex-1"
          : "flex-row items-center justify-between py-3 border-b border-stone-200 dark:border-stone-700"
      }
    >
      <View className={`flex-row items-center ${isColumn ? "gap-1.5" : "gap-3"}`}>
        <View
          className={`rounded-full ${props.colorClass} ${isColumn ? "w-2 h-2" : "w-3 h-3"}`}
        />
        <Text variant={isColumn ? "micro" : "body"}>
          {props.label}
        </Text>
      </View>
      <Text variant="h3">
        {isColumn
          ? `${Math.round(props.value)}${props.unit}`
          : `${Math.round(props.value)} ${props.unit}`}
      </Text>
    </View>
  );
}

