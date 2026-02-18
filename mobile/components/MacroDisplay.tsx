import { View } from "react-native";
import { ThemedText } from "@/components/themed-text";

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
        <ThemedText
          className={
            isColumn ? "text-xs text-stone-500 dark:text-stone-400" : "text-base"
          }
        >
          {props.label}
        </ThemedText>
      </View>
      <ThemedText type="defaultSemiBold">
        {isColumn
          ? `${Math.round(props.value)}${props.unit}`
          : `${Math.round(props.value)} ${props.unit}`}
      </ThemedText>
    </View>
  );
}

