import { Text as RNText, TextProps } from "react-native";

export type TextVariant = "h1" | "h2" | "h3" | "body" | "caption" | "micro";

// Only size/weight — color is applied at component level
const variantClasses: Record<TextVariant, string> = {
  h1: "text-2xl font-bold leading-8",
  h2: "text-xl font-bold",
  h3: "text-base font-semibold leading-6",
  body: "text-base leading-6",
  caption: "text-sm",
  micro: "text-xs",
};

const secondaryVariants: Set<TextVariant> = new Set(["caption", "micro"]);

interface Props extends TextProps {
  variant?: TextVariant;
}

export function Text({ variant = "body", className, ...props }: Props) {
  const color = secondaryVariants.has(variant)
    ? "text-stone-500 dark:text-stone-400"
    : "text-stone-950 dark:text-stone-50";
  const base = `${variantClasses[variant]} ${color}`;
  return (
    <RNText
      className={className ? `${base} ${className}` : base}
      {...props}
    />
  );
}
