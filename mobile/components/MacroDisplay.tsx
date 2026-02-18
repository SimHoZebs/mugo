import { View, Text } from "react-native";

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
        <Text
          className={
            isColumn ? "text-xs text-stone-500 dark:text-stone-400" : "text-base text-stone-950 dark:text-stone-50"
          }
        >
          {props.label}
        </Text>
      </View>
      <Text className="text-base leading-6 font-semibold text-stone-950 dark:text-stone-50">
        {isColumn
          ? `${Math.round(props.value)}${props.unit}`
          : `${Math.round(props.value)} ${props.unit}`}
      </Text>
    </View>
  );
}

