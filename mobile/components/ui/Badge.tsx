import { View, Text, ViewProps } from "react-native";

export type BadgeVariant = "default" | "warning" | "success";

interface Props extends ViewProps {
  label: string;
  variant?: BadgeVariant;
}

// Only variant-specific color classes
const variantClasses: Record<BadgeVariant, { container: string; text: string }> = {
  default: {
    container: "bg-stone-100 dark:bg-stone-800",
    text: "text-stone-600 dark:text-stone-400",
  },
  warning: {
    container: "bg-amber-100 dark:bg-amber-900",
    text: "text-amber-700 dark:text-amber-300",
  },
  success: {
    container: "bg-emerald-100 dark:bg-emerald-900",
    text: "text-emerald-700 dark:text-emerald-300",
  },
};

// Shared layout classes applied directly on the elements
const sharedContainer = "px-2 py-0.5 rounded";
const sharedText = "text-xs";

export function Badge({ label, variant = "default", className, ...props }: Props) {
  const classes = variantClasses[variant];

  return (
    <View
      className={`${sharedContainer} ${classes.container}${className ? ` ${className}` : ""}`}
      {...props}
    >
      <Text className={`${sharedText} ${classes.text}`}>{label}</Text>
    </View>
  );
}
