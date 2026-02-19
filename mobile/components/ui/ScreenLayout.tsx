import { SafeAreaView } from "react-native-safe-area-context";
import { ViewProps } from "react-native";

interface Props extends ViewProps {
  children: React.ReactNode;
}

/**
 * Standard page wrapper that provides:
 * - Full-height background (`bg-stone-50 dark:bg-stone-950`)
 * - Safe area insets
 * - Horizontal page padding (`px-4`)
 */
export function ScreenLayout({ className, children, ...props }: Props) {
  const base = "flex-1 bg-stone-50 dark:bg-stone-950 px-4";
  return (
    <SafeAreaView
      className={className ? `${base} ${className}` : base}
      {...props}
    >
      {children}
    </SafeAreaView>
  );
}
