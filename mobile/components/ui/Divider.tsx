import { View, ViewProps } from "react-native";

export function Divider({ className, ...props }: ViewProps) {
  const base = "h-px bg-stone-200 dark:bg-stone-700";
  return (
    <View
      className={className ? `${base} ${className}` : base}
      {...props}
    />
  );
}
