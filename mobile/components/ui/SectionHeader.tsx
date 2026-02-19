import { View, Pressable, PressableProps, ViewProps } from "react-native";
import { theme } from "@/lib/theme";
import { Text } from "@/components/ui/Text";

interface Props extends ViewProps {
  title: string;
  /** Optional label for the trailing action button */
  actionLabel?: string;
  /** Called when the trailing action is pressed */
  onAction?: PressableProps["onPress"];
}

/**
 * A section header row with a title on the left and an optional
 * action link/button on the right.
 */
export function SectionHeader({ title, actionLabel, onAction, className, ...props }: Props) {
  const base = "flex-row items-center justify-between mb-3";
  return (
    <View className={className ? `${base} ${className}` : base} {...props}>
      <Text variant="h2">{title}</Text>
      {actionLabel && onAction && (
        <Pressable onPress={onAction} className="active:opacity-70">
          <Text className={`text-sm font-medium ${theme.color.primaryText}`}>
            {actionLabel}
          </Text>
        </Pressable>
      )}
    </View>
  );
}
