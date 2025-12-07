import { View } from "react-native";
import { ThemedText } from "@/components/themed-text";
import { ThemedView } from "@/components/themed-view";
import { Macros } from "@/lib/api/conversationAPI.schemas";

interface Props {
  macros?: Macros;
}

interface MacroItemProps {
  label: string;
  value: number;
  unit: string;
  dotColor: string;
}

function MacroItem({ label, value, unit, dotColor }: MacroItemProps) {
  return (
    <View className="items-center">
      <View className="flex-row items-center gap-1.5 mb-0.5">
        <View className={`w-2 h-2 rounded-full ${dotColor}`} />
        <ThemedText className="text-xs text-stone-500 dark:text-stone-400">
          {label}
        </ThemedText>
      </View>
      <ThemedText type="defaultSemiBold">
        {Math.round(value)}{unit}
      </ThemedText>
    </View>
  );
}

export default function TotalMacroPanel({ macros }: Props) {
  const calories = macros?.calories ?? 0;
  const protein = macros?.protein ?? 0;
  const carbs = macros?.carbs ?? 0;
  const fat = macros?.fat ?? 0;

  return (
    <ThemedView className="p-4 w-full border-b border-stone-200 dark:border-stone-700 rounded-lg">
      <ThemedText className="text-xs text-stone-500 dark:text-stone-400 mb-3">
        Today's Total
      </ThemedText>
      <View className="flex-row justify-between">
        <MacroItem 
          label="Calories" 
          value={calories} 
          unit="" 
          dotColor="bg-amber-500" 
        />
        <MacroItem 
          label="Protein" 
          value={protein} 
          unit="g" 
          dotColor="bg-emerald-500" 
        />
        <MacroItem 
          label="Carbs" 
          value={carbs} 
          unit="g" 
          dotColor="bg-blue-500" 
        />
        <MacroItem 
          label="Fat" 
          value={fat} 
          unit="g" 
          dotColor="bg-violet-500" 
        />
      </View>
    </ThemedView>
  );
}
