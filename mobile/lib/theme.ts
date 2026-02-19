// Semantic design tokens — map intent to NativeWind/Tailwind class strings.
// Reference these instead of raw Tailwind classes when building UI.
export const theme = {
  color: {
    // Surfaces
    background: "bg-stone-50 dark:bg-stone-950",
    surface: "bg-white dark:bg-stone-900",
    surfaceSubtle: "bg-stone-100 dark:bg-stone-800",
    border: "border-stone-200 dark:border-stone-700",

    // Text
    textPrimary: "text-stone-950 dark:text-stone-50",
    textSecondary: "text-stone-500 dark:text-stone-400",
    textInverse: "text-white",

    // Brand/Action
    primary: "bg-emerald-500",
    primaryActive: "bg-emerald-600",
    primaryText: "text-emerald-600 dark:text-emerald-400",

    // Semantic
    warning: "bg-amber-100 dark:bg-amber-900",
    warningText: "text-amber-700 dark:text-amber-300",

    // Macros (for NutritionDisplay / MealCard pills)
    calories: "bg-amber-500",
    protein: "bg-emerald-500",
    carbs: "bg-blue-500",
    fat: "bg-rose-500",
  },
  radius: {
    sm: "rounded-lg",
    md: "rounded-xl",
    lg: "rounded-2xl",
    full: "rounded-full",
  },
  text: {
    h1: "text-2xl font-bold leading-8",
    h2: "text-xl font-bold",
    h3: "text-base font-semibold leading-6",
    body: "text-base leading-6",
    caption: "text-sm",
    micro: "text-xs",
  },
} as const;
