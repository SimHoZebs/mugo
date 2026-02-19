import { Text as RNText, TextProps } from "react-native";
import { theme } from "@/lib/theme";

export type TextVariant = "h1" | "h2" | "h3" | "body" | "caption" | "micro";

const variantClasses: Record<TextVariant, string> = {
  h1: `text-2xl font-bold leading-8 ${theme.color.textPrimary}`,
  h2: `text-xl font-bold ${theme.color.textPrimary}`,
  h3: `text-base font-semibold leading-6 ${theme.color.textPrimary}`,
  body: `text-base leading-6 ${theme.color.textPrimary}`,
  caption: `text-sm ${theme.color.textSecondary}`,
  micro: `text-xs ${theme.color.textSecondary}`,
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
