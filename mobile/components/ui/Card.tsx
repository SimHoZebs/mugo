import { View, ViewProps } from "react-native";
import { theme } from "@/lib/theme";

interface Props extends ViewProps {
  /** Remove the default border. Useful for flush surfaces. */
  noBorder?: boolean;
}

export function Card({ noBorder = false, className, children, ...props }: Props) {
  const base = [
    `p-4 ${theme.radius.md} ${theme.color.surface}`,
    noBorder ? "" : `border ${theme.color.border}`,
  ]
    .join(" ")
    .trim();

  return (
    <View className={className ? `${base} ${className}` : base} {...props}>
      {children}
    </View>
  );
}
