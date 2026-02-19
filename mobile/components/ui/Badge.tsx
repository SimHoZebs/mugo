import { View, Text, ViewProps } from "react-native";

export type BadgeVariant = "default" | "warning" | "success";

interface Props extends ViewProps {
  label: string;
  variant?: BadgeVariant;
}

const variantClasses: Record<BadgeVariant, { container: string; text: string }> = {
  default: {
    container: "px-2 py-0.5 bg-stone-100 dark:bg-stone-800 rounded",
    text: "text-xs text-stone-600 dark:text-stone-400",
  },
  warning: {
    container: "px-2 py-0.5 bg-amber-100 dark:bg-amber-900 rounded",
    text: "text-xs text-amber-700 dark:text-amber-300",
  },
  success: {
    container: "px-2 py-0.5 bg-emerald-100 dark:bg-emerald-900 rounded",
    text: "text-xs text-emerald-700 dark:text-emerald-300",
  },
};

export function Badge({ label, variant = "default", className, ...props }: Props) {
  const classes = variantClasses[variant];

  return (
    <View
      className={className ? `${classes.container} ${className}` : classes.container}
      {...props}
    >
      <Text className={classes.text}>{label}</Text>
    </View>
  );
}
