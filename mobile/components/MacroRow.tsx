import { View } from "react-native";
import { ThemedText } from "@/components/themed-text";

interface MacroRowProps {
  label: string;
  value: number;
  unit: string;
  colorClass: string;
}

export function MacroRow({ label, value, unit, colorClass }: MacroRowProps) {
  return (
    <View className="flex-row items-center justify-between py-3 border-b border-stone-200 dark:border-stone-700">
      <View className="flex-row items-center gap-3">
        <View className={`w-3 h-3 rounded-full ${colorClass}`} />
        <ThemedText className="text-base">{label}</ThemedText>
      </View>
      <ThemedText type="defaultSemiBold">
        {Math.round(value)} {unit}
      </ThemedText>
    </View>
  );
}