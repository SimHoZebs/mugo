import { Pressable, Text, PressableProps } from "react-native";

export type ButtonVariant = "primary" | "ghost" | "destructive";

interface Props extends PressableProps {
  variant?: ButtonVariant;
  label: string;
}

const variantClasses: Record<ButtonVariant, { container: string; text: string }> = {
  primary: {
    container:
      "bg-emerald-500 active:bg-emerald-600 rounded-2xl px-4 py-3 items-center justify-center",
    text: "text-white font-semibold text-base",
  },
  ghost: {
    container:
      "bg-stone-200 dark:bg-stone-800 active:opacity-70 rounded-2xl px-4 py-3 items-center justify-center",
    text: "text-stone-950 dark:text-stone-50 font-semibold text-base",
  },
  destructive: {
    container:
      "bg-red-500 active:bg-red-600 rounded-2xl px-4 py-3 items-center justify-center",
    text: "text-white font-semibold text-base",
  },
};

const disabledClasses = {
  container: "bg-stone-300 dark:bg-stone-800 rounded-2xl px-4 py-3 items-center justify-center",
  text: "text-stone-500 font-semibold text-base",
};

export function Button({ variant = "primary", label, disabled, className, ...props }: Props) {
  const classes = disabled ? disabledClasses : variantClasses[variant];

  return (
    <Pressable
      disabled={disabled}
      className={className ? `${classes.container} ${className}` : classes.container}
      {...props}
    >
      <Text className={classes.text}>{label}</Text>
    </Pressable>
  );
}
