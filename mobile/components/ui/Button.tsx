import { Pressable, Text, PressableProps } from "react-native";

export type ButtonVariant = "primary" | "ghost" | "destructive";

interface Props extends PressableProps {
  variant?: ButtonVariant;
  label: string;
}

// Only variant-specific classes (background/active state and text color)
const variantClasses: Record<ButtonVariant, { container: string; text: string }> = {
  primary: {
    container: "bg-emerald-500 active:bg-emerald-600",
    text: "text-white",
  },
  ghost: {
    container: "bg-stone-200 dark:bg-stone-800 active:opacity-70",
    text: "text-stone-950 dark:text-stone-50",
  },
  destructive: {
    container: "bg-red-500 active:bg-red-600",
    text: "text-white",
  },
};

const disabledClasses = {
  container: "bg-stone-300 dark:bg-stone-800",
  text: "text-stone-500",
};

// Shared layout classes applied directly on the elements
const sharedContainer = "rounded-2xl px-4 py-3 items-center justify-center";
const sharedText = "font-semibold text-base";

export function Button({ variant = "primary", label, disabled, className, ...props }: Props) {
  const classes = disabled ? disabledClasses : variantClasses[variant];

  return (
    <Pressable
      disabled={disabled}
      className={`${sharedContainer} ${classes.container}${className ? ` ${className}` : ""}`}
      {...props}
    >
      <Text className={`${sharedText} ${classes.text}`}>{label}</Text>
    </Pressable>
  );
}
