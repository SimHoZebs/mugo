import { View, ViewProps } from "react-native";

interface Props extends ViewProps {
  /** Remove the default border. Useful for flush surfaces. */
  noBorder?: boolean;
}

export function Card({ noBorder = false, className, children, ...props }: Props) {
  const base = [
    "p-4 rounded-xl bg-white dark:bg-stone-900",
    noBorder ? "" : "border border-stone-200 dark:border-stone-700",
  ]
    .join(" ")
    .trim();

  return (
    <View className={className ? `${base} ${className}` : base} {...props}>
      {children}
    </View>
  );
}
