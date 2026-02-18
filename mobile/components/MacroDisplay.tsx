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
  const dot = (
    <View
      className={`rounded-full ${props.colorClass} ${props.variant !== "column" ? "w-3 h-3" : "w-2 h-2"}`}
    />
  );

  return props.variant === "column" ? (
    <View className="items-center flex-1">
      <View className="flex-row items-center gap-1.5">
        {dot}
        <ThemedText className="text-xs text-stone-500 dark:text-stone-400">
          {props.label}
        </ThemedText>
      </View>
      <ThemedText type="defaultSemiBold">
        {Math.round(props.value)}
        {props.unit}
      </ThemedText>
    </View>
  ) : (
    <View className="flex-row items-center justify-between py-3 border-b border-stone-200 dark:border-stone-700">
      <View className="flex-row items-center gap-3">
        {dot}
        <ThemedText className="text-base">{props.label}</ThemedText>
      </View>
      <ThemedText type="defaultSemiBold">
        {Math.round(props.value)} {props.unit}
      </ThemedText>
    </View>
  );
}

