import { Text as RNText, TextProps } from "react-native";

export type TextVariant = "h1" | "h2" | "h3" | "body" | "caption" | "micro";

const variantClasses: Record<TextVariant, string> = {
  h1: "text-2xl font-bold leading-8 text-stone-950 dark:text-stone-50",
  h2: "text-xl font-bold text-stone-950 dark:text-stone-50",
  h3: "text-base font-semibold leading-6 text-stone-950 dark:text-stone-50",
  body: "text-base leading-6 text-stone-950 dark:text-stone-50",
  caption: "text-sm text-stone-500 dark:text-stone-400",
  micro: "text-xs text-stone-500 dark:text-stone-400",
};

interface Props extends TextProps {
  variant?: TextVariant;
}

export function Text({ variant = "body", className, ...props }: Props) {
  const base = variantClasses[variant];
  return (
    <RNText
      className={className ? `${base} ${className}` : base}
      {...props}
    />
  );
}
